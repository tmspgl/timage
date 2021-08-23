package models

import (
	"github.com/google/uuid"
	"time"
)

type Image struct {
	ID         uuid.UUID `json:"id" validate:"required,uuid"`
	SenderID   uuid.UUID `json:"senderID"`
	ReceiverID uuid.UUID `json:"receiverID"`
	Time       time.Time `json:"time"`
}
