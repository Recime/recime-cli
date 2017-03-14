package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Recime/recime-cli/cmd"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

type plugins struct {
	APIKey string
}

func (p *plugins) Add(name string) {
	name = strings.ToLower(name)

	source := fmt.Sprintf("%s/plugin", apiEndpoint)

	uid := cmd.GetUID()

	_pkg := &pkg{
		UID:    uid,
		Name:   name,
		APIKey: p.APIKey,
	}

	jsonBody, err := json.Marshal(p)

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

	s.Start()

	res, err := http.Post(source, "application/json; charset=utf-8", bytes.NewBuffer(jsonBody))

	defer res.Body.Close()

	check(err)

	var data struct {
		Config  cmd.Config `json:"config"`
		Message string     `json:"message"`
	}

	bytes, err := ioutil.ReadAll(res.Body)

	check(err)

	json.Unmarshal(bytes, &data)

	if len(data.Config.Key) > 0 {
		config := cmd.Config{Key: data.Config.Key, Value: data.Config.Value, Source: source}

		config.Save()

		_pkg.save()

	} else {
		red := color.New(color.FgRed).Add(color.Bold)
		red.Printf("\r\nERROR: %s\r\n", data.Message)
		return
	}

	s.Stop()

	console := color.New(color.FgHiMagenta)

	fmt.Println("")

	console.Println("INFO: Plugin added succesfully to your project.")

	fmt.Println("")
}
