package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var grid *tview.Grid

var pkgInfo *tview.TextView
var pkgList *tview.List
var search *tview.InputField

var t *template.Template
var title *tview.TextView

func main() {

	t, _ = template.New("").Parse(tmpl)

	app := tview.NewApplication().
		SetRoot(grid, true).
		EnableMouse(true)

	pkgInfo = tview.NewTextView().
		SetWrap(false).
		SetDynamicColors(true)

	search = tview.NewInputField().
		SetPlaceholder("Search...").
		SetDoneFunc(func(key tcell.Key) {
			pkgSearch()
		})

	title = tview.NewTextView().
		SetText("Package explorer").
		SetTextAlign(tview.AlignCenter)

	initList(pkgList)
	grid = newGrid()
	initGrid(grid)

	if err := app.Run(); err != nil {
		panic(err)
	}

}

func initList(list *tview.List) {
	pkgs, _ := query("")
	pkgList = newList(pkgs)

}

func initGrid(grid *tview.Grid) {
	grid.Clear()
	grid.AddItem(search, 0, 0, 1, 1, 0, 100, false)
	grid.AddItem(title, 0, 1, 1, 1, 0, 100, false)

	grid.AddItem(pkgList, 1, 0, 1, 1, 0, 100, true)
	grid.AddItem(pkgInfo, 1, 1, 1, 1, 0, 100, false)
}

func pkgSearch() {
	query := search.GetText()

	ind := pkgList.FindItems(query, query, false, false)

	if query == "" || query == "*" {
		initList(pkgList)
		return
	}

	if len(ind) == 0 {
		return
	}

	var foundPkgs []string
	for _, i := range ind {

		name, _ := pkgList.GetItemText(i)
		foundPkgs = append(foundPkgs, name)

	}
	newPkgList := newList(foundPkgs)
	pkgList = newPkgList

	initGrid(grid)
	pkgInfo.Clear()
	fmt.Fprint(pkgInfo, "Found ", newPkgList.GetItemCount())

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
