package main

import (
	"flag"
	"fmt"
)

func main() {
	filePath := flag.String("file", "./default.txt", "The filepath to the file you want to operate on.")
	passKey := flag.String("k", "0000", "A key to encrpyt or decrypt your file with...remember this!")
	encryptMode := flag.Bool("e", false, "Encrypt Mode: Encrypt the specified file with the key provided.")
	decryptMode := flag.Bool("d", false, "Decrypt Mode: Attempt to decrypt the file with the key provided.")
	addMode := flag.Bool("a", false, "Add Mode: Decrypt and append lines given as tail arguments to the file, then reencrypt with the key provided.")
	flag.Parse()

	fmt.Printf("Encrypt file at %v with key %v.\n", *filePath, *passKey)
	if len(flag.Args()) > 0 {
		fmt.Printf("We'll be adding lines: %s\n", flag.Args())
	}

	switch {
	case *encryptMode && !*decryptMode && !*addMode:
		fmt.Println("Let's encrpyt!")
	case *decryptMode && !*encryptMode && !*addMode:
		fmt.Println("Let's decrypt!")
	case *addMode && !*encryptMode && !*decryptMode:
		fmt.Println("Add Lines")
	default:
		fmt.Println("Please add a single job flag (-e, -d or -a)")
	}

}
