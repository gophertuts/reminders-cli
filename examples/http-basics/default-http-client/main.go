package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	res, err := http.Get("https://www.google.com")
	if err != nil {
		log.Fatalf("could not make request to google: %v", err)
	}
	bs, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("\nGET REQUEST:\n\n%s\n", string(bs))

	body := bytes.NewReader([]byte(`{"title": "some title"}`))
	res, err = http.Post("http://localhost:8080", "*/*", body)
	if err != nil {
		log.Fatalf("could not make request to localhost: %v", err)
	}
	bs, _ = ioutil.ReadAll(res.Body)
	fmt.Printf("\nPOST REQUEST:\n\n%s\n", string(bs))
}
