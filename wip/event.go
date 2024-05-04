package goui

import (
	"time"
)

type Event interface {
	Time() time.Time
}

// per-pointer identifier?
type MouseEvent interface {
	Event
	X() int
	Y() int
}

type MouseDownEvent interface {
	MouseEvent
	Button() int
}

type MouseUpEvent interface {
	MouseEvent
	Button() int
}



type KeyEvent interface {
	Key() int
}





type EventHub interface {
	Chan() chan State
}

type State interface {
	Apply(*Transition) State
}

type Transition interface {
}

