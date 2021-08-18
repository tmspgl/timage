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

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	bio := r.FormValue("bio")
	image := r.FormValue("image")

	userId := uuid.New()
	user := models.User{
		UserID: userId,
		Username: username,
		Email: email,
		Bio: bio,
		PasswordHash: password,
		Image: &image,
	}
	if err := h.DB.Create(&user).Error; err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)

		return
	}

	RespondWithCreated(w, user)
}