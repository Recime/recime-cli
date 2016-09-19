package main

import "bytes"
import "fmt"
import "os"
import "os/signal"
// import "net/url"
import "io"
import "io/ioutil"

import "encoding/json"
import "bufio"
import "strings"
import "regexp"
import "crypto/md5"

func setValue(data map[string]interface{}, key string, value string) {
    if len(value) > 0 {
        data[key] = strings.TrimRight(value, "\n")
    }
}

func ProcesssInput(in io.Reader) (data map[string]interface{} ){
    scanner := bufio.NewScanner(in)

    asset := MustAsset("data/package.json")

    check(json.Unmarshal(asset, &data))

    fmt.Printf("Title (%s):", data["title"])

    scanner.Scan()

    title := scanner.Text()

    fmt.Printf("Description (%s):", data["description"])

    scanner.Scan()

    desc := scanner.Text()

    author := "Recime Inc."

    fmt.Printf("Author (%s):", author)

    scanner.Scan()

    _author := scanner.Text()

    if (len(_author) > 0) {
      author = scanner.Text()
    }

    email := "hello@recime.ai"

    fmt.Printf("Email (%s):", email)

    scanner.Scan()

    _email := scanner.Text()

    if (len(_email) > 0) {
      email = scanner.Text()
    }

    fmt.Printf("License (%s):", data["license"])

    scanner.Scan()

    license := scanner.Text()

    r, _ := regexp.Compile("[\\s]+")

    author = strings.Trim(author, " ")

    email = r.ReplaceAllString(email, "")

    data["author"] =  author + " " + "<" + email + ">"

    fmt.Println(data["author"])

    setValue(data, "title", title)
    setValue(data, "description", desc)
    setValue(data, "license", license)

    return data
}


func main() {
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)
    go func() {
    <-interrupt
    os.Exit(1)
    }()

    args := os.Args

    fmt.Printf("\r\nRecime Command Line Tool %s \r\n\r\n", VERSION)

    if len(args) == 1 {
        fmt.Println("Usage: recime-cli create")
        return
    }

    command := args[1]

    if command == "deploy" {
        Deploy()
    }

    if command == "create" {
      wd, err := os.Getwd()

      data := ProcesssInput(os.Stdin)

      name := data["title"].(string)

      r, _ := regexp.Compile("[\\s?.$#,()^!&]+")

      normalizedName := r.ReplaceAllString(name, "-")
      normalizedName = strings.ToLower(normalizedName)
      normalizedName = strings.TrimLeft(normalizedName, "_")

      data["name"] = normalizedName

      r, _ = regexp.Compile("[^<>]+")

      author := r.FindAllString(data["author"].(string), -1)

      r, _ = regexp.Compile("[\\s]+")

      _author := author[1]
      _author = r.ReplaceAllString(_author, "")
      _author = strings.ToLower(_author)

      uid := _author + ";" + normalizedName

      // fmt.Println(uid)

      _data := []byte(uid)

      uid = fmt.Sprintf("%x", md5.Sum(_data))

      data["uid"] = uid

      dir, err := os.Getwd()

      check(err)

      path := dir + "/" + name

      if _, err := os.Stat(path); os.IsNotExist(err) {
        si, err := os.Stat(wd)

        check(err)

        err = os.Mkdir(path, si.Mode())

        check(err)
      }

      resources, err := AssetDir("data")

      check(err)

      for key := range resources{
          entry := resources[key]

          asset := MustAsset("data/" + entry)

          if entry == "package.json" {
            asset, err = json.MarshalIndent(data, "", "\t")

            check(err)

            asset = bytes.Replace(asset, []byte("\\u003c"), []byte("<"), -1)
            asset = bytes.Replace(asset, []byte("\\u003e"), []byte(">"), -1)
            asset = bytes.Replace(asset, []byte("\\u0026"), []byte("&"), -1)
          }

          filePath := path + "/" + entry

          err = ioutil.WriteFile(filePath, asset, os.ModePerm)

          check(err)
      }
    }
    fmt.Println("\r\nFor any questions and feedbacks, please reach us at hello@recime.ai. \r\n")
}
