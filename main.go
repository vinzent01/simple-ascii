package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.design/x/clipboard"
)

// Represents the cursor
type Cursor struct {
	Position    rl.Vector2
	Visible     bool
	Active      bool
	BlinkTime   float32
	BlinkPeriod float32
}

func NewCursor() Cursor {
	return Cursor{
		Position:    rl.Vector2{X: 0, Y: 0},
		Visible:     true,
		Active:      false,
		BlinkTime:   0.0,
		BlinkPeriod: 0.25,
	}
}

// / Updates the cursor blink
func (c *Cursor) Update(dt float32, mousePos rl.Vector2, font rl.Font, fontSize float32) {
	c.BlinkTime += dt
	if c.BlinkTime > c.BlinkPeriod {
		c.Visible = !c.Visible
		c.BlinkTime = 0
	}

	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		c.Active = true
		c.Position = world_to_cell_pos(mousePos, font, fontSize)
	} else if rl.IsMouseButtonReleased(rl.MouseButtonRight) {
		c.Active = false
	}

}

func (c *Cursor) Draw(font rl.Font, fontSize float32) {
	if c.Active && c.Visible {
		DrawCursorAtCanvasCell(c.Position, font, fontSize)
	}
}

func (c *Cursor) GetCharPressed() (Char, bool) {
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

func (c *Cursor) GetInputPressed() int32 {
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

// Represents an character in the canvas
type Char struct {
	text     string
	position rl.Vector2
}

// represents the infinite character canvas
type CharCanvas []Char

func load_fonts() map[string]rl.Font {
	fonts_path := "assets/fonts"

	loaded := map[string]rl.Font{}

	entries, err := os.ReadDir(fonts_path)

	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {

		if entry.Type().IsRegular() {
			filename := entry.Name()
			ext := filepath.Ext(entry.Name())
			name := filename[:len(filename)-len(ext)]
			loaded[name] = rl.LoadFont(filepath.Join(fonts_path, filename))
		}
	}

	return loaded

}

func world_to_cell_pos(position rl.Vector2, font rl.Font, font_size float32) rl.Vector2 {
	text_size := rl.MeasureTextEx(font, "A", font_size, 0)

	return rl.Vector2{
		X: float32(
			math.Floor(float64(position.X) / float64(text_size.X)),
		),
		Y: float32(
			math.Floor(float64(position.Y) / float64(text_size.Y)),
		),
	}
}

func cell_to_world_pos(position rl.Vector2, font rl.Font, font_size float32) rl.Vector2 {
	text_size := rl.MeasureTextEx(font, "A", font_size, 0)
	return rl.Vector2Multiply(position, text_size)
}

func draw_canvas(canvas CharCanvas, font rl.Font, font_size float32) {

	for _, char := range canvas {
		DrawCharAtCanvasCell(char, font, font_size, rl.White, rl.Black)
	}
}

func DrawCharAtCanvasCell(char Char, font rl.Font, font_size float32, foreground rl.Color, background rl.Color) {
	// calcula o tamanho do texto
	text_size := rl.MeasureTextEx(font, char.text, float32(font_size), 0)
	world_position := cell_to_world_pos(char.position, font, font_size)

	// desenha o fundo (retângulo)
	rl.DrawRectangleV(world_position, text_size, background)

	// desenha o texto por cima
	rl.DrawTextEx(font, char.text, world_position, font_size, 0, foreground)
}

func DrawCursorAtCanvasCell(cell_position rl.Vector2, font rl.Font, fontsize float32) {

	DrawCharAtCanvasCell(
		Char{
			text:     " ",
			position: cell_position,
		},
		font,
		fontsize,
		rl.White,
		rl.White,
	)
}

func GetCharAtCanvas(char_canvas CharCanvas, cell_position rl.Vector2) (Char, bool) {
	for _, char_entry := range char_canvas {
		if rl.Vector2Equals(char_entry.position, cell_position) {
			return char_entry, true
		}
	}

	return Char{}, false
}

func DeleteCharAtCanvas(char_canvas CharCanvas, to_remove rl.Vector2) CharCanvas {

	new_canvas := CharCanvas{}

	for _, char := range char_canvas {
		if !rl.Vector2Equals(char.position, to_remove) {
			new_canvas = append(new_canvas, char)
		}
	}

	return new_canvas
}

func InserCharAtCanvas(char_canvas CharCanvas, char Char) CharCanvas {
	_, found_char := GetCharAtCanvas(char_canvas, char.position)
	if found_char {
		char_canvas = DeleteCharAtCanvas(char_canvas, char.position)
	}
	char_canvas = append(char_canvas, char)
	return char_canvas
}

func CreateASCIIFromMap(char_canvas CharCanvas) string {
	// Função auxiliar para achar o canto superior esquerdo
	getTopLeftCorner := func() rl.Vector2 {
		minX := float32(math.Inf(1))
		minY := float32(math.Inf(1))

		for _, obj := range char_canvas {
			x, y := obj.position.X, obj.position.Y
			if x < minX {
				minX = x
			}
			if y < minY {
				minY = y
			}
		}

		if minX == float32(math.Inf(1)) {
			minX = 0
		}
		if minY == float32(math.Inf(1)) {
			minY = 0
		}

		return rl.Vector2{X: minX, Y: minY}
	}

	// Função auxiliar para achar o canto inferior direito
	getBottomRightCorner := func() rl.Vector2 {
		maxX := float32(math.Inf(-1))
		maxY := float32(math.Inf(-1))

		for _, obj := range char_canvas {
			x, y := obj.position.X, obj.position.Y
			if x > maxX {
				maxX = x
			}
			if y > maxY {
				maxY = y
			}
		}

		if maxX == float32(math.Inf(-1)) {
			maxX = 0
		}
		if maxY == float32(math.Inf(-1)) {
			maxY = 0
		}

		return rl.Vector2{X: maxX, Y: maxY}
	}

	topLeft := getTopLeftCorner()
	bottomRight := getBottomRightCorner()

	asciiStr := "```\n"

	for y := int(topLeft.Y); y <= int(bottomRight.Y); y++ {
		for x := int(topLeft.X); x <= int(bottomRight.X); x++ {
			pos := rl.Vector2{X: float32(x), Y: float32(y)}
			char, found := GetCharAtCanvas(char_canvas, pos)
			if found {
				asciiStr += char.text
			} else {
				asciiStr += " " // vazio
			}
		}
		asciiStr += "\n"
	}

	asciiStr += "```"
	return asciiStr
}

func Save(ascii string) error {
	saveFolder := "save"

	// Cria a pasta "save" se não existir
	err := os.MkdirAll(saveFolder, os.ModePerm)
	if err != nil {
		return err
	}

	index := 1
	var newFile string

	// Procura o primeiro nome disponível
	for {
		newFile = filepath.Join(saveFolder, fmt.Sprintf("art%d.txt", index))
		if _, err := os.Stat(newFile); os.IsNotExist(err) {
			break
		}
		index++
	}

	// Escreve o conteúdo no arquivo
	return os.WriteFile(newFile, []byte(ascii), 0644)
}

func PasteASCIIAtCanvas(char_canvas CharCanvas, ascii string, startPos rl.Vector2) CharCanvas {
	ascii = strings.ReplaceAll(ascii, "```", "")
	lineList := []string{}
	currentLine := ""

	for _, r := range ascii {
		if r == '\n' {
			lineList = append(lineList, currentLine)
			currentLine = ""
		} else {
			currentLine += string(r)
		}
	}
	if currentLine != "" {
		lineList = append(lineList, currentLine)
	}

	for y, line := range lineList {
		for x, ch := range line {
			if ch != ' ' && ch != '\r' {
				pos := rl.Vector2{
					X: startPos.X + float32(x),
					Y: startPos.Y + float32(y),
				}
				char_canvas = InserCharAtCanvas(char_canvas, Char{
					text:     string(ch),
					position: pos,
				})
			}
		}
	}
	return char_canvas
}

func DrawHelpPanel(screenWidth, screenHeight int32, font rl.Font) {
	helpText := "COMMANDS:\n" +
		"Ctrl + S - saves the art\n" +
		"Ctrl + X - clears the screen\n" +
		"Ctrl + V - pastes content from clipboard to the canvas\n" +
		"Ctrl + H - shows help screen\n" +
		"Typing keys to the canvas\n" +
		"Drag camera: hold SPACE + left mouse button"

	panelfontSize := float32(16.0)
	text_size := rl.MeasureTextEx(font, helpText, panelfontSize, 0)
	panelWidth := int32(text_size.X + 40)
	panelHeight := int32(text_size.Y + 40)

	panelX := (screenWidth - panelWidth) / 2
	panelY := (screenHeight - panelHeight) / 2

	rl.DrawRectangle(panelX, panelY, panelWidth, panelHeight, rl.DarkGray)
	rl.DrawRectangleLines(panelX, panelY, panelWidth, panelHeight, rl.White)
	rl.DrawTextEx(
		font,
		helpText,
		rl.Vector2{
			X: float32(panelX + 10),
			Y: float32(panelY + 10),
		},
		panelfontSize,
		0,
		rl.White,
	)
}

func main() {
	fmt.Println("Hello, world")
	var screen_width int32 = 800
	var screen_height int32 = 450

	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(screen_width, screen_height, "simple ascii")
	rl.SetTargetFPS(60)

	// fonts
	fonts := load_fonts()
	current_font, ok := fonts["Consolas-Regular"]

	var font_size float32 = 24.0
	var ctrl_active bool = false
	var space_active bool = false

	if !ok {
		fmt.Print("default font not found, exiting...")
		os.Exit(0)
	}

	// clipboard
	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	cursor := NewCursor()

	var char_canvas CharCanvas = CharCanvas{}
	camera := rl.Camera2D{}

	camera.Offset = rl.Vector2{X: float32(screen_width) / 2.0, Y: float32(screen_height) / 2.0}
	camera.Target = rl.Vector2{X: 0.0, Y: 0.0}
	camera.Zoom = 1.0

	showHelp := false

	// Load fonts
	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime() // tempo desde o último frame em segundos
		mouse_pos := rl.GetMousePosition()
		mouse_world := rl.GetScreenToWorld2D(mouse_pos, camera)

		// update
		cursor.Update(dt, mouse_world, current_font, font_size)
		ctrl_active = rl.IsKeyDown(rl.KeyLeftControl)
		space_active = rl.IsKeyDown(rl.KeySpace)

		if cursor.Active {
			char, typed := cursor.GetCharPressed()
			keypressed := cursor.GetInputPressed()

			if typed {
				char_canvas = InserCharAtCanvas(char_canvas, char)
				cursor.Position.X += 1
			}

			if keypressed == rl.KeyBackspace {
				char_canvas = DeleteCharAtCanvas(char_canvas, cursor.Position)
				cursor.Position.X -= 1
			}

			if keypressed == rl.KeyEnter {
				cursor.Position.Y += 1
			}
		}

		if ctrl_active && rl.IsKeyPressed(rl.KeyS) {
			// save draw
			ascii := CreateASCIIFromMap(char_canvas)
			Save(ascii)
		} else if ctrl_active && rl.IsKeyPressed(rl.KeyX) {
			char_canvas = CharCanvas{}
		} else if ctrl_active && rl.IsKeyPressed(rl.KeyV) {
			// paste input
			clipboardText := string(clipboard.Read(clipboard.FmtText))
			char_canvas = PasteASCIIAtCanvas(char_canvas, clipboardText, cursor.Position)
		} else if ctrl_active && rl.IsKeyPressed(rl.KeyH) {
			showHelp = !showHelp
		}

		if rl.IsKeyPressed(rl.KeyF11) {
			rl.ToggleFullscreen()
		}

		// camera zoom
		camera.Zoom = float32(
			math.Exp(
				math.Log(float64(camera.Zoom)) + float64(rl.GetMouseWheelMove())*0.1,
			),
		)

		// camera movement
		if space_active && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			mouse_delta := rl.GetMouseDelta()
			camera.Target.X -= mouse_delta.X
			camera.Target.Y -= mouse_delta.Y
		}

		// drawing
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.BeginMode2D(camera)
		draw_canvas(char_canvas, current_font, font_size)
		cursor.Draw(current_font, font_size)
		rl.EndMode2D()

		if showHelp {
			DrawHelpPanel(int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), current_font)
		}

		// draws the help text remember
		help_remember := "Ctrl + h (help)"
		help_size := rl.MeasureTextEx(current_font, help_remember, 16, 0)

		rl.DrawTextEx(
			current_font, help_remember,
			rl.Vector2{
				X: float32(rl.GetScreenWidth()) - help_size.X - 10,
				Y: 0,
			},
			16,
			0,
			rl.Yellow,
		)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
