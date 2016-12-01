package cmd

import "os"
import "os/exec"

func Build() {
	wd, err := os.Getwd()

	check(err)

	cmd := exec.Command("npm", "run", "build")

	cmd.Dir = wd

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

}
