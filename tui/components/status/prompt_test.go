package status

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestPrompt_HandleEvent(t *testing.T) {
	tests := []struct {
		name    string
		model   promptModel
		ev      tcell.Event
		handled bool
	}{
		{
			name: "doesn't handle event without focus",
			model: promptModel{
				hasFocus: false,
			},
			ev:      &tcell.EventKey{},
			handled: false,
		},
		{
			name: "doesn't handle event without keypressHandler",
			model: promptModel{
				hasFocus: true,
			},
			ev:      &tcell.EventKey{},
			handled: false,
		},
		{
			name: "triggers keyPressHandler events (handled)",
			model: promptModel{
				hasFocus: true,
				keyPressHandler: func(ev *tcell.EventKey) bool {
					return true
				},
			},
			ev:      &tcell.EventKey{},
			handled: true,
		},
		{
			name: "triggers keyPressHandler events (unhandled)",
			model: promptModel{
				hasFocus: true,
				keyPressHandler: func(ev *tcell.EventKey) bool {
					return false
				},
			},
			ev:      &tcell.EventKey{},
			handled: false,
		},
		{
			name: "doesn't handle other events",
			model: promptModel{
				hasFocus: true,
				keyPressHandler: func(ev *tcell.EventKey) bool {
					return true
				},
			},
			ev:      &tcell.EventMouse{},
			handled: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Prompt{model: tt.model}
			if got := p.HandleEvent(tt.ev); got != tt.handled {
				t.Errorf("Prompt.HandleEvent() = %v, want %v", got, tt.handled)
			}
		})
	}
}
