package main

import (
	"log"
	"os"
)

func main() {
	url := os.Args[1]
	if url == "" {
		log.Fatal("provide a valid URL")
		return
	}
}
