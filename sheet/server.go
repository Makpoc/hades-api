package sheet

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/makpoc/hades-api/utils"

	gsclient "github.com/makpoc/hades-api/sheet/client"
)

var sheet *gsclient.Sheet

// Init ...
func Init() error {
	spreadsheetID, ok := os.LookupEnv("SHEET_ID")
	if !ok {
		return fmt.Errorf("no SHEET_ID found in environment")
	}
	var err error
	sheet, err = gsclient.New(spreadsheetID)
	if err != nil {
		return fmt.Errorf("failed to create sheet client: %v", err)
	}

	return nil
}

// GetHandleFuncs returns a map with paths and handlers to attach to them
func GetHandleFuncs() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/timezones":            timeZonesHandler,
		"/timezones/{username}": userTimeZoneHandler,
		"/users":                usersHandler,
		"/users/{username}":     userHandler,
	}
}

func timeZonesHandler(res http.ResponseWriter, req *http.Request) {
	result, err := sheet.GetTimeZones()
	if err != nil {
		log.Printf("Failed to get time zones: %v\n", err)
		utils.SendError(res, http.StatusBadRequest, fmt.Errorf("failed to get timezones: %v", err))
		return
	}

	body, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Printf("Failed to marshal json: %v\n", err)
		utils.SendError(res, http.StatusBadRequest, fmt.Errorf("failed to marshal json: %v", err))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(body)
}

func userTimeZoneHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	result, err := sheet.GetTimeZone(username)
	if err != nil {
		log.Printf("Failed to get time zone for user %s. Error was: %v", username, err)
		return
	}

	body, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Printf("Failed to marshal json: %v\n", err)
		utils.SendError(w, http.StatusBadRequest, fmt.Errorf("failed to marshal json: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func usersHandler(res http.ResponseWriter, req *http.Request) {
	result, err := sheet.GetUsers()
	if err != nil {
		log.Printf("Failed to get Users: %v\n", err)
		utils.SendError(res, http.StatusBadRequest, fmt.Errorf("failed to get users: %v", err))
		return
	}

	body, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal json: %v\n", err)
		utils.SendError(res, http.StatusBadRequest, fmt.Errorf("failed to marshal json: %v", err))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(body)
}

func userHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	result, err := sheet.GetUser(username)
	if err != nil {
		log.Printf("Failed to get User %s: %v\n", username, err)
		utils.SendError(res, http.StatusBadRequest, fmt.Errorf("failed to get users: %v", err))
		return
	}

	body, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal json: %v\n", err)
		utils.SendError(res, http.StatusBadRequest, fmt.Errorf("failed to marshal json: %v", err))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(body)
}
