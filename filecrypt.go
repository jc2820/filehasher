package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	filePath := flag.String("f", "", "The filepath to the file you want to operate on.")
	passPhrase := flag.String("k", "secret", "A key to encrpyt or decrypt your file with...remember this!")
	encryptMode := flag.Bool("e", false, "Encrypt Mode: Encrypt the specified file with the key provided.")
	decryptMode := flag.Bool("d", false, "Decrypt Mode: Attempt to decrypt the file with the key provided.")
	addMode := flag.Bool("a", false, "Add Mode: Decrypt and append lines given as tail arguments to the file, then reencrypt with the key provided.")
	readMode := flag.Bool("r", false, "Read Mode: Will read the file given after other operations.")
	flag.Parse()

	if *filePath != "" {
		fmt.Printf("Working on file at %v\n", *filePath)
	} else {
		fmt.Println("Please specify a file to work with.")
		return
	}

	switch {
	case *encryptMode && !*decryptMode && !*addMode:
		fmt.Println("Let's encrypt...")
		if err := encrypt(*filePath, *passPhrase); err != nil {
			fmt.Println(err)
		}
		if err := read(*filePath, *readMode); err != nil {
			fmt.Println(err)
		}
	case *decryptMode && !*encryptMode && !*addMode:
		fmt.Println("Let's decrypt...")
		if err := decrypt(*filePath, *passPhrase); err != nil {
			fmt.Println(err)
		}
		if err := read(*filePath, *readMode); err != nil {
			fmt.Println(err)
		}
	case *addMode && !*encryptMode && !*decryptMode:
		fmt.Println("Add Mode...")
		if err := add(*filePath, *passPhrase, flag.Args()); err != nil {
			fmt.Println(err)
		}
		if err := read(*filePath, *readMode); err != nil {
			fmt.Println(err)
		}
	default:
		if *readMode {
			if err := read(*filePath, *readMode); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Please add a single job flag (-e, -d or -a). -h for more info.")
		}
	}
}

func read(file string, readMode bool) error {
	if readMode {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("Could not read this file: %w", err)
		}
		fmt.Printf("File contents...\n---\n%s\n", data)
	}
	return nil
}

func encrypt(file, secret string) error {
	data, err := os.ReadFile(file)
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

	err = os.WriteFile(file, []byte(ciphertext), 0644)
	if err != nil {
		return fmt.Errorf("Could not write to file: %w", err)
	}

	fmt.Println("File encrypted!")
	return nil
}

func decrypt(file, secret string) error {
	key32 := sha256.Sum256([]byte(secret))
	key := key32[:]

	ciphertext, err := os.ReadFile(file)
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

	err = os.WriteFile(file, []byte(plaintext), 0644)
	if err != nil {
		return fmt.Errorf("Could not write to file: %w", err)
	}

	fmt.Println("File decrypted!")

	return nil
}

func add(file, secret string, args []string) error {
	if err := decrypt(file, secret); err != nil {
		return fmt.Errorf("Error in decryption phase: %w", err)
	}

	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Could not open your file to append to: %w", err)
	}

	for _, v := range args {
		fmt.Printf("Adding line: %s\n", v)
		if _, err := f.Write([]byte(v)); err != nil {
			f.Close()
			return fmt.Errorf("Failed writing line %s: %w", v, err)
		}
		if _, err := f.Write([]byte("\n")); err != nil {
			f.Close()
			return fmt.Errorf("Failed to write new line: %w", err)
		}
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("Failed closing file %w", err)
	}

	if err := encrypt(file, secret); err != nil {
		return fmt.Errorf("Error in re-encryption phase: %w", err)
	}

	fmt.Println("Update complete!")

	return nil
}
