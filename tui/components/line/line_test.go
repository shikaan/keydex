package line

import (
	"reflect"
	"testing"
)

func TestNewPaddedLine(t *testing.T) {
	tests := []struct {
		name string
		line string
		want PaddedLine
	}{
		{"bytes only", "lol", PaddedLine{'l', 'o', 'l'}},
		{"non-ascii", "✅🤖", PaddedLine{'✅', PAD_BYTE, '🤖', PAD_BYTE}},
		{"empty", "", PaddedLine{}},
		{"mixed", "✅I", PaddedLine{'✅', PAD_BYTE, 'I'}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPaddedLine(tt.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaddedLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRune(t *testing.T) {
	tests := []struct {
		name         string
		line         PaddedLine
		x            int
		wantChar     rune
		wantPosition int
	}{
		{"byte", NewPaddedLine("test"), 1, 'e', 1},
		{"non-ascii at the beginning", NewPaddedLine("I🤖"), 1, '🤖', 1},
		{"non-ascii on pad byte", NewPaddedLine("I🤖"), 2, '🤖', 1},
		{"out of bounds (too big)", NewPaddedLine("lol"), 4, EMPTY_CELL, -1},
		{"out of bounds (too small)", NewPaddedLine("lol"), -1, EMPTY_CELL, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChar, gotPosition := GetRune(tt.line, tt.x)
			if gotChar != tt.wantChar {
				t.Errorf("GetRune() got char = %v, want %v", gotChar, tt.wantChar)
			}
			if gotPosition != tt.wantPosition {
				t.Errorf("GetRune() got position = %v, want %v", gotPosition, tt.wantPosition)
			}
		})
	}
}
