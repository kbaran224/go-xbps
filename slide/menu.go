package slide

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func Menu(nextSlide func()) (title string, content tview.Primitive) {

	logoBox := tview.NewTextView().
		SetDoneFunc(func(key tcell.Key) {
			nextSlide()
		})
	fmt.Fprint(logoBox, "In development!")

	return "Menu", logoBox
}
