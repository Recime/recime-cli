package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh/terminal"
)

// Login validates the user
func Login(options map[string]interface{}) {
	scanner := bufio.NewScanner(options["in"].(io.Reader))

	fmt.Printf("Email:")

	scanner.Scan()

	email := scanner.Text()
	email = strings.TrimSpace(email)

	fmt.Printf("Password:")

	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))

	password := string(bytePassword)
	password = strings.TrimSpace(password)

	user := User{Email: email, Password: password}

	jsonBody, err := json.Marshal(user)

	check(err)

	url := options["base"].(string) + "/login"

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

	s.Start()

	r := bytes.NewBuffer(jsonBody)

	resp, err := http.Post(url, "application/json; charset=utf-8", r)

	check(err)

	var result struct {
		User    map[string]interface{} `json:"user"`
		Message string                 `json:"message"`
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)

	check(err)

	json.Unmarshal(bytes, &result)

	s.Stop()

	if result.User != nil {
		verified := result.User["verified"].(bool)
		if !verified {
			fmt.Println("\x1b[31;1mYou must verify your account in order to continue. Please look for the verification email that is sent to you when you signed up for the service.\x1b[0m")
			os.Exit(1)
		}

		saveUser(result.User)
	} else {
		fmt.Printf("\x1b[31;1m\r\n%s\x1b[0m", result.Message)
	}
}

func saveUser(user map[string]interface{}) {
	homeDir, err := homedir.Dir()

	check(err)

	filePath := filepath.Join(".recime", "netrc")

	location := filepath.Join(homeDir, filePath)

	err = os.MkdirAll(filepath.Dir(location), 0755)

	check(err)

	file, err := os.OpenFile(location, os.O_RDONLY|os.O_CREATE, 0600)

	check(err)

	file, err = os.OpenFile(location, os.O_WRONLY|os.O_TRUNC, 0600)

	jsonBody, err := json.Marshal(user)

	check(err)

	file.Write(jsonBody)

	fmt.Println("")

	color := color.New(color.FgHiMagenta)

	fmt.Println("")

	fmt.Print("Logged in as: ")

	color.Print(user["email"])

	fmt.Println("")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
