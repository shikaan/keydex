// package form
// 
// import (
// 	"github.com/gdamore/tcell/v2"
// 	"github.com/gdamore/tcell/v2/views"
// )
// 
// type Input struct {
//   textArea views.TextArea
//   label views.SimpleStyledText
// 
//   views.BoxLayout
// }
// 
// func (i *Input) HandleEvent(ev tcell.Event) bool {
//   switch ev := ev.(type) {	
//   case *tcell.EventKey:
// 		if ev.Key() == tcell.KeyRune {
// 		  i.ReplaceCharAtCursor(ev.Rune())	
// 			return true
// 		}
// 	}
// 	return i.BoxLayout.HandleEvent(ev)
// }
// 
// func (i *Input) ReplaceCharAtCursor(char rune) {
//   m := i.textArea.GetModel() 
//   x, y, on, visible := m.GetCursor()
// 
//   if on && visible {
//     println(m.runes[0])
//   }
// }
// 
// func NewInput() *Input {
//   input := &Input{ }
//   input.SetOrientation(1)
// 
//   ta := views.NewTextArea()
//   ta.Init()
//   ta.SetContent("lol")
//   ta.EnableCursor(true)
// 
//   t := views.NewSimpleStyledText()
//   t.SetText("label")
// 
//   input.AddWidget(t, 0)
//   input.AddWidget(ta, 2)
// 
//   input.textArea = *ta
//   input.label = *t
// 
//   return input
// }
