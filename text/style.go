package text

import (
	"image/color"
	"reflect"
	"runtime"

	pbtext "github.com/gomatcha/matcha/proto/text"
)

// Alignment represents a text alignment.
type Alignment int

const (
	AlignmentLeft Alignment = iota
	AlignmentRight
	AlignmentCenter
	AlignmentJustified
)

func (a Alignment) MarshalProtobuf() pbtext.TextAlignment {
	return pbtext.TextAlignment(a)
}

// StrikethroughStyle represents a text strikethrough style.
type StrikethroughStyle int

const (
	StrikethroughStyleNone StrikethroughStyle = iota
	StrikethroughStyleSingle
	StrikethroughStyleDouble
	StrikethroughStyleThick
	StrikethroughStyleDotted
	StrikethroughStyleDashed
)

func (a StrikethroughStyle) MarshalProtobuf() pbtext.StrikethroughStyle {
	return pbtext.StrikethroughStyle(a)
}

// StrikethroughStyle represents a text underline style.
type UnderlineStyle int

const (
	UnderlineStyleNone UnderlineStyle = iota
	UnderlineStyleSingle
	UnderlineStyleDouble
	UnderlineStyleThick
	UnderlineStyleDotted
	UnderlineStyleDashed
)

func (a UnderlineStyle) MarshalProtobuf() pbtext.UnderlineStyle {
	return pbtext.UnderlineStyle(a)
}

func DefaultFont(size float64) *Font {
	if runtime.GOOS == "android" {
		return FontWithName("sans-serif", size)
	} else if runtime.GOOS == "darwin" {
		return FontWithName("HelveticaNeue", size)
	}
	return &Font{}
}

func DefaultBoldFont(size float64) *Font {
	if runtime.GOOS == "android" {
		return FontWithName("sans-serif-bold", size)
	} else if runtime.GOOS == "darwin" {
		return FontWithName("HelveticaNeue-Bold", size)
	}
	return &Font{}
}

func DefaultItalicFont(size float64) *Font {
	if runtime.GOOS == "android" {
		return FontWithName("sans-serif-italic", size)
	} else if runtime.GOOS == "darwin" {
		return FontWithName("HelveticaNeue-Italic", size)
	}
	return &Font{}
}

func FontWithName(name string, size float64) *Font {
	return &Font{
		name: name,
		size: size,
	}
}

// StrikethroughStyle represents a text font.
type Font struct {
	name string // Postscript name
	size float64
}

func (f *Font) MarshalProtobuf() *pbtext.Font {
	return &pbtext.Font{
		Family: f.name,
		Size:   f.size,
	}
}

// StrikethroughStyle represents how text is wrapped.
type Wrap int

const (
	WrapNone Wrap = iota
	WrapWord
	WrapCharacter
)

func (a Wrap) MarshalProtobuf() pbtext.TextWrap {
	return pbtext.TextWrap(a)
}

// Truncation represents how text is truncated to fit within the bounds.
type Truncation int

const (
	TruncationNone Truncation = iota
	TruncationStart
	TruncationMiddle
	TruncationEnd
)

func (a Truncation) MarshalProtobuf() pbtext.Truncation {
	return pbtext.Truncation(a)
}

type styleKey int

const (
	styleKeyAlignment styleKey = iota
	styleKeyStrikethroughStyle
	styleKeyStrikethroughColor
	styleKeyUnderlineStyle
	styleKeyUnderlineColor
	styleKeyFont
	styleKeyHyphenation
	styleKeyLineHeightMultiple
	styleKeyMaxLines // Deprecated
	styleKeyTextColor
	styleKeyWrap
	styleKeyTruncation
	styleKeyTruncationString
)

// Style holds a group of text formatting options.
type Style struct {
	attributes map[styleKey]interface{}
	cleared    map[styleKey]bool
}

func (f *Style) map_() map[styleKey]interface{} {
	return f.attributes
}

func (f *Style) clear(k styleKey) {
	if f.cleared == nil || f.attributes == nil {
		f.attributes = map[styleKey]interface{}{}
		f.cleared = map[styleKey]bool{}
	}

	delete(f.attributes, k)
	f.cleared[k] = true
}

func (f *Style) get(k styleKey) interface{} {
	v, ok := f.attributes[k]
	if ok {
		return v
	}
	switch k {
	case styleKeyAlignment:
		return AlignmentLeft
	case styleKeyStrikethroughStyle:
		return StrikethroughStyleNone
	case styleKeyStrikethroughColor:
		return color.Gray{0}
	case styleKeyUnderlineStyle:
		return UnderlineStyleNone
	case styleKeyUnderlineColor:
		return color.Gray{0}
	case styleKeyFont:
		return DefaultFont(14)
	case styleKeyHyphenation:
		return float64(0.0)
	case styleKeyLineHeightMultiple:
		return float64(1.0)
	case styleKeyMaxLines:
		return 0
	case styleKeyTextColor:
		return color.Gray{0}
	case styleKeyWrap:
		return WrapWord
	case styleKeyTruncation:
		return TruncationNone
	case styleKeyTruncationString:
		return "…"
	}
	return nil
}

func (f *Style) set(k styleKey, v interface{}) {
	if f.cleared == nil || f.attributes == nil {
		f.attributes = map[styleKey]interface{}{}
		f.cleared = map[styleKey]bool{}
	}

	f.attributes[k] = v
	delete(f.cleared, k)
}

func (f *Style) Equal(f2 *Style) bool {
	return reflect.DeepEqual(f, f2)
}

func (f *Style) Copy() *Style {
	if f == nil {
		return nil
	}

	c := &Style{}
	if f.attributes != nil {
		c.attributes = map[styleKey]interface{}{}
		for k, v := range f.attributes {
			c.attributes[k] = v
		}
	}
	if f.cleared != nil {
		c.cleared = map[styleKey]bool{}
		for k, v := range f.cleared {
			c.cleared[k] = v
		}
	}
	return c
}

// Applies the styels from u to f.
func (f *Style) Update(u *Style) {
	for k, v := range u.attributes {
		f.attributes[k] = v
	}
	for k := range u.cleared {
		delete(f.attributes, k)
	}
}

func (f *Style) MarshalProtobuf() *pbtext.TextStyle {
	if f == nil {
		f = &Style{}
	}

	s := &pbtext.TextStyle{
		TextAlignment:      f.get(styleKeyAlignment).(Alignment).MarshalProtobuf(),
		StrikethroughStyle: f.get(styleKeyStrikethroughStyle).(StrikethroughStyle).MarshalProtobuf(),
		UnderlineStyle:     f.get(styleKeyUnderlineStyle).(UnderlineStyle).MarshalProtobuf(),
		Hyphenation:        f.get(styleKeyHyphenation).(float64),
		LineHeightMultiple: f.get(styleKeyLineHeightMultiple).(float64),
		Wrap:               f.get(styleKeyWrap).(Wrap).MarshalProtobuf(),
		Truncation:         f.get(styleKeyTruncation).(Truncation).MarshalProtobuf(),
		TruncationString:   f.get(styleKeyTruncationString).(string),
	}
	{
		font := f.get(styleKeyFont).(*Font)
		s.FontName = font.name
		s.FontSize = font.size
	}
	{
		r, g, b, a := f.get(styleKeyTextColor).(color.Color).RGBA()
		s.HasTextColor = true
		s.TextColorRed = r
		s.TextColorGreen = g
		s.TextColorBlue = b
		s.TextColorAlpha = a
	}
	{
		r, g, b, a := f.get(styleKeyStrikethroughColor).(color.Color).RGBA()
		s.HasStrikethroughColor = true
		s.StrikethroughColorRed = r
		s.StrikethroughColorGreen = g
		s.StrikethroughColorBlue = b
		s.StrikethroughColorAlpha = a
	}
	{
		r, g, b, a := f.get(styleKeyUnderlineColor).(color.Color).RGBA()
		s.HasUnderlineColor = true
		s.UnderlineColorRed = r
		s.UnderlineColorGreen = g
		s.UnderlineColorBlue = b
		s.UnderlineColorAlpha = a
	}
	return s
}

func (f *Style) Alignment() Alignment {
	return f.get(styleKeyAlignment).(Alignment)
}

func (f *Style) SetAlignment(v Alignment) {
	f.set(styleKeyAlignment, v)
}

func (f *Style) ClearAlignment() {
	f.clear(styleKeyAlignment)
}

func (f *Style) StrikethroughStyle() StrikethroughStyle {
	return f.get(styleKeyStrikethroughStyle).(StrikethroughStyle)
}

func (f *Style) SetStrikethroughStyle(v StrikethroughStyle) {
	f.set(styleKeyStrikethroughStyle, v)
}

func (f *Style) ClearStrikethroughStyle() {
	f.clear(styleKeyStrikethroughStyle)
}

func (f *Style) StrikethroughColor() color.Color {
	return f.get(styleKeyStrikethroughColor).(color.Color)
}

func (f *Style) SetStrikethroughColor(v color.Color) {
	f.set(styleKeyStrikethroughColor, v)
}

func (f *Style) ClearStrikethroughColor() {
	f.clear(styleKeyStrikethroughColor)
}

func (f *Style) UnderlineStyle() UnderlineStyle {
	return f.get(styleKeyUnderlineStyle).(UnderlineStyle)
}

func (f *Style) SetUnderlineStyle(v UnderlineStyle) {
	f.set(styleKeyUnderlineStyle, v)
}

func (f *Style) ClearUnderlineStyle() {
	f.clear(styleKeyUnderlineStyle)
}

func (f *Style) UnderlineColor() color.Color {
	return f.get(styleKeyUnderlineColor).(color.Color)
}

func (f *Style) SetUnderlineColor(v color.Color) {
	f.set(styleKeyUnderlineColor, v)
}

func (f *Style) ClearUnderlineColor() {
	f.clear(styleKeyUnderlineColor)
}

func (f *Style) Font() *Font {
	return f.get(styleKeyFont).(*Font)
}

func (f *Style) SetFont(v *Font) {
	f.set(styleKeyFont, v)
}

func (f *Style) ClearFont() {
	f.clear(styleKeyFont)
}

func (f *Style) Hyphenation() float64 {
	return f.get(styleKeyHyphenation).(float64)
}

func (f *Style) SetHyphenation(v float64) {
	f.set(styleKeyHyphenation, v)
}

func (f *Style) ClearHyphenation() {
	f.clear(styleKeyHyphenation)
}

func (f *Style) LineHeightMultiple() float64 {
	return f.get(styleKeyLineHeightMultiple).(float64)
}

func (f *Style) SetLineHeightMultiple(v float64) {
	f.set(styleKeyLineHeightMultiple, v)
}

func (f *Style) ClearLineHeightMultiple() {
	f.clear(styleKeyLineHeightMultiple)
}

func (f *Style) TextColor() color.Color {
	return f.get(styleKeyTextColor).(color.Color)
}

func (f *Style) SetTextColor(v color.Color) {
	f.set(styleKeyTextColor, v)
}

func (f *Style) ClearTextColor() {
	f.clear(styleKeyTextColor)
}

func (f *Style) Wrap() Wrap {
	return f.get(styleKeyWrap).(Wrap)
}

func (f *Style) SetWrap(v Wrap) {
	f.set(styleKeyWrap, v)
}

func (f *Style) ClearWrap() {
	f.clear(styleKeyWrap)
}

func (f *Style) Truncation() Truncation {
	return f.get(styleKeyTruncation).(Truncation)
}

func (f *Style) SetTruncation(v Truncation) {
	f.set(styleKeyTruncation, v)
}

func (f *Style) ClearTruncation() {
	f.clear(styleKeyTruncation)
}

func (f *Style) TruncationString() string {
	return f.get(styleKeyTruncationString).(string)
}

func (f *Style) SetTruncationString(v string) {
	f.set(styleKeyTruncationString, v)
}

func (f *Style) ClearTruncationString() {
	f.clear(styleKeyTruncationString)
}
