package tui

type Focusable interface {
  SetFocus(on bool)
  HasFocus() bool
}

