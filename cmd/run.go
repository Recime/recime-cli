package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Recime/recime-cli/util"
	"github.com/howeyc/fsnotify"
	"github.com/mitchellh/go-homedir"
)

// ExecuteInDir executes command in a given directory
func ExecuteInDir(args []string, wd string, config []Config) {
	cmd := exec.Command(args[0], args[1])

	cmd.Dir = wd

	if config != nil {

		env := os.Environ()

		for _, c := range config {
			env = append(env, fmt.Sprintf("%s=%s", c.Key, c.Value))
		}

		cmd.Env = env
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}

//WatchForChanges watch file for changes
func WatchForChanges(dir string, targetDir string) {
	watcher, err := fsnotify.NewWatcher()

	check(err)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if !ev.IsAttrib() {
					fmt.Println("INFO: File change event.")

					util.CopyDir(dir, targetDir)

					Build(targetDir)
				}
			case err := <-watcher.Error:
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(dir)

	check(err)
}

//Run runs the bot in a local node server.
func Run(options map[string]interface{}) {
	url := options["url"].(string)
	base := options["base"].(string)
	uid := options["uid"].(string)
	watch := options["watch"].(bool)

	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	version := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	home, err := homedir.Dir()

	check(err)

	home = filepath.ToSlash(home) + "/recime-cli"

	fileName = fmt.Sprintf("%s/recime-%s.zip", home, version)

	_, err = os.Stat(home)

	if os.IsNotExist(err) {
		err = os.Mkdir(home, os.ModePerm)
		check(err)
	}

	Download(url, fileName)

	target := home

	util.Unzip(fileName, target)

	fmt.Println("INFO: Preparing.")

	templateDir := target + "/recime-template-" + version

	wd, err := os.Getwd()

	botDir := templateDir + "/" + uid

	check(err)

	fmt.Println("INFO: Deploying Bot.")

	ExecuteInDir([]string{"npm", "install"}, templateDir, nil)

	fmt.Println("INFO: Starting.")

	Build(wd)

	util.CopyDir(filepath.ToSlash(wd), botDir)

	if watch {
		WatchForChanges(filepath.ToSlash(wd), botDir)
	}

	config := []Config{Config{Key: "BOT_UNIQUE_ID", Value: uid}}
	config = append(config, Config{Key: "BASE_URL", Value: base})

	// Add config user config
	reader, _ := OpenConfig(wd)

	vars := GetConfigVars(reader)

	for key, value := range vars {
		config = append(config, Config{Key: key, Value: value.(string)})
	}

	ExecuteInDir([]string{"npm", "start"}, templateDir, config)
}

// Download downloads url to a file name
func Download(url string, fileName string) {
	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)

	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
}
