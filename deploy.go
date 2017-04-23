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

	"github.com/Recime/recime-cli/cmd"
	"github.com/Recime/recime-cli/util"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"

	bar "gopkg.in/cheggaaa/pb.v1"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"errors"

	pb "github.com/Recime/recime-cli/pb"
)

const (
	address   = "deployer.recime.io"
	port      = 3000
	bucket    = "recime-io"
	singedURL = apiEndpoint + "/signedurl"
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

// Resource contains the bucket information.
type resource struct {
	Key    string `json:"key"`
	Bucket string `json:"bucket"`
}

type bot struct {
	ID       string       `json:"uid"`
	Type     string       `json:"fileType"`
	Name     string       `json:"name"`
	Title    string       `json:"title"`
	Author   string       `json:"author"`
	Desc     string       `json:"description"`
	Version  string       `json:"version"`
	Owner    string       `json:"owner"`
	Config   []cmd.Config `json:"config"`
	Icon     []byte       `json:"icon"`
	Resource *resource    `json:"resource"`
}

type deployer struct {
	ID          string
	UserID      string
	Environment map[string]string
}

// Prepare prepares the bot for deploy.
func (d *deployer) Prepare() {
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
		UserId:      d.UserID,
		BotId:       d.ID,
		Environment: d.Environment,
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

// Deploy deploys the bot in aws lambda.
func (d *deployer) Deploy(b bot) []byte {
	uid := b.ID

	jsonBody, err := json.Marshal(b)

	check(err)

	url := fmt.Sprintf("%s/bot/register/%s", apiEndpoint, uid)

	r := bytes.NewBuffer(jsonBody)

	resp, err := http.Post(url, "application/json; charset=utf-8", r)

	check(err)

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)

	check(err)

	return bytes
}

// UploadIcon uploads the icon from bot folder.
func (d *deployer) UploadIcon() {
	wd, _ := os.Getwd()

	icon, size := readFile(fmt.Sprintf("%s/icon.png", wd))

	source := fmt.Sprintf("%s/bot/icon", apiEndpoint)

	jsonBody, err := json.Marshal(&bot{
		ID: d.ID,
	})

	response := sendRequest(source, bytes.NewBuffer(jsonBody))

	signedURL := response["url"].(string)

	reader := bytes.NewReader(icon)

	req, err := http.NewRequest("PUT", signedURL, reader)

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

	if cmd.Build(botdir) != nil {
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

	err = ioutil.WriteFile(pkgFilePath, cmd.MarshalIndent(data), os.ModePerm)

	check(err)
}

// SendRequest sends POST request
func sendRequest(url string, body io.Reader) map[string]interface{} {
	res, err := http.Post(url, "application/json; charset=utf-8", body)

	check(err)

	var data map[string]interface{}

	bytes, err := ioutil.ReadAll(res.Body)

	check(err)

	json.Unmarshal(bytes, &data)

	defer res.Body.Close()

	return data
}

func printError(msg string) {
	if len(msg) > 0 {
		console := color.New(color.FgHiRed)
		message := fmt.Sprintf("FATAL: %s", msg)
		console.Println(message)
	}
}

func guard(uid string, owner string) {
	source := fmt.Sprintf("%s/bot/status", apiEndpoint)

	jsonBody, err := json.Marshal(&bot{
		ID:      uid,
		Owner:   owner,
		Version: Version,
	})

	res, err := http.Post(source, "application/json; charset=utf-8", bytes.NewBuffer(jsonBody))

	check(err)

	var data struct {
		ID      string `json:"uid"`
		Message string `json:"message"`
	}

	bytes, err := ioutil.ReadAll(res.Body)

	defer res.Body.Close()

	check(err)

	json.Unmarshal(bytes, &data)

	if len(data.Message) > 0 {
		printError(data.Message)

		fmt.Println("")

		os.Exit(1)
	}
}

// Deploy deploys the bot
func Deploy() {
	uid := cmd.GetUID()
	user, err := cmd.GetStoredUser()

	guard(uid, user.Email)

	fmt.Println("Creating bot package to deploy into \"Recime\" cloud.")

	pkgPath, err := preparePackage(uid)

	if err != nil {
		return
	}

	buffer, size := readFile(pkgPath)

	fileType := http.DetectContentType(buffer)

	var config []cmd.Config

	wd, err := os.Getwd()

	check(err)

	_config := cmd.Config{}

	env := make(map[string]string)

	// open config.json
	reader, err := _config.Open(wd)

	if err == nil {
		env = _config.Get(reader)
		for key, value := range env {
			config = append(config, cmd.Config{Key: key, Value: value})
		}
	}

	var data map[string]interface{}

	buff, err := ioutil.ReadFile(wd + "/package.json")

	check(err)

	if err := json.Unmarshal(buff, &data); err != nil {
		panic(err)
	}

	r := &resource{
		Key: fmt.Sprintf("bot/%s", uid),
	}

	jsonBody, err := json.Marshal(r)

	response := sendRequest(singedURL, bytes.NewBuffer(jsonBody))

	signedURL := response["url"].(string)

	bar := bar.New(len(buffer))

	bar.ShowCounters = false

	bar.Prefix("Uploading: ")

	bar.Format("[## ]")

	bar.Start()

	proxy := NewReader(buffer, bar)

	req, err := http.NewRequest("PUT", signedURL, proxy)

	req.ContentLength = size

	check(err)

	resp, err := http.DefaultClient.Do(req)

	check(err)

	_, err = ioutil.ReadAll(resp.Body)

	check(err)

	defer resp.Body.Close()

	b := bot{
		Author:  data["author"].(string),
		ID:      uid,
		Type:    fileType,
		Version: Version,
		Owner:   user.Email,
		Config:  config,
		Name:    data["name"].(string),
	}

	d := &deployer{
		ID:          b.ID,
		UserID:      user.ID,
		Environment: env,
	}

	d.Prepare()

	fmt.Println("")
	fmt.Println("Registering the bot and creating the API endpoint.")

	d.UploadIcon()

	if title, ok := data["title"].(string); ok {
		b.Title = title
	}

	if desc, ok := data["description"].(string); ok {
		b.Desc = desc
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

	s.Start()

	bytes := d.Deploy(b)

	s.Stop()

	var result struct {
		Name    string `json:"name"`
		ID      string `json:"uid"`
		Message string `json:"message"`
		URI     string `json:"uri"`
	}

	json.Unmarshal(bytes, &result)

	if len(result.ID) > 0 {
		console := color.New(color.FgHiMagenta)

		fmt.Println("")

		console.Println(result.URI)

		fmt.Println("")

		fmt.Println("INFO: Success!")

		return
	}

	printError(result.Message)
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
