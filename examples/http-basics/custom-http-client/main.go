package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type customHTTPClient struct {
	client http.Client
}

type reqBody struct {
	Title string `json:"title"`
}

func main() {
	httpClient := customHTTPClient{}
	reqBody := reqBody{
		Title: "some title",
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal("could not marshal body")
	}
	body := bytes.NewReader(bodyBytes)
	//body := bytes.NewReader([]byte(`{"title": "some title"}`))
	req, _ := http.NewRequest(http.MethodPatch, "http://localhost:8080", body)
	res, err := httpClient.client.Do(req)
	if err != nil {
		log.Fatalf("could not make request: %v", err)
	}
	bs, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(bs))
}
