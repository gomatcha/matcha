package settings

import (
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/scrollview"
)

type CellularView struct {
	view.Embed
	app *App
}

func NewCellularView(ctx *view.Context, key string, app *App) *CellularView {
	if v, ok := ctx.Prev(key).(*CellularView); ok {
		return v
	}
	return &CellularView{Embed: ctx.NewEmbed(key), app: app}
}

func (v *CellularView) Build(ctx *view.Context) view.Model {
	l := &table.Layouter{}
	chlds := []view.View{}

	scrollView := scrollview.New(ctx, "b")
	scrollView.ContentLayouter = l
	scrollView.ContentChildren = chlds

	return view.Model{
		Children: []view.View{scrollView},
		Painter:  &paint.Style{BackgroundColor: backgroundColor},
	}
}
