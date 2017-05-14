package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Recime/recime-cli/shared"
	"github.com/fatih/color"

	"strings"
)

type platform struct {
}

func (p *platform) set(key string, value string) {
	if len(value) > 0 {
		config := shared.Config{Key: key, Value: value, Source: apiEndpoint}
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
	key := fmt.Sprintf("RECIME_%v_ACCESS_TOKEN", strings.ToUpper(name))

	switch strings.ToLower(name) {
	case "facebook":
		p.processInput(key, "Page access token:")
	case "telegram":
		p.processInput(key, "Telegram access key:")
	case "wechat":
		p.processInput(key, "WeChat access token:")
	case "slack":
		p.processInput(key, "Slack oAuth access token:")
	case "viber":
		p.processInput(key, "Viber authentication token:")
	default:
		panic("ERROR: Unsupported Platform.")
	}
}
