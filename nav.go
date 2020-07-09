package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/rivo/tview"
)

var pkgInfo *tview.TextView
var t *template.Template

func main() {

	t, _ = template.New("").Parse(tmpl)

	app := tview.NewApplication()

	pkgs, _ := query("")
	menu := newList(pkgs)

	pkgInfo = tview.NewTextView().
		SetWrap(false).
		SetDynamicColors(true)

	search := newPrimitive("Search", true)
	title := newPrimitive("Title", true)

	grid := newGrid()
	grid.AddItem(search, 0, 0, 1, 1, 0, 100, false)
	grid.AddItem(title, 0, 1, 1, 1, 0, 100, false)

	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, true)
	grid.AddItem(pkgInfo, 1, 1, 1, 1, 0, 100, false)

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}

}

func newGrid() *tview.Grid {
	return tview.NewGrid().
		SetRows(1, 0).
		SetColumns(30, 0).
		SetBorders(true)
}

func newPrimitive(text string, align bool) tview.Primitive {
	prim := tview.NewTextView().
		SetText(text)

	if align {
		prim.SetTextAlign(tview.AlignCenter)

	}
	return prim
}

func newList(text []string) *tview.List {
	list := tview.NewList()

	for _, v := range text {
		f := newInfo(v)
		list.AddItem(v, "", rune('>'), f)
	}

	return list
}

func newInfo(pkgName string) func() {
	return func() {
		pkg, _ := info(pkgName)
		var buf bytes.Buffer
		_ = t.Execute(&buf, pkg)
		pkgInfo.Clear()
		fmt.Fprint(pkgInfo, buf.String())
		pkgInfo.ScrollToBeginning()
	}
}
