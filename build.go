package main

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/fatih/color"
)

func execute(cmd *exec.Cmd) error {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		red := color.New(color.FgHiRed)
		red.Println(fmt.Sprint(err) + ": " + stderr.String())
	}

	cyan := color.New(color.FgHiMagenta)
	cyan.Println(out.String())

	return err
}

// Build builds the bot.
func Build(dir string) error {
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = dir

	err := execute(cmd)

	return err
}
