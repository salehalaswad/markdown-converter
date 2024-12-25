package main

import (
	"fmt"
	"os"
)

func main() {
	filepath := "./test.md"
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("error opening the file", err)
		return
	}
	defer file.Close()
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println(string(data))
}
