package main

import (
	"context"
	"math/rand"
	"time"

	"fortio.org/terminal/ansipixels/tcolor"
)

type streak struct {
	x, y  int
	char  rune
	color tcolor.RGBColor
}

type matrix struct {
	maxX, maxY int
	streaks    chan (streak)
}

func (m *matrix) newStreak(ctx context.Context, speedDividend int) {
	s := streak{0, rand.Intn(m.maxY), rune(rand.Intn(128) + 8), BrightGreen}
	speed := rand.Intn(100)
	timeBetween := max(time.Duration(speed*int(time.Millisecond)/speedDividend), 100)
	go func() {
		ticker := time.NewTicker(timeBetween)
		defer func() {
			ticker.Stop()
		}()
		m.streaks <- s
		for {
			select {
			case <-ticker.C:
				s.x++
				if s.x >= m.maxX {
					return
				}
				s.char = rune(rand.Intn(128) + 8) // + 8 because we don't wanna hit the bell character
				for s.char == '\n' || s.char == '\r' || s.char == ' ' {
					s.char = rune(rand.Intn(128) + 8) // + 8 because we don't wanna hit the bell character
				}
				m.streaks <- s
			case <-ctx.Done():
				return
			}
		}
	}()
}
