package printer

import (
	"fmt"
	"github.com/DusanKasan/snaek/internal/app/snaek"
)

type Terminal struct{}

func (p Terminal) Arena(arena snaek.Arena) error {
	fmt.Print("\033[H\033[2J")
	fmt.Println("░░░░░░░░░░░░░░░░░░░░░░░░")

	dim := int64(arena.Dimension)

	for y := int64(0); y < dim; y += 2 {
		fmt.Print("░░")
		for x := int64(0); x < dim; x++ {
			var i = y*dim + x
			switch {
			case arena.State[i] == snaek.FieldFood:
				fmt.Print("⁰")
			case arena.State[i+dim] == snaek.FieldFood:
				fmt.Print("ₒ")
			case arena.State[i] == snaek.FieldSnake && arena.State[i+dim] == snaek.FieldSnake:
				fmt.Print("█")
			case arena.State[i] == snaek.FieldSnake:
				fmt.Print("▀")
			case arena.State[i+dim] == snaek.FieldSnake:
				fmt.Print("▄")
			default:
				fmt.Print(" ")
			}
		}
		fmt.Print("░░\n")
	}
	fmt.Println("░░░░░░░░░░░░░░░░░░░░░░░░")
	fmt.Println("Exit: ESC")
	return nil
}

func (p Terminal) GameOver() error {
	fmt.Println("game over")
	return nil
}

func (p Terminal) Error(err error) error {
	fmt.Println(err)
	return nil
}
