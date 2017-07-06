package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Recime/recime-cli/lib"
	"github.com/Recime/recime-cli/shared"

	"github.com/howeyc/fsnotify"
	"github.com/mitchellh/go-homedir"
)

//Run runs the bot in a local node server.
func Run(watch bool) {
	uid := UID{}

	wd, err := os.Getwd()

	check(err)

	id := uid.Get(wd)

	tokens := strings.Split(template, "/")
	fileName := tokens[len(tokens)-1]
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	fileName = fmt.Sprintf("recime-bot-template-%s", fileName)

	home, err := homedir.Dir()

	check(err)

	home = fmt.Sprintf("%v/recime-cli-%v", filepath.ToSlash(home), Version)

	zipName := fmt.Sprintf("%s/%s.zip", home, fileName)

	_, err = os.Stat(home)

	if os.IsNotExist(err) {
		err = os.Mkdir(home, os.ModePerm)
		check(err)
	}

	h := &httpClient{}
	h.download(template, zipName)

	util.Unzip(zipName, home)

	templatedir := fmt.Sprintf("%s/%s", home, fileName)

	botdir := fmt.Sprintf("%s/%s", templatedir, id)

	util.CopyDir(filepath.ToSlash(wd), botdir)

	fmt.Println("INFO: Installing Dependencies...")

	sh := &shell{}

	pkg := &pkg{}

	pkg.sync(botdir, templatedir)

	sh.execute(templatedir, "install")
	sh.execute(botdir, "install")

	fmt.Println("INFO: Building...")

	if Build(botdir) != nil {
		return
	}

	fmt.Println("INFO: Starting...")

	if watch {
		watchDir(filepath.ToSlash(wd), botdir)
	}

	config := []shared.Config{shared.Config{Key: "BOT_UNIQUE_ID", Value: id}}

	_config := shared.Config{}
	// Add config user config
	reader, _ := _config.Open(wd)

	vars := _config.Get(reader)

	createFBPersistentMenu(vars["RECIME_FACEBOOK_ACCESS_TOKEN"])
	createFBGettingStarted(vars["RECIME_FACEBOOK_ACCESS_TOKEN"])

	syncConfigVars(id, renewToken().ID)

	for key, value := range vars {
		config = append(config, shared.Config{Key: key, Value: value})
	}

	sh = &shell{
		config: config,
	}

	sh.execute(templatedir, "start")
}

//watchDir watch file for changes
func watchDir(dir string, targetDir string) {
	watcher, err := fsnotify.NewWatcher()

	check(err)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				fileInfo, err := os.Stat(ev.Name)

				if !ev.IsAttrib() {
					if os.IsNotExist(err) {
						components := strings.Split(filepath.ToSlash(ev.Name), "/")
						name := components[len(components)-1]

						os.Remove(fmt.Sprintf("%s/%s", targetDir, string(name)))
					} else {
						targetFile := fmt.Sprintf("%s/%s", targetDir, fileInfo.Name())

						fmt.Println("INFO: File change event.")

						util.CopyFile(ev.Name, targetFile)
						Build(targetDir)

						fmt.Println("----")
					}
				}
			case err := <-watcher.Error:
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(dir)

	check(err)
}
