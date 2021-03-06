package slide

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gdamore/tcell"
	"github.com/kbaran224/go-xbps/xbps"
	"github.com/rivo/tview"
)

var grid = tview.NewGrid()

var pkgInfo = tview.NewTextView()
var pkgList *tview.List
var search = tview.NewInputField()
var ph = tview.NewTextView()

var t *template.Template

// Install searches and installs pkgs from a list
func Install(nextSlide func()) (title string, content tview.Primitive) {

	t, _ = template.New("").Parse(xbps.Tmpl)

	pkgInfo.SetWrap(false).
		SetDynamicColors(true)

	search.SetPlaceholder("Search...").
		SetDoneFunc(func(key tcell.Key) {
			pkgSearch(key)
		})

	initList(pkgList)

	ph.SetText(fmt.Sprintf("%d packages", pkgList.GetItemCount())).
		SetTextAlign(tview.AlignCenter)

	grid.SetRows(1, 0).
		SetColumns(30, 0).
		SetBorders(true)

	initGrid(grid)

	return "Install", grid

}

// init a list with all packages
func initList(list *tview.List) {
	pkgs, _ := xbps.Query("")
	pkgList = newList(pkgs)

}

// init (or reinit) whole grid
func initGrid(grid *tview.Grid) {
	grid.Clear()
	grid.AddItem(search, 0, 0, 1, 1, 0, 100, false)
	grid.AddItem(ph, 0, 1, 1, 1, 0, 100, false)

	grid.AddItem(pkgList, 1, 0, 1, 1, 0, 100, true)
	grid.AddItem(pkgInfo, 1, 1, 1, 1, 0, 100, false)
}

// search for pkg with input field
func pkgSearch(key tcell.Key) {

	switch key {
	case tcell.KeyDown:
		return
	case tcell.KeyEnter:
		query := search.GetText()

		ind := pkgList.FindItems(query, query, false, false)

		// show every pkg if none is searched
		if query == "" || query == "*" {
			initList(pkgList)
			ph.SetText(fmt.Sprintf("%d packages", pkgList.GetItemCount()))
			initGrid(grid)
			return
		}

		if len(ind) == 0 {
			ph.SetText(fmt.Sprintf("%d packages", pkgList.GetItemCount()))
			return
		}

		var foundPkgs []string
		for _, i := range ind {

			name, _ := pkgList.GetItemText(i)
			foundPkgs = append(foundPkgs, name)

		}

		// reinit whole grid with found pkgs
		newPkgList := newList(foundPkgs)
		pkgList = newPkgList

		initGrid(grid)
		pkgInfo.Clear()
		ph.SetText(fmt.Sprintf("%d packages", pkgList.GetItemCount()))

	}
}

// create a list for given pkgs
func newList(text []string) *tview.List {
	list := tview.NewList()

	for _, v := range text {
		f := newInfo(v)
		list.AddItem(v, "", rune('>'), f)
	}

	return list
}

// create pkg info for given pkg using template
func newInfo(pkgName string) func() {
	return func() {
		pkg, _ := xbps.Info(pkgName)
		var buf bytes.Buffer
		_ = t.Execute(&buf, pkg)
		pkgInfo.Clear()
		fmt.Fprint(pkgInfo, buf.String())
		pkgInfo.ScrollToBeginning()
	}
}
