package cmd

import "fmt"
import "os"
import "os/exec"


func Install() (error){
  fmt.Println("INFO: Installing package dependencies")

  wd, err := os.Getwd()

  check(err)

  cmd := exec.Command("npm", "install")

  cmd.Dir = wd

  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  return cmd.Run()
}
