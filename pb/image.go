package pb

import (
	"image"
	"image/color"
)

func ColorEncode(c color.Color) *Color {
	if c == nil {
		return nil
	}
	r, g, b, a := c.RGBA()
	return &Color{
		Red:   r,
		Green: g,
		Blue:  b,
		Alpha: a,
	}
}

func ImageEncode(img image.Image) *Image {
	if img == nil {
		return nil
	}
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			newImg.Set(x, y, img.At(x, y))
		}
	}

	return &Image{
		Width:  int64(bounds.Max.X - bounds.Min.X),
		Height: int64(bounds.Max.Y - bounds.Min.Y),
		Stride: int64(newImg.Stride),
		Data:   newImg.Pix,
	}
}

func ImageDecode(img *Image) *image.RGBA {
	return &image.RGBA{
		Pix:    img.Data,
		Stride: int(img.Stride),
		Rect:   image.Rect(0, 0, int(img.Width), int(img.Height)),
	}
}
