package utils

import (
	"fmt"
	"io/ioutil"
	"os"
)

type PaddingFunc func(*[]byte, int)

type UnpaddingFunc func(*[]byte)

func ReadStringHex(filename string) string {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error in reading", filename, err)
	}
	return string(f)
}

func WriteStringHex(filename string, msg string) {
	err := ioutil.WriteFile(filename, []byte(msg), 0666)
	if err != nil {
		fmt.Println("Error in writing", filename, err)
	}
}

func ReadBytesHex(filename string) []byte {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("File reading error: ", err)
		return nil
	}
	defer file.Close()

	var tmp []byte
	_, _ = fmt.Fscanf(file, "%X", &tmp)
	return tmp
}

func WriteBytesHex(filename string, msg []byte) {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Open file err =", err)
		return
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%02X", msg)
	if err != nil {
		return
	}
	return
}

func DumpWords(note string, in []uint32) {
	fmt.Printf("%s", note)
	for i, v := range in {
		if i%4 == 0 {
			fmt.Printf("\nword[%02d]: %.8X ", i/4, v)
		} else {
			fmt.Printf("%.8X ", v)
		}
	}
	fmt.Println("\n")
}

func DumpBytes(note string, in []byte) {
	fmt.Printf("%s", note)
	for i, v := range in {
		if i%16 == 0 {
			fmt.Printf("\nblock[%d]: %02X", i/16, v)
		} else {
			if i%4 == 0 {
				fmt.Printf(" %02X", v)
			} else {
				fmt.Printf("%02X", v)
			}
		}
	}
	fmt.Println("\n")
}

//
func PaddingZeros(in *[]byte, blockLen int) {
	for len(*in)%blockLen != 0 {
		*in = append(*in, 0x00)
	}
}

func UnpaddingZeros(in *[]byte) {
	for (*in)[len(*in)-1] == 0x00 {
		*in = (*in)[:len(*in)-1]
	}
}
