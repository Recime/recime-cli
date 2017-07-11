package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/Recime/recime-cli/shared"

	"strings"
)

type platform struct {
}

func (p *platform) set(key string, value string) error {
	if len(value) > 0 {
		config := shared.Config{Key: key, Value: value, Source: apiEndpoint}
		config.Save()
	} else {
		return errors.New("Invalid config property")
	}
	return nil
}

func (p *platform) processInput(key string, title string) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(fmt.Sprintf("%s (Press \"Enter\" to continue):", title))

	scanner.Scan()

	return p.set(key, scanner.Text())
}

// Prepare prepares the bot for deploy.
func (p *platform) install(name string) {
	key := fmt.Sprintf("RECIME_%v_ACCESS_TOKEN", strings.ToUpper(name))

	var err error

	switch strings.ToLower(name) {
	case "facebook":
		err = p.processInput(key, "Page access token")
	case "telegram":
		err = p.processInput(key, "Telegram access key")
	case "wechat":
		err = p.processInput(key, "WeChat access token")
	case "slack":
		{
			fmt.Println("Please enter Slack credentials")
			m := map[string]string{
				"RECIME_SLACK_CLIENT_ID":     "Client ID",
				"RECIME_SLACK_CLIENT_SECRET": "Client Secret",
			}
			for key, value := range m {
				err = p.processInput(key, value)
				if err != nil {
					break
				}
			}
		}
	case "viber":
		err = p.processInput(key, "Viber Auth token:")
	default:
		panic("ERROR: Unsupported Platform.")
	}

	if err != nil {
		printError(err.Error())
		return
	}

	fmt.Println("")
	fmt.Println("INFO: Platform Configured Successfully. \r\nPlease do \"recime-cli deploy\" for changes to take effect.")
}
