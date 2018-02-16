package utils

import (
	"log"
	"net/http"
)

// SendError logs the error and sends the message back as an response
func SendError(w http.ResponseWriter, status int, err error) {
	log.Println(err)
	w.WriteHeader(status) // perhaps handle this nicer
	w.Write([]byte("Something's wrong: " + err.Error()))
}
