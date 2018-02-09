package main

import (
	"github.com/yencrypt/yencrypt_server/encryptserver"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) <= 1 {
		log.Print("No port given in args")
		return
	}
	strPort := os.Args[1]
	port, err := strconv.Atoi(strPort)
	if err != nil {
		log.Fatal("Invalid port given, shutting down...")
	}

	server := encryptserver.YServer{}
	err = server.StartServer(port)
	if err != nil {
		log.Fatalf("Could not initialize service: %v", err)
	}
}
