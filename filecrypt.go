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
	filePath := flag.String("f", "./cryptfile.txt", "The filepath to the file you want to operate on.")
	passPhrase := flag.String("k", "secret", "A key to encrpyt or decrypt your file with...remember this!")
	encryptMode := flag.Bool("e", false, "Encrypt Mode: Encrypt the specified file with the key provided.")
	decryptMode := flag.Bool("d", false, "Decrypt Mode: Attempt to decrypt the file with the key provided.")
	addMode := flag.Bool("a", false, "Add Mode: Decrypt and append lines given as tail arguments to the file, then reencrypt with the key provided.")
	readMode := flag.Bool("r", false, "Read Mode: Will read the file given before and after other operations.")
	flag.Parse()

	fmt.Printf("Working on file at %v\n", *filePath)

	read(*filePath, "Read before:", *readMode)
	switch {
	case *encryptMode && !*decryptMode && !*addMode:
		fmt.Println("Let's encrypt...")
		err := encrypt(*filePath, *passPhrase)
		if err != nil {
			fmt.Println(err)
		}
	case *decryptMode && !*encryptMode && !*addMode:
		fmt.Println("Let's decrypt...")
		err := decrypt(*filePath, *passPhrase)
		if err != nil {
			fmt.Println(err)
		}
	case *addMode && !*encryptMode && !*decryptMode:
		fmt.Println("Add Mode...")
		err := add(*filePath, *passPhrase, flag.Args())
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("Please add a single job flag (-e, -d or -a). See help -h for more info.")
	}
	read(*filePath, "Read after:", *readMode)

}

func read(file, location string, readMode bool) error {
	if readMode {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("Could not read this file: %w", err)
		}
		fmt.Println(location)
		fmt.Printf("%s\n", data)
	}
	return nil
}

func encrypt(file, secret string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Could not read this file: %w", err)
	}
	plaintext := []byte(data)

	key32 := sha256.Sum256([]byte(secret))
	key := key32[:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)

	err = ioutil.WriteFile(file, []byte(ciphertext), 0777)
	if err != nil {
		return fmt.Errorf("Could not write to file: %w", err)
	}

	fmt.Println("File encrypted!")
	return nil
}

func decrypt(file, secret string) error {
	key32 := sha256.Sum256([]byte(secret))
	key := key32[:]

	ciphertext, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Could not read file: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("Error: %w", err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	err = ioutil.WriteFile(file, []byte(plaintext), 0777)
	if err != nil {
		return fmt.Errorf("Could not write to file: %w", err)
	}

	fmt.Println("File decrypted!")

	return nil
}

func add(file, secret string, args []string) error {
	err := decrypt(file, secret)
	if err != nil {
		return fmt.Errorf("Error in decryption phase: %w", err)
	}
	for _, v := range args {
		fmt.Printf("Adding line: %s\n", v)
	}
	err = encrypt(file, secret)
	if err != nil {
		return fmt.Errorf("Error in re-encryption phase: %w", err)
	}

	fmt.Println("Update complete!")

	return nil
}
