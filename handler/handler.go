package handler

import (
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Handler struct {
	DB *gorm.DB
}

func HandleBodyDecode(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		RespondWithError(w, http.StatusBadRequest, err)
		return false
	}
	defer r.Body.Close()

	return true
}

func RespondWithError(w http.ResponseWriter, code int, err error) {
	log.Println(err)
	RespondWithCode(w, code, map[string]string{"error": err.Error()})
}

func RespondWithSuccess(w http.ResponseWriter, v interface{}) {
	RespondWithCode(w, http.StatusOK, v)
}

func RespondWithCreated(w http.ResponseWriter, v interface{}) {
	RespondWithCode(w, http.StatusCreated, v)
}

func RespondWithNotFound(w http.ResponseWriter, v interface{}) {
	RespondWithCode(w, http.StatusNotFound, v)
}

func RespondWithCode(w http.ResponseWriter, code int, v interface{}) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return
	}
}

func RespondEmptyWithCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func RespondWithEmptyArray(w http.ResponseWriter) {
	json.NewEncoder(w).Encode([]string{})
}

func RespondWithFile(w http.ResponseWriter, file []byte) {
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(file)
}

//func RespondWithImage(w http.ResponseWriter, image image.Image) {
//	json.NewEncoder(w).Encode(image)
//}
