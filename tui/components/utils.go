package components

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// Clamps an integer between minVale and maxValue (both included)
func clamp(n, minValue, maxValue int) int {
	return max(min(n, maxValue), minValue)
}

// Breaks a text in lines in a platform independent way
func GetLines(text string) []string {
	lines := []string{}

	for _, line := range strings.Split(text, "\n") {
		lines = append(lines, strings.TrimSuffix(line, "\r"))
	}

	return lines
}

type PaddedLine = []rune

const PAD_BYTE = 0

func NewPaddedLine(line string) PaddedLine {
	cells := PaddedLine{}
	for _, char := range line {
		width := runewidth.RuneWidth(char)

		cells = append(cells, char)
		// Pad rune with PAD_BYTE, in case the rune is longer than one cell
		for i := 1; i < width; i++ {
			cells = append(cells, PAD_BYTE)
		}
	}
	return cells
}

// Stably retuns a rune and its position from a PaddedLine,
// skipping all the PAD_BYTEs. Returns EMPTY_CELL and -1 for invalid positions
func GetRune(line PaddedLine, x int) (char rune, position int) {
	if x >= len(line) {
		return EMPTY_CELL, -1
	}

	for j := x; j >= 0; j-- {
		if line[j] != PAD_BYTE {
			return line[j], j
		}
	}

	return EMPTY_CELL, -1
}
