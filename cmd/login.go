package cmd

import "fmt"

import "bufio"
import "bytes"
import "path/filepath"
import "io"
import "io/ioutil"
import "encoding/json"
import "net/http"
import "os"
import "strings"
import "syscall"
import "time"

import "golang.org/x/crypto/ssh/terminal"
import "github.com/briandowns/spinner"
import "github.com/mitchellh/go-homedir"

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
	fmt.Printf("INFO: User Verification Successful.")
	fmt.Println("")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
