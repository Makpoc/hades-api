package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/makpoc/hades-api/utils"
)

func UsersHandler(res http.ResponseWriter, req *http.Request) {
	result, err := sheet.GetUsers()
	if err != nil {
		log.Printf("Failed to get Users: %v\n", err)
		utils.SendError(res, http.StatusBadRequest, fmt.Errorf("failed to get users: %v", err))
		return
	}

	sendResponse(result, res, req)
}

func UserHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	result, err := sheet.GetUser(username)
	if err != nil {
		log.Printf("Failed to get User %s: %v\n", username, err)
		utils.SendError(res, http.StatusBadRequest, fmt.Errorf("failed to get users: %v", err))
		return
	}

	sendResponse(result, res, req)
}
