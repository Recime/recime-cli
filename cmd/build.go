package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/fatih/color"
)

func execute(cmd *exec.Cmd) error {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	red := color.New(color.FgMagenta)

	if err != nil {
		red.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	red.Println(out.String())

	return err
}

func buildES6(dir string) error {
	var data map[string]interface{}

	path := fmt.Sprintf("%s/.babelrc", dir)

	buff, err := ioutil.ReadFile(path)

	check(err)

	if err := json.Unmarshal(buff, &data); err != nil {
		panic(err)
	}

	args := []string{dir, "-d", dir}

	for key, value := range data {
		rt := reflect.TypeOf(value)

		switch rt.Kind() {
		case reflect.Slice:
			arr := value.([]interface{})

			values := make([]string, 0)

			for _, v := range arr {
				values = append(values, v.(string))
			}

			args = append(args, fmt.Sprintf("--%s=%s", key, strings.Join(values, ",")))
		default:
			args = append(args, fmt.Sprintf("--%s=%s", key, value.(string)))
		}
	}

	cmd := exec.Command("node_modules/.bin/babel", args...)
	cmd.Dir = dir

	return execute(cmd)
}

// Build builds the bot.
func Build(dir string) error {
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = dir

	err := execute(cmd)

	if _, err := os.Stat(fmt.Sprintf("%s/.babelrc", dir)); err == nil {
		return buildES6(dir)
	}

	return err
}
