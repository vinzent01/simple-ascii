package main

import rl "github.com/gen2brain/raylib-go/raylib"

type CursorSelection struct {
	rect      Rect
	selecting bool
}

// Represents the keyboard and cursor
type KeyBoardCursor struct {
	Position         rl.Vector2
	Visible          bool
	Active           bool
	BlinkTime        float32
	BlinkPeriod      float32
	CurrentSelection CursorSelection
}

func NewCursor() KeyBoardCursor {
	return KeyBoardCursor{
		Position:         rl.Vector2{X: 0, Y: 0},
		Visible:          true,
		Active:           false,
		BlinkTime:        0.0,
		BlinkPeriod:      0.25,
		CurrentSelection: CursorSelection{},
	}
}

// / Updates the cursor blink
func (c *KeyBoardCursor) Update(dt float32, mousePos rl.Vector2, font rl.Font, fontSize float32) {
	c.BlinkTime += dt
	if c.BlinkTime > c.BlinkPeriod {
		c.Visible = !c.Visible
		c.BlinkTime = 0
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		c.Active = true
		c.Position = world_to_cell_pos(mousePos, font, fontSize)
		c.CurrentSelection.rect.Start = c.Position
		c.CurrentSelection.rect.End = c.Position
		c.CurrentSelection.selecting = true

	} else if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		c.Position = world_to_cell_pos(mousePos, font, fontSize)
		c.CurrentSelection.rect.End = c.Position
	} else if rl.IsKeyPressed(rl.KeyEscape) {
		c.Active = false
		c.CurrentSelection = CursorSelection{} // limpa tudo
	}

	if rl.IsKeyDown(rl.KeySpace) {
		c.CurrentSelection.selecting = false
	}

}

func (c *KeyBoardCursor) Draw(font rl.Font, fontSize float32) {
	if c.Active && c.Visible {
		DrawCursorAtCanvasCell(c.Position, font, fontSize)
	}
}

func (c *KeyBoardCursor) GetCharPressed() (Char, bool) {
	typed_key := rl.GetCharPressed()

	// Check if a character was pressed
	if typed_key > 0 {
		// Only add if it's a printable character
		if (typed_key >= 32) && (typed_key <= 126) {
			new_char := Char{
				text:     string(typed_key),
				position: c.Position,
			}

			return new_char, true
		}
	}

	return Char{}, false
}

func (c *KeyBoardCursor) GetInputPressed() int32 {
	typed_key := rl.GetKeyPressed()

	if typed_key == rl.KeyUp {
		c.Position.Y -= 1
	} else if typed_key == rl.KeyDown {
		c.Position.Y += 1

	} else if typed_key == rl.KeyLeft {
		c.Position.X -= 1

	} else if typed_key == rl.KeyRight {
		c.Position.X += 1
	}

	return typed_key
}
