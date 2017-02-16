package cmd

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

//Config bot configuration
type Config struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Source string
}

// Get gets stored config.
func (c *Config) Get(reader io.Reader) map[string]interface{} {
	dat, _ := ioutil.ReadAll(reader)

	var config map[string]interface{}

	if len(dat) > 0 {
		json.Unmarshal(dat, &config)
	} else {
		config = make(map[string]interface{})
	}

	return config
}

// Open opens config from working directory.
func (c *Config) Open(wd string) (io.Reader, error) {
	path := filepath.Join(wd, filepath.Join(".recime", "config.json"))

	reader, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)

	return reader, err
}

// Save saves the config to disk.
func (c *Config) Save() {
	wd, err := os.Getwd()

	check(err)

	data := make(map[string]interface{})

	_filepath := filepath.Join(".recime", "config.json")

	target := filepath.Join(wd, _filepath)

	reader, err := c.Open(wd)

	if err != nil {
		err = os.MkdirAll(filepath.Dir(target), 0755)

		check(err)

		_, err := os.Stat(target)

		if os.IsNotExist(err) {
			os.Create(target)
		}
	} else {
		data = c.Get(reader)
	}

	data[c.Key] = c.Value

	file, err := os.OpenFile(target, os.O_WRONLY|os.O_TRUNC, 0600)

	check(err)

	jsonBody, err := json.Marshal(data)

	check(err)

	file.Write(jsonBody)

}
