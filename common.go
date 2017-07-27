package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Recime/recime-cli/shared"
	"github.com/briandowns/spinner"
)

func syncConfigVars(uid string, token string) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/%s", botBaseURL, uid)

	req, err := http.NewRequest("GET", url, nil)

	if len(token) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

	s.Start()

	res, err := client.Do(req)

	check(err)

	s.Stop()

	defer res.Body.Close()

	dat, err := ioutil.ReadAll(res.Body)

	check(err)

	var result bot

	json.Unmarshal(dat, &result)

	for key, value := range result.Config {
		if strings.Index(key, "SYSTEM_") != -1 {
			config := shared.Config{Key: key, Value: value, Source: apiEndpoint}
			config.Save()
		}
	}
}

func renewToken() shared.Token {
	token := shared.Token{Source: apiEndpoint}
	return token.Renew()
}
