package main

import "fmt"
import "os"
// import "net/url"
import "io/ioutil"
import "encoding/json"
// import "net/http"
import "strings"
import "os/exec"

func Deploy() {
  wd, err := os.Getwd()

  var data map[string]interface{}

  buff, err := ioutil.ReadFile(wd + "/package.json")

  check(err)

  if err := json.Unmarshal(buff, &data); err != nil {
    panic(err)
  }

  fmt.Println("Preparing.")

  // resp, err := http.PostForm("http://localhost:8081/module/deploy", v)
  //
  // check(err)
  //
  // defer resp.Body.Close()

  name := data["name"].(string)

  _name := strings.ToLower(name)
  _name = strings.Replace(_name, " ", "-", -1)

  // temp, err :=  ioutil.TempDir("", _name)
  //
  // check(err)
  //
  // err = CopyDir(wd, temp)
  //
  // check(err)

  cmd := exec.Command("npm", "install")

  cmd.Dir = wd

  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  err = cmd.Run()

  check(err)

  fmt.Println("Creating Archive.")
  //


  temp, err :=  ioutil.TempDir("", _name)

  check(err)

  dest := temp + "/.recime"

  err = os.Mkdir(dest, os.ModePerm)

  check(err)

  err = CopyDir(wd, dest)

  // fmt.Println(temp)

  Archive(dest, temp+ "/" + name + ".zip")

  // os.RemoveAll(dest + "/")

  fmt.Println("Done.")

  // body, err := ioutil.ReadAll(resp.Body)

  // fmt.Println(string(body))
}
