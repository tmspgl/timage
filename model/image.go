package models

import (
	"github.com/google/uuid"
	"time"
)

type Image struct {
    ImageID     uuid.UUID   `json:"id" validate:"required,uuid" gorm:"type:uuid;primary_key;"`
    Path    	string      `json:"path"`
    Time    	time.Time   `json:"time" gorm:"required"`
    SenderID	uuid.UUID	//`json:"sender" gorm:"required;association_foreignKey:UserID;"`
    ReceiverID	uuid.UUID	//`json:"receiver" gorm:"required;association_foreignKey:UserID;"`
}



