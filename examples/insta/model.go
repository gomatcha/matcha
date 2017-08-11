package insta

import (
	"math/rand"

	golorem "github.com/drhodes/golorem"
	"gomatcha.io/matcha/view/stackview"
)

type App struct {
	Stack *stackview.Stack
	Posts []*Post
}

func NewApp() *App {
	posts := []*Post{}
	for i := 0; i < 5; i++ {
		posts = append(posts, GeneratePost())
	}

	return &App{
		Stack: &stackview.Stack{},
		Posts: posts,
	}
}

type Post struct {
	Id int64
	// UserId       int64
	UserName     string
	UserImageURL string
	Location     string
	ImageURL     string
	LikeCount    int
	// CommentIds   []int64
	Comments []Comment
}

func GeneratePost() *Post {
	return &Post{
		Id:           0,
		UserName:     "User name",
		UserImageURL: "https://unsplash.it/50/50/?random",
		Location:     golorem.Word(5, 15),
		ImageURL:     "https://unsplash.it/200/200/?random",
		LikeCount:    rand.Intn(500),
		Comments:     nil,
	}
}

type Comment struct {
	Id       int64
	UserId   int64
	UserName string
	Text     string
}

type User struct {
	Id       int64
	Name     string
	ImageURL string
}
