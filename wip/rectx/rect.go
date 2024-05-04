package rectx

import (
	"image"
)

type Rectangle image.Rectangle

func (r Rectangle) SplitH(y int) (top, bottom Rectangle) {
	t, b := SplitH(r, y)
	return Rectangle(t), Rectangle(b)
}

func (r Rectangle) SplitV(x int) (left, right Rectangle) {
	l, r := SplitV(r, x)
	return Rectangle(l), Rectangle(r)
}

