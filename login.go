package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"syscall"

	"github.com/Recime/recime-cli/shared"
	"golang.org/x/crypto/ssh/terminal"
)

// Login validates the user
func Login(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	fmt.Printf("Email:")

	scanner.Scan()

	email := scanner.Text()
	email = strings.TrimSpace(email)

	fmt.Printf("Paste your auth token from :")

	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))

	password := string(bytePassword)
	password = strings.TrimSpace(password)

	t := shared.Token{Source: apiEndpoint}
	t.Lease(email, password)
}
