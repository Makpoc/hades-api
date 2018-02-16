package main

import (
	"log"

	"github.com/makpoc/hades-api/server"
)

func main() {
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
