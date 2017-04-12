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

	pkg := &pkg{
		UID:  uid,
		Name: name,
	}

	jsonBody, err := json.Marshal(pkg)

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

	s.Start()

	res, err := http.Post(source, "application/json; charset=utf-8", bytes.NewBuffer(jsonBody))

	defer res.Body.Close()

	check(err)

	var data struct {
		Name         string       `json:"name"`
		Dependencies []dependency `json:"dependencies"`
		Message      string       `json:"message"`
	}

	bytes, err := ioutil.ReadAll(res.Body)

	check(err)

	json.Unmarshal(bytes, &data)

	if len(data.Name) > 0 {
		key := fmt.Sprintf("%s_API_KEY", strings.ToUpper(data.Name))

		config := cmd.Config{Key: key, Value: p.APIKey, Source: source}

		config.Save()

		pkg.save(data.Dependencies)

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

func (p *plugins) Remove(name string) {
	source := fmt.Sprintf("%s/plugin", apiEndpoint)

	uid := cmd.GetUID()

	pkg := &pkg{
		UID:  uid,
		Name: name,
	}

	jsonBody, err := json.Marshal(pkg)

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

	s.Start()

	reader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("DELETE", source, reader)
	req.Header.Set("Content-Type", "application/json")

	check(err)

	res, err := http.DefaultClient.Do(req)

	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)

	check(err)

	s.Stop()

	pkg.remove()

	config := cmd.Config{
		Key: fmt.Sprintf("%s_API_KEY", strings.ToUpper(name)),
	}

	config.Remove()

	console := color.New(color.FgHiMagenta)
	console.Println("\r\nINFO: Plugin removed succesfully from your project.\r\n")
}
