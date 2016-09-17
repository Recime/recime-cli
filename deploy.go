package main

import "fmt"
import "os"

import "bytes"

import "io"
import "io/ioutil"

import "encoding/json"
import "os/exec"

import "net/http"

import "gopkg.in/cheggaaa/pb.v1"
import "github.com/briandowns/spinner"
import "time"


type Bot struct{
    Id string `json:"uid"`
    Type string `json:"fileType"`
}

func SendRequest(url string, body io.Reader) (string){
    res, err := http.Post(url, "application/json; charset=utf-8", body)

    check(err)

    var result struct {
        Url string `json:"url"`
    }

    bytes, err := ioutil.ReadAll(res.Body)

    json.Unmarshal(bytes, &result)

    defer res.Body.Close()

    // fmt.Println(string(res.Body))

    return result.Url
}


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

  fmt.Println("INFO: Installing modules.")

  cmd := exec.Command("npm", "install")

  cmd.Dir = wd

  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  err = cmd.Run()

  check(err)

  fmt.Println("INFO: Compressing.")

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

  fileInfo, _ := file.Stat()

  var size = fileInfo.Size()

  buffer := make([]byte, size)

  // read file content to buffer
  file.Read(buffer)

  fmt.Println("INFO: Preparing to upload.")

  url := BASE_URL + "/signed-url"

  fileType := http.DetectContentType(buffer)

  bot := Bot { Id : uid, Type : fileType, }

  jsonBody, err := json.Marshal(bot)

  check(err)

  signedUrl := SendRequest(url, bytes.NewBuffer(jsonBody))

  bar := pb.New(len(buffer)).SetUnits(pb.U_BYTES)

  bar.Format("[## ]")

  bar.Start()

  proxy := NewReader(buffer, bar)

  req, err := http.NewRequest("PUT", signedUrl, proxy)

  req.ContentLength = size

  check(err)

  // bar.Finish()

  resp, err := http.DefaultClient.Do(req)

  check(err)

  defer resp.Body.Close()

  dat, err := ioutil.ReadAll(resp.Body)

  check(err)

  fmt.Println(string(dat))

  if len(dat) == 0 {
    fmt.Println("INFO: Finalizing.")
  }

  url = BASE_URL + "/module/deploy/" + name

  s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)  // Build our new spinner

  s.Start()

  r := bytes.NewBuffer(jsonBody)

  resp, err = http.Post(url, "application/json; charset=utf-8", r)

  check(err)

  var result struct {
      Name string `json:"name"`
      Id string `json:"uid"`
      Message string `json:message`
  }

  defer resp.Body.Close()

  bytes, err := ioutil.ReadAll(resp.Body)

  check(err)

  json.Unmarshal(bytes, &result)

  s.Stop()

  time.Sleep(time.Millisecond * 10)

  if len(result.Name) > 0 {
    fmt.Println("=> " + BASE_URL + "/bot/" + result.Name)
    fmt.Println("INFO: Publish Successful")
    fmt.Println("For any questions and feedbacks, please reach us at hello@recime.ai.")
    return
  }
  fmt.Println("\x1b[31;1mFatal: Publish Failed!!!\x1b[0m")
}
