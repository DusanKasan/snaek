package main

import (
	"context"
	"github.com/DusanKasan/snaek/internal/app/snaek"
	"github.com/DusanKasan/snaek/internal/app/snaek/printer"
	"github.com/DusanKasan/snaek/internal/lib/terminal"
	"log"
	"time"
)



func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	keys, err := terminal.Keystrokes(ctx)
	if err != nil {
		panic(err)
	}

	game := snaek.NewGame(20, time.Second/4, printer.Terminal{})
	errs := game.Start(ctx)

	for{
		select {
		case err := <-errs:
			if err != nil {
				log.Panic(err)
			}

			cancelFunc()
			return
		case key := <-keys:
			switch key {
			case terminal.KeyArrowRight:
				game.Move(snaek.DirectionRight)
			case terminal.KeyArrowDown:
				game.Move(snaek.DirectionDown)
			case terminal.KeyArrowLeft:
				game.Move(snaek.DirectionLeft)
			case terminal.KeyArrowUp:
				game.Move(snaek.DirectionUp)
			case terminal.KeyEsc:
				cancelFunc()
				return
			}
		}
	}
}
