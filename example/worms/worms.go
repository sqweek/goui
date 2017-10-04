package main

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"time"

	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
	"github.com/sqweek/goui"
	"github.com/sqweek/goui/wdedrv"
)

type World struct {
	width, height int
}

func main() {
	go func() {
		wdemain()
		wde.Stop()
	}()
	wde.Run()
}

type Direction int

const (
	N Direction = iota; NE; E; SE; S; SW; W; NW
)

func (d Direction) dx() int {
	switch d {
	case NE, E, SE:
		return 1
	case SW, W, NW:
		return -1
	}
	return 0
}

func (d Direction) dy() int {
	switch d {
	case NW, N, NE:
		return -1
	case SE, S, SW:
		return 1
	}
	return 0
}

func (d Direction) Advance(pt image.Point) image.Point {
	return image.Point{pt.X + d.dx(), pt.Y + d.dy()}
}

func (d Direction) Pt() image.Point {
	return image.Point{d.dx(), d.dy()}
}

func (d Direction) ReflectX() Direction {
	if d % 4 == 0 {
		return d /* 0/4 (N/S) remain unchanged */
	}
	return 8 - d    /* 1 <-> 7, 2 <-> 6, 3 <-> 5 */
}

func (d Direction) ReflectY() Direction {
	if d % 4 == 2 {
		return d /* 2/6 (E/W) remain unchanged */
	}
	seg := 4*((d-1)/4)
	return (seg + 4) - (d - seg)   /* 0 <-> 4, 1 <-> 3, 5 <-> 7 */
}

func (d Direction) Turn(way int) Direction {
	n := (int(d) + way) % 8
	if n < 0 {
		n += 8
	}
	return Direction(n)
}

var world = World{320, 240}
var WORMS = 8192

func wdemain() {
	rand.Seed(time.Now().UnixNano())
	w, err := wde.NewWindow(world.width, world.height)
	if err != nil {
		panic(err)
	}
	w.SetTitle("Worms")
	w.Show()
	painter := goui.MakePainter(wdedrv.Make(w))
	for i := 0; i < WORMS; i++ {
		go worm(painter)
	}
	world.randPt()
	go painter.Loop()
	events: for ei := range w.EventChan() {
		switch e := ei.(type) {
		case wde.KeyEvent:
			println(e.Key)
		case wde.CloseEvent:
			break events
		}
	}
}

func (w World) randPt() image.Point {
	return image.Pt(rand.Intn(w.width), rand.Intn(w.height))
}

func (w World) ContainsX(x int) bool {
	return x >= 0 && x < w.width
}

func (w World) ContainsY(y int) bool {
	return y >= 0 && y < w.height
}

type Worm struct {
	pts []image.Point
	dir Direction
	//speed
}

func randWorm(world World) Worm {
	worm := Worm{
		pts: make([]image.Point, 0, rand.Intn(8)+3),
		dir: Direction(rand.Intn(8)),
	}
	worm.pts = append(worm.pts, world.randPt())
	for len(worm.pts) < cap(worm.pts) {
		worm.advance()
	}
	return worm
}

func (w *Worm) advance() {
	if len(w.pts) < cap(w.pts) {
		w.pts = append(w.pts, image.Point{0, 0})
	}
	n := len(w.pts)
	copy(w.pts[1:n], w.pts[0:n-1])
	w.pts[0] = w.dir.Advance(w.pts[1])
	reflect := false
	if !world.ContainsX(w.pts[0].X) {
		w.dir, reflect = w.dir.ReflectX(), true
	}
	if !world.ContainsY(w.pts[0].Y) {
		w.dir, reflect = w.dir.ReflectY(), true
	}
	if reflect {
		w.pts[0] = w.dir.Advance(w.pts[1])
	}
}

func (w *Worm) turn(way int) {
	w.dir = w.dir.Turn(way)
}

type DrawWormCmd struct {
	w Worm
	col color.Color
}

func (d DrawWormCmd) Bounds() image.Rectangle {
	r := image.Rectangle{Min: d.w.pts[0], Max: d.w.pts[0]}
	for _, pt := range d.w.pts {
		if pt.X < r.Min.X {
			r.Min.X = pt.X
		} else if pt.X + 1 > r.Max.X {
			r.Max.X = pt.X + 1
		}
		if pt.Y < r.Min.Y {
			r.Min.Y = pt.Y
		} else if pt.Y + 1 > r.Max.Y {
			r.Max.Y = pt.Y + 1
		}
	}
	return r
}

func (d DrawWormCmd) Partial() bool {
	return true
}

func (d DrawWormCmd) Paint(dst draw.Image) {
	r := d.Bounds()
	draw.DrawMask(dst, r, d, r.Min, d, r.Min, draw.Over)
}

func (d DrawWormCmd) ColorModel() color.Model {
	return color.RGBAModel
}

func (d DrawWormCmd) At(x, y int) color.Color {
	for _, pt := range d.w.pts {
		if pt.X == x && pt.Y == y {
			return d.col
		}
	}
	return color.Transparent
}

func randCol() color.Color {
	bits := 1 + rand.Intn(6)
	c := []uint8{0x88, 0x88, 0x88}
	for i := range c {
		if bits & (1 << uint(i)) != 0 {
			c[i] += 0x33 * uint8(2*(float32(rand.Intn(5)) - 1.5))
		}
	}
	return color.NRGBA{c[0], c[1], c[2], 0x44}
}

func worm(painter goui.Painter) {
	worm := randWorm(world)
	col := randCol()
	speed := time.Millisecond * time.Duration(15 + rand.Intn(36))
	for {
		//println(worm.pts[0].X, worm.pts[1].Y, worm.dir.dx(), worm.dir.dy())
		worm.advance()
		time.Sleep(speed)
		worm.turn(rand.Intn(3) - 1)
		painter.Queue(DrawWormCmd{worm, col})
	}	
}
