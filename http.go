package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

var secret string

func useHTTP() {

	secret = getEnvPropOrDefault("secret", "")
	port := getEnvPropOrDefault("port", "8080")
	var portInt, err = strconv.Atoi(port)
	if err != nil || portInt < 1 || portInt > 65535 {
		fmt.Printf("Cannot use %s as port - not a valid port\n", port)
		os.Exit(1)
	}

	http.HandleFunc("/map", auth(imageHandler))

	fmt.Printf("Starting server on port: %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	coord := r.URL.Query().Get("coords")
	err := generateImage(coord)
	if err != nil {
		returnError(w, err)
		return
	}

	var path = getImagePath(coord)

	img, err := os.Open(path)
	if err != nil {
		returnError(w, err)
		return
	}
	defer img.Close()

	w.Header().Set("Content-Type", "image/png")
	_, err = io.Copy(w, img)
	if err != nil {
		returnError(w, err)
	}
}

func returnError(w http.ResponseWriter, err error) {
	fmt.Println(err)
	w.WriteHeader(http.StatusInternalServerError) // perhaps handle this nicer
	w.Write([]byte("Something's wrong: " + err.Error()))
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		querySecret := r.URL.Query().Get("secret")
		if querySecret != "" && querySecret == secret {
			next(w, r)
			return
		}
		returnError(w, fmt.Errorf("unauthorized"))
		return
	}
}
