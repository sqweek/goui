package goui

import (
	"image"
	"image/draw"
)

type DrawCmd interface {
	Bounds() image.Rectangle
	Partial() bool
	Paint(dst draw.Image)
}

type DrawImg struct {
	r image.Rectangle
	src image.Image
	sp image.Point
	mask image.Image
	mp image.Point
	op draw.Op
}

func (d DrawImg) Bounds() image.Rectangle {
	return d.r
}

func (d DrawImg) Partial() bool {
	return false
}

func (d DrawImg) Paint(dst draw.Image) {
	draw.DrawMask(dst, d.r, d.src, d.sp, d.mask, d.mp, d.op)
}

type PaintDrv interface {
	Img() draw.Image
	Flip() // TODO include bounds of updated area?
}

func MakePainter(drv PaintDrv) Painter {
	return Painter{screen: drv, cmds: make(chan DrawCmd)}
}

type Painter struct {
	screen PaintDrv 
	cmds chan DrawCmd
	pending []DrawCmd
	pBounds image.Rectangle
}

func (p *Painter) Queue(cmd DrawCmd) {
	p.cmds <- cmd
}

func (p *Painter) Loop() {
	// TODO ¿split coalescing into separate goroutine with configurable refresh rate, and change p.cmds to []DrawCmd?
	for {
		cmd, ok := <-p.cmds
		if !ok {
			return
		}
		p.add(cmd)
		for {
			if p.drain() == 0 {
				break
			}
		}
		dst := p.screen.Img()
		p.exec(dst)
		p.screen.Flip()
	}
}

func (p *Painter) drain() (n int) {
	for {
		select {
		case cmd := <-p.cmds:
			p.add(cmd)
			n++
		default:
			return n
		}
	}
}

// adds a command to the pending list, possibly rendering previous commands superflous
func (p *Painter) add(cmd DrawCmd) {
	first := (p.pending == nil)
	if !first && !cmd.Partial() && cmd.Bounds().Overlaps(p.pBounds) {
		// Clear previous commands which are obscured
		for i, c := range p.pending {
			if c != nil && c.Bounds().In(cmd.Bounds()) {
				p.pending[i] = nil
			}
		}
	}
	p.pending = append(p.pending, cmd)
	if first {
		p.pBounds = cmd.Bounds()
	} else {
		p.pBounds = p.pBounds.Union(cmd.Bounds())
	}
}

// executes pending commands and clears pending list
func (p *Painter) exec(dst draw.Image) {
	for _, cmd := range p.pending {
		if cmd != nil {
			cmd.Paint(dst)
		}
	}
	p.pending = nil
	p.pBounds = image.ZR // does this need to be retained for flip()?
}

