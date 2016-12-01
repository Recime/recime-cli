package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Recime/recime-cli/util"
)

// ExecuteInDir executes command in a given directory
func ExecuteInDir(args []string, wd string) {
	cmd := exec.Command(args[0], args[1])

	cmd.Dir = wd

	if len(args) == 3 {
		env := os.Environ()
		env = append(env, fmt.Sprintf("BOT_UNIQUE_ID=%s", args[2]))
		cmd.Env = env
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}

// UserHomeDir returns user home directory
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

//Run runs the bot in a local node server.
func Run(options map[string]interface{}) {
	url := options["url"].(string)
	uid := options["uid"].(string)

	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	version := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	homedir := filepath.ToSlash(UserHomeDir()) + "/recime-app"

	fileName = fmt.Sprintf("%s/recime-%s.zip", homedir, version)

	Download(url, fileName)

	target := homedir

	check(Unzip(fileName, target))

	fmt.Println("INFO: Preparing.")

	templateDir := target + "/recime-template-" + version

	fmt.Println(templateDir)

	wd, err := os.Getwd()

	check(err)

	botDir := templateDir + "/" + uid

	util.CopyDir(filepath.ToSlash(wd), botDir)

	fmt.Println("INFO: Deploying Bot.")

	ExecuteInDir([]string{"npm", "install"}, templateDir)

	fmt.Println("INFO: Starting.")

	ExecuteInDir([]string{"npm", "start", uid}, templateDir)
}

// Download downloads url to a file name
func Download(url string, fileName string) {
	fmt.Println("Downloading", url)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}

// Unzip unzips a given archive to a target
func Unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)

	if err != nil {
		return err
	}

	if err = os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}
