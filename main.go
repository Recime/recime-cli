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
	"regexp"
	"strings"
	"time"

	"github.com/Recime/recime-cli/cmd"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// PrintStatus outputs formatted status.
func PrintStatus(status string) {
	console := color.New(color.FgHiMagenta)

	console.Print("INFO: ")

	console = color.New(color.FgHiCyan)

	console.Println(fmt.Sprintf("%v", status))
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		fmt.Println("")
		os.Exit(1)
	}()

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

	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Scaffolds the bot from an interactive prompt",
		Long:  `Scaffolds the necessary files required for the bot to work correctly in Recime cloud from an interactive prompt`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				folder := strings.Join(args, " ")
				Create(folder)
			} else {
				fmt.Println("\n\rUSAGE: recime-cli create [folderName]\n\r")
			}
		},
	}

	var cmdDeploy = &cobra.Command{
		Use:   "deploy",
		Short: "Deploys the bot to Recime cloud",
		Long:  "Prepares and deploys to bot to Recime cloud. Installs the node modules defined in package.json, validates provides the endpoint to test the bot",
		Run: func(_ *cobra.Command, args []string) {
			wd, err := os.Getwd()

			check(err)

			if cmd.Build(wd) != nil {
				os.Exit(1)
			}

			Deploy()
		},
	}

	var cmdInstall = &cobra.Command{
		Use:   "install",
		Short: "Installs the dependencies",
		Long:  "Installs the required dependencies for the bot to work in Recime cloud",
		Run: func(_ *cobra.Command, args []string) {
			cmd.Install()
		},
	}

	var cmdBuild = &cobra.Command{
		Use:   "build",
		Short: "Builds the bot module",
		Long:  "Builds the bot module. Uses the build script from pacakge.json",
		Run: func(_ *cobra.Command, args []string) {
			wd, err := os.Getwd()

			check(err)

			cmd.Build(wd)
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
				fmt.Println("INFO: Invalid Number of Arguments.")
			}
		},
	}

	cmdConfig.AddCommand(cmdConfigAdd)

	var watch bool

	var cmdRun = &cobra.Command{
		Use:   "run",
		Short: "Runs the bot locally",
		Long:  "Runs the bot locally",
		Run: func(_ *cobra.Command, args []string) {
			// install any dependencies
			cmd.Install()

			// execute run Command
			options := map[string]interface{}{
				"url":   appTemplateURL,
				"uid":   cmd.GetUID(),
				"base":  baseURL,
				"watch": watch,
			}
			cmd.Run(options)
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

	rootCmd.AddCommand(cmdInstall)
	rootCmd.AddCommand(cmdBuild)
	rootCmd.AddCommand(cmdConfig)
	rootCmd.AddCommand(cmdLogin)
	rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdDeploy)
	rootCmd.AddCommand(cmdRun)
	rootCmd.AddCommand(cmdPlatform)

	rootCmd.Execute()

	fmt.Println("")
	fmt.Println("For any questions and feedback, please reach us at hello@recime.io.")
	fmt.Println("")
}
