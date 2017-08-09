package main

import (
	"bufio"
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
	}
	return nil
}

func (p *platform) processInput(key string, title string) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(fmt.Sprintf("%s (Press \"Enter\" to continue):", title))

	scanner.Scan()

	return p.set(key, scanner.Text())
}

type keyMap struct {
	m    map[string]string
	keys []string
}

func (k *keyMap) set(key string, value string) {
	k.m[key] = value
	k.keys = append(k.keys, key)
}

// Prepare prepares the bot for deploy.
func (p *platform) install(name string) {
	key := fmt.Sprintf("RECIME_%v_ACCESS_TOKEN", strings.ToUpper(name))

	var err error

	switch strings.ToLower(name) {
	case "facebook":
		{
			fmt.Println("Please enter your facebook app settings")

			k := &keyMap{m: make(map[string]string)}

			k.set("RECIME_FACEBOOK_APP_ID", "App ID")
			k.set("RECIME_FACEBOOK_APP_SECRET", "App Secret")
			k.set("RECIME_FACEBOOK_ACCESS_TOKEN", "Page access token")

			for _, key := range k.keys {
				err = p.processInput(key, k.m[key])
			}
		}
	case "telegram":
		err = p.processInput(key, "Telegram access key")
	case "wechat":
		err = p.processInput(key, "WeChat access token")
	case "slack":
		{
			fmt.Println("Please enter Slack credentials")

			k := &keyMap{m: make(map[string]string)}

			k.set("RECIME_SLACK_CLIENT_ID", "Client ID")
			k.set("RECIME_SLACK_CLIENT_SECRET", "Client Secret")

			for _, key := range k.keys {
				err = p.processInput(key, k.m[key])
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
