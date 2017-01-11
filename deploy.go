package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"path/filepath"

	"github.com/Recime/recime-cli/cmd"
	"github.com/Recime/recime-cli/util"
	"github.com/briandowns/spinner"
	pb "gopkg.in/cheggaaa/pb.v1"
)

type Bot struct {
	Id      string `json:"uid"`
	Type    string `json:"fileType"`
	Version string `json:"version"`
	Owner   string `json:"owner"`
	Config []cmd.Config `json:"config"`
}

func SendRequest(url string, body io.Reader) string {
	res, err := http.Post(url, "application/json; charset=utf-8", body)

	check(err)

	var result struct {
		Url string `json:"url"`
	}

	bytes, err := ioutil.ReadAll(res.Body)

	json.Unmarshal(bytes, &result)

	defer res.Body.Close()

	// fmt.Println(string(res.Body))

	return result.Url
}

// Deploy deploys the bot with the given uid
func Deploy(uid string) {
	wd, err := os.Getwd()

	var data map[string]interface{}

	buff, err := ioutil.ReadFile(wd + "/package.json")

	check(err)

	if err := json.Unmarshal(buff, &data); err != nil {
		panic(err)
	}

	name := data["name"].(string)

	fmt.Println("INFO: Compressing.")

	temp, err := ioutil.TempDir("", name)

	check(err)

	dest := filepath.ToSlash(temp) + "/" + uid

	err = os.Mkdir(dest, os.ModePerm)

	check(err)

	err = util.CopyDir(wd, dest)

	filePath := temp + "/" + name + ".zip"

	Archive(dest, filePath)

	file, err := os.Open(filePath)

	defer file.Close()

	fileInfo, _ := file.Stat()

	var size = fileInfo.Size()

	buffer := make([]byte, size)

	// read file content to buffer
	file.Read(buffer)

	fmt.Println("INFO: Preparing to upload.")

	url := BaseURL + "/signed-url"

	fileType := http.DetectContentType(buffer)

	user, err := cmd.GetStoredUser()

	var config []cmd.Config

	vars := cmd.GetConfigVars(wd)

	for key, value := range vars {
		config = append(config, cmd.Config{ Key : key, Value : value.(string) })
	}

	bot := Bot{
		Id: uid, 
		Type: fileType, 
		Version: Version, 
		Owner: user.Email, 
		Config : config,
	}

	jsonBody, err := json.Marshal(bot)

	check(err)

	signedUrl := SendRequest(url, bytes.NewBuffer(jsonBody))

	bar := pb.New(len(buffer)).SetUnits(pb.U_BYTES)

	bar.Format("[## ]")

	bar.Start()

	proxy := NewReader(buffer, bar)

	req, err := http.NewRequest("PUT", signedUrl, proxy)

	req.ContentLength = size

	check(err)

	// bar.Finish()

	resp, err := http.DefaultClient.Do(req)

	check(err)

	dat, err := ioutil.ReadAll(resp.Body)

	check(err)

	defer resp.Body.Close()

	fmt.Println(string(dat))

	if len(dat) == 0 {
		fmt.Println("INFO: Finalizing.")
	}

	url = BaseURL + "/module/deploy/" + uid

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

	s.Start()

	r := bytes.NewBuffer(jsonBody)

	resp, err = http.Post(url, "application/json; charset=utf-8", r)

	check(err)

	var result struct {
		Name    string `json:"name"`
		Id      string `json:"uid"`
		Message string `json:message`
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)

	check(err)

	json.Unmarshal(bytes, &result)

	s.Stop()

	if len(result.Name) > 0 {
		fmt.Println("\r\n=> " + BaseURL + "/bot/" + result.Name + "\r\n")
		fmt.Println("INFO: Bot publish successful.")
		return
	}

	fmt.Println("\x1b[31;1mFatal: Publish Failed!!!\x1b[0m")
}
