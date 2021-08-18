package handler

import (
    "errors"
    "fmt"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "timage.flomas.net/model"
    "time"

    "gorm.io/gorm"
)

var latest = models.Image{}

func (h *Handler) RetrieveImages(w http.ResponseWriter, _ *http.Request) {
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
    if err := h.DB.First(&image, "image_id = ?", vars["imageId"]).Error; errors.Is(err, gorm.ErrRecordNotFound) {
        RespondEmptyWithCode(w, http.StatusNotFound)

        return
    } else if err != nil {
        RespondWithError(w, http.StatusInternalServerError, err)

        return
    }
    err := checkSendDate(image.Time)
    if err != nil {
        RespondWithError(w, http.StatusBadRequest, err)

        return
    }
    fileBytes, err := ioutil.ReadFile(image.Path)
    if err != nil {
        panic(err)
    }
    RespondWithFile(w, fileBytes)
}

func checkSendDate(t time.Time) error {
    if t.Before(time.Now()) {
        return errors.New("send date is not reached yet")
    }
    return nil
}
func (h *Handler) GetLatest(w http.ResponseWriter, r *http.Request) {
    fileBytes, err := ioutil.ReadFile(latest.Path)
    if err != nil {
        RespondWithError(w, http.StatusUnprocessableEntity, err)

        return
    }
    RespondWithFile(w, fileBytes)
}

func (h *Handler) StartImageFetch() () {
    ticker := time.NewTicker(10* time.Second)
    go func() {
        for {
            select {
            case t := <-ticker.C:
                h.SendLatest()
                fmt.Println("FETCH NOW", t)
            }
        }
    }()
}

func (h *Handler) SendLatest() (w http.ResponseWriter, r *http.Request) {
    c := time.Now()
    var a []models.Image
    if err := h.DB.Find(&a, "time < ?", c).Error; errors.Is(err, gorm.ErrRecordNotFound) {
        log.Fatal(err)
    }
    if len(a) == 0 {
        x := time.Now()
        diff := c.Sub(x)
        fmt.Print("difference is ")
        fmt.Println(diff)
        return nil, nil
    }
    for _, i := range a {
        // send image here
        image := models.ImageStore{
            ImageID:    i.ImageID,
            Path:       i.Path,
            Time:       i.Time,
            SenderID:   i.SenderID,
            ReceiverID: i.ReceiverID,
        }
        fmt.Println(image.Path)
        if err := h.DB.Create(&image).Error; err != nil {
            log.Fatal(err)

            return
        } else { //TODO only delete when other created
            if err := h.DB.Delete(&i).Error; errors.Is(err, gorm.ErrRecordNotFound) {
                log.Fatal(err)
            }
        }
        _, err := ioutil.ReadFile(i.Path)
        if err != nil {
            panic(err)
        }
        latest = i
        //RespondWithFile(w, fileBytes)
    }
    x := time.Now()
    diff := c.Sub(x)
    fmt.Print("difference is ")
    fmt.Println(diff)
    return nil, nil

}

func (h *Handler) CreateImage(w http.ResponseWriter, r *http.Request) {
    file, _, err := r.FormFile("photo")
    if err != nil {
       RespondWithError(w, http.StatusBadRequest, err)

       return
    }
    date := r.FormValue("date")
    imageTime, err := time.Parse(time.RFC3339, date)
    if err != nil {
        RespondWithError(w, http.StatusBadRequest, err)

        return
    }
    receiver := models.User{}
    if err := h.DB.First(&receiver, "user_id = ?", r.FormValue("receiver")).Error; errors.Is(err, gorm.ErrRecordNotFound) {
        RespondWithError(w, http.StatusBadRequest, err)

        return
    }
    sender:= models.User{}
    if err:= h.DB.First(&sender, "user_id = ?", r.FormValue("sender")).Error; errors.Is(err, gorm.ErrRecordNotFound) {
        RespondWithError(w, http.StatusBadRequest, err)

        return
    }
   imageId := uuid.New()
   path := "./photo/" + imageId.String()
   image := models.Image{ImageID: imageId, Path: path, Time: imageTime, SenderID: sender.UserID, ReceiverID: receiver.UserID}
   if err := h.DB.Create(&image).Error; err != nil {
       RespondWithError(w, http.StatusInternalServerError, err)

       return
   }
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
    RespondWithCreated(w, image)
}

