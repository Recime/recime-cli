package main

import "fmt"
import "os"

// import "bufio"

import "bytes"

import "io"
import "io/ioutil"

import "encoding/json"
// import "net/http"
// import "strings"
import "os/exec"

import "net/http"
import "mime/multipart"
import fs "path/filepath"


func Deploy() {

  wd, err := os.Getwd()

  var data map[string]interface{}

  buff, err := ioutil.ReadFile(wd + "/package.json")

  check(err)

  if err := json.Unmarshal(buff, &data); err != nil {
    panic(err)
  }

  name := data["name"].(string)

  uid := data["uid"].(string)

  fmt.Println("Installing modules")

  cmd := exec.Command("npm", "install")

  cmd.Dir = wd

  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  err = cmd.Run()

  check(err)

  fmt.Println("Compressing")

  temp, err :=  ioutil.TempDir("", name)

  check(err)

  dest := temp + "/" + uid

  err = os.Mkdir(dest, os.ModePerm)

  check(err)

  err = CopyDir(wd, dest)

  filePath := temp + "/" + name + ".zip"

  Archive(dest, filePath)

  file, err := os.Open(filePath)

  defer file.Close()

  body := &bytes.Buffer{}

  writer := multipart.NewWriter(body)

  params := map[string]string {
      "uid" : uid,
  }

  part, err := writer.CreateFormFile("file", fs.Base(filePath))

  check(err)

  _, err = io.Copy(part, file)

  for key, val := range params {
      _ = writer.WriteField(key, val)
  }
  err = writer.Close()

  check(err)

  fmt.Println("Preparing to Upload")

  req, err := http.NewRequest("POST", "http://localhost:8081/module/deploy/" + name, body)

  req.Header.Set("Content-Type", writer.FormDataContentType())

  check(err)

  resp, err := http.DefaultClient.Do(req)

  check(err)

  dat, err := ioutil.ReadAll(resp.Body)

  // for {
  //     line, err := reader.ReadBytes('\n')
  //     line = bytes.TrimSpace(line)
  //
  //     if err != nil {
  //       break
  //     }
  //
      fmt.Println(string(dat))
  // }

  defer resp.Body.Close()
}
