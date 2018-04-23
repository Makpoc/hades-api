package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/makpoc/hades-api/sheet/gsclient"
	"github.com/makpoc/hades-api/utils"
)

var sheet *gsclient.Sheet

func InitSheet() error {
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

func sendResponse(content interface{}, res http.ResponseWriter, req *http.Request) {
	body, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal json: %v\n", err)
		utils.SendError(res, http.StatusBadRequest, fmt.Errorf("failed to marshal json: %v", err))
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(body)
}
