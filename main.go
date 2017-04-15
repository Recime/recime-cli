// Copyright 2017 The Recime Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Recime/recime-cli/cmd"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// PrintStatus outputs formatted status.
func PrintStatus(status string) {
	pattern := regexp.MustCompile(`[a-z1-9A-Z.]+`)
	if pattern.MatchString(status) {
		fmt.Print("INFO: ")
		fmt.Println(fmt.Sprintf("%v", status))
	} else {
		fmt.Print(status)
	}
}

func checkMainFolder() {
	home, err := homedir.Dir()

	check(err)

	home = filepath.ToSlash(home) + "/recime-cli"

	_, err = os.Stat(home)

	if os.IsNotExist(err) {
		err = os.Mkdir(home, os.ModePerm)
		check(err)
	}
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		fmt.Println("")
		os.Exit(1)
	}()

	checkMainFolder()

	var cmdLogin = &cobra.Command{
		Use:   "login",
		Short: "Logs into your Recime account",
		Long:  `Logs into your Recime account. You need to create and verify your account from https://recime.io in order to get started.`,
		Run: func(_ *cobra.Command, args []string) {
			options := map[string]interface{}{
				"in":   os.Stdin,
				"base": apiEndpoint,
			}
			cmd.Login(options)
		},
	}

	var lang string

	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Scaffolds the bot from an interactive prompt",
		Long:  `Scaffolds the necessary files required for the bot to work correctly in Recime cloud from an interactive prompt`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				folder := strings.Join(args, " ")
				Create(folder, lang)
			} else {
				fmt.Println("\n\rUSAGE: recime-cli create [folderName|.]\n\r")
			}
		},
	}

	cmdCreate.PersistentFlags().StringVarP(&lang, "lang", "l", "es6", "Specifies the language of the template.")

	var cmdDeploy = &cobra.Command{
		Use:   "deploy",
		Short: "Deploys the bot to Recime cloud",
		Long:  "Prepares and deploys to bot to Recime cloud. Installs the node modules defined in package.json, validates provides the endpoint to test the bot",
		Run: func(_ *cobra.Command, args []string) {
			Deploy()
		},
	}

	var cmdConfig = &cobra.Command{
		Use:   "config",
		Short: "Initializes bot config vars",
		Long:  "Add or edit bot config vars that will be acccessed via `process.env` from within bot module",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("\r\nUSAGE: recime-cli config set NAME=Joe Smith")
			}
		},
	}
	var cmdConfigAdd = &cobra.Command{
		Use:   "set",
		Short: "Sets a new or existing config var",
		Long:  "Sets a new or existing config var",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 1 {
				pattern := regexp.MustCompile(`[a-zA-Z][1-9a-zA-Z_]+=[0-9a-zA-Z]+`)

				if pattern.MatchString(args[0]) {
					pair := strings.Split(args[0], "=")

					config := cmd.Config{Key: pair[0], Value: pair[1], Source: apiEndpoint}

					config.Save()

					fmt.Println("\r\nINFO: Config Var Set Successfully.")
				} else {
					red := color.New(color.FgRed).Add(color.Bold)
					red.Println("\r\nERROR: Invalid Key-Value Pair!")
				}

			} else {
				red := color.New(color.FgRed).Add(color.Bold)
				red.Println("\r\nERROR: Missing arguments. USAGE: recime-cli config set KEY_NAME=value")
			}
		},
	}

	var cmdConfigRemove = &cobra.Command{
		Use:   "unset",
		Short: "Removes a config var",
		Long:  "Removes a config var",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 1 {
				config := cmd.Config{
					Key: args[0],
				}

				if config.Remove() != nil {
					console := color.New(color.FgHiMagenta)
					console.Println("\r\nINFO: Config key removed successfully.\r\n")
				}
			} else {
				red := color.New(color.FgRed).Add(color.Bold)
				red.Println("\r\nERROR: Missing arguments. USAGE: recime-cli config unset KEY_NAME")
			}
		},
	}

	cmdConfig.AddCommand(cmdConfigAdd)
	cmdConfig.AddCommand(cmdConfigRemove)

	var cmdPlugin = &cobra.Command{
		Use:   "plugins",
		Short: "Installs third-party plugins",
		Long:  "Install third-party plugins for bot to take advantage",
		Run: func(_ *cobra.Command, args []string) {

		},
	}

	var apikey string

	var cmdPluginAdd = &cobra.Command{
		Use:   "add",
		Short: "Adds a new Plugin",
		Long:  "Adds a new Plugin",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 1 && len(apikey) > 0 {
				p := &plugins{APIKey: apikey}
				p.Add(args[0])
			} else {
				red := color.New(color.FgRed).Add(color.Bold)
				red.Println("\r\nERROR: Missing arguments. USAGE: recime-cli plugins add pluginname --apikey=[KEY]")
			}
		},
	}

	var cmdPluginRemove = &cobra.Command{
		Use:   "remove",
		Short: "Remove a Plugin",
		Long:  "Removes a Plugin",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 1 {
				p := &plugins{}
				p.Remove(args[0])
			} else {
				red := color.New(color.FgRed).Add(color.Bold)
				red.Println("\r\nERROR: Missing arguments. USAGE: recime-cli plugins remove pluginname")
			}
		},
	}

	cmdPluginAdd.PersistentFlags().StringVarP(&apikey, "apikey", "k", "", "Sets the api key")

	cmdPlugin.AddCommand(cmdPluginAdd)
	cmdPlugin.AddCommand(cmdPluginRemove)

	var watch bool

	var cmdRun = &cobra.Command{
		Use:   "run",
		Short: "Runs the bot locally",
		Long:  "Runs the bot locally",
		Run: func(_ *cobra.Command, args []string) {
			Run(watch)
		},
	}

	cmdRun.PersistentFlags().BoolVarP(&watch, "watch", "w", false, "Watches the bot folder for changes")

	var rootCmd = &cobra.Command{
		Use: "recime-cli",
		Long: fmt.Sprintf(`Recime Command Line Interface
Version %v
Copyright %d Recime, Inc.
%s`,
			Version,
			time.Now().Year(),
			SiteURL,
		),
	}

	var cmdPlatform = &cobra.Command{
		Use:   "platform",
		Short: "Initializes Platform",
		Long:  "Initializes the bot to be used in the target platform",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 0 || (len(args) == 1 && args[0] != "config") {
				red := color.New(color.FgRed).Add(color.Bold)
				red.Println("\r\nUSAGE: recime-cli platform config facebook")
			}
		},
	}

	var cmdPlaformConfig = &cobra.Command{
		Use:   "config",
		Short: "Initializes a new platform",
		Long:  "Initializes a new platform",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 1 {
				p := &platform{}
				p.install(args[0])
			} else {
				red := color.New(color.FgRed).Add(color.Bold)
				red.Println("\r\nERROR: Missing Platform. USAGE: recime-cli platform config facebook")
			}
		},
	}

	cmdPlatform.AddCommand(cmdPlaformConfig)

	rootCmd.AddCommand(cmdConfig)
	rootCmd.AddCommand(cmdLogin)
	rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdDeploy)
	rootCmd.AddCommand(cmdRun)
	rootCmd.AddCommand(cmdPlatform)
	rootCmd.AddCommand(cmdPlugin)

	rootCmd.Execute()

	fmt.Println("")
	fmt.Println("For any questions and feedback, please reach us at hello@recime.io.")
	fmt.Println("")
}
