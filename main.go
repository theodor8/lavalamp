package main

import (
	"log"
	"math"
	"math/rand/v2"
	"time"

	"github.com/gdamore/tcell/v2"
)

type ball struct {
	x, y   float64
	vx, vy float64
	size   float64
}

type lava struct {
	balls   []ball
	gravity float64
}

func (l *lava) brightness(x, y float64) float64 {
	var d float64
	for _, ball := range l.balls {
		d += ball.size / math.Pow(math.Pow(ball.x-x, 2)+math.Pow(ball.y-y, 2), 0.65)
	}
	return math.Min(1.0, d)
}

func (l *lava) update() {
	for i1 := range l.balls {
		for i2 := i1 + 1; i2 < len(l.balls); i2++ {
			b1 := &l.balls[i1]
			b2 := &l.balls[i2]
			dx := b2.x - b1.x
			dy := b2.y - b1.y
			d := math.Hypot(dx, dy)
			dx /= d
			dy /= d
			d = math.Max(d, 0.1)
			d2 := d * d
			f1 := l.gravity * b2.size / d2
			f2 := l.gravity * b1.size / d2
			b1.vx += dx * f1
			b1.vy += dy * f1
			b2.vx -= dx * f2
			b2.vy -= dy * f2
		}
	}
	w, h := s.Size()
	for i := range l.balls {
		b := &l.balls[i]
		b.vx = math.Max(-2, math.Min(2, b.vx))
		b.vy = math.Max(-2, math.Min(2, b.vy))
		b.x += b.vx
		b.y += b.vy
		if b.x < -float64(w) || b.x >= 2.0*float64(w) {
			b.vx = -b.vx
		}
		if b.y < -2.0*float64(h) || b.y >= 4.0*float64(h) {
			b.vy = -b.vy
		}
	}
}

func pollEvents() {
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				close(quit)
				return
			}
		}
	}
}

func drawScreen() {
	const gl = '▄'
	w, h := s.Size()
	for x := range w {
		for y := range h {
			y1 := y * 2
			y2 := y1 + 1
			st := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlack)
			b := l.brightness(float64(x), float64(y1))
			st = st.Background(tcell.NewRGBColor(int32(b*255), 0, 0))
			b = l.brightness(float64(x), float64(y2))
			st = st.Foreground(tcell.NewRGBColor(int32(b*255), 0, 0))
			s.SetContent(x, y, gl, nil, st)
		}
	}
	s.Show()
}

var quit chan struct{}
var s tcell.Screen
var l *lava

func main() {

	var err error
	s, err = tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	s.SetStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset))

	s.Clear()
	s.Show()

	l = &lava{
		balls:   []ball{},
		gravity: 0.3,
	}
	w, h := s.Size()
	for range 15 {
		b := ball{
			x:    rand.Float64() * 3.0 * float64(w) - float64(w),
			y:    rand.Float64() * 3.0 * float64(h*2) - float64(h*2),
			vx:   rand.Float64()*2 - 1,
			vy:   rand.Float64()*2 - 1,
			size: rand.Float64()*6 + 9,
		}
		l.balls = append(l.balls, b)
	}
	quit = make(chan struct{})
	go pollEvents()

	go func() {
		for {
			l.update()
			drawScreen()
			time.Sleep(time.Millisecond * 16)
		}
	}()

	<-quit
	s.Fini()
}
