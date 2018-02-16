package hmap

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/makpoc/hades-api/utils"
)

// GetHandleFuncs returns a map with paths and handlers to attach to them
func GetHandleFuncs() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/map": imageHandler,
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
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	var path = getImagePath(coordsArray)

	img, err := os.Open(path)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer img.Close()

	w.Header().Set("Content-Type", "image/jpeg")
	_, err = io.Copy(w, img)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err)
	}
}
