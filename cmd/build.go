package cmd

import (
	"bytes"
	"os/exec"

	"github.com/fatih/color"
)

// Build builds the bot.
func Build(dir string) error {

	// if _, err := os.Stat(fmt.print); err == nil {
	// 	// path/to/whatever exists
	// }

	cmd := exec.Command("npm", "run", "build")

	cmd.Dir = dir

	var out bytes.Buffer

	cmd.Stdout = &out

	err := cmd.Run()

	red := color.New(color.FgMagenta)
	red.Println(out.String())

	return err
}
