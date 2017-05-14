package shared

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

//User Recime User
type User struct {
	Email    string   `json:"email"`
	ID       string   `json:"_id"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Company  string   `json:"company"`
	Config   []Config `json:"config"`
}

// Guard validates the account against recime cloud
func (u *User) Guard(user User) {
	if user.Email == "" {
		console := color.New(color.FgHiRed)
		console.Println("User is not logged in. Please run \"recime-cli login\" to get started.")
		fmt.Println("")
		os.Exit(1)
	}
}
