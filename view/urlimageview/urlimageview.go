// Package urlimageview implements a view which loads and displays an image.
package urlimageview

import (
	"context"
	"image"
	"image/color"
	"net/http"

	"gomatcha.io/matcha"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

// Layouter that returns the child's layout
type layouter struct{}

func (l layouter) Layout(ctx *layout.Context) (layout.Guide, []layout.Guide) {
	g := layout.Guide{Frame: layout.Rect{Max: ctx.MaxSize}}
	gs := []layout.Guide{}
	for i := 0; i < ctx.ChildCount; i++ {
		f := ctx.LayoutChild(i, ctx.MinSize, ctx.MaxSize)
		gs = append(gs, f)
		g.Frame = f.Frame
	}
	return g, gs
}

func (l layouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

func (l layouter) Unnotify(id comm.Id) {
	// no-op
}

type View struct {
	view.Embed
	PaintStyle *paint.Style
	ResizeMode view.ImageResizeMode
	URL        string
	Path       string
	ImageTint  color.Color
	stage      view.Stage
	// Image request
	cancelFunc context.CancelFunc
	image      image.Image
	err        error
}

// New returns either the previous View in ctx with matching key, or a new View if none exists.
func New() *View {
	return &View{}
}

// Build implements view.View.
func (v *View) Build(ctx *view.Context) view.Model {
	v.reload()

	chl := view.NewImageView()
	chl.ResizeMode = v.ResizeMode
	chl.Image = v.image
	chl.ImageTint = v.ImageTint

	var painter paint.Painter
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return view.Model{
		Painter:  painter,
		Layouter: layouter{},
		Children: []view.View{chl},
	}
}

func (v *View) Update(v2 view.View) {
	if v2.(*View).URL != v.URL {
		v.cancel()
		v.image = nil
		v.err = nil
	}
	view.CopyFields(v, v2)
}

// Lifecycle implements view.View.
func (v *View) Lifecycle(from, to view.Stage) {
	v.stage = to
	v.reload()
}

func (v *View) reload() {
	if v.stage < view.StageMounted {
		v.cancel()
		return
	}

	if v.URL == "" || v.cancelFunc != nil || v.image != nil || v.err != nil {
		return
	}

	c, cancelFunc := context.WithCancel(context.Background())
	v.cancelFunc = cancelFunc
	go func(url string) {
		image, err := loadImageURL(url)

		matcha.MainLocker.Lock()
		defer matcha.MainLocker.Unlock()

		select {
		case <-c.Done():
		default:
			v.cancelFunc()
			v.cancelFunc = nil
			v.image = image
			v.err = err
			v.Signal()
		}
	}(v.URL)
}

func (v *View) cancel() {
	if v.cancelFunc != nil {
		v.cancelFunc()
		v.cancelFunc = nil
	}
}

func loadImageURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	return img, err
}
