package main

import "math"

type ball struct {
	x, y   float64
	vx, vy float64
	size   float64
}

type lava struct {
	balls     []ball
	gravity   float64
	maxVel    float64
	intensity float64
	w, h int
}

func (l *lava) brightness(x, y float64) float64 {
	var d float64
	for _, ball := range l.balls {
		d += ball.size / math.Pow(math.Pow(ball.x-x, 2)+math.Pow(ball.y-y, 2), 1.0-l.intensity)
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
	for i := range l.balls {
		b := &l.balls[i]
		b.vx = math.Max(-l.maxVel, math.Min(l.maxVel, b.vx))
		b.vy = math.Max(-l.maxVel, math.Min(l.maxVel, b.vy))
		b.x += b.vx
		b.y += b.vy
		if b.x < 0 {
			b.vx += 0.05
		} else if b.x >= float64(l.w) {
			b.vx -= 0.05
		}
		if b.y < 0 {
			b.vy += 0.05
		} else if b.y >= float64(l.h) {
			b.vy -= 0.05
		}
	}
}
