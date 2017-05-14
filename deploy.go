// Copyright 2017 The Recime Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless 261d by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"path/filepath"

	"github.com/Recime/recime-cli/lib"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"

	bar "gopkg.in/cheggaaa/pb.v1"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"errors"

	pb "github.com/Recime/recime-cli/deployment"

	"github.com/Recime/recime-cli/shared"
)

const (
	address    = "ec2-52-36-237-55.us-west-2.compute.amazonaws.com"
	port       = 3000
	bucket     = "recime-io"
	botBaseURL = apiEndpoint + "/bots"
)

// PrintStatus outputs formatted status.
func printRemoteStatus(status string) {
	pattern := regexp.MustCompile(`[a-z1-9A-Z.]+`)
	if pattern.MatchString(status) {
		fmt.Println(fmt.Sprintf("remote ---> %v", status))
	} else {
		fmt.Print(status)
	}
}

type bot struct {
	Name    string            `json:"name"`
	Title   string            `json:"title"`
	Author  string            `json:"author"`
	Desc    string            `json:"description"`
	Version string            `json:"version"`
	Owner   string            `json:"owner"`
	Config  map[string]string `json:"config"`
}

type deployer struct {
	ID    string
	Token string
}

// Deploy deployes the bot.
func (d *deployer) Deploy() {
	target := fmt.Sprintf("%s:%v", address, port)

	connection, err := grpc.Dial(
		target,
		grpc.WithBackoffMaxDelay(10*time.Second),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{Time: 5 * time.Second}),
		grpc.WithInsecure())

	if err != nil {
		fmt.Println(fmt.Sprintf("\x1b[31;1mFatal: %v\x1b[0m", err))
	}

	defer connection.Close()

	// Creates a new CustomerClient
	client := pb.NewDeployerClient(connection)

	deployRequest := &pb.DeployRequest{
		Token: d.Token,
		BotId: d.ID,
	}

	stream, err := client.Deploy(context.Background(), deployRequest)

	if err != nil {
		fmt.Println(fmt.Sprintf("\x1b[31;1mFatal: %v\x1b[0m", err))
		os.Exit(1)
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

	s.Stop()

	failed := false

	for {
		resp, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(fmt.Sprintf("\x1b[31;1mFatal: %v\x1b[0m", err))
			os.Exit(1)
		}

		if resp.Code == 0 {
			s.Stop()
			printRemoteStatus(resp.Message)
		}

		if resp.Code > 0 {
			fmt.Println("")
			printError(resp.Message)
			fmt.Println("")
			failed = true
			break
		}

		s.Start()
	}

	s.Stop()

	if failed {
		fmt.Println("For any questions and feedback, please reach us at hello@recime.io.")
		fmt.Println("")

		os.Exit(1)
	}

}

// Register registers the bot.
func (d *deployer) UpdateMetadata(b bot) []byte {
	uid := d.ID

	jsonBody, err := json.Marshal(b)

	check(err)

	url := fmt.Sprintf("%s/%s", botBaseURL, uid)

	dat := sendRequest(url, d.Token, bytes.NewBuffer(jsonBody))

	return dat
}

func (d *deployer) printMetadata() {
	uid := d.ID

	client := &http.Client{}

	url := fmt.Sprintf("%s/%s", botBaseURL, uid)

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", d.Token))

	resp, err := client.Do(req)

	check(err)

	defer resp.Body.Close()

	dat, err := ioutil.ReadAll(resp.Body)

	var result struct {
		Region string `json:"region"`
		ID     string `json:"id"`
	}

	json.Unmarshal(dat, &result)

	if len(result.ID) > 0 {
		console := color.New(color.FgHiMagenta)

		fmt.Println("")

		console.Println(fmt.Sprintf("https://%s-bot.recime.io/%s/v1", result.Region, result.ID))

		fmt.Println("")

		fmt.Println("INFO: Success!")
	}
}

// UploadIcon uploads the icon from bot folder.
func (d *deployer) UploadIcon() {
	wd, _ := os.Getwd()

	icon, size := readFile(fmt.Sprintf("%s/icon.png", wd))

	requestURL := fmt.Sprintf("%v/%v/uploads/icon", botBaseURL, d.ID)

	dat := sendRequest(requestURL, d.Token, nil)

	var result struct {
		URL string `json:"url"`
	}

	check(json.Unmarshal(dat, &result))

	fmt.Println(result.URL)

	reader := bytes.NewReader(icon)

	req, err := http.NewRequest("PUT", result.URL, reader)

	req.ContentLength = size

	check(err)

	http.DefaultClient.Do(req)
}

func preparePackage(uid string) (string, error) {
	temp, err := ioutil.TempDir("", "recime-cli")

	check(err)

	dest := fmt.Sprintf("%s/bin", temp)

	err = os.Mkdir(filepath.ToSlash(dest), os.ModePerm)

	fileName := filepath.ToSlash(fmt.Sprintf("%s/%s.zip", dest, uid))

	h := &httpClient{}
	h.download(template, fileName)

	check(util.Unzip(fileName, dest))

	tokens := strings.Split(template, "/")

	templatedir := tokens[len(tokens)-1]
	templatedir = strings.TrimSuffix(templatedir, filepath.Ext(templatedir))
	templatedir = fmt.Sprintf("recime-bot-template-%s", templatedir)

	wd, err := os.Getwd()

	check(err)

	bindir := filepath.ToSlash(fmt.Sprintf("%s/%s", dest, templatedir))
	botdir := filepath.ToSlash(fmt.Sprintf("%s/%s", bindir, uid))

	_ = util.CopyDir(wd, botdir)

	sh := &shell{}

	sh.execute(botdir, "install")

	if Build(botdir) != nil {
		return "", errors.New("Build failed")
	}

	pkg := &pkg{}
	pkg.sync(botdir, bindir)

	removeScript(botdir)

	pkgdir := filepath.ToSlash(fmt.Sprintf("%s/%s", dest, uid))

	util.CopyDir(bindir, pkgdir)

	zip := filepath.ToSlash(fmt.Sprintf("%s/%s.zip", temp, uid))

	util.Zip(pkgdir, zip)

	return zip, nil
}

func removeScript(dir string) {
	var data map[string]interface{}

	pkgFilePath := fmt.Sprintf("%s/package.json", dir)

	buff, err := ioutil.ReadFile(pkgFilePath)

	check(err)

	if err := json.Unmarshal(buff, &data); err != nil {
		panic(err)
	}

	delete(data, "scripts")

	err = ioutil.WriteFile(pkgFilePath, marshalIndent(data), os.ModePerm)

	check(err)
}

// SendRequest sends POST request
func sendRequest(url string, token string, body io.Reader) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)

	if len(token) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

	s.Start()

	res, err := client.Do(req)

	check(err)

	defer res.Body.Close()

	s.Stop()

	bytes, err := ioutil.ReadAll(res.Body)

	check(err)

	color := color.New(color.FgHiRed)

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return bytes
	}

	switch res.StatusCode {
	case 400:
		color.Println("Bad Request.")
	case 401:
		color.Println("Unauthorized. Invalid or expired token. Please do \"recime-cli login\" and try again.")
	case 405:
		color.Println("The operation is not allowed in your subscription.")
	}
	os.Exit(1)

	return nil
}

func printError(msg string) {
	if len(msg) > 0 {
		console := color.New(color.FgHiRed)
		message := fmt.Sprintf("FATAL: %s", msg)
		console.Println(message)
	}
}

// Deploy deploys the bot
func Deploy() {
	uid := UID{}

	id := uid.Get()

	token := renewToken()

	fmt.Println("Creating bot package to deploy into \"Recime\" cloud.")

	pkgPath, err := preparePackage(id)

	if err != nil {
		color := color.New(color.FgHiRed)
		color.Println(err)
		return
	}

	buffer, size := readFile(pkgPath)

	wd, err := os.Getwd()

	check(err)

	_config := shared.Config{}

	cfg := make(map[string]string)

	// open config.json
	reader, err := _config.Open(wd)

	if err == nil {
		cfg = _config.Get(reader)
	}

	var data map[string]interface{}

	buff, err := ioutil.ReadFile(wd + "/package.json")

	check(err)

	if err := json.Unmarshal(buff, &data); err != nil {
		panic(err)
	}

	fmt.Println("Updating metadata information.")

	_bot := bot{
		Author:  data["author"].(string),
		Version: Version,
		Owner:   token.Email,
		Config:  cfg,
		Name:    data["name"].(string),
	}

	if title, ok := data["title"].(string); ok {
		_bot.Title = title
	}

	if desc, ok := data["description"].(string); ok {
		_bot.Desc = desc
	}

	d := &deployer{
		ID:    id,
		Token: token.ID,
	}

	d.UpdateMetadata(_bot)

	botUploadRequestURL := fmt.Sprintf("%v/%v/uploads/bot", botBaseURL, uid)

	dat := sendRequest(botUploadRequestURL, token.ID, nil)

	var uploadResult struct {
		URL string `json:"url"`
	}

	check(json.Unmarshal(dat, &uploadResult))

	bar := bar.New(len(buffer))

	bar.ShowCounters = false

	bar.Prefix("Uploading: ")

	bar.Format("[## ]")

	bar.Start()

	proxy := NewReader(buffer, bar)

	req, err := http.NewRequest("PUT", uploadResult.URL, proxy)

	req.ContentLength = size

	check(err)

	resp, err := http.DefaultClient.Do(req)

	check(err)

	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)

	check(err)

	fmt.Println("")
	fmt.Println("Updating  \"icon.png\" from source folder.")

	d.UploadIcon()

	fmt.Println("")
	fmt.Println("")

	fmt.Println("Starting deployment...")
	fmt.Println("")

	d.Deploy()

	d.printMetadata()
}

func readFile(path string) ([]byte, int64) {
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

func renewToken() shared.Token {
	token := shared.Token{Source: apiEndpoint}
	return token.Renew()
}
