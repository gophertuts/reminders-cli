package main

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

type idsFlag []string

func (ids idsFlag) String() string {
	return strings.Join(ids, ",")
}

func (ids *idsFlag) Set(v string) error {
	*ids = append(*ids, v)
	return nil
}

type person struct {
	name string
	born time.Time
}

func (p *person) String() string {
	return fmt.Sprintf("person %s was born: %v", p.name, p.born)
}

func (p *person) Set(v string) error {
	p.name = v
	p.born = time.Now()
	return nil
}

func main() {
	var ids idsFlag
	var p person

	flag.Var(&ids, "ids", "list of ids")
	flag.Var(&p, "person", "person name")
	flag.Parse()

	fmt.Println(ids)
	fmt.Println(p.name)
	fmt.Println(p.born)
}
