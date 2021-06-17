package main

import (
	"io/ioutil"
	"log"
)

func main() {
	bytes, err := ioutil.ReadFile("data")
	if err != nil {
		log.Fatal(err)
	}
	print(bytes[0])
}
