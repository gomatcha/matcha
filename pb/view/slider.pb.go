// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gomatcha.io/matcha/pb/view/slider.proto

package view

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Slider struct {
	Value    float64 `protobuf:"fixed64,1,opt,name=value" json:"value,omitempty"`
	MaxValue float64 `protobuf:"fixed64,2,opt,name=maxValue" json:"maxValue,omitempty"`
	MinValue float64 `protobuf:"fixed64,3,opt,name=minValue" json:"minValue,omitempty"`
	Enabled  bool    `protobuf:"varint,4,opt,name=enabled" json:"enabled,omitempty"`
}

func (m *Slider) Reset()                    { *m = Slider{} }
func (m *Slider) String() string            { return proto.CompactTextString(m) }
func (*Slider) ProtoMessage()               {}
func (*Slider) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{0} }

func (m *Slider) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *Slider) GetMaxValue() float64 {
	if m != nil {
		return m.MaxValue
	}
	return 0
}

func (m *Slider) GetMinValue() float64 {
	if m != nil {
		return m.MinValue
	}
	return 0
}

func (m *Slider) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

type SliderEvent struct {
	Value float64 `protobuf:"fixed64,1,opt,name=value" json:"value,omitempty"`
}

func (m *SliderEvent) Reset()                    { *m = SliderEvent{} }
func (m *SliderEvent) String() string            { return proto.CompactTextString(m) }
func (*SliderEvent) ProtoMessage()               {}
func (*SliderEvent) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{1} }

func (m *SliderEvent) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func init() {
	proto.RegisterType((*Slider)(nil), "matcha.view.Slider")
	proto.RegisterType((*SliderEvent)(nil), "matcha.view.SliderEvent")
}

func init() { proto.RegisterFile("gomatcha.io/matcha/pb/view/slider.proto", fileDescriptor4) }

var fileDescriptor4 = []byte{
	// 192 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4f, 0xcf, 0xcf, 0x4d,
	0x2c, 0x49, 0xce, 0x48, 0xd4, 0xcb, 0xcc, 0xd7, 0x87, 0xb0, 0xf4, 0x0b, 0x92, 0xf4, 0xcb, 0x32,
	0x53, 0xcb, 0xf5, 0x8b, 0x73, 0x32, 0x53, 0x52, 0x8b, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85,
	0xb8, 0xa1, 0xca, 0x40, 0x32, 0x4a, 0x05, 0x5c, 0x6c, 0xc1, 0x60, 0x49, 0x21, 0x11, 0x2e, 0xd6,
	0xb2, 0xc4, 0x9c, 0xd2, 0x54, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xc6, 0x20, 0x08, 0x47, 0x48, 0x8a,
	0x8b, 0x23, 0x37, 0xb1, 0x22, 0x0c, 0x2c, 0xc1, 0x04, 0x96, 0x80, 0xf3, 0xc1, 0x72, 0x99, 0x79,
	0x10, 0x39, 0x66, 0xa8, 0x1c, 0x94, 0x2f, 0x24, 0xc1, 0xc5, 0x9e, 0x9a, 0x97, 0x98, 0x94, 0x93,
	0x9a, 0x22, 0xc1, 0xa2, 0xc0, 0xa8, 0xc1, 0x11, 0x04, 0xe3, 0x2a, 0x29, 0x73, 0x71, 0x43, 0x6c,
	0x74, 0x2d, 0x4b, 0xcd, 0x2b, 0xc1, 0x6e, 0xad, 0x93, 0x35, 0x97, 0x54, 0x66, 0xbe, 0x1e, 0xdc,
	0x47, 0x50, 0xaa, 0x20, 0x09, 0xec, 0x68, 0x27, 0x8e, 0x80, 0x24, 0x88, 0x11, 0x51, 0x2c, 0x20,
	0xfe, 0x22, 0x26, 0x1e, 0x5f, 0xb0, 0x82, 0xb0, 0xcc, 0xd4, 0xf2, 0x80, 0xa4, 0x24, 0x36, 0xb0,
	0x3f, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x80, 0x55, 0x06, 0xec, 0x12, 0x01, 0x00, 0x00,
}