package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Recime/recime-cli/cmd"
)

type shell struct {
}

func (c *shell) execute(args []string, wd string, config []cmd.Config) {
	cmd := exec.Command(args[0], args[1])

	cmd.Dir = wd

	if config != nil {

		env := os.Environ()

		for _, c := range config {
			env = append(env, fmt.Sprintf("%s=%s", c.Key, c.Value))
		}

		cmd.Env = env
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
