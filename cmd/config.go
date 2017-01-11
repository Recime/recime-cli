package cmd

import "fmt"
import "encoding/json"
import "path/filepath"
import "io/ioutil"
import "os"

//Config user configuration
type Config struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// IsExist checks for a file in a given path.
func IsExist(path string){
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		check(err)
	}
}

// GetConfigVars returns config map for a given path.
func GetConfigVars(wd string) map[string]interface{} {
	path := filepath.Join(wd, filepath.Join(".recime", "config.json"))
	
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)

	check(err)

	dat, err := ioutil.ReadAll(file)

	check(err)

	var config map[string]interface{}

	if len(dat) > 0 {
		json.Unmarshal(dat, &config)
	} else{
		config = make(map[string]interface{})
	}
	return config
}

// SetConfig adds / edits a config var.
func SetConfig(config Config){
	wd, err := os.Getwd()

	check(err)

	pkgPath := wd + "/package.json"

	IsExist(pkgPath)

	path := filepath.Join(wd, filepath.Join(".recime", "config.json"))
	
	os.MkdirAll(path, 0755)

	data := GetConfigVars(wd)

	data[config.Key] = config.Value

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0600)

	check(err)

	jsonBody, err := json.Marshal(data)

	check(err)

	file.Write(jsonBody)

	fmt.Println("\r\nINFO: Config Vars Set Successfully.\r\n")
}
