package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Recime/recime-cli/shared"
)

type shell struct {
	config []shared.Config
}

func (sh *shell) execute(wd string, arg ...string) {
	cmd := exec.Command("npm", arg...)

	cmd.Dir = wd

	if sh.config != nil {

		env := os.Environ()

		for _, c := range sh.config {
			env = append(env, fmt.Sprintf("%s=%s", c.Key, c.Value))
		}

		cmd.Env = env
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
