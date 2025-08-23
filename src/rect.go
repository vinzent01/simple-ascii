package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Abstração de retângulo
type Rect struct {
	Start rl.Vector2
	End   rl.Vector2
}

// Itera sobre todas as posições X inteiras dentro do retângulo
func (r Rect) EachX(f func(x float32)) {
	for x := int(r.Start.X); x <= int(r.End.X); x++ {
		f(float32(x))
	}
}

// Itera sobre todas as posições Y inteiras dentro do retângulo
func (r Rect) EachY(f func(y float32)) {
	for y := int(r.Start.Y); y <= int(r.End.Y); y++ {
		f(float32(y))
	}
}

// Verifica se um ponto está dentro do retângulo
func (r Rect) Contains(v rl.Vector2) bool {
	return v.X >= r.Start.X && v.X <= r.End.X && v.Y >= r.Start.Y && v.Y <= r.End.Y
}

// Itera sobre todas as posições inteiras dentro do retângulo
func (r Rect) Each(f func(pos rl.Vector2)) {
	for y := int(r.Start.Y); y <= int(r.End.Y); y++ {
		for x := int(r.Start.X); x <= int(r.End.X); x++ {
			f(rl.Vector2{X: float32(x), Y: float32(y)})
		}
	}
}

func (r *Rect) Normalize() {
	startX := math.Min(float64(r.Start.X), float64(r.End.X))
	startY := math.Min(float64(r.Start.Y), float64(r.End.Y))
	endX := math.Max(float64(r.Start.X), float64(r.End.X))
	endY := math.Max(float64(r.Start.Y), float64(r.End.Y))

	r.Start = rl.Vector2{X: float32(startX), Y: float32(startY)}
	r.End = rl.Vector2{X: float32(endX), Y: float32(endY)}
}
