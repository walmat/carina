package util

import "image/color"

var (
	ColorSuccess = color.RGBA{
		R: 135,
		G: 179,
		B: 141,
		A: 255,
	}

	ColorWarning = color.RGBA{
		R: 255,
		G: 200,
		B: 87,
		A: 255,
	}

	ColorFailed = color.RGBA{
		R: 242,
		G: 110,
		B: 134,
		A: 255,
	}
)

func NewRGB(r, g, b uint8) color.RGBA {
	return color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	}
}
