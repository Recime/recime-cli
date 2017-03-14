package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Recime/recime-cli/cmd"
)

type pkg struct {
	Name   string `json:"name"`
	UID    string `json:"uid"`
	APIKey string `json:"apikey"`
}

func (p *pkg) get(reader io.Reader) map[string]interface{} {
	dat, _ := ioutil.ReadAll(reader)

	var config map[string]interface{}

	if len(dat) > 0 {
		json.Unmarshal(dat, &config)
	} else {
		config = make(map[string]interface{})
	}

	return config
}

func (p *pkg) open(wd string) (io.Reader, error) {
	path := filepath.Join(wd, filepath.Join(".recime", "plugins.json"))

	reader, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)

	return reader, err
}

func (p *pkg) sync(source string, dest string) {
	reader, _ := p.open(source)

	if reader != nil {
		data := p.get(reader)

		var meta map[string]interface{}

		fp := dest + "/package.json"

		buff, err := ioutil.ReadFile(fp)

		check(err)

		if err := json.Unmarshal(buff, &meta); err != nil {
			panic(err)
		}

		deps := meta["dependencies"].(map[string]interface{})

		for key, value := range data {
			deps[key] = value
		}

		ioutil.WriteFile(fp, cmd.MarshalIndent(meta), os.ModePerm)
	}

}

func (p *pkg) save() {
	data := make(map[string]interface{})

	wd, err := os.Getwd()

	check(err)

	_filepath := filepath.Join(".recime", "plugins.json")

	target := filepath.Join(wd, _filepath)

	reader, err := p.open(wd)

	if err != nil {
		err = os.MkdirAll(filepath.Dir(target), 0755)

		check(err)

		_, err := os.Stat(target)

		if os.IsNotExist(err) {
			os.Create(target)
		}
	} else {
		data = p.get(reader)
	}

	data[p.Name] = "latest"

	file, err := os.OpenFile(target, os.O_WRONLY|os.O_TRUNC, 0600)

	check(err)

	jsonBody, err := json.Marshal(data)

	check(err)

	file.Write(jsonBody)
}
