package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"./shared"
)

// Login validates the user
func Login(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	fmt.Printf("Email:")

	scanner.Scan()

	email := scanner.Text()
	email = strings.TrimSpace(email)

	fmt.Printf("Password:")

	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))

	password := string(bytePassword)
	password = strings.TrimSpace(password)

	t := shared.Token{Source: apiEndpoint}
	t.Lease(email, password)
}
