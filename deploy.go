package main

import "fmt"
import "os"

import "bufio"

import "bytes"

import "io"
import "io/ioutil"

import "encoding/json"
// import "net/http"
// import "strings"
import "os/exec"

import "net/http"

// import "mime/multipart"
// import fs "path/filepath"

// import "gopkg.in/cheggaaa/pb.v1"
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

  // body := &bytes.Buffer{}
  //
  // writer := multipart.NewWriter(body)
  //
  // params := map[string]string {
  //     "uid" : uid,
  // }
  //
  // part, err := writer.CreateFormFile("file", fs.Base(filePath))
  //
  // check(err)
  //
  // _, err = io.Copy(part, file)
  //
  // for key, val := range params {
  //     _ = writer.WriteField(key, val)
  // }
  // err = writer.Close()
  //
  // check(err)

  //read file content to buffer

  fileInfo, _ := file.Stat()

  var size = fileInfo.Size()

  buffer := make([]byte, size)

  // read file content to buffer
  file.Read(buffer)

  fmt.Println("INFO: Preparing to upload.")

  url := "http://recimedev-env.us-west-1.elasticbeanstalk.com/signed-url"

  fileType := http.DetectContentType(buffer)

  bot := Bot { Id : uid, Type : fileType, }

  jsonBody, err := json.Marshal(bot)

  check(err)

  signedUrl := SendRequest(url, bytes.NewBuffer(jsonBody))

  // fmt.Println(signedUrl)

  // url = "http://localhost:8081/module/deploy/" + name

  fileReader := bytes.NewReader(buffer)

  // bar := pb.New(int(size)).SetUnits(pb.U_BYTES)
  //
  // bar.Start()

  // _ = bar.NewProxyReader(fileReader)


  // proxy := &ioprogress.Reader{
  //   Reader: fileReader,
  //   Size:   size,
  // }

  req, err := http.NewRequest("PUT", signedUrl, fileReader)

  s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)  // Build our new spinner

  s.Start()

  check(err)

  // bar.Finish()

  resp, err := http.DefaultClient.Do(req)

  check(err)

  defer resp.Body.Close()

  dat, err := ioutil.ReadAll(resp.Body)

  s.Stop()

  if len(dat) == 0 {
    fmt.Println("INFO: Finalizing.")
  }


  //bar.Finish()

  url = "http://recimedev-env.us-west-1.elasticbeanstalk.com/module/deploy/" + name

  s.Start()

  resp, err = http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(jsonBody))

  check(err)

  reader := bufio.NewReader(resp.Body)

  for {
      line, err := reader.ReadBytes('\n')
      line = bytes.TrimSpace(line)

      fmt.Println(string(line))

      if err == io.EOF {
        s.Stop()
        break
      }
  }

  defer resp.Body.Close()

}
