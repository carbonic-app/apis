package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

var (
	port = "80"
)

func main() {
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	log.Printf("Starting gRPC Server at :%s", port)
	if err := RunServer(lis); err != nil {
		log.Fatal(err)
	}
}
