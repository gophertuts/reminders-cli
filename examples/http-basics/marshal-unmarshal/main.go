package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	numbers := []int{1, 2, 3}
	grades := map[string]int{
		"Math":    8,
		"English": 10,
	}
	type person struct {
		Name string `json:"name"` // Export the field
	}
	p1 := person{Name: "Steve"}

	// Marshal
	bs, _ := json.Marshal(numbers)
	fmt.Println(string(bs))

	bs, _ = json.Marshal(grades)
	fmt.Println(string(bs))

	bs, _ = json.Marshal(p1)
	fmt.Println(string(bs))

	// Unmarshal
	var p2 person
	json.Unmarshal([]byte(`{"name": "John"}`), &p2)
	fmt.Println(p2.Name)
}
