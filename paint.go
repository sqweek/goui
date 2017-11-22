package goui

import (
	"image"
	"image/draw"
	"runtime"
)

type DrawCmd interface {
	Bounds() image.Rectangle
	Partial() bool
	Paint(dst draw.Image)
}

type flushCmd struct{}
func (f flushCmd) Bounds() image.Rectangle { return image.ZR }
func (f flushCmd) Partial() bool { return true }
func (f flushCmd) Paint(dst draw.Image) { }

var Flush = flushCmd{}

func Draw(r image.Rectangle, src image.Image, sp image.Point, op draw.Op) DrawCmd {
	return DrawImg{r, src, sp, op}
}

func DrawMask(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) DrawCmd {
	return DrawImgMask{DrawImg{r, src, sp, op}, mask, mp}
}

type DrawImg struct {
	r image.Rectangle
	src image.Image
	sp image.Point
	op draw.Op
}

func (d DrawImg) Bounds() image.Rectangle {
	return d.r
}

func (d DrawImg) Partial() bool {
	return d.op == draw.Over
}

func (d DrawImg) Paint(dst draw.Image) {
	draw.Draw(dst, d.r, d.src, d.sp, d.op)
}

type DrawImgMask struct {
	DrawImg
	mask image.Image
	mp image.Point
}

func (d DrawImgMask) Partial() bool {
	return true
}

func (d DrawImgMask) Paint(dst draw.Image) {
	draw.DrawMask(dst, d.r, d.src, d.sp, d.mask, d.mp, d.op)
}

type PaintDrv interface {
	Img() draw.Image
	Flip() // TODO include bounds of updated area?
}

func MakePainter(drv PaintDrv) Painter {
	return Painter{screen: drv, queue: &CoalescingQueue{}, cmds: make(chan DrawCmd)}
}

type Painter struct {
	screen PaintDrv
	queue PaintQueue
	cmds chan DrawCmd
}

type PaintQueue interface {
	Add(cmd DrawCmd)
	Drain() ([]DrawCmd, image.Rectangle)
}


func (p *Painter) Queue(cmd DrawCmd) {
	p.cmds <- cmd
}

func (p *Painter) Loop() {
	// TODO ¿split coalescing into separate goroutine with configurable refresh rate, and change p.cmds to []DrawCmd?
	runtime.LockOSThread()
	for {
		cmd, ok := <-p.cmds
		if !ok {
			return
		}
		if cmd == Flush {
			cmds, _ := p.queue.Drain()
			dst := p.screen.Img()
			p.exec(cmds, dst)
			p.screen.Flip()
		} else {
			p.queue.Add(cmd)
		}
	}
}

func (p *Painter) exec(cmds []DrawCmd, dst draw.Image) {
	for _, cmd := range cmds {
		if cmd != nil {
			cmd.Paint(dst)
		}
	}
}

type CoalescingQueue struct {
	pending []DrawCmd
	pBounds image.Rectangle
}

// adds a command to the pending list, possibly rendering previous commands superflous
func (q *CoalescingQueue) Add(cmd DrawCmd) {
	first := (q.pending == nil)
	if !first && !cmd.Partial() && cmd.Bounds().Overlaps(q.pBounds) {
		// Clear previous commands which are obscured
		for i, c := range q.pending {
			if c != nil && c.Bounds().In(cmd.Bounds()) {
				q.pending[i] = nil
			}
		}
	}
	q.pending = append(q.pending, cmd)
	if first {
		q.pBounds = cmd.Bounds()
	} else {
		q.pBounds = q.pBounds.Union(cmd.Bounds())
	}
}

// executes pending commands and clears pending list
func (q *CoalescingQueue) Drain() (cmds []DrawCmd, bounds image.Rectangle) {
	cmds, bounds = q.pending, q.pBounds
	q.pending = nil
	q.pBounds = image.ZR
	return
}
