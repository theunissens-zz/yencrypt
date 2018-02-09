package main

import (
	"fmt"
	"github.com/yencrypt/yencrypt_client/client"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		log.Print("Please provide arguments: action id plaintextfile")
		return
	}
	var err error
	client.Port, err = strconv.Atoi(os.Args[1])
	if err != nil {
		log.Printf("Invalid port: %s", os.Args[1])
		return
	}
	action := os.Args[2]

	if strings.Compare(action, "encrypt") == 0 {
		encrypt()
	} else if strings.Compare(action, "decrypt") == 0 {
		decrypt()
	} else {
		log.Printf("Invalid action: %s", action)
	}

}

func encrypt() {
	id := os.Args[3]
	plainTextFile := os.Args[4]

	plainText, err := readPlainTextFromFile(plainTextFile)
	if err != nil {
		log.Printf("Unable to read file: %v", err)
		return
	}

	client := client.YClient{}
	key, err := client.Store([]byte(id), plainText)
	if err != nil {
		log.Printf("An unexpected error occurred: %v", err)
		return
	}
	fmt.Printf("Key: %s", string(key[0:]))
	return
}

func decrypt() {
	id := os.Args[3]
	key := os.Args[4]

	client := client.YClient{}
	plainText, err := client.Retrieve([]byte(id), []byte(key))
	if err != nil {
		log.Printf("An unexpected error occurred: %v", err)
		return
	}
	fmt.Printf("Original plaintext: %s", string(plainText[0:]))
}

func readPlainTextFromFile(file string) ([]byte, error) {
	plaintext, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
