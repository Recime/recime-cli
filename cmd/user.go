package cmd

// import "fmt"
// import "bufio"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mitchellh/go-homedir"
)

//Config user configuration
type Config struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//User Recime User
type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Company  string `json:"company"`
}

// GetStoredUser fetches the stored user
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

//GetUserConfig Gets user configuration.
func GetUserConfig(options map[string]interface{}) map[string][]Config {
	user, err := GetStoredUser()

	body := User{Email: user.Email}

	jsonBody, err := json.Marshal(body)

	check(err)

	url := options["base"].(string) + "/api/config"

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

	s.Start()

	r := bytes.NewBuffer(jsonBody)

	resp, err := http.Post(url, "application/json; charset=utf-8", r)

	var result map[string][]Config

	if err == nil {
		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)

		check(err)

		json.Unmarshal(bytes, &result)
	}

	s.Stop()

	return result
}

// Guard validates the account against recime cloud
func Guard(user User) {
	if user.Email == "" {
		fmt.Println("\x1b[31;1mInvalid account. Please run \"recime-cli init\" to get started.\x1b[0m")
		os.Exit(1)
	}
}
