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

	for _, char := range canvas.chars {
		DrawCharAtCanvasCell(char, font, font_size, rl.White, rl.Black)
	}
}

func draw_selected_canvas(canvas CharCanvas, selection CursorSelection, font rl.Font, font_size float32) {
	selection.rect.Normalize()

	selection.rect.EachY(func(y float32) {
		selection.rect.EachX(func(x float32) {
			pos := rl.Vector2{X: x, Y: y}
			charAt, found := canvas.GetCharAt(pos)

			var draw_char Char

			if found {
				draw_char = charAt
			} else {
				draw_char = Char{
					text:     " ",
					position: pos,
				}
			}

			DrawCharAtCanvasCell(draw_char, font, font_size, rl.Black, rl.White)
		})
	})

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

func paste_clipboard_to_canvas(canvas *CharCanvas, ascii string, startPos rl.Vector2) {
	text := strings.ReplaceAll(ascii, "```", "")

	text_width := TextWidth(text)
	text_height := TextHeight(text)

	lines := strings.Split(text, "\n")

	startPos.X -= float32(text_width / 2)
	startPos.Y -= float32(text_height / 2)

	for y, line := range lines {
		for x, ch := range line {
			if ch != ' ' && ch != '\r' {
				pos := rl.Vector2{
					X: startPos.X + float32(x),
					Y: startPos.Y + float32(y),
				}

				canvas.InsertChar(Char{
					text:     string(ch),
					position: pos,
				})
			}
		}
	}
}

func DrawHelpPanel(screenWidth, screenHeight int32, font rl.Font) {
	helpText := "COMMANDS:\n" +
		"Ctrl + S - saves the art\n" +
		"Left Click + drag - selects an area\n" +
		"Ctrl + X - cut selection\n" +
		"Ctrl + C - copy selection\n" +
		"Ctrl + V - pastes content from clipboard to the canvas\n" +
		"Ctrl + H - shows help screen\n" +
		"you can type keys to the canvas\n" +
		"Drag camera: hold SPACE + left mouse button\n"

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
	var drag_active bool = false

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

	var canvas CharCanvas = CharCanvas{}
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
		drag_active = rl.IsKeyDown(rl.KeySpace) && rl.IsMouseButtonDown(rl.MouseButtonLeft) || rl.IsMouseButtonDown(rl.MouseButtonMiddle)

		// keyboard inputs
		if ctrl_active && rl.IsKeyPressed(rl.KeyS) {
			// save draw
			text := canvas.ToText()
			ascii := TextToASCII(text)
			Save(ascii)
		} else if cursor.CurrentSelection.selecting {
			if ctrl_active && rl.IsKeyPressed(rl.KeyV) {
				// paste input
				clipboardText := string(clipboard.Read(clipboard.FmtText))
				paste_clipboard_to_canvas(&canvas, clipboardText, cursor.Position)
				cursor.CurrentSelection.selecting = false
			} else if ctrl_active && rl.IsKeyPressed(rl.KeyC) {
				// copy
				selected_text := canvas.TextFromSelection(cursor.CurrentSelection)

				tmpClipboard := selected_text
				cursor.CurrentSelection.selecting = false

				clipboard.Write(clipboard.FmtText, []byte(tmpClipboard))
			} else if ctrl_active && rl.IsKeyPressed(rl.KeyX) {
				selected_text := canvas.TextFromSelection(cursor.CurrentSelection)

				tmpClipboard := selected_text

				canvas.ClearSelection(cursor.CurrentSelection)
				cursor.CurrentSelection.selecting = false

				clipboard.Write(clipboard.FmtText, []byte(tmpClipboard))
			}
		} else if ctrl_active && rl.IsKeyPressed(rl.KeyZero) {
			canvas.chars = []Char{}
		} else if ctrl_active && rl.IsKeyPressed(rl.KeyH) {
			showHelp = !showHelp
		}

		if cursor.Active {
			char, typed := cursor.GetCharPressed()
			keypressed := cursor.GetInputPressed()

			if typed {
				canvas.InsertChar(char)
				cursor.Position.X += 1
			}

			if keypressed == rl.KeyBackspace {
				canvas.DeleteCharAt(cursor.Position)
				cursor.Position.X -= 1
			}

			if keypressed == rl.KeyEnter {
				cursor.Position.Y += 1
			}

			if keypressed != 0 {
				if keypressed != rl.KeyLeftControl && keypressed != rl.KeyLeftShift {
					cursor.CurrentSelection.selecting = false // desativa seleção ao digitar
				}
			}

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
		if drag_active {
			mouse_delta := rl.GetMouseDelta()
			camera.Target.X -= mouse_delta.X
			camera.Target.Y -= mouse_delta.Y
		}

		// drawing
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.BeginMode2D(camera)

		draw_canvas(canvas, current_font, font_size)

		if cursor.CurrentSelection.selecting {
			draw_selected_canvas(canvas, cursor.CurrentSelection, current_font, font_size)
		}

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
