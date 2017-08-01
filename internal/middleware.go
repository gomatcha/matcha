package internal

import "sync"

var middlewaresMu sync.Mutex
var middlewares = []func() interface{}{}

// RegisterMiddleware adds v to the list of default middleware that Root starts with.
func RegisterMiddleware(v func() interface{}) {
	middlewaresMu.Lock()
	defer middlewaresMu.Unlock()

	middlewares = append(middlewares, v)
}

func Middlewares() []func() interface{} {
	middlewaresMu.Lock()
	defer middlewaresMu.Unlock()

	return middlewares
}
