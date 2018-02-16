package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/makpoc/hades-api/utils"
)

// auth provides authorization layer based on secret in query parameter
func auth(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var secret string
		var ok bool
		if secret, ok = os.LookupEnv("secret"); !ok {
			// server was not configured with secret
			h.ServeHTTP(w, r)
			return
		}

		querySecret := r.URL.Query().Get("secret")
		if querySecret != "" && querySecret == secret {
			h.ServeHTTP(w, r)
			return
		}
		utils.SendError(w, http.StatusForbidden, fmt.Errorf("unauthorized"))
		return
	})
}

// timeLogger wraps a handler and logs the time for its execution
func timeLogger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		defer log.Printf("It took %s to respond to TimeZone request", time.Since(startTime))
		h.ServeHTTP(res, req)
	})
}
