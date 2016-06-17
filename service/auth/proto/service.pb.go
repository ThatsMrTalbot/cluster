// Code generated by protoc-gen-go.
// source: proto/service.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	proto/service.proto

It has these top-level messages:
	User
	Token
	AuthRequest
	RefreshRequest
	Response
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type TokenType int32

const (
	TokenType_Auth    TokenType = 0
	TokenType_Refresh TokenType = 1
)

var TokenType_name = map[int32]string{
	0: "Auth",
	1: "Refresh",
}
var TokenType_value = map[string]int32{
	"Auth":    0,
	"Refresh": 1,
}

func (x TokenType) String() string {
	return proto1.EnumName(TokenType_name, int32(x))
}
func (TokenType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type User struct {
	UID         string   `protobuf:"bytes,1,opt,name=UID,json=uID" json:"UID,omitempty"`
	Permissions []string `protobuf:"bytes,2,rep,name=Permissions,json=permissions" json:"Permissions,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto1.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Token struct {
	Type   TokenType `protobuf:"varint,1,opt,name=type,enum=proto.TokenType" json:"type,omitempty"`
	Expiry int64     `protobuf:"varint,2,opt,name=expiry" json:"expiry,omitempty"`
	User   *User     `protobuf:"bytes,3,opt,name=user" json:"user,omitempty"`
}

func (m *Token) Reset()                    { *m = Token{} }
func (m *Token) String() string            { return proto1.CompactTextString(m) }
func (*Token) ProtoMessage()               {}
func (*Token) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Token) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type AuthRequest struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *AuthRequest) Reset()                    { *m = AuthRequest{} }
func (m *AuthRequest) String() string            { return proto1.CompactTextString(m) }
func (*AuthRequest) ProtoMessage()               {}
func (*AuthRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type RefreshRequest struct {
	Token string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
}

func (m *RefreshRequest) Reset()                    { *m = RefreshRequest{} }
func (m *RefreshRequest) String() string            { return proto1.CompactTextString(m) }
func (*RefreshRequest) ProtoMessage()               {}
func (*RefreshRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type Response struct {
	Token   string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
	Refresh string `protobuf:"bytes,2,opt,name=refresh" json:"refresh,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto1.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func init() {
	proto1.RegisterType((*User)(nil), "proto.User")
	proto1.RegisterType((*Token)(nil), "proto.Token")
	proto1.RegisterType((*AuthRequest)(nil), "proto.AuthRequest")
	proto1.RegisterType((*RefreshRequest)(nil), "proto.RefreshRequest")
	proto1.RegisterType((*Response)(nil), "proto.Response")
	proto1.RegisterEnum("proto.TokenType", TokenType_name, TokenType_value)
}

var fileDescriptor0 = []byte{
	// 272 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x90, 0xc1, 0x4b, 0xc3, 0x30,
	0x14, 0xc6, 0xad, 0x6d, 0xb7, 0xee, 0x05, 0x46, 0x79, 0x8a, 0x14, 0x2f, 0x96, 0x20, 0x32, 0x3c,
	0x4c, 0x98, 0x37, 0x6f, 0xc2, 0x3c, 0xec, 0x26, 0x61, 0xfb, 0x03, 0xa6, 0xbe, 0x61, 0x91, 0x35,
	0x31, 0x49, 0xd5, 0xfd, 0xf7, 0x26, 0x69, 0x56, 0x2f, 0x9e, 0xca, 0xf7, 0x7e, 0x7d, 0xdf, 0xf7,
	0xbd, 0xc0, 0x99, 0xd2, 0xd2, 0xca, 0x3b, 0x43, 0xfa, 0xab, 0x79, 0xa5, 0x79, 0x50, 0x98, 0x87,
	0x0f, 0x7f, 0x80, 0x6c, 0xe3, 0x00, 0x96, 0x90, 0x6e, 0x56, 0xcb, 0x2a, 0xa9, 0x93, 0xd9, 0x44,
	0xa4, 0xdd, 0x6a, 0x89, 0x35, 0xb0, 0x67, 0xd2, 0xfb, 0xc6, 0x98, 0x46, 0xb6, 0xa6, 0x3a, 0xad,
	0x53, 0x47, 0x98, 0xfa, 0x1b, 0xf1, 0x1d, 0xe4, 0x6b, 0xf9, 0x41, 0x2d, 0x5e, 0x43, 0x66, 0x0f,
	0x8a, 0xc2, 0xf6, 0x74, 0x51, 0xf6, 0x09, 0xf3, 0xc0, 0xd6, 0x6e, 0x2e, 0x02, 0xc5, 0x0b, 0x18,
	0xd1, 0x8f, 0x6a, 0xf4, 0xc1, 0x79, 0x25, 0xb3, 0x54, 0x44, 0x85, 0x57, 0x90, 0x75, 0xae, 0x42,
	0x95, 0xba, 0x29, 0x5b, 0xb0, 0xb8, 0xed, 0x5b, 0x89, 0x00, 0xf8, 0x13, 0xb0, 0xc7, 0xce, 0xbe,
	0x0b, 0xfa, 0xec, 0xc8, 0x58, 0xbc, 0x84, 0xc2, 0x8f, 0xdb, 0xed, 0x9e, 0x62, 0xdf, 0x41, 0x7b,
	0xa6, 0xb6, 0xc6, 0x7c, 0x4b, 0xfd, 0x16, 0x52, 0x1c, 0x3b, 0x6a, 0x7e, 0x03, 0x53, 0x41, 0x3b,
	0x4d, 0x66, 0x70, 0x3a, 0x87, 0xdc, 0xfa, 0x92, 0xd1, 0xa6, 0x17, 0xee, 0x49, 0x0a, 0x41, 0x46,
	0xb9, 0x0b, 0xe9, 0xff, 0x3f, 0xb0, 0x82, 0xb1, 0xee, 0x9d, 0x62, 0xc8, 0x51, 0xde, 0x72, 0x98,
	0x0c, 0x67, 0x63, 0x01, 0x99, 0xef, 0x5d, 0x9e, 0x20, 0x83, 0x71, 0x8c, 0x2e, 0x93, 0x97, 0x51,
	0x38, 0xf0, 0xfe, 0x37, 0x00, 0x00, 0xff, 0xff, 0x72, 0x4c, 0x0b, 0xc7, 0x97, 0x01, 0x00, 0x00,
}