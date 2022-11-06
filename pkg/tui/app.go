package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/kpcli/pkg/clipboard"
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
    key := f[0]
    value := f[1]
    // Do not print empty fields
    if value == "" {
      continue 
    }

    inputType := InputTypeText
    if key == kdbx.PASSWORD_KEY {
      inputType = InputTypePassword
    }

    fieldOptions := &FieldOptions{Label: key , InitialValue: value, InputType: inputType}
    field := NewField(fieldOptions)

    field.input.OnKeyPress(func(ev *tcell.EventKey) bool {
      if ev.Key() == tcell.KeyEnter {
        clipboard.Write(value)
      }
      return true
    })

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
