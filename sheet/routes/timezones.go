package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/makpoc/hades-api/utils"
)

func TimeZonesHandler(res http.ResponseWriter, req *http.Request) {
	result, err := sheet.GetTimeZones()
	if err != nil {
		log.Printf("Failed to get time zones: %v\n", err)
		utils.SendError(res, http.StatusBadRequest, fmt.Errorf("failed to get timezones: %v", err))
		return
	}

	sendResponse(result, res, req)
}

func TimeZoneHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	result, err := sheet.GetTimeZone(username)
	if err != nil {
		log.Printf("Failed to get time zone for user %s. Error was: %v", username, err)
		return
	}

	sendResponse(result, res, req)
}
