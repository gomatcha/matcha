package text

import (
	"image/color"

	"gomatcha.io/matcha/pb"
	pbtext "gomatcha.io/matcha/pb/text"
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

// StrikethroughStyle represents a text font.
type Font struct {
	Family string
	Face   string
	Size   float64
}

func (f Font) MarshalProtobuf() *pbtext.Font {
	return &pbtext.Font{
		Family: f.Family,
		Face:   f.Face,
		Size:   f.Size,
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
		return Font{
			Family: "Helvetica Neue",
			Size:   14,
		}
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

func (f *Style) copy() *Style {
	c := &Style{
		attributes: map[styleKey]interface{}{},
		cleared:    map[styleKey]bool{},
	}
	for k, v := range f.attributes {
		c.attributes[k] = v
	}
	for k, v := range f.cleared {
		c.attributes[k] = v
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

	return &pbtext.TextStyle{
		TextAlignment:      f.get(styleKeyAlignment).(Alignment).MarshalProtobuf(),
		StrikethroughStyle: f.get(styleKeyStrikethroughStyle).(StrikethroughStyle).MarshalProtobuf(),
		StrikethroughColor: pb.ColorEncode(f.get(styleKeyStrikethroughColor).(color.Color)),
		UnderlineStyle:     f.get(styleKeyUnderlineStyle).(UnderlineStyle).MarshalProtobuf(),
		UnderlineColor:     pb.ColorEncode(f.get(styleKeyUnderlineColor).(color.Color)),
		Font:               f.get(styleKeyFont).(Font).MarshalProtobuf(),
		Hyphenation:        f.get(styleKeyHyphenation).(float64),
		LineHeightMultiple: f.get(styleKeyLineHeightMultiple).(float64),
		TextColor:          pb.ColorEncode(f.get(styleKeyTextColor).(color.Color)),
		Wrap:               f.get(styleKeyWrap).(Wrap).MarshalProtobuf(),
		Truncation:         f.get(styleKeyTruncation).(Truncation).MarshalProtobuf(),
		TruncationString:   f.get(styleKeyTruncationString).(string),
	}
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

func (f *Style) Font() Font {
	return f.get(styleKeyFont).(Font)
}

func (f *Style) SetFont(v Font) {
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
