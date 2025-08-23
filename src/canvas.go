package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Represents an character in the canvas
type Char struct {
	text     string
	position rl.Vector2
}

// represents the infinite character canvas
type CharCanvas struct {
	chars     []Char
	font      rl.Font
	font_size float32
}

// Retorna o retângulo que cobre toda a área ocupada pelo CharCanvas
func (canvas CharCanvas) VisibleRect() Rect {
	if len(canvas.chars) == 0 {
		return Rect{Start: rl.Vector2{X: 0, Y: 0}, End: rl.Vector2{X: 0, Y: 0}}
	}
	minX := float32(math.Inf(1))
	minY := float32(math.Inf(1))
	maxX := float32(math.Inf(-1))
	maxY := float32(math.Inf(-1))
	for _, obj := range canvas.chars {
		x, y := obj.position.X, obj.position.Y
		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
	}
	return Rect{
		Start: rl.Vector2{X: minX, Y: minY},
		End:   rl.Vector2{X: maxX, Y: maxY},
	}
}

func (canvas *CharCanvas) DeleteCharAt(to_remove rl.Vector2) {
	new_list := []Char{}

	for _, char := range canvas.chars {
		if !rl.Vector2Equals(char.position, to_remove) {
			new_list = append(new_list, char)
		}
	}

	canvas.chars = new_list
}

func (canvas *CharCanvas) ClearSelection(selection CursorSelection) {
	selection.rect.Normalize()

	selection.rect.EachY(func(y float32) {
		selection.rect.EachX(func(x float32) {
			pos := rl.Vector2{X: x, Y: y}
			canvas.DeleteCharAt(pos)
		})
	})
}

func (canvas *CharCanvas) GetCharAt(cell_position rl.Vector2) (Char, bool) {
	for _, char_entry := range canvas.chars {
		if rl.Vector2Equals(char_entry.position, cell_position) {
			return char_entry, true
		}
	}

	return Char{}, false
}

func (canvas *CharCanvas) InsertChar(char Char) {
	_, found_char := canvas.GetCharAt(char.position)

	if found_char {
		canvas.DeleteCharAt(char.position)
	}

	canvas.chars = append(canvas.chars, char)
}

func (canvas *CharCanvas) TextFromSelection(selection CursorSelection) string {
	selection.rect.Normalize()

	str := ""
	selection.rect.EachY(func(y float32) {
		selection.rect.EachX((func(x float32) {
			pos := rl.Vector2{X: float32(x), Y: float32(y)}
			char, found := canvas.GetCharAt(pos)

			if found {
				str += char.text
			} else {
				str += " "
			}
		}))

		str += "\n"
	})

	return str
}

func (canvas *CharCanvas) ToText() string {
	str := ""
	canvas_rect := canvas.VisibleRect()

	canvas_rect.EachY(func(y float32) {
		canvas_rect.EachX((func(x float32) {
			pos := rl.Vector2{X: float32(x), Y: float32(y)}
			char, found := canvas.GetCharAt(pos)

			if found {
				str += char.text
			} else {
				str += " "
			}
		}))

		str += "\n"
	})

	return str
}
