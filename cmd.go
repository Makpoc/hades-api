package main

import (
	"fmt"
	"os"
)

func useCmd() {
	var coord string
	if len(os.Args) > 2 {
		coord = os.Args[2]
	}
	err := generateImage(coord)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Operation successful. Check output/ folder for your image")
}
