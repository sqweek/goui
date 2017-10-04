package wdedrv

import (
	"image/draw"
	"github.com/skelterjohn/go.wde"
)

func Make(win wde.Window) WdeDrv {
	return WdeDrv{win}
}

// WdeDrv implements goui.PaintDrv interface
type WdeDrv struct {
	win wde.Window
}

func (d WdeDrv) Img() draw.Image {
	return d.win.Screen()
}

func (d WdeDrv) Flip() {
	d.win.FlushImage()
}
