package goui

type Unblockable interface {
	Dest() chan<- interface{}
	Unblock() func(src <-chan interface{}, dst chan<- interface{})
}

func coalesce(src, dst chan interface{}) {
	for {
		/* wait for one event */
		latest, ok := <-src
		if !ok {
			return
		}
		/* drain as many as possible */
		for {
			select {
			case l, ok := <-src:
				if !ok {
					dst <- latest
					return
				}
				latest = l
			default:
				break
			}
		}
		/* keep draining until receiver is ready */
		for {
			select {
			case latest = <-src:
			case dst <- latest:
				break
			}
		}
	}
}

func Coalesce(src chan interface{}) chan interface{} {
	dst := make(chan interface{})
	go coalesce(src, dst)
	return dst
}


func broadcast(src chan interface{}, add, rem chan Unblockable) {
	dsts := make(map[Unblockable]chan interface{})
	for {
		select {
		case dst := <-add:
			buf := make(chan interface{})
			dsts[dst] = buf
			go dst.Unblock()(buf, dst.Dest())
		case dst := <-rem:
			close(dsts[dst])
			delete(dsts, dst)
		case ev := <-src:
			for _, c := range dsts {
				c <- ev
			}
		}
	}
}


