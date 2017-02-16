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

import "fmt"
import "os"
import "io"
import "io/ioutil"

import "encoding/json"
import "bufio"
import "strings"
import "regexp"

import "github.com/Recime/recime-cli/cmd"
import "path/filepath"

// Create Generates the bot
func Create(folder string) {
	user, err := cmd.GetStoredUser()

	cmd.Guard(user)

	wd, err := os.Getwd()

	data := processsInput(os.Stdin)

	data["author"] = fmt.Sprintf("%s <%s>", user.Company, user.Email)

	name := data["title"].(string)

	r, _ := regexp.Compile("[\\s?.$#,()^!&]+")

	normalizedName := r.ReplaceAllString(name, "-")
	normalizedName = strings.ToLower(normalizedName)
	normalizedName = strings.TrimLeft(normalizedName, "_")

	data["name"] = normalizedName

	check(err)

	path := filepath.ToSlash(folder)

	if !filepath.IsAbs(path) {
		path = filepath.Join(wd, path)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		si, err := os.Stat(wd)

		check(err)

		err = os.Mkdir(path, si.Mode())

		check(err)
	}

	resources, err := AssetDir("data")

	check(err)

	for key := range resources {
		entry := resources[key]

		asset := MustAsset("data/" + entry)

		if entry == "package.json" {
			asset = cmd.MarshalIndent(data)
		}

		filePath := path + "/" + entry

		err = ioutil.WriteFile(filePath, asset, os.ModePerm)

		check(err)
	}

	fmt.Println("Bot Created Successfully.")

}

func setValue(data map[string]interface{}, key string, value string) {
	if len(value) > 0 {
		data[key] = strings.TrimRight(value, "\n")
	}
}

func processsInput(in io.Reader) (data map[string]interface{}) {
	scanner := bufio.NewScanner(in)

	asset := MustAsset("data/package.json")

	check(json.Unmarshal(asset, &data))

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

	return data
}
