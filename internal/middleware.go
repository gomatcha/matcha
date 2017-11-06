package internal

import (
	"sync"

	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
)

var middlewaresMu sync.Mutex
var middlewares = []func(m *MiddlewareRoot) interface{}{}

// RegisterMiddleware adds v to the list of default middleware that Root starts with.
func RegisterMiddleware(v func(m *MiddlewareRoot) interface{}) {
	middlewaresMu.Lock()
	defer middlewaresMu.Unlock()

	middlewares = append(middlewares, v)
}

func Middlewares() []func(m *MiddlewareRoot) interface{} {
	middlewaresMu.Lock()
	defer middlewaresMu.Unlock()

	return middlewares
}

type MiddlewareRoot struct {
	LayoutGuideFunc func(path int64) layout.Guide
	PaintStyleFunc  func(path int64) paint.Style
}

// func (m *MiddlewareRoot) viewModel(path []Id) view.Model {
// }

func (m *MiddlewareRoot) LayoutGuide(path int64) layout.Guide {
	return m.LayoutGuideFunc(path)
}

func (m *MiddlewareRoot) PaintStyle(path int64) paint.Style {
	return m.PaintStyleFunc(path)
}
