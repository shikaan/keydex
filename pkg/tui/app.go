package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/kpcli/pkg/kdbx"
)

var App = &views.Application{}
var Screen, _ = tcell.NewScreen()

type root struct {
  views.Panel
}

func (r *root) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape || ev.Name() == "Ctrl+C"{
			App.Quit()
			return true
		}
  }
    return r.Panel.HandleEvent(ev)
}

func OpenEntryEditor(e kdbx.Entry) error {
  r := &root{}
  r.SetStyle(tcell.StyleDefault)

  title := NewTitle(e.Name)
  r.SetTitle(title)

  main := NewForm()
  hasSetFocus := false

  for _, f := range e.Fields {
    // Do not print empty fields
    if f[1] == "" {
      continue 
    }

    inputType := InputTypeText
    if f[0] == kdbx.PASSWORD_KEY {
      inputType = InputTypePassword
    }

    fieldOptions := &FieldOptions{Label: f[0] , InitialValue: f[1], InputType: inputType}
    field := NewField(fieldOptions)
    main.AddWidget(field, 0)
  
    if !hasSetFocus {
      field.SetFocus(true)
      hasSetFocus = true
    }
  }

  flex := NewResponsive(views.Horizontal)
  flex.SetContent(main)
  r.SetContent(flex)

  status := NewStatus()
  r.SetStatus(status)

  App.SetRootWidget(r)
  App.SetScreen(Screen)

  return App.Run()
}
