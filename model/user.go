package models

import (
	"github.com/google/uuid"
)

type User struct {
	UserID          uuid.UUID	    `json:"id" validate:"required.uuid" gorm:"primary_key"`
	Username     	string  		`json:"username" gorm:"required;"`
	Email        	string  		`json:"email" gorm:"required;"`
	Bio          	string  		`json:"bio;size:1024"`
	Image        	*string 		`json:"image,omitempty"`
	PasswordHash	string  		`gorm:"required;"`
	CreatedAt		int				`json:"created_at" gorm:"autoCreateTime"`
	ImagesSent		[]Image			`json:"images" gorm:"foreignKey:SenderID"`
}