// Copyright 2017 The Recime Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
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
	"time"

	"path/filepath"

	"github.com/Recime/recime-cli/cmd"
	"github.com/Recime/recime-cli/util"
	"github.com/briandowns/spinner"

	bar "gopkg.in/cheggaaa/pb.v1"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/Recime/recime-cli/pb"
)

const (
	address   = "agent.recime.io"
	port      = 3000
	bucket    = "recime-io"
	template  = "https://github.com/Recime/recime-lambda-package-template/releases/download/1.0.1/package.zip"
	singedURL = BaseURL + "/signed-url"
)

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
	ID       string
	Resource *resource
}

// Prepare prepares the bot for deploy.
func (d *deployer) Prepare() {
	target := fmt.Sprintf("%s:%v", address, port)

	connection, err := grpc.Dial(target, grpc.WithInsecure())

	if err != nil {
		fmt.Println(fmt.Sprintf("\x1b[31;1mFatal: %v\x1b[0m", err))
	}

	defer connection.Close()

	// Creates a new CustomerClient
	client := pb.NewDeployerClient(connection)

	deployRequest := &pb.DeployRequest{
		Resource: &pb.Resource{
			Bucket: bucket,
			Key:    fmt.Sprintf("bot/%s", d.ID),
		},
	}

	stream, err := client.Deploy(context.Background(), deployRequest)

	if err != nil {
		fmt.Println(fmt.Sprintf("\x1b[31;1mFatal: %v\x1b[0m", err))
		os.Exit(1)
	}

	var r *pb.Resource

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

	for {
		resp, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(fmt.Sprintf("\x1b[31;1mFatal: %v\x1b[0m", err))
			os.Exit(1)
		}

		if len(resp.Status) > 1 {
			s.Stop()
			fmt.Println(fmt.Sprintf("INFO: %v", resp.Status))
		}

		s.Start()

		if resp.Resource != nil {
			r = resp.Resource
		}
	}

	s.Stop()

	if len(r.Key) == 0 {
		fmt.Println("\x1b[31;1mFatal: Deploy Failed!!!\x1b[0m")
		os.Exit(1)
	}

	d.Resource = &resource{
		Key:    r.Key,
		Bucket: r.Bucket,
	}
}

// Deploy deploys the bot in aws lambda.
func (d *deployer) Deploy(b bot) []byte {
	uid := b.ID

	b.Resource = &resource{
		Bucket: d.Resource.Bucket,
		Key:    d.Resource.Key,
	}

	jsonBody, err := json.Marshal(b)

	check(err)

	url := fmt.Sprintf("%s/module/deploy/%s", BaseURL, uid)

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

	source := fmt.Sprintf("%s/bot/icon", BaseURL)

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

func prepareLambdaPackage(uid string) string {
	temp, err := ioutil.TempDir("", "recime-cli")

	check(err)

	dest := fmt.Sprintf("%s/bin", temp)

	err = os.Mkdir(filepath.ToSlash(dest), os.ModePerm)

	fileName := filepath.ToSlash(fmt.Sprintf("%s/%s.zip", dest, uid))

	cmd.Download(template, fileName)

	target := filepath.ToSlash(fmt.Sprintf("%s/%s", dest, uid))

	check(util.Unzip(fileName, target))

	wd, err := os.Getwd()

	check(err)

	botDir := filepath.ToSlash(fmt.Sprintf("%s/%s", target, uid))

	_ = util.CopyDir(wd, botDir)

	removeScript(botDir)

	pkg := filepath.ToSlash(fmt.Sprintf("%s/%s.zip", temp, uid))

	util.Zip(target, pkg)

	return pkg
}

func removeScript(dir string) {
	var data map[string]interface{}

	pkgFilePath := fmt.Sprintf("%s/package.json", dir)

	buff, err := ioutil.ReadFile(pkgFilePath)

	check(err)

	if err := json.Unmarshal(buff, &data); err != nil {
		panic(err)
	}

	data["scripts"] = make([]interface{}, 0)

	err = ioutil.WriteFile(pkgFilePath, cmd.MarshalIndent(data), os.ModePerm)

	check(err)
}

// SendRequest sends POST request
func sendRequest(url string, body io.Reader) map[string]interface{} {
	res, err := http.Post(url, "application/json; charset=utf-8", body)

	check(err)

	var data map[string]interface{}

	bytes, err := ioutil.ReadAll(res.Body)

	json.Unmarshal(bytes, &data)

	defer res.Body.Close()

	return data
}

// Deploy deploys the bot
func Deploy() {
	uid := cmd.GetUID()

	fmt.Println("INFO: Preparing Package.")

	pkgPath := prepareLambdaPackage(uid)

	buffer, size := readFile(pkgPath)

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

	fmt.Println("INFO: Uploading.")

	r := &resource{
		Key: fmt.Sprintf("bot/%s", uid),
	}

	jsonBody, err := json.Marshal(r)

	response := sendRequest(singedURL, bytes.NewBuffer(jsonBody))

	signedURL := response["url"].(string)

	bar := bar.New(len(buffer)).SetUnits(bar.U_BYTES)

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

	d := &deployer{
		ID: uid,
	}

	d.Prepare()

	fmt.Println("INFO: Finishing.")

	d.UploadIcon()

	b := bot{
		Author:  data["author"].(string),
		ID:      uid,
		Type:    fileType,
		Version: Version,
		Owner:   user.Email,
		Config:  config,
		Name:    data["name"].(string),
	}

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
	}

	json.Unmarshal(bytes, &result)

	if len(result.ID) > 0 {
		fmt.Println("\r\n=> " + BaseURL + "/bot/" + result.ID + "\r\n")
		fmt.Println("INFO: Success!")
		return
	}

	if len(result.Message) > 0 {
		message := fmt.Sprintf("INFO: %s", result.Message)
		fmt.Println(message)
	}

	fmt.Println("\x1b[31;1mFatal: Deploy Failed!!!\x1b[0m")
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
