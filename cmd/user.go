package cmd

// import "fmt"
// import "bufio"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/mitchellh/go-homedir"
)

//User Recime User
type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Company  string `json:"company"`
	Config	 []Config `json:"config"`
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

// Guard validates the account against recime cloud
func Guard(user User) {
	if user.Email == "" {
		fmt.Println("\x1b[31;1mInvalid account. Please run \"recime-cli init\" to get started.\x1b[0m")
		os.Exit(1)
	}
}
