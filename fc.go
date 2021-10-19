package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
)

func main() {
	filePath := flag.String("f", "./default", "The filepath to the file you want to operate on.")
	passKey := flag.String("k", "passphrasewhichneedstobe32bytes!", "A key to encrpyt or decrypt your file with...remember this!")
	encryptMode := flag.Bool("e", false, "Encrypt Mode: Encrypt the specified file with the key provided.")
	decryptMode := flag.Bool("d", false, "Decrypt Mode: Attempt to decrypt the file with the key provided.")
	addMode := flag.Bool("a", false, "Add Mode: Decrypt and append lines given as tail arguments to the file, then reencrypt with the key provided.")
	flag.Parse()

	fmt.Printf("Working on file at %v with key %v.\n", *filePath, *passKey)
	if len(flag.Args()) > 0 {
		fmt.Printf("We')ll be adding lines: %s\n", flag.Args())
	}

	switch {
	case *encryptMode && !*decryptMode && !*addMode:
		fmt.Println("Let's encrpyt!")
		encrypt(*filePath, *passKey)
	case *decryptMode && !*encryptMode && !*addMode:
		fmt.Println("Let's decrypt!")
	case *addMode && !*encryptMode && !*decryptMode:
		fmt.Println("Add Lines")
	default:
		fmt.Println("Please add a single job flag (-e, -d or -a)")
	}

}

func encrypt(file, key string) {
	plainText := []byte(file)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	cipherText := aesgcm.Seal(nonce, nonce, plainText, nil)
	fmt.Printf("%x\n", cipherText)
}
