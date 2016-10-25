package cmd

import "fmt"
import "os"
import "os/exec"

func Install() {
	fmt.Println("INFO: Installing package dependencies")

	wd, err := os.Getwd()

	check(err)

	cmd := exec.Command("npm", "install")

	cmd.Dir = wd

	cmd.Stdout = os.Stdout

	cmd.Run()
}
