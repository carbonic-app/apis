package main

import (
	"log"
)

func main() {
	if err := RunServer(); err != nil {
		log.Fatal(err)
	}
}
