// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gomatcha.io/matcha/pb/view/android/stackview.proto

package android

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import matcha "gomatcha.io/matcha/pb"
import matcha_text "gomatcha.io/matcha/pb/text"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type StackChildView struct {
	ScreenId int64 `protobuf:"varint,3,opt,name=screenId" json:"screenId,omitempty"`
}

func (m *StackChildView) Reset()                    { *m = StackChildView{} }
func (m *StackChildView) String() string            { return proto.CompactTextString(m) }
func (*StackChildView) ProtoMessage()               {}
func (*StackChildView) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *StackChildView) GetScreenId() int64 {
	if m != nil {
		return m.ScreenId
	}
	return 0
}

type StackView struct {
	Children []*StackChildView `protobuf:"bytes,1,rep,name=children" json:"children,omitempty"`
}

func (m *StackView) Reset()                    { *m = StackView{} }
func (m *StackView) String() string            { return proto.CompactTextString(m) }
func (*StackView) ProtoMessage()               {}
func (*StackView) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *StackView) GetChildren() []*StackChildView {
	if m != nil {
		return m.Children
	}
	return nil
}

type StackBar struct {
	Title            string                  `protobuf:"bytes,1,opt,name=title" json:"title,omitempty"`
	StyledTitle      *matcha_text.StyledText `protobuf:"bytes,6,opt,name=styledTitle" json:"styledTitle,omitempty"`
	Subtitle         string                  `protobuf:"bytes,3,opt,name=subtitle" json:"subtitle,omitempty"`
	StyledSubtitle   *matcha_text.StyledText `protobuf:"bytes,7,opt,name=styledSubtitle" json:"styledSubtitle,omitempty"`
	Color            *matcha.Color           `protobuf:"bytes,4,opt,name=color" json:"color,omitempty"`
	Items            []*StackBarItem         `protobuf:"bytes,5,rep,name=items" json:"items,omitempty"`
	BackButtonHidden bool                    `protobuf:"varint,2,opt,name=backButtonHidden" json:"backButtonHidden,omitempty"`
}

func (m *StackBar) Reset()                    { *m = StackBar{} }
func (m *StackBar) String() string            { return proto.CompactTextString(m) }
func (*StackBar) ProtoMessage()               {}
func (*StackBar) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *StackBar) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *StackBar) GetStyledTitle() *matcha_text.StyledText {
	if m != nil {
		return m.StyledTitle
	}
	return nil
}

func (m *StackBar) GetSubtitle() string {
	if m != nil {
		return m.Subtitle
	}
	return ""
}

func (m *StackBar) GetStyledSubtitle() *matcha_text.StyledText {
	if m != nil {
		return m.StyledSubtitle
	}
	return nil
}

func (m *StackBar) GetColor() *matcha.Color {
	if m != nil {
		return m.Color
	}
	return nil
}

func (m *StackBar) GetItems() []*StackBarItem {
	if m != nil {
		return m.Items
	}
	return nil
}

func (m *StackBar) GetBackButtonHidden() bool {
	if m != nil {
		return m.BackButtonHidden
	}
	return false
}

type StackBarItem struct {
	Title       string                  `protobuf:"bytes,1,opt,name=title" json:"title,omitempty"`
	Icon        *matcha.ImageOrResource `protobuf:"bytes,3,opt,name=icon" json:"icon,omitempty"`
	IconTint    *matcha.Color           `protobuf:"bytes,2,opt,name=iconTint" json:"iconTint,omitempty"`
	Disabled    bool                    `protobuf:"varint,4,opt,name=disabled" json:"disabled,omitempty"`
	OnPressFunc string                  `protobuf:"bytes,5,opt,name=onPressFunc" json:"onPressFunc,omitempty"`
}

func (m *StackBarItem) Reset()                    { *m = StackBarItem{} }
func (m *StackBarItem) String() string            { return proto.CompactTextString(m) }
func (*StackBarItem) ProtoMessage()               {}
func (*StackBarItem) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *StackBarItem) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *StackBarItem) GetIcon() *matcha.ImageOrResource {
	if m != nil {
		return m.Icon
	}
	return nil
}

func (m *StackBarItem) GetIconTint() *matcha.Color {
	if m != nil {
		return m.IconTint
	}
	return nil
}

func (m *StackBarItem) GetDisabled() bool {
	if m != nil {
		return m.Disabled
	}
	return false
}

func (m *StackBarItem) GetOnPressFunc() string {
	if m != nil {
		return m.OnPressFunc
	}
	return ""
}

type StackEvent struct {
	Id []int64 `protobuf:"varint,1,rep,packed,name=id" json:"id,omitempty"`
}

func (m *StackEvent) Reset()                    { *m = StackEvent{} }
func (m *StackEvent) String() string            { return proto.CompactTextString(m) }
func (*StackEvent) ProtoMessage()               {}
func (*StackEvent) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *StackEvent) GetId() []int64 {
	if m != nil {
		return m.Id
	}
	return nil
}

func init() {
	proto.RegisterType((*StackChildView)(nil), "matcha.view.android.StackChildView")
	proto.RegisterType((*StackView)(nil), "matcha.view.android.StackView")
	proto.RegisterType((*StackBar)(nil), "matcha.view.android.StackBar")
	proto.RegisterType((*StackBarItem)(nil), "matcha.view.android.StackBarItem")
	proto.RegisterType((*StackEvent)(nil), "matcha.view.android.StackEvent")
}

func init() { proto.RegisterFile("gomatcha.io/matcha/pb/view/android/stackview.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 461 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xe5, 0xa4, 0x6e, 0xdd, 0x31, 0x04, 0xb4, 0x20, 0x61, 0x45, 0x1c, 0x5c, 0x57, 0x48,
	0xe6, 0x8f, 0x6c, 0x29, 0x1c, 0x10, 0xa7, 0x0a, 0x57, 0x20, 0x22, 0x15, 0x11, 0x6d, 0x2a, 0x0e,
	0xdc, 0x6c, 0xef, 0xa8, 0x5d, 0xe1, 0xec, 0x46, 0xeb, 0x4d, 0x5b, 0x5e, 0x87, 0xc7, 0xe0, 0xc8,
	0x93, 0x21, 0xcf, 0x26, 0x56, 0x4a, 0x53, 0x2e, 0x5e, 0xcf, 0x37, 0xbf, 0x99, 0xdd, 0xf9, 0x34,
	0x30, 0xb9, 0xd0, 0x8b, 0xd2, 0xd6, 0x97, 0x65, 0x26, 0x75, 0xee, 0xfe, 0xf2, 0x65, 0x95, 0x5f,
	0x49, 0xbc, 0xce, 0x4b, 0x25, 0x8c, 0x96, 0x22, 0x6f, 0x6d, 0x59, 0xff, 0xe8, 0x94, 0x6c, 0x69,
	0xb4, 0xd5, 0xec, 0xc9, 0xba, 0x82, 0xa4, 0x35, 0x34, 0x3e, 0xda, 0xdd, 0x48, 0x2e, 0xca, 0x0b,
	0x74, 0x75, 0xe3, 0x17, 0xbb, 0x11, 0x8b, 0x37, 0x96, 0x3e, 0x0e, 0x4b, 0xde, 0xc0, 0x68, 0xde,
	0xdd, 0x78, 0x7a, 0x29, 0x1b, 0xf1, 0x4d, 0xe2, 0x35, 0x1b, 0x43, 0xd0, 0xd6, 0x06, 0x51, 0x4d,
	0x45, 0x34, 0x8c, 0xbd, 0x74, 0xc8, 0xfb, 0x38, 0x39, 0x83, 0x43, 0xa2, 0x09, 0x3c, 0x81, 0xa0,
	0xee, 0xaa, 0x0c, 0xaa, 0xc8, 0x8b, 0x87, 0x69, 0x38, 0x39, 0xce, 0x76, 0x3c, 0x36, 0xbb, 0xdd,
	0x9f, 0xf7, 0x45, 0xc9, 0x9f, 0x01, 0x04, 0x94, 0x2c, 0x4a, 0xc3, 0x9e, 0x82, 0x6f, 0xa5, 0x6d,
	0x30, 0xf2, 0x62, 0x2f, 0x3d, 0xe4, 0x2e, 0x60, 0xef, 0x21, 0x6c, 0xed, 0xcf, 0x06, 0xc5, 0x39,
	0xe5, 0xf6, 0x63, 0x2f, 0x0d, 0x27, 0xcf, 0x36, 0xd7, 0xd0, 0x1c, 0x73, 0x97, 0xc7, 0x1b, 0xcb,
	0xb7, 0x59, 0x9a, 0x63, 0x55, 0xb9, 0x9e, 0x43, 0xea, 0xd9, 0xc7, 0xec, 0x04, 0x46, 0x0e, 0x9d,
	0x6f, 0x88, 0x83, 0xff, 0x77, 0xfe, 0x07, 0x67, 0xc7, 0xe0, 0xd7, 0xba, 0xd1, 0x26, 0xda, 0xa3,
	0xba, 0x87, 0x9b, 0xba, 0xd3, 0x4e, 0xe4, 0x2e, 0xc7, 0xde, 0x81, 0x2f, 0x2d, 0x2e, 0xda, 0xc8,
	0x27, 0x77, 0x8e, 0xee, 0x77, 0xa7, 0x28, 0xcd, 0xd4, 0xe2, 0x82, 0x3b, 0x9e, 0xbd, 0x82, 0xc7,
	0x55, 0xa7, 0xae, 0xac, 0xd5, 0xea, 0xb3, 0x14, 0x02, 0x55, 0x34, 0x88, 0xbd, 0x34, 0xe0, 0x77,
	0xf4, 0xe4, 0xb7, 0x07, 0x0f, 0xb6, 0x7b, 0xdc, 0x63, 0xe4, 0x6b, 0xd8, 0x93, 0xb5, 0x56, 0xe4,
	0xc4, 0xd6, 0x9c, 0xd3, 0x6e, 0x63, 0xbe, 0x1a, 0x8e, 0xad, 0x5e, 0x99, 0x1a, 0x39, 0x41, 0xec,
	0x25, 0x04, 0xdd, 0x79, 0x2e, 0x95, 0xa5, 0x7b, 0xef, 0x0c, 0xd8, 0xa7, 0x3b, 0x97, 0x85, 0x6c,
	0xcb, 0xaa, 0x41, 0x41, 0x5e, 0x04, 0xbc, 0x8f, 0x59, 0x0c, 0xa1, 0x56, 0x33, 0x83, 0x6d, 0xfb,
	0x69, 0xa5, 0xea, 0xc8, 0xa7, 0xf7, 0x6c, 0x4b, 0xc9, 0x73, 0x00, 0x7a, 0xfb, 0xc7, 0x2b, 0x54,
	0x96, 0x8d, 0x60, 0x20, 0x05, 0xad, 0xd2, 0x90, 0x0f, 0xa4, 0x28, 0xce, 0x20, 0x91, 0x3a, 0xeb,
	0xf7, 0x78, 0x7d, 0x2c, 0xab, 0x5b, 0x06, 0x16, 0xe1, 0xac, 0xea, 0x77, 0xf2, 0xfb, 0xc1, 0x5a,
	0xfd, 0x35, 0x78, 0xf4, 0x85, 0xf0, 0x0f, 0x2e, 0x9e, 0x15, 0xd5, 0x3e, 0x2d, 0xfc, 0xdb, 0xbf,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x49, 0xd4, 0x32, 0x13, 0x85, 0x03, 0x00, 0x00,
}