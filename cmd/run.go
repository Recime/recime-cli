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

const (
	template = "https://s3-us-west-2.amazonaws.com/recime-io/package-with-container.zip"
)

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
func Run(source string, watch bool) {
	uid := GetUID()

	tokens := strings.Split(template, "/")
	fileName := tokens[len(tokens)-1]
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	home, err := homedir.Dir()

	check(err)

	home = filepath.ToSlash(home) + "/recime-cli"

	zipName := fmt.Sprintf("%s/%s.zip", home, fileName)

	_, err = os.Stat(home)

	if os.IsNotExist(err) {
		err = os.Mkdir(home, os.ModePerm)
		check(err)
	}

	Download(template, zipName)

	templateDir := fmt.Sprintf("%s/recime-bot-template", home)

	util.Unzip(zipName, templateDir)

	wd, err := os.Getwd()

	check(err)

	botDir := templateDir + "/" + uid

	fmt.Println("INFO: Deploying Bot...")

	Build(wd)

	util.CopyDir(filepath.ToSlash(wd), botDir)

	fmt.Println("INFO: Installing Dependencies...")

	installCmd := []string{"npm", "install"}

	runCmd(installCmd, templateDir, nil)
	runCmd(installCmd, botDir, nil)

	fmt.Println("INFO: Starting...")

	if watch {
		WatchForChanges(filepath.ToSlash(wd), botDir)
	}

	config := []Config{Config{Key: "BOT_UNIQUE_ID", Value: uid}}
	config = append(config, Config{Key: "BASE_URL", Value: source})

	_config := Config{}
	// Add config user config
	reader, _ := _config.Open(wd)

	vars := _config.Get(reader)

	for key, value := range vars {
		config = append(config, Config{Key: key, Value: value})
	}

	runCmd([]string{"npm", "start"}, templateDir, config)
}

func runCmd(args []string, wd string, config []Config) {
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
