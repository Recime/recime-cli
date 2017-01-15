package cmd

import "fmt"
import "net/http"
import "encoding/json"
import "path/filepath"
import "io"
import "io/ioutil"
import "os"
import "time"
import "github.com/briandowns/spinner"

//Config user configuration
type Config struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Source string
}

func OpenConfig(wd string) (io.Reader, error){
	path := filepath.Join(wd, filepath.Join(".recime", "config.json"))

	reader, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)

	return reader, err
}

// GetConfigVars returns config map for a given path.
func GetConfigVars(reader io.Reader) map[string]interface{} {
	dat, err := ioutil.ReadAll(reader)

	check(err)

	var config map[string]interface{}

	if len(dat) > 0 {
		json.Unmarshal(dat, &config)
	} else{
		config = make(map[string]interface{})
	}
	return config
}

func Sync(source string) []Config{
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

	s.Start()

	url := fmt.Sprintf("%s/api/bot/config/%s", source, GetUID())

	res, err := http.Get(url)

	check(err)

	var result struct {
		Config []Config `json:"config"`
	}

	bytes, err := ioutil.ReadAll(res.Body)

	json.Unmarshal(bytes, &result)

	s.Stop()

	defer res.Body.Close()

	return result.Config
}

// SetConfig adds / edits a config var.
func SaveConfig(config Config){
	wd, err := os.Getwd()

	check(err)

	data := make(map[string]interface{})

	_filepath := filepath.Join(".recime", "config.json")

	target := filepath.Join(wd, _filepath)

	reader, err  := OpenConfig(wd)

	if err != nil {
		_config := Sync(config.Source)

		for _, element := range _config {
				data[element.Key] = element.Value
		}

		err = os.MkdirAll(filepath.Dir(target), 0755)

		check(err)

		_, err := os.Stat(target)

		if os.IsNotExist(err) {
				os.Create(target)
		}
	} else {
		data = GetConfigVars(reader)
	}

	data[config.Key] = config.Value

	file, err := os.OpenFile(target, os.O_WRONLY|os.O_TRUNC, 0600)

	check(err)

	jsonBody, err := json.Marshal(data)

	check(err)

	file.Write(jsonBody)

	fmt.Println("\r\nINFO: Config Variables Set Successfully.\r\n")
}
