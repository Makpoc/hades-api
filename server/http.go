package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/makpoc/hades-api/hmap"
	"github.com/makpoc/hades-api/sheet"

	"github.com/gorilla/mux"

	"github.com/makpoc/hades-api/utils"
)

var secret string

// Start starts the server
func Start() error {
	r := mux.NewRouter().StrictSlash(true)
	s := r.PathPrefix("/api/v1").Subrouter()

	port := utils.GetEnvPropOrDefault("API_PORT", "8080")
	var portInt, err = strconv.Atoi(port)
	if err != nil || portInt < 1 || portInt > 65535 {
		fmt.Printf("Cannot use %s as port - not a valid port\n", port)
		os.Exit(1)
	}

	registerRoutes(s, hmap.GetHandleFuncs())
	registerRoutes(s, sheet.GetHandleFuncs())
	if err = sheet.Init(); err != nil {
		log.Printf("Failed to initialize Google Sheet client: %v", err)
		return err
	}

	fmt.Printf("Starting server on port: %s\n", port)
	return http.ListenAndServe(":"+port, s)
}

func registerRoutes(r *mux.Router, routes map[string]http.HandlerFunc) {
	for p, h := range routes {
		log.Printf("Adding %s route", p)
		r.HandleFunc(p, auth(timeLogger(h)))
	}
}
