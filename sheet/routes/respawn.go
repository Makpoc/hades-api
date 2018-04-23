package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/makpoc/hades-api/utils"
)

func RespawnsHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		getRespawns(res, req)
	} else if req.Method == http.MethodPost {
		postRespawn(res, req)
	}
}

func getRespawns(res http.ResponseWriter, req *http.Request) {
	result, err := sheet.GetRespawnTimes()
	if err != nil {
		log.Printf("Failed to get Users: %v\n", err)
		utils.SendError(res, http.StatusBadRequest, fmt.Errorf("failed to get respawn times: %v", err))
		return
	}

	sendResponse(result, res, req)
}

func postRespawn(res http.ResponseWriter, req *http.Request) {
	user := req.PostFormValue("user")
	if user == "" {
		user = "unknown"
	}
	affiliation := req.PostFormValue("affiliation")
	if affiliation == "" {
		affiliation = "unknown"
	}
	respawnsIn, err := strconv.Atoi(req.PostFormValue("respawnsInH"))
	if err != nil {
		utils.SendError(res, http.StatusBadRequest, fmt.Errorf("Failed to parse respawnsIn parameter. It needs to be an integer representing time duration"))
	}

	err = sheet.AddRespawnTime(user, affiliation, time.Now().Add(time.Duration(respawnsIn)*time.Hour))
	if err != nil {
		utils.SendError(res, http.StatusBadRequest, fmt.Errorf("Failed to add respawn info to sheet: %s", err))
	}
}

func RespawnHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	result, err := sheet.GetRespawnTime(username)
	if err != nil {
		log.Printf("Failed to get respawn time for %s: %v\n", username, err)
		utils.SendError(res, http.StatusBadRequest, fmt.Errorf("failed to get respawn time: %v", err))
		return
	}

	sendResponse(result, res, req)
}
