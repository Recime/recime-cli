package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"

	"strings"

	"github.com/Recime/recime-cli/cmd"
)

type platform struct {
}

// Prepare prepares the bot for deploy.
func (p *platform) install(name string) {
	switch strings.ToLower(name) {
	case "facebook":
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("Page Access Token:")

		scanner.Scan()

		token := scanner.Text()

		if len(token) > 0 {
			config := cmd.Config{Key: "FACEBOOK_PAGE_ACCESS_TOKEN", Value: token, Source: apiEndpoint}

			cmd.SaveConfig(config)

			fmt.Println("")
			fmt.Println("INFO: Platform Configured Successfully. \r\nPlease do \"recime-cli deploy\" for changes to take effect.")
		} else {
			red := color.New(color.FgRed).Add(color.Bold)

			red.Println("ERROR: Invalid Page Token! Please verify input and try again.")
		}
	default:
		panic("INFO: Unsupported Platform.")
	}

}
