package main

import (
	"flag"
	"fmt"
)

func main() {
	filePath := flag.String("file", "./default.txt", "The filepath to the file you want to operate on.")
	flag.Parse()
	fmt.Println(*filePath)
}
