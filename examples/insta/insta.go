package insta

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/app"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/button"
	"gomatcha.io/matcha/view/scrollview"
	"gomatcha.io/matcha/view/stackview"
	"gomatcha.io/matcha/view/textview"
	"gomatcha.io/matcha/view/urlimageview"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/insta New", func() *view.Root {
		app := NewApp()

		v := stackview.New()
		v.Stack = app.Stack
		v.Stack.SetViews(NewRootView(app))
		return view.NewRoot(v)
	})
}

type RootView struct {
	view.Embed
	app *App
}

func NewRootView(app *App) *RootView {
	return &RootView{
		Embed: view.NewEmbed(app),
		app:   app,
	}
}

func (v *RootView) Build(ctx *view.Context) view.Model {
	l := &table.Layouter{}

	for _, i := range v.app.Posts {
		postView := NewPostView(i)
		l.Add(postView, nil)
	}

	scrollView := scrollview.New()
	scrollView.ContentChildren = l.Views()
	scrollView.ContentLayouter = l

	return view.Model{
		Children: []view.View{scrollView},
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

// func (v *RootView) StackBar(ctx *view.Context) *stackview.Bar {
// 	return &stackview.Bar{Title: "Settings Example"}
// }

type PostView struct {
	view.Embed
	post *Post
}

func NewPostView(p *Post) *PostView {
	return &PostView{
		Embed: view.NewEmbed(p),
		post:  p,
	}
}

func (v *PostView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	header := NewPostHeaderView()
	header.Title = v.post.UserName
	header.ImageURL = v.post.UserImageURL
	headerGuide := l.Add(header, func(s *constraint.Solver) {
		s.TopEqual(l.Top())
		s.LeftEqual(l.Left())
		s.RightEqual(l.Right())
	})

	image := urlimageview.New()
	image.URL = v.post.ImageURL
	imageGuide := l.Add(image, func(s *constraint.Solver) {
		s.TopEqual(headerGuide.Bottom())
		s.LeftEqual(l.Left())
		s.RightEqual(l.Right())
		s.HeightEqual(l.Width())
	})

	buttons := NewPostButtonsView()
	buttonsGuide := l.Add(buttons, func(s *constraint.Solver) {
		s.TopEqual(imageGuide.Bottom())
		s.LeftEqual(l.Left())
		s.RightEqual(l.Right())
	})

	l.Solve(func(s *constraint.Solver) {
		s.Top(0)
		s.BottomEqual(buttonsGuide.Bottom())
		s.WidthEqual(l.MinGuide().Width())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Blue},
	}
}

type PostHeaderView struct {
	view.Embed
	Title    string
	ImageURL string
}

func NewPostHeaderView() *PostHeaderView {
	return &PostHeaderView{}
}

func (v *PostHeaderView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(55)
	})

	imageView := urlimageview.New()
	imageView.PaintStyle = &paint.Style{CornerRadius: 16, BackgroundColor: colornames.Gray}
	imageView.URL = v.ImageURL
	g := l.Add(imageView, func(s *constraint.Solver) {
		s.Width(32)
		s.Height(32)
		s.Left(10)
		s.CenterYEqual(l.CenterY())
	})

	titleView := textview.New()
	titleView.MaxLines = 1
	titleView.String = v.Title
	l.Add(titleView, func(s *constraint.Solver) {
		s.LeftEqual(g.Right().Add(10))
		s.CenterYEqual(l.CenterY())
		s.RightEqual(l.Right().Add(-10))
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

type PostButtonsView struct {
	view.Embed
	Liked           bool
	Bookmarked      bool
	OnTouchLike     func(bool)
	OnTouchComment  func()
	OnTouchShare    func()
	OnTouchBookmark func(bool)
}

func NewPostButtonsView() *PostButtonsView {
	return &PostButtonsView{}
}

func (v *PostButtonsView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(50)
	})

	likeButton := button.NewImageButton()
	likeButton.Image = app.MustLoadImage("Like")
	likeButton.OnPress = func() {
		if v.OnTouchLike != nil {
			v.Liked = !v.Liked
			v.OnTouchLike(v.Liked)
			v.Signal()
		}
	}
	likeGuide := l.Add(likeButton, func(s *constraint.Solver) {
		s.CenterYEqual(l.CenterY())
		s.Left(13)
	})

	commentButton := button.NewImageButton()
	commentButton.Image = app.MustLoadImage("Comment")
	commentButton.OnPress = func() {
		if v.OnTouchComment != nil {
			v.OnTouchComment()
		}
	}
	commentGuide := l.Add(commentButton, func(s *constraint.Solver) {
		s.CenterYEqual(l.CenterY())
		s.LeftEqual(likeGuide.Right().Add(13))
	})

	shareButton := button.NewImageButton()
	shareButton.Image = app.MustLoadImage("Share")
	shareButton.OnPress = func() {
		if v.OnTouchShare != nil {
			v.OnTouchShare()
		}
	}
	l.Add(shareButton, func(s *constraint.Solver) {
		s.CenterYEqual(l.CenterY())
		s.LeftEqual(commentGuide.Right().Add(13))
	})

	bookmarkButton := button.NewImageButton()
	bookmarkButton.Image = app.MustLoadImage("Bookmark")
	bookmarkButton.OnPress = func() {
		if v.OnTouchBookmark != nil {
			v.Bookmarked = !v.Bookmarked
			v.OnTouchBookmark(v.Bookmarked)
			v.Signal()
		}
	}
	l.Add(bookmarkButton, func(s *constraint.Solver) {
		s.CenterYEqual(l.CenterY())
		s.RightEqual(l.Right().Add(-13))
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}
