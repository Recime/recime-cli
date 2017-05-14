package shared

import (
	"encoding/json"
	"errors"
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
func (c *Config) Get(reader io.Reader) map[string]string {
	dat, _ := ioutil.ReadAll(reader)

	var config map[string]string

	if len(dat) > 0 {
		json.Unmarshal(dat, &config)
	} else {
		config = make(map[string]string)
	}

	return config
}

func (c *Config) getTargetPath(wd string) string {
	_filepath := filepath.Join(".recime", "config.json")

	path := filepath.Join(wd, _filepath)

	return path
}

// Open opens config from working directory.
func (c *Config) Open(wd string) (io.Reader, error) {
	path := filepath.Join(wd, filepath.Join(".recime", "config.json"))

	reader, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)

	return reader, err
}

// Remove remvoes the config.
func (c *Config) Remove() error {
	wd, err := os.Getwd()

	check(err)

	reader, _ := c.Open(wd)

	data := c.Get(reader)

	if len(data[c.Key]) > 0 {
		delete(data, c.Key)
		c.saveDataToFile(data)
		return nil
	}

	return errors.New("Failed to unset the config var")
}

// Save saves the config to disk.
func (c *Config) Save() {
	wd, err := os.Getwd()

	check(err)

	data := make(map[string]string)

	path := c.getTargetPath(wd)

	reader, err := c.Open(wd)

	if err != nil {
		err = os.MkdirAll(filepath.Dir(path), 0755)

		check(err)

		_, err := os.Stat(path)

		if os.IsNotExist(err) {
			os.Create(path)
		}
	} else {
		data = c.Get(reader)
	}

	data[c.Key] = c.Value

	c.saveDataToFile(data)
}

func (c *Config) saveDataToFile(data map[string]string) {
	wd, _ := os.Getwd()
	path := c.getTargetPath(wd)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0600)

	check(err)

	jsonBody, err := json.Marshal(data)

	check(err)

	file.Write(jsonBody)
}

