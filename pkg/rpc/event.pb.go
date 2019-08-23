// Code generated by protoc-gen-go. DO NOT EDIT.
// source: event.proto

package rpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ConsumeRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConsumeRequest) Reset()         { *m = ConsumeRequest{} }
func (m *ConsumeRequest) String() string { return proto.CompactTextString(m) }
func (*ConsumeRequest) ProtoMessage()    {}
func (*ConsumeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_event_a1be43c7c699c898, []int{0}
}
func (m *ConsumeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsumeRequest.Unmarshal(m, b)
}
func (m *ConsumeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsumeRequest.Marshal(b, m, deterministic)
}
func (dst *ConsumeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsumeRequest.Merge(dst, src)
}
func (m *ConsumeRequest) XXX_Size() int {
	return xxx_messageInfo_ConsumeRequest.Size(m)
}
func (m *ConsumeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsumeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ConsumeRequest proto.InternalMessageInfo

// 事件的详细信息
type ConsumeResponse struct {
	// 事件id
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// 事件类型
	EventType string `protobuf:"bytes,2,opt,name=eventType,proto3" json:"eventType,omitempty"`
	// 聚合id
	AggId string `protobuf:"bytes,3,opt,name=aggId,proto3" json:"aggId,omitempty"`
	// 聚合类型
	AggType string `protobuf:"bytes,4,opt,name=aggType,proto3" json:"aggType,omitempty"`
	// 事件创建时间
	Create *timestamp.Timestamp `protobuf:"bytes,5,opt,name=create,proto3" json:"create,omitempty"`
	// 事件内容
	Data                 string   `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConsumeResponse) Reset()         { *m = ConsumeResponse{} }
func (m *ConsumeResponse) String() string { return proto.CompactTextString(m) }
func (*ConsumeResponse) ProtoMessage()    {}
func (*ConsumeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_event_a1be43c7c699c898, []int{1}
}
func (m *ConsumeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsumeResponse.Unmarshal(m, b)
}
func (m *ConsumeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsumeResponse.Marshal(b, m, deterministic)
}
func (dst *ConsumeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsumeResponse.Merge(dst, src)
}
func (m *ConsumeResponse) XXX_Size() int {
	return xxx_messageInfo_ConsumeResponse.Size(m)
}
func (m *ConsumeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsumeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ConsumeResponse proto.InternalMessageInfo

func (m *ConsumeResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ConsumeResponse) GetEventType() string {
	if m != nil {
		return m.EventType
	}
	return ""
}

func (m *ConsumeResponse) GetAggId() string {
	if m != nil {
		return m.AggId
	}
	return ""
}

func (m *ConsumeResponse) GetAggType() string {
	if m != nil {
		return m.AggType
	}
	return ""
}

func (m *ConsumeResponse) GetCreate() *timestamp.Timestamp {
	if m != nil {
		return m.Create
	}
	return nil
}

func (m *ConsumeResponse) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*ConsumeRequest)(nil), "rpc.ConsumeRequest")
	proto.RegisterType((*ConsumeResponse)(nil), "rpc.ConsumeResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ConsumerClient is the client API for Consumer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConsumerClient interface {
	Consume(ctx context.Context, in *ConsumeRequest, opts ...grpc.CallOption) (Consumer_ConsumeClient, error)
}

type consumerClient struct {
	cc *grpc.ClientConn
}

func NewConsumerClient(cc *grpc.ClientConn) ConsumerClient {
	return &consumerClient{cc}
}

func (c *consumerClient) Consume(ctx context.Context, in *ConsumeRequest, opts ...grpc.CallOption) (Consumer_ConsumeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Consumer_serviceDesc.Streams[0], "/rpc.Consumer/Consume", opts...)
	if err != nil {
		return nil, err
	}
	x := &consumerConsumeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Consumer_ConsumeClient interface {
	Recv() (*ConsumeResponse, error)
	grpc.ClientStream
}

type consumerConsumeClient struct {
	grpc.ClientStream
}

func (x *consumerConsumeClient) Recv() (*ConsumeResponse, error) {
	m := new(ConsumeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ConsumerServer is the server API for Consumer service.
type ConsumerServer interface {
	Consume(*ConsumeRequest, Consumer_ConsumeServer) error
}

func RegisterConsumerServer(s *grpc.Server, srv ConsumerServer) {
	s.RegisterService(&_Consumer_serviceDesc, srv)
}

func _Consumer_Consume_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConsumeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConsumerServer).Consume(m, &consumerConsumeServer{stream})
}

type Consumer_ConsumeServer interface {
	Send(*ConsumeResponse) error
	grpc.ServerStream
}

type consumerConsumeServer struct {
	grpc.ServerStream
}

func (x *consumerConsumeServer) Send(m *ConsumeResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _Consumer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Consumer",
	HandlerType: (*ConsumerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Consume",
			Handler:       _Consumer_Consume_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "event.proto",
}

func init() { proto.RegisterFile("event.proto", fileDescriptor_event_a1be43c7c699c898) }

var fileDescriptor_event_a1be43c7c699c898 = []byte{
	// 235 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8f, 0xc1, 0x4e, 0xc3, 0x30,
	0x0c, 0x86, 0x49, 0xb7, 0x75, 0xcc, 0x93, 0x06, 0x32, 0x3b, 0x44, 0x15, 0x12, 0x53, 0x4f, 0x3b,
	0x65, 0xa8, 0x5c, 0xb8, 0xc3, 0x85, 0x6b, 0xb5, 0x17, 0xc8, 0x5a, 0x13, 0x55, 0xa2, 0x4d, 0x48,
	0x52, 0x24, 0x9e, 0x8c, 0xd7, 0x43, 0x72, 0x5b, 0x50, 0x6f, 0xf6, 0x17, 0xc7, 0xfe, 0x7e, 0xd8,
	0xd2, 0x17, 0x75, 0x51, 0x39, 0x6f, 0xa3, 0xc5, 0x85, 0x77, 0x55, 0xf6, 0x60, 0xac, 0x35, 0x1f,
	0x74, 0x62, 0x74, 0xe9, 0xdf, 0x4f, 0xb1, 0x69, 0x29, 0x44, 0xdd, 0xba, 0x61, 0x2a, 0xbf, 0x85,
	0xdd, 0x8b, 0xed, 0x42, 0xdf, 0x52, 0x49, 0x9f, 0x3d, 0x85, 0x98, 0xff, 0x08, 0xb8, 0xf9, 0x43,
	0xc1, 0xd9, 0x2e, 0x10, 0xee, 0x20, 0x69, 0x6a, 0x29, 0x0e, 0xe2, 0xb8, 0x29, 0x93, 0xa6, 0xc6,
	0x7b, 0xd8, 0xf0, 0xa9, 0xf3, 0xb7, 0x23, 0x99, 0x30, 0xfe, 0x07, 0xb8, 0x87, 0x95, 0x36, 0xe6,
	0xad, 0x96, 0x0b, 0x7e, 0x19, 0x1a, 0x94, 0xb0, 0xd6, 0xc6, 0xf0, 0x8f, 0x25, 0xf3, 0xa9, 0xc5,
	0x02, 0xd2, 0xca, 0x93, 0x8e, 0x24, 0x57, 0x07, 0x71, 0xdc, 0x16, 0x99, 0x1a, 0xac, 0xd5, 0x64,
	0xad, 0xce, 0x93, 0x75, 0x39, 0x4e, 0x22, 0xc2, 0xb2, 0xd6, 0x51, 0xcb, 0x94, 0x57, 0x71, 0x5d,
	0xbc, 0xc2, 0xf5, 0x28, 0xee, 0xf1, 0x19, 0xd6, 0x63, 0x8d, 0x77, 0xca, 0xbb, 0x4a, 0xcd, 0x53,
	0x66, 0xfb, 0x39, 0x1c, 0x72, 0xe6, 0x57, 0x8f, 0xe2, 0x92, 0xf2, 0xd5, 0xa7, 0xdf, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x9f, 0x9c, 0xe7, 0x31, 0x4d, 0x01, 0x00, 0x00,
}