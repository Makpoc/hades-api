package hmap

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/makpoc/hades-api/utils"
)

// GetHandleFuncs returns a map with paths and handlers to attach to them
func GetHandleFuncs() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/map": mapHandler,
	}
}

func mapHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getMapHandler(w, r)
	case http.MethodPost:
		uploadMapHandler(w, r)
	}
}

func getMapHandler(w http.ResponseWriter, r *http.Request) {
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

func uploadMapHandler(w http.ResponseWriter, r *http.Request) {

	var picTypeExt = map[string]string{
		"labels":     "png",
		"screenshot": "jpeg",
	}

	var created bool
	for picType, picExt := range picTypeExt {
		file, handle, err := r.FormFile(picType)
		if err != nil {
			// not part of the request
			continue
		}
		defer file.Close()

		fmt.Println(handle.Header.Get("Content-Type"))
		//		mimeType := handle.Header.Get("Content-Type")
		//		if mimeType != "image/jpeg" && mimeType != "image/png" {
		//			jsonResponse(w, http.StatusBadRequest, fmt.Sprintf("The format file is not valid - %s", mimeType))
		//			return
		//		}
		err = saveFile(w, file, fmt.Sprintf("%s.%s", picType, picExt))
		if err != nil {
			jsonResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// at least one file was created
		created = true
	}

	if !created {
		jsonResponse(w, http.StatusBadRequest, fmt.Sprintf("Request did not contains any known pictures!"))
	}

	err := cleanOutputFolder()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete output folder: %v", err))
		return
	}
	jsonResponse(w, http.StatusCreated, "File(s) uploaded successfully!.")
}

func saveFile(w http.ResponseWriter, file multipart.File, fileName string) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file %s.", fileName)
	}

	if len(data) == 0 {
		return fmt.Errorf("file is empty")
	}

	fmt.Println(len(data))

	err = ioutil.WriteFile("res/"+fileName, data, 0666)
	if err != nil {
		return fmt.Errorf("failed to file file %s.", fileName)
	}

	return nil
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	fmt.Println(code, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, fmt.Sprintf("{\"code\": %d, \"message\": \"%s\"}", code, message))
}

func cleanOutputFolder() error {
	if err := os.RemoveAll("output/"); err != nil {
		return err
	}
	if err := os.Mkdir("output/", 0777); err != nil {
		return err
	}
	return nil
}
