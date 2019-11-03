package snaek

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type err string

func (e err) Error() string {
	return string(e)
}

const (
	ErrUnableToMoveInReverseDirection err = "unable to move in reverse direction"
	ErrObstacleHit                    err = "obstacle hit"
)

type field int8

const (
	FieldEmpty field = iota
	FieldSnake
	FieldFood
)

type direction int8

const (
	DirectionRight direction = iota
	DirectionDown
	DirectionLeft
	DirectionUp
)

type Printer interface {
	Arena(Arena) error
	GameOver() error
	Error(error) error
}

type Arena struct {
	// x = y = Dimension
	Dimension int16
	// x = i/Dimension
	// y = i%Dimension
	State map[int64]field
	// fields filled by Snake starting from head
	Snake []int64
}

func New(dimension int16) Arena {
	rand.Seed(time.Now().UnixNano())

	arena := Arena{
		Dimension: dimension,
		State:     map[int64]field{},
	}

	var centre = int64(arena.Dimension) * int64(arena.Dimension) / 2
	arena.State[centre] = FieldSnake
	arena.State[centre+1] = FieldSnake
	arena.Snake = []int64{centre + 1, centre}
	arena.spawnFood()

	return arena
}

func (a *Arena) spawnFood() {
	var size = int64(a.Dimension) * int64(a.Dimension)
	var empty = size - int64(len(a.Snake))

	if empty == 0 {
		// you win -> emit
		return
	}

	if empty < 0 {
		panic("empty less than 0")
	}

	foodslot := rand.Intn(int(empty))
	for i := int64(0); i < size; i++ {
		if a.State[i] == FieldEmpty {
			foodslot--

		}

		if foodslot == 0 {
			a.State[i] = FieldFood
			return
		}
	}
}

func (a *Arena) Move(d direction) error {
	var head = a.Snake[0]
	var dim = int64(a.Dimension)
	var newHead int64

	switch d {
	case DirectionRight:
		if head%dim+1 == dim {
			return ErrObstacleHit
		}
		newHead = head + 1
	case DirectionDown:
		if head/dim == dim-1 {
			return ErrObstacleHit
		}
		newHead = head + dim
	case DirectionLeft:
		if head%dim == 0 {
			return ErrObstacleHit
		}
		newHead = head - 1
	case DirectionUp:
		if head/dim == 0 {
			return ErrObstacleHit
		}
		newHead = head - dim
	default:
		return fmt.Errorf("invalid direction: %v", d)
	}

	if newHead == a.Snake[1] {
		return ErrUnableToMoveInReverseDirection
	}

	if a.State[newHead] == FieldSnake {
		return ErrObstacleHit
	}

	var grow = a.State[newHead] == FieldFood
	a.Snake = append([]int64{newHead}, a.Snake...)
	a.State[newHead] = FieldSnake

	if !grow {
		a.State[a.Snake[len(a.Snake)-1]] = FieldEmpty
		a.Snake = a.Snake[0 : len(a.Snake)-1]
	} else {
		a.spawnFood()
	}

	return nil
}

func NewGame(dimension int16, tick time.Duration, printer Printer) Game {
	ticker := time.NewTicker(tick)
	return Game{
		printer: printer,
		arena:   New(dimension),
		ticker:  *ticker,
	}
}

type Game struct {
	printer   Printer
	arena     Arena
	direction direction
	ticker    time.Ticker
	mux       sync.Mutex
}

func (g *Game) Move(d direction) {
	g.mux.Lock()
	g.direction = d
	g.mux.Unlock()
}

func (g *Game) Start(ctx context.Context) chan error {
	var errs = make(chan error)
	go func() {
		for {
			select {
			case <-g.ticker.C:
				g.mux.Lock()
				err := g.arena.Move(g.direction)
				if err == ErrUnableToMoveInReverseDirection {
					err = g.arena.Move((g.direction + 2) % 4)
				}

				switch err {
				case nil:
					if err := g.printer.Arena(g.arena); err != nil {
						errs <- err
					}
				case ErrObstacleHit:
					if err := g.printer.GameOver(); err != nil {
						errs <- err
					}
					close(errs)
					return
				default:
					if e := g.printer.Error(err); e != nil {
						// TODO: Wrap
						err = fmt.Errorf("%v: %v", e, err)
					}
					errs <- err
					close(errs)
					return
				}
				g.mux.Unlock()
			case <-ctx.Done():
				close(errs)
				return
			}
		}
	}()
	return errs
}
