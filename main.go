package main

import (
	"os"
)

func main() {

	var ui string
	if len(os.Args) > 1 {
		ui = os.Args[1]
	}

	if ui == "cmd" {
		useCmd()
	} else if ui == "http" {
		useHTTP()
	}
}

func getEnvPropOrDefault(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
