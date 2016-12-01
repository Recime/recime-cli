// Copyright 2016 The Recime Inc. All rights reserved.
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
	"time"

	"strings"

	"github.com/Recime/recime-cli/cmd"
	"github.com/spf13/cobra"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		fmt.Println("")
		os.Exit(1)
	}()

	var cmdInit = &cobra.Command{
		Use:   "init",
		Short: "Initializes your Recime account",
		Long:  `Initializes the CLI with your Recime account. You need to create and verify your account from https://recime.ai in order to get started.`,
		Run: func(_ *cobra.Command, args []string) {
			options := map[string]interface{}{
				"in":   os.Stdin,
				"base": BaseURL,
			}
			cmd.Init(options)
		},
	}

	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Scaffolds the bot from an interactive prompt",
		Long:  `Scaffolds the necessary files required for the bot to work correctly in Recime cloud from an interactive prompt`,
		Run: func(cmd *cobra.Command, args []string) {
			folder := "."
			if len(args) > 0 {
				folder = strings.Join(args, " ")
			}
			Create(folder)
		},
	}

	var cmdDeploy = &cobra.Command{
		Use:   "deploy",
		Short: "Deploys the bot to Recime cloud",
		Long:  "Prepares and deploys to bot to Recime cloud. Installs the node modules defined in package.json, validates provides the endpoint to test the bot",
		Run: func(_ *cobra.Command, args []string) {
			// install any dependencies
			cmd.Install()

			// execute run Command

			cmd.Build()

			uid := cmd.Prepare()

			Deploy(uid)
		},
	}

	var cmdInstall = &cobra.Command{
		Use:   "install",
		Short: "Installs the dependencies",
		Long:  "Installs the requried dependencies for the bot to work in Recime cloud",
		Run: func(_ *cobra.Command, args []string) {
			cmd.Install()
		},
	}

	var cmdBuild = &cobra.Command{
		Use:   "build",
		Short: "Builds the bot module",
		Long:  "Builds the bot module. Uses the build script from pacakge.json",
		Run: func(_ *cobra.Command, args []string) {
			cmd.Build()
		},
	}

	var cmdRun = &cobra.Command{
		Use:   "run",
		Short: "Runs the bot locally",
		Long:  "Runs the bot locally",
		Run: func(_ *cobra.Command, args []string) {
			options := map[string]interface{}{
				"url": AppTemplateURL,
				"uid": cmd.Prepare(),
			}
			cmd.Run(options)
		},
	}

	var rootCmd = &cobra.Command{
		Use: "recime-cli",
		Long: fmt.Sprintf(`Recime Command Line Interface
Version %v
Copyright %d Recime, Inc.
https://recime.ai`,
			Version,
			time.Now().Year(),
		),
	}

	rootCmd.AddCommand(cmdInstall)
	rootCmd.AddCommand(cmdBuild)
	rootCmd.AddCommand(cmdInit)
	rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdDeploy)
	rootCmd.AddCommand(cmdRun)

	rootCmd.Execute()

	fmt.Println("")
	fmt.Println("For any questions and feedback, please reach us at hello@recime.ai.")
	fmt.Println("")
}
