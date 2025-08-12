package main

import (
	"flag"
	"log"
	"math/rand/v2"
	"time"

	"github.com/gdamore/tcell/v2"
)

func pollEvents(s tcell.Screen, l *lava) {
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			w, h := ev.Size()
			l.w, l.h = w, h*2
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				close(quit)
				return
			}
		}
	}
}

func drawScreen(s tcell.Screen, l *lava) {
	const gl = '▄'
	w, h := s.Size()
	for x := range w {
		for y := range h {
			y1 := y * 2
			y2 := y1 + 1
			st := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlack)
			b := int32(l.brightness(float64(x), float64(y1)) * 255)
			st = st.Background(tcell.NewRGBColor(b, 0, 255-b/2))
			b = int32(l.brightness(float64(x), float64(y2)) * 255)
			st = st.Foreground(tcell.NewRGBColor(b, 0, 255-b/2))
			s.SetContent(x, y, gl, nil, st)
		}
	}
	s.Show()
}

var quit chan struct{}

func main() {

	l := &lava{
		balls: []ball{},
	}

	flag.Float64Var(&l.intensity, "i", 0.5, "intensity of the glow")
	flag.Float64Var(&l.gravity, "g", 0.2, "gravity force strength")
	flag.Float64Var(&l.maxVel, "m", 1.0, "maximum velocity of the balls")
	ballsSize := flag.Float64("s", 5.0, "size of the balls")
	numBalls := flag.Int("n", 5, "number of balls")
	flag.Parse()

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	s.SetStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset))

	s.Clear()
	s.Show()

	w, h := s.Size()
	for range *numBalls {
		b := ball{
			x:    rand.Float64() * float64(w),
			y:    rand.Float64() * float64(h*2),
			vx:   rand.Float64()*2 - 1,
			vy:   rand.Float64()*2 - 1,
			size: rand.Float64()**ballsSize + *ballsSize,
		}
		l.balls = append(l.balls, b)
	}
	quit = make(chan struct{})
	go pollEvents(s, l)

	go func() {
		for {
			l.update()
			drawScreen(s, l)
			time.Sleep(time.Millisecond * 16)
		}
	}()

	<-quit
	s.Fini()
}
