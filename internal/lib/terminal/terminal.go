package terminal

import (
	"context"
	"github.com/nsf/termbox-go"
	"sync"
)

type Key termbox.Key

const (
	KeyArrowRight = Key(termbox.KeyArrowRight)
	KeyArrowDown = Key(termbox.KeyArrowDown)
	KeyArrowLeft = Key(termbox.KeyArrowLeft)
	KeyArrowUp = Key(termbox.KeyArrowUp)
	KeyEsc = Key(termbox.KeyEsc)
)

func Keystrokes(ctx context.Context) (chan Key, error) {
	if err := termbox.Init(); err != nil {
		return nil, err
	}

	var done bool
	var mux sync.Mutex
	go func() {
		select {
		case <-ctx.Done():
			mux.Lock()
			done = true
			mux.Unlock()
		}
	}()

	var keys = make(chan Key)
	go func() {
		for {
			event := termbox.PollEvent()
			if event.Type != termbox.EventKey {
				continue
			}

			mux.Lock()
			if done {
				termbox.Close()
				close(keys)
				return
			}
			mux.Unlock()

			keys <- Key(event.Key)
		}
	}()

	return keys, nil
}
