package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
)

func main() {
	filePath := flag.String("f", "./default", "The filepath to the file you want to operate on.")
	passPhrase := flag.String("k", "secret", "A key to encrpyt or decrypt your file with...remember this!")
	encryptMode := flag.Bool("e", false, "Encrypt Mode: Encrypt the specified file with the key provided.")
	decryptMode := flag.Bool("d", false, "Decrypt Mode: Attempt to decrypt the file with the key provided.")
	addMode := flag.Bool("a", false, "Add Mode: Decrypt and append lines given as tail arguments to the file, then reencrypt with the key provided.")
	flag.Parse()

	fmt.Printf("Working on file at %v\n", *filePath)
	if len(flag.Args()) > 0 {
		fmt.Printf("We')ll be adding lines: %s\n", flag.Args())
	}

	switch {
	case *encryptMode && !*decryptMode && !*addMode:
		fmt.Println("Let's encrpyt!")
		encrypt(*filePath, *passPhrase)
	case *decryptMode && !*encryptMode && !*addMode:
		fmt.Println("Let's decrypt!")
	case *addMode && !*encryptMode && !*decryptMode:
		fmt.Println("Add Lines")
	default:
		fmt.Println("Please add a single job flag (-e, -d or -a). See help -h for more info.")
	}

}

func encrypt(file, secret string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}
	plaintext := []byte(data)

	key32 := sha256.Sum256([]byte(secret))
	key := key32[:]

	block, err := aes.NewCipher(key)
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

	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)

	err = ioutil.WriteFile(file, []byte(ciphertext), 0777)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("File encrypted!")
}
