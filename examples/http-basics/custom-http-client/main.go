package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type customHTTPClient struct {
	client http.Client
}

func main() {
	httpClient := customHTTPClient{}
	body := bytes.NewReader([]byte(`{"title": "some title"}`))
	req, _ := http.NewRequest(http.MethodPatch, "http://localhost:8080", body)
	res, err := httpClient.client.Do(req)
	if err != nil {
		log.Fatalf("could not make request: %v", err)
	}
	bs, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(bs))
}
