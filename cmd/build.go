package cmd

import "os"
import "os/exec"

func Build(dir string) {
	cmd := exec.Command("npm", "run", "build")

	cmd.Dir = dir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

}
