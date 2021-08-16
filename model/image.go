package models

import (
	"github.com/google/uuid"
	"image"
	"image/color"
	"time"
)

type Image struct {
    ID      uuid.UUID   `json:"imageId"`
    Path    string      `json:"path"`
    Time    time.Time      `json:"time"`
}

func (i Image) ColorModel() color.Model {
	panic("implement me")
}

func (i Image) Bounds() image.Rectangle {
	panic("implement me")
}

func (i Image) At(x, y int) color.Color {
	panic("implement me")
}


