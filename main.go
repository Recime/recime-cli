package main

import "fmt"
import "os"
import "os/signal"
import "time"

import  "github.com/spf13/cobra"

func main() {
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)
    go func() {
    <-interrupt
    os.Exit(1)
    }()

    var cmdCreate = &cobra.Command {
      Use : "create",
      Short: "Scafolds the bot from an interactive prompt",
      Long : "Scalfolds the nesseary files required for the bot module to work correctly in Recime cloud via an interactive prompt",
      Run : func (cmd *cobra.Command, args []string){
          Create()
      },
    }

    var cmdDeploy = &cobra.Command {
      Use : "deploy",
      Short: "Deploys the bot module to Recime cloud",
      Long : "Prepares and deploys to bot to Recime cloud. Installs the node modules defined in package.json, validates provides the endpoint to test the bot",
      Run : func (cmd *cobra.Command, args []string){
          Deploy()
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

    rootCmd.AddCommand(cmdCreate)
    rootCmd.AddCommand(cmdDeploy)

    rootCmd.Execute()

    fmt.Println("\r\nFor any questions and feedbacks, please reach us at hello@recime.ai. \r\n")
}
