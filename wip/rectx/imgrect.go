package rectx

import (
	"image"
)

func SplitH(r image.Rectangle, y int) (top, bottom image.Rectangle) {
	if y < r.Min.Y {
		bottom = r
	} else if y >= r.Max.Y {
		top = r
	} else {
		top = image.Rectangle{r.Min, image.Point{r.Max.X, y}}
		bottom = image.Rectangle{image.Point{r.Min.X, y}, r.Max}
	}
	return
}

func SplitV(r image.Rectangle, x int) (left, right image.Rectangle) {
	if x < r.Min.X {
		right = r
	} else if x >= r.Max.X {
		left = r
	} else {
		left = image.Rectangle{r.Min, image.Point{x, r.Max.Y}}
		right = image.Rectangle{image.Point{x, r.Min.Y}, r.Max}
	}
	return
}
