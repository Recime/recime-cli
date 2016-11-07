package cmd

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Prepare prepares the bot to deploy.
func Prepare() string {
	wd, err := os.Getwd()

	var data map[string]interface{}

	buff, err := ioutil.ReadFile(wd + "/package.json")

	check(err)

	if err := json.Unmarshal(buff, &data); err != nil {
		panic(err)
	}

	user, err := GetStoredUser()

	check(err)

	name := data["name"].(string)

	uid := CreateUID(name, user.Email)

	// update the author
	data["author"] = fmt.Sprintf("%s <%s>", user.Company, user.Email)

	fmt.Println("INFO: Preparing.")

	asset := MarshalIndent(data)

	writePath := wd + "/package.json"

	err = ioutil.WriteFile(writePath, asset, os.ModePerm)

	check(err)

	return uid
}

//CreateUID creates md5 hash with bot name and author
func CreateUID(name string, author string) string {
	uid := author + ";" + name

	_data := []byte(uid)

	uid = fmt.Sprintf("%x", md5.Sum(_data))

	return uid
}

// MarshalIndent marshals the given object
func MarshalIndent(data map[string]interface{}) []byte {
	asset, err := json.MarshalIndent(data, "", "\t")

	check(err)

	asset = bytes.Replace(asset, []byte("\\u003c"), []byte("<"), -1)
	asset = bytes.Replace(asset, []byte("\\u003e"), []byte(">"), -1)
	asset = bytes.Replace(asset, []byte("\\u0026"), []byte("&"), -1)

	return asset
}
