package main

import "fmt"
import "os"
import "os/signal"
import "time"

import c "./cmd"

import  "github.com/spf13/cobra"

func main() {
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)
    go func() {
    <-interrupt
    os.Exit(1)
    }()

    var cmdInit = &cobra.Command {
      Use : "init",
      Short: "Initializes your Recime account",
      Long: `Initializes the CLI with your Recime account. You need to create and verify your account from https://recime.ai in order to get started.`,
      Run : func (cmd *cobra.Command, args []string){
          options := map[string]interface{}{
            "in" : os.Stdin,
            "base" : BASE_URL,
          }
          c.Init(options)
      },
    }

    var cmdCreate = &cobra.Command {
      Use : "create",
      Short: "Scaffolds the bot from an interactive prompt",
      Long : `Scaffolds the necessary files required for the bot to work correctly in Recime cloud from an interactive prompt`,
      Run : func (cmd *cobra.Command, args []string){
          Create(ValidateUser())
      },
    }

    var cmdDeploy = &cobra.Command {
      Use : "deploy",
      Short: "Deploys the bot to Recime cloud",
      Long : "Prepares and deploys to bot to Recime cloud. Installs the node modules defined in package.json, validates provides the endpoint to test the bot",
      Run : func (cmd *cobra.Command, args []string){
        Deploy(ValidateUser())
      },
    }

    var rootCmd = &cobra.Command{
      Use : "recime-cli",
      Long: fmt.Sprintf(`Recime Command Line Interface
Version %v
Copyright %d Recime, Inc.
https://recime.ai`,
        VERSION,
        time.Now().Year(),
      ),
    }

    rootCmd.AddCommand(cmdInit)
    rootCmd.AddCommand(cmdCreate)
    rootCmd.AddCommand(cmdDeploy)

    rootCmd.Execute()

    fmt.Println("\r\nFor any questions and feedbacks, please reach us at hello@recime.ai. \r\n")
}

func ValidateUser()(User) {
  user, err := GetStoredUser()

  if err != nil {
        fmt.Println("\x1b[31;1mInvalid account. Please run \"recime-cli init\" to get started.\x1b[0m")
        os.Exit(1)
  }
  return user
}
