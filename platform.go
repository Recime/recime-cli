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

func (p *platform) set(key string, value string) {
	if len(value) > 0 {
		config := cmd.Config{Key: key, Value: value, Source: apiEndpoint}
		config.Save()

		fmt.Println("")
		fmt.Println("INFO: Platform Configured Successfully. \r\nPlease do \"recime-cli deploy\" for changes to take effect.")
	} else {
		red := color.New(color.FgRed).Add(color.Bold)
		red.Println("ERROR: Invalid Page Token! Please verify input and try again.")
	}
}

// Prepare prepares the bot for deploy.
func (p *platform) install(name string) {
	switch strings.ToLower(name) {
	case "facebook":
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("Page access token:")

		scanner.Scan()

		p.set("TELEGRAM_ACCESS_TOKEN", scanner.Text())
	case "telegram":
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("Telegram access key:")

		scanner.Scan()

		p.set("TELEGRAM_ACCESS_TOKEN", scanner.Text())
	default:
		panic("INFO: Unsupported Platform.")
	}

}
