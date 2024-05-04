package rectx

type Rects []Rectangle

func (rex *Rects) Include(area image.Rectangle) bool {
	if len(rex) == 0 {
		*rex = append(rex, area)
		return true
	} else {
		for i, r := range rex {
			if r.Contains(area) {
				return false
			}
		}
	}
}

func (r *Rects) Exclude(area image.Rectangle) {
}










type Occlusion struct {
	r image.Rectangle
	visible image.Rectangle // bounds of visible area
	partials []image.Rectangle
}

type Rects []image.Rectangle


const (
	TopLeft = iota
	TopRight
	BottomRight
	BottomLeft
)

func corners(r image.Rectangle) [4]image.Point {
	max := r.Max.Sub(image.Point{1, 1})
	return [4]image.Point{r.Min, image.Point{max.X, r.Min.Y}, max, image.Point{r.Min.X, max.Y}}
}


func subtract(large, small image.Rect) (image.Rect, bool) {
	result, changed := large, false
	if small.Min.X == large.Min.X && small.Max.X == large.Max.X {
		if small.Min.Y <= large.Max.Y {
			result.Max.Y, changed = small.Min.Y - 1, true
		}
		if small.Max.Y >= large.Min.Y {
			result.Min.Y, changed = small.Max.Y + 1, true
		}
	}
	if small.Min.Y == large.Min.Y && small.Max.Y == large.Max.Y {
		if small.Min.X <= large.Max.X {
			result.Max.X, changed = small.Min.X - 1, true
		}
		if small.Max.X >= large.Min.X {
			result.Min.X, changed = small.Max.X + 1, true
		}
	}
	return result, changed
}

func join(r1

func (o *Occlusion) Occlude(r image.Rectangle) {
	if area := r.Intersect(o.visible); !area.Empty() {
		if sub, changed := subtract(o.visible, area); changed {
			o.visible = sub
			return
		}
		for i, r2 := range o.partials {
			if area.In(r2) {
				return // this area already accounted for
			}
		}
		o.partials = append(o.partials, area)
		
	}
}
