package form

type Focusable interface {
  SetFocus(on bool)
  HasFocus() bool
}

