package shared

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
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

type user struct {
	Email string `json:"email"`
}

type api struct {
	Key string `json:"apiKey"`
}

// Token defines the token
type Token struct {
	Source   string
	ID       string `json:"token"`
	ExpireAt int64  `json:"expireAt"`
	User     user   `json:"user"`
}

func (t *user) currentUser(source string, token string) string {
	client := &http.Client{}

	url := fmt.Sprintf("%s/user", source)

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	resp, err := client.Do(req)

	check(err)

	defer resp.Body.Close()

	dat, err := ioutil.ReadAll(resp.Body)

	var result struct {
		Email string `json:"email"`
	}

	json.Unmarshal(dat, &result)

	return result.Email

}

// Lease leases a new token
func (t *Token) Lease(apiKey string) {
	body := api{
		Key: apiKey,
	}

	jsonBody, err := json.Marshal(body)

	check(err)

	endpoint := fmt.Sprintf("%v/token/api-key", t.Source)

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

	s.Start()

	r := bytes.NewBuffer(jsonBody)

	resp, err := http.Post(endpoint, "application/json; charset=utf-8", r)

	check(err)

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)

	check(err)

	result := Token{}

	json.Unmarshal(bytes, &result)

	u := user{}

	result.User.Email = u.currentUser(t.Source, result.ID)

	s.Stop()

	fmt.Println("")

	if resp.StatusCode == 200 {
		t.save(result)

		color := color.New(color.FgHiMagenta)

		fmt.Println("")

		fmt.Print("Logged in as: ")

		color.Print(result.User.Email)

		fmt.Println("")
	} else {
		color := color.New(color.FgHiRed)

		switch resp.StatusCode {
		case 400:
			color.Println("Missing/Invalid argument.")
		case 401:
			color.Println("Invalid or expired API key. Please check that you have pasted it correctly.")
		case 403:
			color.Println("Account is not verified.")
		default:
			color.Println("Opps...There had been some issues, please try again later.")
		}
	}
}

func (t *Token) read() Token {
	dir, _ := homedir.Dir()

	path := filepath.Join(dir, filepath.Join(".recime", "netrc"))

	var result Token

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return result
	}

	reader, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)

	check(err)

	dat, _ := ioutil.ReadAll(reader)

	if len(dat) > 0 {
		json.Unmarshal(dat, &result)
	}

	return result
}

func (t *Token) save(result Token) {
	homeDir, err := homedir.Dir()

	check(err)

	filePath := filepath.Join(".recime", "netrc")

	location := filepath.Join(homeDir, filePath)

	err = os.MkdirAll(filepath.Dir(location), 0755)

	check(err)

	file, err := os.OpenFile(location, os.O_RDONLY|os.O_CREATE, 0600)

	check(err)

	file, err = os.OpenFile(location, os.O_WRONLY|os.O_TRUNC, 0600)

	jsonBody, err := json.Marshal(result)

	check(err)

	file.Write(jsonBody)
}

//Validate validates the token
func (t *Token) Validate() (*Token, error) {
	homeDir, err := homedir.Dir()

	if err != nil {
		return t, err
	}

	filePath := filepath.Join(".recime", "netrc")

	location := filepath.Join(homeDir, filePath)

	file, err := os.OpenFile(location, os.O_RDONLY|os.O_CREATE, 0600)

	if err != nil {
		return t, err
	}

	dat, err := ioutil.ReadAll(file)

	if len(dat) > 0 {
		json.Unmarshal(dat, &t)
	}

	return t, err
}

// Renew renews a token
func (t *Token) Renew() Token {
	endpoint := fmt.Sprintf("%v/token", t.Source)

	client := &http.Client{}

	req, err := http.NewRequest("PUT", endpoint, nil)

	check(err)

	token := t.read()

	// Renew after 20 days
	secs := time.Now().AddDate(0, 0, 20).Unix()

	if len(token.ID) > 0 && token.ExpireAt < secs {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.ID))

		resp, err := client.Do(req)

		check(err)

		defer resp.Body.Close()

		dat, err := ioutil.ReadAll(resp.Body)

		if len(dat) > 0 {
			json.Unmarshal(dat, &token)
			t.save(token)
		}

		return token
	}

	return token
}
