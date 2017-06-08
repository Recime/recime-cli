package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type httpClient struct {
}

// Download downloads url to a file name
func (h *httpClient) download(url string, fileName string) {
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

	_, err = io.Copy(output, response.Body)

	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
}

func (h *httpClient) post(url string, data map[string]interface{}) []byte {
	jsonBody, err := json.Marshal(data)

	check(err)

	r := bytes.NewBuffer(jsonBody)

	resp, err := http.Post(url, "application/json; charset=utf-8", r)

	check(err)

	defer resp.Body.Close()

	dat, _ := ioutil.ReadAll(resp.Body)

	return dat
}
