package utils

import (
	"fmt"
	"io/ioutil"
)

func ReadHex(filename string) string {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error in reading", filename, err)
	}
	return string(f)
}

func WriteHex(filename string, msg string) {
	err := ioutil.WriteFile(filename, []byte(msg), 0666)
	if err != nil {
		fmt.Println("Error in writing", filename, err)
	}
}
