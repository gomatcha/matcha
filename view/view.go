// Package view provides the component library. See
// https://gomatcha.io/guide/view/ for more details.
package view

import (
	"fmt"
	"reflect"
	"sync"

	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
)

type Id int64

type View interface {
	Build(Context) Model
	Lifecycle(from, to Stage)
	ViewKey() interface{}
	Update(View)
	comm.Notifier
}

type Option interface {
	OptionKey() string
}

var embedUpdate bool

// Embed is a convenience struct that provides a default implementation of View. It also wraps a comm.Relay.
type Embed struct {
	key       interface{}
	Key       interface{}
	mu        sync.Mutex
	relay     comm.Relay
	lastField struct{}
}

func NewEmbed(key interface{}) Embed {
	return Embed{
		key: key,
	}
}

// Build is an empty implementation of View's Build method.
func (e *Embed) Build(ctx Context) Model {
	return Model{}
}

func (e *Embed) ViewKey() interface{} {
	return struct {
		A interface{}
		B interface{}
	}{A: e.key, B: e.Key}
}

// Lifecycle is an empty implementation of View's Lifecycle method.
func (e *Embed) Lifecycle(from, to Stage) {
	// no-op
}

// Update is an empty implementation of View's Update method.
func (e *Embed) Update(v View) {
	embedUpdate = true
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

// Copy all public fields from src to dst, that aren't 'Embed'.
func CopyFields(dst, src View) {
	va := reflect.ValueOf(dst).Elem()
	vb := reflect.ValueOf(src).Elem()
	for i := 0; i < va.NumField(); i++ {
		fa := va.Field(i)
		if fa.CanSet() && va.Type().Field(i).Name != "Embed" {
			fa.Set(vb.Field(i))
		}
	}
}

type Stage int

const (
	// Views start in StageDead.
	StageDead Stage = iota
	// Views enter StageMounted before the first Build() call and exit it after the last one.
	StageMounted
	// StageVisible marks views that are in the view hierarchy and visible.
	// TODO(KD): StageVisible is not used.
	StageVisible
)

// EntersStage returns true if from﹤s and to≥s.
func EntersStage(from, to, s Stage) bool {
	return from < s && to >= s
}

// ExitsStage returns true if from≥s and to﹤s.
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
	NativeViewState []byte
	NativeOptions   map[string][]byte
	NativeFuncs     map[string]interface{}
}

// WithPainter wraps the view v, and replaces its Model.Painter with p.
func WithPainter(v View, p paint.Painter) View {
	return &painterView{view: v, painter: p}
}

type painterView struct {
	view    View
	painter paint.Painter
}

func (v *painterView) ViewKey() interface{} {
	return struct {
		A interface{}
		B interface{}
	}{A: v.view.ViewKey(), B: internal.ReflectName(v.view)}
}

func (v *painterView) Lifecycle(from, to Stage) {
	v.view.Lifecycle(from, to)
}

func (v *painterView) Update(v2 View) {
	v.view.Update(v2.(*painterView).view)
}

func (v *painterView) Build(ctx Context) Model {
	m := v.view.Build(ctx)
	m.Painter = v.painter
	return m
}

func (v *painterView) Notify(f func()) comm.Id {
	id := v.view.Notify(f)
	return id
}

func (v *painterView) Unnotify(id comm.Id) {
	v.view.Unnotify(id)
}

// WithOptions wraps the view v, and adds the given options to its Model.Options.
func WithOptions(v View, opts ...Option) View {
	return &optionsView{view: v, options: opts}
}

type optionsView struct {
	view    View
	options []Option
}

func (v *optionsView) ViewKey() interface{} {
	return struct {
		A interface{}
		B interface{}
	}{A: v.view.ViewKey(), B: internal.ReflectName(v.view)}
}

func (v *optionsView) Lifecycle(from, to Stage) {
	v.view.Lifecycle(from, to)
}

func (v *optionsView) Update(v2 View) {
	v.view.Update(v2.(*optionsView).view)
}

func (v *optionsView) Build(ctx Context) Model {
	m := v.view.Build(ctx)
	fmt.Println("build", m.Options, v.options, append(m.Options, v.options...))
	m.Options = append(m.Options, v.options...)
	return m
}

func (v *optionsView) Notify(f func()) comm.Id {
	id := v.view.Notify(f)
	return id
}

func (v *optionsView) Unnotify(id comm.Id) {
	v.view.Unnotify(id)
}
