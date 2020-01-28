package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

type R []byte

func (r R) Read(bs []byte) (int, error) {
	for i, b := range r {
		bs[i] = b
	}
	return len(r), io.EOF
}

type W []byte

func (w W) Write(bs []byte) (int, error) {
	for i, b := range bs {
		w[i] = b
	}
	return len(w), io.EOF
}

func main() {
	r := R(`hello`)
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bs))

	w := make(W, len(r))
	io.Copy(w, r)
	fmt.Println(string(w))
}
