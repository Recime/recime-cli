package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Recime/recime-cli/shared"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

//Defines UID
type UID struct {
}

//Create creates md5 hash with bot name and author
func (u *UID) Create(name string, author string) string {
	uid := author + ";" + name

	_data := []byte(uid)

	uid = fmt.Sprintf("%x", md5.Sum(_data))

	return uid
}

// Get gets the uid for the package.
func (u *UID) Get() string {
	home, _ := homedir.Dir()

	var data map[string]interface{}

	buff, err := ioutil.ReadFile(home + "/package.json")

	check(err)

	if err := json.Unmarshal(buff, &data); err != nil {
		panic(err)
	}

	t := shared.Token{}

	token, err := t.Validate()

	if len(token.ID) > 0 {
		console := color.New(color.FgHiRed)
		console.Println("User is not logged in. Please run \"recime-cli login\" to get started.")
		fmt.Println("")
		os.Exit(1)
	}

	check(err)

	name := data["name"].(string)

	uid := u.Create(name, t.Email)

	return uid
}
