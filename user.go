package main

// import "fmt"
// import "bufio"

import "io/ioutil"
import "os"

import "path/filepath"
import "encoding/json"

import "github.com/mitchellh/go-homedir"

type User struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Company string `json:"company"`
}

func GetStoredUser() (User, error) {
	var user User
	homeDir, err := homedir.Dir()

	if err != nil {
		return user, err
	}

	filePath := filepath.Join(".recime", "netrc")

	location := filepath.Join(homeDir, filePath)

	file, err := os.OpenFile(location, os.O_RDONLY|os.O_CREATE, 0600)

	if err != nil {
		return user, err
	}

	dat, err := ioutil.ReadAll(file)

	if len(dat) > 0 {
		json.Unmarshal(dat, &user)
	}

	return user, err
}
