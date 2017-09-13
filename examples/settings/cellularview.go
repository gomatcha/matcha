package settings

import (
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

type CellularView struct {
	view.Embed
	app *App
}

func NewCellularView(app *App) *CellularView {
	return &CellularView{
		Embed: view.NewEmbed(app),
		app:   app,
	}
}

func (v *CellularView) Build(ctx view.Context) view.Model {
	l := &table.Layouter{}
	chlds := []view.View{}

	scrollView := view.NewScrollView()
	scrollView.ContentLayouter = l
	scrollView.ContentChildren = chlds

	return view.Model{
		Children: []view.View{scrollView},
		Painter:  &paint.Style{BackgroundColor: backgroundColor},
	}
}
