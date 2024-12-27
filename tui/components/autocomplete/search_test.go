package autocomplete

import (
	"reflect"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/shikaan/keydex/tui/components/line"
)

func Test_searchModel_GetCell(t *testing.T) {
	type fields struct {
		content       string
		runes         line.PaddedLine
		width         int
		x             int
		style         tcell.Style
		hasFocus      bool
		changeHandler func(ev tcell.Event) bool
		focusHandler  func() bool
	}
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   rune
		want1  tcell.Style
		want2  []rune
		want3  int
	}{
		{
			name: "Valid printable character",
			fields: fields{
				content: "hello",
				runes:   line.NewPaddedLine("hello"),
				style:   tcell.StyleDefault,
			},
			args:  args{x: 1, y: 0},
			want:  'e',
			want1: tcell.StyleDefault,
			want2: nil,
			want3: 1,
		},
		{
			name: "Valid printable non-ascii character",
			fields: fields{
				content: "hello",
				runes:   line.NewPaddedLine("h🤖llo"),
				style:   tcell.StyleDefault,
			},
			args:  args{x: 1, y: 0},
			want:  '🤖',
			want1: tcell.StyleDefault,
			want2: nil,
			want3: 2,
		},
		{
			name: "Out of bounds x",
			fields: fields{
				content: "hello",
				runes:   line.NewPaddedLine("hello"),
				style:   tcell.StyleDefault,
			},
			args:  args{x: 10, y: 0},
			want:  line.EMPTY_CELL,
			want1: tcell.StyleDefault,
			want2: nil,
			want3: 1,
		},
		{
			name: "Out of bounds y",
			fields: fields{
				content: "hello",
				runes:   line.NewPaddedLine("hello"),
				style:   tcell.StyleDefault,
			},
			args:  args{x: 1, y: 1},
			want:  line.EMPTY_CELL,
			want1: tcell.StyleDefault,
			want2: nil,
			want3: 1,
		},
		{
			name: "Non-printable character",
			fields: fields{
				content: "\x00",
				runes:   line.NewPaddedLine("\x00"),
				style:   tcell.StyleDefault,
			},
			args:  args{x: 0, y: 0},
			want:  line.EMPTY_CELL,
			want1: tcell.StyleDefault,
			want2: nil,
			want3: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &searchModel{
				content:       tt.fields.content,
				runes:         tt.fields.runes,
				width:         tt.fields.width,
				x:             tt.fields.x,
				style:         tt.fields.style,
				hasFocus:      tt.fields.hasFocus,
				changeHandler: tt.fields.changeHandler,
				focusHandler:  tt.fields.focusHandler,
			}
			got, got1, got2, got3 := m.GetCell(tt.args.x, tt.args.y)
			if got != tt.want {
				t.Errorf("searchModel.GetCell() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("searchModel.GetCell() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("searchModel.GetCell() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("searchModel.GetCell() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}
func TestSearch_HandleEvent(t *testing.T) {
	type fields struct {
		model *searchModel
	}
	tests := []struct {
		name        string
		fields      fields
		ev          tcell.Event
		wantHandled bool
		wantContent string
	}{
		{
			name: "handles rune key event",
			fields: fields{
				model: &searchModel{
					content:  "hello",
					runes:    line.NewPaddedLine("hello"),
					x:        1,
					style:    tcell.StyleDefault,
					hasFocus: true,
				},
			},
			ev:          tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
			wantHandled: true,
			wantContent: "haello",
		},
		{
			name: "handles backspace key event",
			fields: fields{
				model: &searchModel{
					content:  "hello",
					runes:    line.NewPaddedLine("hello"),
					x:        1,
					style:    tcell.StyleDefault,
					hasFocus: true,
				},
			},
			ev:          tcell.NewEventKey(tcell.KeyBackspace, 0, tcell.ModNone),
			wantHandled: true,
			wantContent: "ello",
		},
		{
			name: "doesn't handle event without focus",
			fields: fields{
				model: &searchModel{
					content:  "hello",
					runes:    line.NewPaddedLine("hello"),
					x:        1,
					style:    tcell.StyleDefault,
					hasFocus: false,
				},
			},
			ev:          tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
			wantHandled: false,
			wantContent: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Search{
				model: tt.fields.model,
			}
			if got := s.HandleEvent(tt.ev); got != tt.wantHandled {
				t.Errorf("Search.HandleEvent() = %v, want %v", got, tt.wantHandled)
			}
			if got := s.GetContent(); got != tt.wantContent {
				t.Errorf("Search.GetContent() = %v, want %v", got, tt.wantContent)
			}
		})
	}
}
