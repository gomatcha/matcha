package view

import (
	"sync"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
)

type View interface {
	Build(*Context) Model
	Lifecycle(from, to Stage)
	Id() matcha.Id
	comm.Notifier
}

type Option interface {
	OptionsKey() string
}

// Embed is a convenience struct that provides a default implementation of View. It also wraps a comm.Relay.
type Embed struct {
	mu    sync.Mutex
	id    matcha.Id
	relay comm.Relay
}

// NewEmbed creates a new Embed with the given Id.
func NewEmbed(id matcha.Id) Embed {
	return Embed{id: id}
}

// Build is an empty implementation of View's Build method.
func (e *Embed) Build(ctx *Context) Model {
	return Model{}
}

// Id returns the id passed into NewEmbed
func (e *Embed) Id() matcha.Id {
	return e.id
}

// Lifecycle is an empty implementation of View's Lifecycle method.
func (e *Embed) Lifecycle(from, to Stage) {
	// no-op
}

// Notify calls Notify(id) on the underlying comm.Relay.
func (e *Embed) Notify(f func()) comm.Id {
	return e.relay.Notify(f)
}

// Unnotify calls Unnotify(id) on the underlying comm.Relay.
func (e *Embed) Unnotify(id comm.Id) {
	e.relay.Unnotify(id)
}

// Subscribe calls Subscribe(n) on the underlying comm.Relay.
func (e *Embed) Subscribe(n comm.Notifier) {
	e.relay.Subscribe(n)
}

// Unsubscribe calls Unsubscribe(n) on the underlying comm.Relay.
func (e *Embed) Unsubscribe(n comm.Notifier) {
	e.relay.Unsubscribe(n)
}

// Update calls Signal() on the underlying comm.Relay.
func (e *Embed) Signal() {
	e.relay.Signal()
}

type Stage int

const (
	// StageDead marks views that are not attached to the view hierarchy.
	StageDead Stage = iota
	// StageMounted marks views that are in the view hierarchy but not visible.
	StageMounted
	// StageVisible marks views that are in the view hierarchy and visible.
	StageVisible
)

// EntersStage returns true if from<s and to>=s.
func EntersStage(from, to, s Stage) bool {
	return from < s && to >= s
}

// ExitsStage returns true if from>=s and to<s.
func ExitsStage(from, to, s Stage) bool {
	return from >= s && to < s
}

// Model describes the view and its children.
type Model struct {
	Children []View
	Layouter layout.Layouter
	Painter  paint.Painter
	Options  []Option

	NativeViewName  string
	NativeViewState proto.Message
	NativeValues    map[string]proto.Message
	NativeFuncs     map[string]interface{}
}

// WithPainter wraps the view v, and replaces its Model.Painter with p.
func WithPainter(v View, p paint.Painter) View {
	return &painterView{View: v, painter: p}
}

type painterView struct {
	View
	painter paint.Painter
}

func (v *painterView) Build(ctx *Context) Model {
	m := v.View.Build(ctx)
	m.Painter = v.painter
	return m
}

// WithOptions wraps the view v, and adds the given options to its Model.Options.
func WithOptions(v View, opts []Option) View {
	return &optionsView{View: v, options: opts}
}

type optionsView struct {
	View
	options []Option
}

func (v *optionsView) Build(ctx *Context) Model {
	m := v.View.Build(ctx)
	m.Options = append(m.Options, v.options...)
	return m
}
