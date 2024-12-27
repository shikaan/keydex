package line

import (
	"github.com/mattn/go-runewidth"
)

type PaddedLine = []rune

const PAD_BYTE = 0
const EMPTY_CELL = 0

// Returns a slice of runes where each item represents a logical position
// in the line. Runes whose length is more than 1, will be padded with PAD_BYTES
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

// Stably retuns a rune and its position from a PaddedLine, skipping all the
// PAD_BYTEs. Returns EMPTY_CELL and -1 for invalid positions
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
