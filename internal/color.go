package internal

import (
	"image"
	"image/color"

	_ "github.com/gomatcha/matcha/internal/device"
)

func TintColor(img image.Image, color color.Color) image.Image {
	return img
}
