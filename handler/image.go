package handler

import (
    "errors"
    "fmt"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "timage.flomas.net/model"
    "time"

    "gorm.io/gorm"
)

func (h *Handler) RetrieveImages(w http.ResponseWriter, r *http.Request) {
    Images := &[]models.Image{}

    if err := h.DB.Find(&Images).Error; errors.Is(err, gorm.ErrRecordNotFound) {
        RespondWithEmptyArray(w)

        return
    } else if err != nil {
    	RespondWithError(w, http.StatusInternalServerError, err)

    	return
    }

    RespondWithSuccess(w, Images)
}

func (h *Handler) RetrieveImage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    image := models.Image{}
    if err := h.DB.First(&image, "id = ?", vars["imageId"]).Error; errors.Is(err, gorm.ErrRecordNotFound) {
        RespondEmptyWithCode(w, http.StatusNotFound)

        return
    } else if err != nil {
        RespondWithError(w, http.StatusInternalServerError, err)

        return
    }
    if image.Time.After(time.Now()) {
        //TODO proper error handling
        err := errors.New("IMAGE NOT YET SENT. CANNOT BE RETRIEVED")
        RespondWithError(w, http.StatusBadRequest, err)

        return
    }
    fileBytes, err := ioutil.ReadFile(image.Path)
    if err != nil {
        panic(err)
    }

    RespondWithFile(w, fileBytes)
}

func (h *Handler) CreateImage(w http.ResponseWriter, request *http.Request) {
   err := request.ParseMultipartForm(32 << 20) // maxMemory 32MB
   if err != nil {
       RespondWithError(w, http.StatusBadRequest, err)

      return
   }
   file, header, err := request.FormFile("photo")
   fmt.Println(header.Header.Get("Content-Type"))
   if err != nil {
       RespondWithError(w, http.StatusBadRequest, err)

       return
   }
   date := request.FormValue("date")
   time, err := time.Parse(time.RFC3339, date)
   if err != nil {
       RespondWithError(w, http.StatusBadRequest, err)
   }
   imageId := uuid.New()
   path := "./photo/" + imageId.String()

   image := models.Image{ID: imageId, Path: path, Time: time}
   h.DB.Create(&image)

   tmpfile, err := os.Create("./photo/" + imageId.String())
   defer tmpfile.Close()
   if err != nil {
       RespondWithError(w, http.StatusInternalServerError, err)

       return
   }
   _, err = io.Copy(tmpfile, file)
   if err != nil {
       RespondWithError(w, http.StatusInternalServerError, err)

       return
   }

    ResponseWithCreated(w, image)
}