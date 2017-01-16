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

	"github.com/recime/recime-cli/cmd"
	"github.com/recime/recime-cli/util"
	"github.com/briandowns/spinner"
	pb "gopkg.in/cheggaaa/pb.v1"
	"github.com/jhoonb/archivex"
)

type Bot struct {
	Id      string `json:"uid"`
	Type    string `json:"fileType"`
	Name  string `json:"name"`
	Title string `json:"title"`
	Author string `json:"author"`
	Desc  string `json:"description"`
	Version string `json:"version"`
	Owner   string `json:"owner"`
	Config []cmd.Config `json:"config"`
}

type Package struct {
	FileType string `json:"fileType"`
	Id string `json:"uid"`
	Method string `json:"method"`
}

func PrepareLambdaPackage(uid string) string{
	jsonBody, err := json.Marshal(Package {
		FileType : "application/octet-stream",
		Id : "package.zip",
		Method : "getObject",
	})

	check(err)

	source := fmt.Sprintf("%s/signed-url", BaseURL)
	signedUrl := SendHTTPRequest(source, bytes.NewBuffer(jsonBody))
	
	temp, err := ioutil.TempDir("", "recime-cli")

	check(err)

	dest := filepath.ToSlash(temp) + "/bin"

	err = os.Mkdir(dest, os.ModePerm)

	fileName := fmt.Sprintf("%s/%s.zip", dest, uid)

	cmd.Download(signedUrl, fileName)

	target := fmt.Sprintf("%s/%s", dest, uid)

	check(util.Unzip(fileName, target))

	wd, err := os.Getwd()

	check(err)

	botDir := fmt.Sprintf("%s/%s", target, uid)

	_ = util.CopyDir(wd, botDir)

	pkg := fmt.Sprintf("%s/%s.zip", temp, uid)

	zip := new(archivex.ZipFile)
    zip.Create(pkg)
    zip.AddAll(target, true)
    zip.Close()

	return pkg
}

func SendHTTPRequest(url string, body io.Reader) string {
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
func Deploy() {
	uid := cmd.GetUID()

	fmt.Println("INFO: Preparing Package.")

	pkgPath := PrepareLambdaPackage(uid)

	file, err := os.Open(pkgPath)

	defer file.Close()

	fileInfo, _ := file.Stat()

	var size = fileInfo.Size()

	buffer := make([]byte, size)

	// // read file content to buffer
	file.Read(buffer)

	url := BaseURL + "/signed-url"

	fileType := http.DetectContentType(buffer)

	user, err := cmd.GetStoredUser()

	var config []cmd.Config

	wd, err := os.Getwd()

	check(err)
	
	// Add config user config
	reader, err := cmd.OpenConfig(wd)

	if reader != nil {
		cfg := cmd.GetConfigVars(reader)
		for key, value := range cfg {
			config = append(config, cmd.Config{ Key : key, Value : value.(string) })
		}
	}

	var data map[string]interface{}

	buff, err := ioutil.ReadFile(wd + "/package.json")

	check(err)

	if err := json.Unmarshal(buff, &data); err != nil {
		panic(err)
	}

	bot := Bot{
		Author : data["author"].(string),
		Id: uid,
		Type: fileType,
		Version: Version,
		Owner: user.Email,
		Config : config,
		Name : data["name"].(string),
		Desc : data["description"].(string),
		Title : data["title"].(string),
	}

	jsonBody, err := json.Marshal(bot)

	check(err)

	fmt.Println("INFO: Uploading.")

	signedUrl := SendHTTPRequest(url, bytes.NewBuffer(jsonBody))

	bar := pb.New(len(buffer)).SetUnits(pb.U_BYTES)

	bar.Format("[## ]")

	bar.Start()

	proxy := NewReader(buffer, bar)

	req, err := http.NewRequest("PUT", signedUrl, proxy)

	req.ContentLength = size

	check(err)

	resp, err := http.DefaultClient.Do(req)

	check(err)

	dat, err := ioutil.ReadAll(resp.Body)

	check(err)

	defer resp.Body.Close()

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

	if len(result.Id) > 0{
		fmt.Println("\r\n=> " + BaseURL + "/bot/" + result.Id + "\r\n")
		fmt.Println("INFO: Success!")
		return
	}

	if len(result.Message) > 0{
		message:= fmt.Sprintf("INFO: %s", result.Message)
		fmt.Println(message)
	}

	fmt.Println("\x1b[31;1mFatal: Deploy Failed!!!\x1b[0m")
}
