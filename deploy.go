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
	"archive/zip"
	"strings"

	"github.com/Recime/recime-cli/cmd"
	"github.com/Recime/recime-cli/util"
	"github.com/briandowns/spinner"
	pb "gopkg.in/cheggaaa/pb.v1"
)

type Bot struct {
	Id      string       `json:"uid"`
	Type    string       `json:"fileType"`
	Name    string       `json:"name"`
	Title   string       `json:"title"`
	Author  string       `json:"author"`
	Desc    string       `json:"description"`
	Version string       `json:"version"`
	Owner   string       `json:"owner"`
	Config  []cmd.Config `json:"config"`
	Icon    []byte       `json:"icon"`
}

type Package struct {
	FileType string `json:"fileType"`
	Id       string `json:"uid"`
	Method   string `json:"method"`
}

func PrepareLambdaPackage(uid string) string {
	jsonBody, err := json.Marshal(Package{
		FileType: "application/octet-stream",
		Id:       "package.zip",
		Method:   "getObject",
	})

	check(err)

	source := fmt.Sprintf("%s/signed-url", BaseURL)
	response := SendRequest(source, bytes.NewBuffer(jsonBody))

	signedURL := response["url"].(string)

	temp, err := ioutil.TempDir("", "recime-cli")

	check(err)

	dest := fmt.Sprintf("%s/bin", temp)

	err = os.Mkdir(dest, os.ModePerm)

	fileName := filepath.ToSlash(fmt.Sprintf("%s/%s.zip", dest, uid))

	cmd.Download(signedURL, fileName)

	target := filepath.ToSlash(fmt.Sprintf("%s/%s", dest, uid))

	check(util.Unzip(fileName, target))

	wd, err := os.Getwd()

	check(err)

	botDir := filepath.ToSlash(fmt.Sprintf("%s/%s", target, uid))

	_ = util.CopyDir(wd, botDir)

	pkg := filepath.ToSlash(fmt.Sprintf("%s/%s.zip", temp, uid))

	Zip(target, pkg)

	// zip := new(archivex.ZipFile)
	// zip.Create(pkg)
	// zip.AddAll(target, true)
	// zip.Close()

	return pkg
}

func Zip(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string

	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			path = filepath.ToSlash(path)
			source = filepath.ToSlash(source)

			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		header.Name = filepath.ToSlash(header.Name)

		writer, err := archive.CreateHeader(header)

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

// SendRequest sends POST request
func SendRequest(url string, body io.Reader) map[string]interface{} {
	res, err := http.Post(url, "application/json; charset=utf-8", body)

	check(err)

	var data map[string]interface{}

	bytes, err := ioutil.ReadAll(res.Body)

	json.Unmarshal(bytes, &data)

	defer res.Body.Close()

	return data
}

// Deploy deploys the bot with the given uid
func Deploy() {
	uid := cmd.GetUID()

	fmt.Println("INFO: Preparing Package.")

	pkgPath := PrepareLambdaPackage(uid)

	buffer, size := ReadFile(pkgPath)

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
			config = append(config, cmd.Config{Key: key, Value: value.(string)})
		}
	}

	var data map[string]interface{}

	buff, err := ioutil.ReadFile(wd + "/package.json")

	check(err)

	if err := json.Unmarshal(buff, &data); err != nil {
		panic(err)
	}

	bot := Bot{
		Author:  data["author"].(string),
		Id:      uid,
		Type:    fileType,
		Version: Version,
		Owner:   user.Email,
		Config:  config,
		Name:    data["name"].(string),
	}

	if title, ok := data["title"].(string); ok {
		bot.Title = title
	}

	if desc, ok := data["description"].(string); ok {
		bot.Desc = desc
	}

	jsonBody, err := json.Marshal(bot)

	check(err)

	fmt.Println("INFO: Uploading.")

	response := SendRequest(url, bytes.NewBuffer(jsonBody))

	signedURL := response["url"].(string)

	bar := pb.New(len(buffer)).SetUnits(pb.U_BYTES)

	bar.Format("[## ]")

	bar.Start()

	proxy := NewReader(buffer, bar)

	req, err := http.NewRequest("PUT", signedURL, proxy)

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

	UploadIcon(uid)

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

	if len(result.Id) > 0 {
		fmt.Println("\r\n=> " + BaseURL + "/bot/" + result.Id + "\r\n")
		fmt.Println("INFO: Success!")
		return
	}

	if len(result.Message) > 0 {
		message := fmt.Sprintf("INFO: %s", result.Message)
		fmt.Println(message)
	}

	fmt.Println("\x1b[31;1mFatal: Deploy Failed!!!\x1b[0m")
}

func ReadFile(path string) ([]byte, int64) {
	file, err := os.Open(path)

	check(err)

	defer file.Close()

	fileInfo, _ := file.Stat()

	var size = fileInfo.Size()

	buffer := make([]byte, size)

	// // read file content to buffer
	file.Read(buffer)

	return buffer, size
}

// UploadIcon uploads the icon from bot folder.
func UploadIcon(id string) {
	wd, _ := os.Getwd()

	icon, size := ReadFile(fmt.Sprintf("%s/icon.png", wd))

	source := fmt.Sprintf("%s/bot/icon", BaseURL)

	bot := Bot{
		Id: id,
	}

	jsonBody, err := json.Marshal(bot)

	response := SendRequest(source, bytes.NewBuffer(jsonBody))

	signedURL := response["url"].(string)

	reader := bytes.NewReader(icon)

	req, err := http.NewRequest("PUT", signedURL, reader)

	req.ContentLength = size

	check(err)

	http.DefaultClient.Do(req)
}
