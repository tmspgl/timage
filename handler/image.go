package handler

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"timage.flomas.net/model"
	"time"

	"gorm.io/gorm"
)

var latest = models.Image{}

func (h *Handler) RetrieveAllImages(w http.ResponseWriter, _ *http.Request) {
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
	path := "./photo/" + image.ID.String()
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	RespondWithFile(w, fileBytes)
}

func (h *Handler) GetLatest(w http.ResponseWriter, r *http.Request) {
	path := "./photo/" + latest.ID.String()
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		RespondWithError(w, http.StatusUnprocessableEntity, err)

		return
	}
	RespondWithFile(w, fileBytes)
}

func (h *Handler) StartImageFetch(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case _ = <-ticker.C:
				h.SendLatest()
			}
		}
	}()
}

func (h *Handler) sendImage() {

}

func (h *Handler) SendLatest() (w http.ResponseWriter, r *http.Request) {
	c := time.Now()
	var timers []models.Timer
	if err := h.DB.Find(&timers, "time < ?", c).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if len(timers) == 0 {
		return nil, nil
	}
	for _, timer := range timers {
		image := models.Image{}
		if err := h.DB.Find(&image, "id = ?", timer.ID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println(err)

			return //TODO
		}
		// TODO sendImage()
		fmt.Println(image)
		if err := h.DB.Delete(&timer, "id = ?", timer.ID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Fatal(err)
			//retry. make sure image doesn't get send multiple times
		}
		_, err := ioutil.ReadFile("./photo/" + image.ID.String())
		if err != nil {
			panic(err)
		}
		fmt.Println("Image was sent")
		latest = image //TODO delete. Just to show updated image on base_url.
	}

	//RespondWithFile(w, fileBytes)
	x := time.Now()
	diff := c.Sub(x)
	fmt.Print("difference is ")
	fmt.Println(diff)
	return nil, nil
}

func (h *Handler) CreateImage(w http.ResponseWriter, r *http.Request) {

	date := r.FormValue("date")
	imageTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err)

		return
	}
	receiver := models.User{}
	if err := h.DB.First(&receiver, "id = ?", r.FormValue("receiver")).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		RespondWithError(w, http.StatusBadRequest, err)

		return
	}
	sender := models.User{}
	if err := h.DB.First(&sender, "id = ?", r.FormValue("sender")).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		RespondWithError(w, http.StatusBadRequest, err)

		return
	}
	imageId := uuid.New()
	image := models.Image{ID: imageId, Time: imageTime, SenderID: sender.ID, ReceiverID: receiver.ID}
	if err := h.DB.Create(&image).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)

		return
	}
	timer := models.Timer{ID: imageId, Time: imageTime}
	if err := h.DB.Create(&timer).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)

		return
	}
	file, _, err := r.FormFile("photo")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err)

		return
	}
	path := "./photo/" + imageId.String()
	err = SaveImageOnDisk(path, file)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
	}

	RespondWithCreated(w, image)
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	h.DB.Preload("ImageSent").Preload("ImagesReceived").First(&user, "id = ?", "d06281ab-49dc-426a-8747-f98c4725c9e3")
	RespondWithSuccess(w, user)
}

func (h *Handler) Test2(w http.ResponseWriter, r *http.Request) {
	image := &[]models.Image{}
	h.DB.Preload("Timer").Find(&image, "receiver_id = ?", "d06281ab-49dc-426a-8747-f98c4725c9e3")
	RespondWithSuccess(w, image)
}

func checkSendDate(t time.Time) error {
	if t.Before(time.Now()) {
		return errors.New("send date is not reached yet")
	}
	return nil
}

func SaveImageOnDisk(path string, file multipart.File) error {
	tmpfile, err := os.Create(path)
	defer tmpfile.Close()
	if err != nil {

		return err
	}
	_, err = io.Copy(tmpfile, file)
	if err != nil {

		return err
	}
	return nil
}
