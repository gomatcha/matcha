package internal

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/layout"
	pb "gomatcha.io/matcha/pb/layout"
	pbtext "gomatcha.io/matcha/pb/text"
	"gomatcha.io/matcha/text"
)

type StyledText struct {
	text  *text.Text
	style *text.Style
}

func NewStyledText(t *text.Text) *StyledText {
	return &StyledText{
		text:  t,
		style: &text.Style{},
	}
}

func (st *StyledText) Text() *text.Text {
	return st.text
}

func (st *StyledText) At(a int) *text.Style {
	return nil
}

func (st *StyledText) Set(s *text.Style, start, end int) {
	st.style = s
}

func (st *StyledText) Update(s *text.Style, start, end int) {
}

func (st *StyledText) Size(min layout.Point, max layout.Point, maxLines int) layout.Point {
	if st.text.String() == "" {
		st = &StyledText{
			text:  text.New("A"),
			style: st.style,
		}
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

	pointData := bridge.Bridge().Call("sizeForAttributedString:maxLines:", bridge.Bytes(data), bridge.Int64(int64(maxLines))).ToInterface().([]byte)
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
	return &pbtext.StyledText{
		Text:  st.text.MarshalProtobuf(),
		Style: st.style.MarshalProtobuf(),
	}
}
