package models

import (
	"github.com/google/uuid"
	"time"
)

type Timer struct {
	ID   uuid.UUID `json:"id" validate:"required,uuid"`
	Time time.Time `json:"time"`
}
