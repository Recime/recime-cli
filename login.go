package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/Recime/recime-cli/shared"
)

// Login validates the user
func Login(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	fmt.Println("Paste your api key from \"https://console.recime.io/my-account\":")

	scanner.Scan()

	apiKey := scanner.Text()
	apiKey = strings.TrimSpace(apiKey)

	t := shared.Token{Source: apiEndpoint}
	t.Lease(apiKey)
}
