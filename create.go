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
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Recime/recime-cli/lib"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	botBaseURL = apiEndpoint + "/bots"
)

type bot struct {
	Name    string            `json:"name"`
	Title   string            `json:"title"`
	Author  string            `json:"author"`
	Desc    string            `json:"description"`
	Version string            `json:"version"`
	Owner   string            `json:"owner"`
	Config  map[string]string `json:"config"`
}

// Create Generates the bot
func Create(folder string, lang string) {
	h := &httpClient{}

	home, err := homedir.Dir()

	check(err)

	home = fmt.Sprintf("%v/recime-cli-%v", filepath.ToSlash(home), Version)

	fileName := filepath.ToSlash(fmt.Sprintf("%s/recime-bot-%s-template.zip", home, lang))

	tokens := strings.Split(botTemplateURL(lang), "/")
	templateDir := tokens[len(tokens)-1]
	templateDir = strings.TrimSuffix(templateDir, filepath.Ext(templateDir))
	templateDir = fmt.Sprintf("%s/recime-bot-%s-template-%s", home, strings.ToLower(lang), templateDir)

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		fmt.Println("INFO: Downloading template...")

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

		s.Start()

		h.download(botTemplateURL(lang), fileName)

		s.Stop()
	}

	botDir := filepath.ToSlash(folder)

	wd, err := os.Getwd()

	if !filepath.IsAbs(botDir) {
		botDir = filepath.Join(wd, botDir)
	}

	if _, err := os.Stat(botDir); os.IsNotExist(err) {
		si, err := os.Stat(wd)

		check(err)

		err = os.Mkdir(botDir, si.Mode())

		check(err)
	}

	util.Unzip(fileName, home)

	var data map[string]interface{}

	pkgFilePath := fmt.Sprintf("%s/package.json", templateDir)

	buff, err := ioutil.ReadFile(pkgFilePath)

	check(err)

	if err := json.Unmarshal(buff, &data); err != nil {
		panic(err)
	}

	readFromStdin(data)

	// data["author"] = fmt.Sprintf("%s <%s>", user.Company, user.Email)

	name := data["title"].(string)

	r, _ := regexp.Compile("[\\s?.$#,()^!&]+")

	normalizedName := r.ReplaceAllString(name, "-")
	normalizedName = strings.ToLower(normalizedName)
	normalizedName = strings.TrimLeft(normalizedName, "_")

	data["name"] = normalizedName

	filePath := fmt.Sprintf("%s/package.json", templateDir)

	err = ioutil.WriteFile(filePath, marshalIndent(data), os.ModePerm)

	check(err)

	util.CopyDir(templateDir, botDir)

	// sync to dashboard.
	token := renewToken()

	uid := UID{}
	id := uid.Get(botDir)

	b := bot{
		Author:  data["author"].(string),
		Name:    data["name"].(string),
		Config:  make(map[string]string),
		Owner:   token.User.Email,
		Version: data["version"].(string),
	}

	if title, ok := data["title"].(string); ok {
		b.Title = title
	}

	if desc, ok := data["description"].(string); ok {
		b.Desc = desc
	}

	d := &deployer{
		ID:    id,
		Token: token.ID,
	}

	fmt.Println("Creating the experience for your dashboard...")

	d.UpdateMetadata(b)

	fmt.Println("Uploading  \"icon.png\" from source folder...")

	d.UploadIcon(botDir)

	console := color.New(color.FgHiMagenta)

	console.Println("")
	console.Println("INFO: Bot project is created successfully.")
	console.Println("")
}

func botTemplateURL(lang string) string {
	switch lang {
	case "typescript":
		return typescriptBotTemplate
	default:
		return es6BotTemplate
	}
}

func setValue(data map[string]interface{}, key string, value string) {
	if len(value) > 0 {
		data[key] = strings.TrimRight(value, "\n")
	}
}

func readFromStdin(data map[string]interface{}) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Title (%s):", data["title"])

	scanner.Scan()

	title := scanner.Text()

	fmt.Printf("Description (%s):", data["description"])

	scanner.Scan()

	desc := scanner.Text()

	fmt.Printf("License (%s):", data["license"])

	scanner.Scan()

	license := scanner.Text()

	setValue(data, "title", title)
	setValue(data, "description", desc)
	setValue(data, "license", license)
}
