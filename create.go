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

func SetValue(data map[string]interface{}, key string, value string) {
	if len(value) > 0 {
		data[key] = strings.TrimRight(value, "\n")
	}
}

func ProcesssInput(in io.Reader) (data map[string]interface{}) {
	scanner := bufio.NewScanner(in)

	res := &Resource{}

	asset := res.Get("data/package.json")

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

	SetValue(data, "title", title)
	SetValue(data, "description", desc)
	SetValue(data, "license", license)

	return data
}

// Create Generates the bot
func Create(folder string) {
	user, err := cmd.GetStoredUser()

	cmd.Guard(user)

	wd, err := os.Getwd()

	data := ProcesssInput(os.Stdin)

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

	res := &Resource{}

	resources, err := res.GetDir("data")

	check(err)

	for key := range resources {
		entry := resources[key]

		asset := res.Get("data/" + entry)

		if entry == "package.json" {
			asset = cmd.MarshalIndent(data)
		}

		filePath := path + "/" + entry

		err = ioutil.WriteFile(filePath, asset, os.ModePerm)

		check(err)
	}

	fmt.Println("Bot Created Successfully.")

}
