package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id" validate:"required.uuid" gorm:"primary_key"`
	Username       string    `json:"username" gorm:"required;"`
	Email          string    `json:"email" gorm:"required;"`
	Bio            string    `json:"bio"`
	UserImage      uuid.UUID `json:"-"`
	PasswordHash   string    `json:"-" gorm:"required;"`
	CreatedAt      int       `json:"-" gorm:"autoCreateTime"`
	ImageSent      []Image   `json:"sentImages" gorm:"foreignKey:SenderID;"`
	ImagesReceived []Image   `json:"receivedImages" gorm:"foreignKey:ReceiverID;"`
}
