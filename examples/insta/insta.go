package insta

import (
	"fmt"
	"image/color"
	"time"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/app"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/touch"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/ios/stackview"
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

	scrollView := view.NewScrollView()
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

	image := NewPostImageView()
	image.ImageURL = v.post.ImageURL
	image.OnDoubleTouch = func() {
		v.post.Liked = true
		v.Signal()
	}
	imageGuide := l.Add(image, func(s *constraint.Solver) {
		s.TopEqual(headerGuide.Bottom())
		s.LeftEqual(l.Left())
		s.RightEqual(l.Right())
		s.HeightEqual(l.Width())
	})

	buttons := NewPostButtonsView()
	buttons.Liked = v.post.Liked
	buttons.Bookmarked = v.post.Bookmarked
	buttons.LikeCount = v.post.LikeCount
	buttons.OnTouchLike = func(a bool) {
		fmt.Println("Like", a)
		v.post.Liked = a
		v.Signal()
	}
	buttons.OnTouchComment = func() {
		fmt.Println("Comment")
	}
	buttons.OnTouchShare = func() {
		fmt.Println("Share")
	}
	buttons.OnTouchBookmark = func(a bool) {
		fmt.Println("Bookmark", a)
		v.post.Bookmarked = a
		v.Signal()
	}
	buttonsGuide := l.Add(buttons, func(s *constraint.Solver) {
		s.TopEqual(imageGuide.Bottom())
		s.LeftEqual(l.Left())
		s.RightEqual(l.Right())
	})

	comments := NewCommentsView()
	comments.Comments = v.post.Comments
	commentsGuide := l.Add(comments, func(s *constraint.Solver) {
		s.TopEqual(buttonsGuide.Bottom())
		s.LeftEqual(l.Left())
		s.RightEqual(l.Right())
	})

	l.Solve(func(s *constraint.Solver) {
		s.Top(0)
		s.BottomEqual(commentsGuide.Bottom())
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
		s.Height(60)
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

	titleView := view.NewTextView()
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

type PostImageView struct {
	view.Embed
	ImageURL      string
	OnDoubleTouch func()
	showHeart     bool
}

func NewPostImageView() *PostImageView {
	return &PostImageView{}
}

func (v *PostImageView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	image := urlimageview.New()
	image.URL = v.ImageURL
	l.Add(image, func(s *constraint.Solver) {
		s.WidthEqual(l.Width())
		s.HeightEqual(l.Height())
	})

	if v.showHeart {
		heart := view.NewImageView()
		heart.Image = app.MustLoadImage("Heart")
		heart.ResizeMode = view.ImageResizeModeCenter
		heart.PaintStyle = &paint.Style{
			ShadowRadius: 10,
			ShadowOffset: layout.Pt(0, 3),
			ShadowColor:  color.RGBA{0, 0, 0, 128},
		}
		l.Add(heart, func(s *constraint.Solver) {
			s.CenterXEqual(l.CenterX())
			s.CenterYEqual(l.CenterY())
		})
	}

	tap := &touch.TapRecognizer{
		Count: 2,
		OnTouch: func(e *touch.TapEvent) {
			v.showHeart = true
			v.Signal()
			time.AfterFunc(time.Second, func() {
				v.showHeart = false
				v.Signal()
			})

			if v.OnDoubleTouch != nil {
				v.OnDoubleTouch()
			}
		},
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Options: []view.Option{
			touch.RecognizerList{tap},
		},
	}
}

type PostButtonsView struct {
	view.Embed
	LikeCount       int
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
		s.Height(65)
	})

	likeButton := view.NewImageButton()
	if v.Liked {
		likeButton.Image = app.MustLoadImage("LikeFilled")
	} else {
		likeButton.Image = app.MustLoadImage("Like")
	}
	likeButton.OnPress = func() {
		if v.OnTouchLike != nil {
			v.Liked = !v.Liked
			v.OnTouchLike(v.Liked)
			v.Signal()
		}
	}
	likeGuide := l.Add(likeButton, func(s *constraint.Solver) {
		s.Top(13)
		s.Left(13)
	})

	commentButton := view.NewImageButton()
	commentButton.Image = app.MustLoadImage("Comment")
	commentButton.OnPress = func() {
		if v.OnTouchComment != nil {
			v.OnTouchComment()
		}
	}
	commentGuide := l.Add(commentButton, func(s *constraint.Solver) {
		s.Top(13)
		s.LeftEqual(likeGuide.Right().Add(13))
	})

	shareButton := view.NewImageButton()
	shareButton.Image = app.MustLoadImage("Share")
	shareButton.OnPress = func() {
		if v.OnTouchShare != nil {
			v.OnTouchShare()
		}
	}
	l.Add(shareButton, func(s *constraint.Solver) {
		s.Top(13)
		s.LeftEqual(commentGuide.Right().Add(13))
	})

	bookmarkButton := view.NewImageButton()
	if v.Bookmarked {
		bookmarkButton.Image = app.MustLoadImage("BookmarkFilled")
	} else {
		bookmarkButton.Image = app.MustLoadImage("Bookmark")
	}

	bookmarkButton.OnPress = func() {
		if v.OnTouchBookmark != nil {
			v.Bookmarked = !v.Bookmarked
			v.OnTouchBookmark(v.Bookmarked)
			v.Signal()
		}
	}
	l.Add(bookmarkButton, func(s *constraint.Solver) {
		s.Top(13)
		s.RightEqual(l.Right().Add(-13))
	})

	likeTextView := view.NewTextView()
	likeTextView.String = fmt.Sprintf("%v Likes", v.LikeCount)
	l.Add(likeTextView, func(s *constraint.Solver) {
		s.Top(50)
		s.LeftEqual(l.Left().Add(13))
		s.RightEqual(l.Right().Add(-13))
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

type CommentsView struct {
	view.Embed
	Comments []*Comment
}

func NewCommentsView() *CommentsView {
	return &CommentsView{}
}

func (v *CommentsView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	topGuide := l.Top().Add(13)

	for _, i := range v.Comments {
		textView := view.NewTextView()
		textView.String = i.UserName + i.Text
		textGuide := l.Add(textView, func(s *constraint.Solver) {
			s.TopEqual(topGuide)
			s.LeftEqual(l.Left().Add(13))
			s.RightEqual(l.Right().Add(-13))
		})

		topGuide = textGuide.Top()
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
	}
}
