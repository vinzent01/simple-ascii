package main

import "strings"

// Extrai o width (maior comprimento de linha entre minX e maxX com caracteres não espaço)
func TextWidth(ascii string) int {
	lines := strings.Split(ascii, "\n")
	minX, maxX := -1, -1
	for _, line := range lines {
		for x, ch := range line {
			if ch != ' ' && ch != '\r' {
				if minX == -1 || x < minX {
					minX = x
				}
				if maxX == -1 || x > maxX {
					maxX = x
				}
			}
		}
	}
	if minX == -1 || maxX == -1 {
		return 0
	}
	return (maxX - minX) + 1
}

// Extrai o height (número de linhas entre minY e maxY com pelo menos um caractere não espaço)
func TextHeight(ascii string) int {
	lines := strings.Split(ascii, "\n")
	minY, maxY := -1, -1
	for y, line := range lines {
		hasChar := false
		for _, ch := range line {
			if ch != ' ' && ch != '\r' {
				hasChar = true
				break
			}
		}
		if hasChar {
			if minY == -1 || y < minY {
				minY = y
			}
			if maxY == -1 || y > maxY {
				maxY = y
			}
		}
	}
	if minY == -1 || maxY == -1 {
		return 0
	}
	return (maxY - minY) + 1
}

func TextToASCII(text string) string {
	asciiStr := "```\n"

	asciiStr += text

	asciiStr += "```"

	return asciiStr
}
