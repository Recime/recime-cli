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

func (p *platform) processInput(key string, title string) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(title)

	scanner.Scan()

	p.set(key, scanner.Text())
}

// Prepare prepares the bot for deploy.
func (p *platform) install(name string) {
	switch strings.ToLower(name) {
	case "facebook":
		p.processInput("FACEBOOK_PAGE_ACCESS_TOKEN", "Page access token:")
	case "telegram":
		p.processInput("TELEGRAM_ACCESS_TOKEN", "Telegram access key:")
	case "wechat":
		p.processInput("WECHAT_ACCESS_TOKEN", "WeChat access token:")
	case "slack":
		p.processInput("SLACK_ACCESS_TOKEN", "Slack access token:")
	case "viber":
		p.processInput("VIBER_ACCESS_TOKEN", "Viber authentication token:")
	default:
		panic("ERROR: Unsupported Platform.")
	}
}
