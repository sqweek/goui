package wdedrv

import (
	"image"
	"image/draw"
	"github.com/skelterjohn/go.wde"
)

func Make(win wde.Window) *WdeDrv {
	return &WdeDrv{win: win}
}

// WdeDrv implements goui.PaintDrv interface
type WdeDrv struct {
	win wde.Window
	buf *image.RGBA
	Direct bool // if true, double-buffering is not used
}

func (d *WdeDrv) Img() draw.Image {
	s := d.win.Screen()
	if d.Direct {
		return s
	}
	if d.buf == nil || d.buf.Bounds() != s.Bounds() {
		d.buf = image.NewRGBA(s.Bounds())
	}
	return d.buf
}

func (d WdeDrv) Flip() {
	if d.buf != nil {
		d.win.Screen().CopyRGBA(d.buf, d.buf.Bounds())
	}
	d.win.FlushImage()
}
