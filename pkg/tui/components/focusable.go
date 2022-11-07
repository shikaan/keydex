package components

type Focusable interface {
  SetFocus(on bool)
  HasFocus() bool
}

