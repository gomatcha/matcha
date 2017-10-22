package text

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout"
	pb "gomatcha.io/matcha/proto/layout"
	pbtext "gomatcha.io/matcha/proto/text"
)

type styleRange struct {
	index int
	style *Style
}

type StyledText struct {
	text   *Text
	styles []styleRange
}

func NewStyledText(str string, s *Style) *StyledText {
	st := &StyledText{
		text: New(str),
		styles: []styleRange{
			styleRange{index: 0, style: s},
		},
	}
	return st
}

func (st *StyledText) Equal(st2 *StyledText) bool {
	return reflect.DeepEqual(st, st2)
}

func (st *StyledText) Copy() *StyledText {
	if st == nil {
		return nil
	}

	s := make([]styleRange, len(st.styles))
	for idx, i := range st.styles {
		s[idx] = styleRange{
			index: i.index,
			style: i.style.Copy(),
		}
	}
	c := &StyledText{
		text:   st.text.Copy(),
		styles: s,
	}
	return c
}

// returns null if a is outside of text range.
func (st *StyledText) At(a int) *Style {
	if a >= st.text.runeCount || a < 0 {
		return nil
	}

	var style *Style
	for _, i := range st.styles {
		if i.index > a {
			break
		}
		style = i.style
	}
	return style.Copy()
}

func (st *StyledText) Set(s *Style, start, end int) {
	st.update(func(prev *Style) *Style {
		return s.Copy()
	}, start, end)
}

func (st *StyledText) Update(s *Style, start, end int) {
	st.update(func(prev *Style) *Style {
		prev = prev.Copy()
		prev.Update(s)
		return prev
	}, start, end)
}

func (st *StyledText) update(f func(*Style) *Style, start, end int) {
	styles := []styleRange{}
	for idx, i := range st.styles {
		// Calculate the range for the style. rangeMin and rangeMax are inclusive.
		rangeMin := i.index
		rangeMax := 0
		if idx == len(st.styles)-1 {
			rangeMax = st.text.runeCount - 1
		} else {
			rangeMax = st.styles[idx+1].index - 1
		}

		if rangeMax < start || rangeMin > end {
			// If range does not overlap with start/end, add as normal.
			styles = append(styles, i)
		} else if rangeMin < start {
			if rangeMax <= end {
				styles = append(styles, styleRange{index: rangeMin, style: i.style})
				styles = append(styles, styleRange{index: start, style: f(i.style)})
			} else if rangeMax > end {
				styles = append(styles, styleRange{index: rangeMin, style: i.style})
				styles = append(styles, styleRange{index: start, style: f(i.style)})
				styles = append(styles, styleRange{index: end + 1, style: i.style.Copy()})
			}
		} else if rangeMin == start {
			if rangeMax <= end {
				styles = append(styles, styleRange{index: rangeMin, style: f(i.style)})
			} else if rangeMax > end {
				styles = append(styles, styleRange{index: start, style: f(i.style)})
				styles = append(styles, styleRange{index: end + 1, style: i.style})
			}
		} else if rangeMin > start {
			if rangeMax <= end {
				// ignore
			} else if rangeMax > end {
				styles = append(styles, styleRange{index: end + 1, style: i.style})
			}
		}
	}
	st.styles = styles
}

func (st *StyledText) Size(min layout.Point, max layout.Point, maxLines int) layout.Point {
	if st.text.String() == "" {
		st = NewStyledText("A", nil)
	}

	sizeFunc := &pbtext.SizeFunc{
		Text:    st.MarshalProtobuf(),
		MinSize: min.MarshalProtobuf(),
		MaxSize: max.MarshalProtobuf(),
	}
	data, err := proto.Marshal(sizeFunc)
	if err != nil {
		return layout.Pt(0, 0)
	}

	var pointData []byte
	if runtime.GOOS == "android" {
		pointData = bridge.Bridge("").Call("sizeForStyledText", bridge.Bytes(data), bridge.Int64(int64(maxLines))).ToInterface().([]byte)
	} else if runtime.GOOS == "darwin" {
		pointData = bridge.Bridge("").Call("sizeForAttributedString:maxLines:", bridge.Bytes(data), bridge.Int64(int64(maxLines))).ToInterface().([]byte)
	}
	pbpoint := &pb.Point{}
	err = proto.Unmarshal(pointData, pbpoint)
	if err != nil {
		fmt.Println("StyledText.Size(): Decode error", err)
		return layout.Pt(0, 0)
	}
	return layout.Pt(pbpoint.X, pbpoint.Y)
}

func (st *StyledText) MarshalProtobuf() *pbtext.StyledText {
	if st == nil {
		return nil
	}

	styles := []*pbtext.TextStyle{}
	for _, i := range st.styles {
		style := i.style.MarshalProtobuf()
		style.Index = int64(i.index)
		styles = append(styles, style)
	}

	return &pbtext.StyledText{
		Text:   st.text.MarshalProtobuf(),
		Styles: styles,
	}
}
