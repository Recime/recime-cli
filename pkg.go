package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

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

func (p *pkg) getFilePath() string {
	return filepath.Join(".recime", "plugins.json")
}

func (p *pkg) getTargetPath(wd string) string {
	target := filepath.Join(wd, p.getFilePath())

	return target
}

func (p *pkg) open(wd string) (io.Reader, error) {
	path := filepath.Join(wd, filepath.Join(".recime", "plugins.json"))

	reader, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)

	return reader, err
}

func (p *pkg) read(wd string) map[string]interface{} {
	data := make(map[string]interface{})
	reader, err := p.open(wd)

	if err != nil {
		targetPath := p.getTargetPath(wd)
		err = os.MkdirAll(filepath.Dir(targetPath), 0755)

		check(err)

		_, err := os.Stat(targetPath)

		if os.IsNotExist(err) {
			os.Create(targetPath)
		}
	} else {
		data = p.get(reader)
	}

	return data
}

func (p *pkg) remove() {
	wd, err := os.Getwd()

	check(err)

	data := p.read(wd)

	delete(data, p.Name)

	p.saveDataToDisk(data)
}

func (p *pkg) sync(source string, dest string) {
	reader, _ := p.open(source)

	if reader != nil {
		data := p.get(reader)

		var meta map[string]interface{}

		fp := fmt.Sprintf("%s/package.json", dest)

		buff, err := ioutil.ReadFile(fp)

		check(err)

		if err := json.Unmarshal(buff, &meta); err != nil {
			panic(err)
		}

		deps := meta["dependencies"].(map[string]interface{})

		for key, value := range data {
			deps[key] = value
		}

		ioutil.WriteFile(fp, MarshalIndent(meta), os.ModePerm)
	}

}

func (p *pkg) save(dependencies []dependency) {
	wd, err := os.Getwd()

	check(err)

	data := p.read(wd)

	for _, dep := range dependencies {
		data[dep.Name] = dep.Version
	}

	p.saveDataToDisk(data)
}

func (p *pkg) saveDataToDisk(data map[string]interface{}) {
	wd, err := os.Getwd()

	check(err)

	targetPath := p.getTargetPath(wd)

	file, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_TRUNC, 0600)

	check(err)

	jsonBody, err := json.Marshal(data)

	check(err)

	file.Write(jsonBody)
}
