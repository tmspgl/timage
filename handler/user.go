package handler

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	models "timage.flomas.net/model"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
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

// TODO
func (h *Handler) validateValue(key string, value string) error {
	switch key {
	case "username":
		var a []models.User
		err := h.DB.Find(&a, "username = ?", key).Error
		errors.Is(err, gorm.ErrRecordNotFound)
		if err == gorm.ErrRecordNotFound {
			return nil
		} else {
			return errors.New("username already existing")
		}
	}
	return nil
}

func (h *Handler) GetUserImage(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// TODO validateValues()
	username := r.FormValue("username")
	//err := h.validateValue("username", username)
	//if err != nil {
	//	RespondWithError(w, http.StatusConflict, err)
	//
	//	return
	//}
	password := r.FormValue("password")
	email := r.FormValue("email")
	bio := r.FormValue("bio")

	userId := uuid.New()
	imageId := uuid.New()
	user := models.User{
		ID:           userId,
		Username:     username,
		Email:        email,
		Bio:          bio,
		PasswordHash: password,
		UserImage:    imageId,
	}
	if err := h.DB.Create(&user).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)

		return
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err)

		return
	}
	err = SaveImageOnDisk("./userImage/"+imageId.String(), file)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
	}
	RespondWithCreated(w, user)
}

func (h *Handler) EditUserImage(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) EditUserInfo(w http.ResponseWriter, r *http.Request) {
}
