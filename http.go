package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	coords := strings.ToLower(r.URL.Query().Get("coords"))

	var coordsArray []string
	if coords != "" {
		coordsArray = strings.Split(coords, ",")
	}
	err := generateImage(coordsArray)
	if err != nil {
		returnError(w, http.StatusBadRequest, err)
		return
	}

	var path = getImagePath(coordsArray)

	img, err := os.Open(path)
	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
		return
	}
	defer img.Close()

	w.Header().Set("Content-Type", "image/jpeg")
	_, err = io.Copy(w, img)
	if err != nil {
		returnError(w, http.StatusInternalServerError, err)
	}
}

func returnError(w http.ResponseWriter, status int, err error) {
	fmt.Println(err)
	w.WriteHeader(status) // perhaps handle this nicer
	w.Write([]byte("Something's wrong: " + err.Error()))
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		querySecret := r.URL.Query().Get("secret")
		if querySecret != "" && querySecret == secret {
			next(w, r)
			return
		}
		returnError(w, http.StatusForbidden, fmt.Errorf("unauthorized"))
		return
	}
}
