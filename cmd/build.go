package cmd

import "os"
import "os/exec"

func Build() error {
	wd, err := os.Getwd()

	check(err)

	cmd := exec.Command("npm", "run", "build")

	cmd.Dir = wd

	cmd.Stdout = os.Stdout

	return cmd.Run()

}
