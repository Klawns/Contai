package main

import (
	"log"

	"contai/internal/server"
)

func main() {
	server, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
