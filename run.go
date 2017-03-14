package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Recime/recime-cli/cmd"
	"github.com/Recime/recime-cli/util"

	"github.com/howeyc/fsnotify"
	"github.com/mitchellh/go-homedir"
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

					cmd.Build(targetDir)
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
func Run(watch bool) {
	uid := cmd.GetUID()

	tokens := strings.Split(template, "/")
	fileName := tokens[len(tokens)-1]
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	fileName = fmt.Sprintf("recime-bot-template-%s", fileName)

	home, err := homedir.Dir()

	check(err)

	home = filepath.ToSlash(home) + "/recime-cli"

	zipName := fmt.Sprintf("%s/%s.zip", home, fileName)

	_, err = os.Stat(home)

	if os.IsNotExist(err) {
		err = os.Mkdir(home, os.ModePerm)
		check(err)
	}

	download(template, zipName)

	util.Unzip(zipName, home)

	templatedir := fmt.Sprintf("%s/%s", home, fileName)

	wd, err := os.Getwd()

	check(err)

	botdir := fmt.Sprintf("%s/%s", templatedir, uid)

	fmt.Println("INFO: Deploying Bot...")

	cmd.Build(wd)

	util.CopyDir(filepath.ToSlash(wd), botdir)

	fmt.Println("INFO: Installing Dependencies...")

	installCmd := []string{"npm", "install"}

	shell := &shell{}

	pkg := &pkg{}

	pkg.sync(botdir, templatedir)

	shell.execute(installCmd, templatedir, nil)
	shell.execute(installCmd, botdir, nil)

	fmt.Println("INFO: Starting...")

	if watch {
		WatchForChanges(filepath.ToSlash(wd), botdir)
	}

	config := []cmd.Config{cmd.Config{Key: "BOT_UNIQUE_ID", Value: uid}}
	config = append(config, cmd.Config{Key: "BASE_URL", Value: baseURL})

	_config := cmd.Config{}
	// Add config user config
	reader, _ := _config.Open(wd)

	vars := _config.Get(reader)

	for key, value := range vars {
		config = append(config, cmd.Config{Key: key, Value: value})
	}

	shell.execute([]string{"npm", "start"}, templatedir, config)
}
