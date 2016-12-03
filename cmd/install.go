package cmd

import "fmt"
import "os"
import "os/exec"

// PrintError prints error message
func PrintError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("\n\r> Error: %s\n", err.Error()))
	}
}

// Install installs dependencies.
func Install() {
	fmt.Println("INFO: Installing Dependencies.")

	wd, err := os.Getwd()

	check(err)

	cmd := exec.Command("npm1", "install")

	cmd.Dir = wd

	cmd.Stdout = os.Stdout

	err = cmd.Run()

	PrintError(err)
}
