// Code generated by protoc-gen-go.
// source: contracts.proto
// DO NOT EDIT!

/*
Package deployment is a generated protocol buffer package.

It is generated from these files:
	contracts.proto

It has these top-level messages:
	DeployRequest
	Response
*/
package deployment

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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

type DeployRequest struct {
	Token       string            `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
	BotId       string            `protobuf:"bytes,2,opt,name=botId" json:"botId,omitempty"`
	Environment map[string]string `protobuf:"bytes,3,rep,name=environment" json:"environment,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *DeployRequest) Reset()                    { *m = DeployRequest{} }
func (m *DeployRequest) String() string            { return proto.CompactTextString(m) }
func (*DeployRequest) ProtoMessage()               {}
func (*DeployRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DeployRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *DeployRequest) GetBotId() string {
	if m != nil {
		return m.BotId
	}
	return ""
}

func (m *DeployRequest) GetEnvironment() map[string]string {
	if m != nil {
		return m.Environment
	}
	return nil
}

type Response struct {
	Code    int32  `protobuf:"varint,1,opt,name=code" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Response) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Response) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*DeployRequest)(nil), "deployment.DeployRequest")
	proto.RegisterType((*Response)(nil), "deployment.Response")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Deployer service

type DeployerClient interface {
	Deploy(ctx context.Context, in *DeployRequest, opts ...grpc.CallOption) (Deployer_DeployClient, error)
}

type deployerClient struct {
	cc *grpc.ClientConn
}

func NewDeployerClient(cc *grpc.ClientConn) DeployerClient {
	return &deployerClient{cc}
}

func (c *deployerClient) Deploy(ctx context.Context, in *DeployRequest, opts ...grpc.CallOption) (Deployer_DeployClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Deployer_serviceDesc.Streams[0], c.cc, "/deployment.Deployer/Deploy", opts...)
	if err != nil {
		return nil, err
	}
	x := &deployerDeployClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Deployer_DeployClient interface {
	Recv() (*Response, error)
	grpc.ClientStream
}

type deployerDeployClient struct {
	grpc.ClientStream
}

func (x *deployerDeployClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Deployer service

type DeployerServer interface {
	Deploy(*DeployRequest, Deployer_DeployServer) error
}

func RegisterDeployerServer(s *grpc.Server, srv DeployerServer) {
	s.RegisterService(&_Deployer_serviceDesc, srv)
}

func _Deployer_Deploy_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DeployRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DeployerServer).Deploy(m, &deployerDeployServer{stream})
}

type Deployer_DeployServer interface {
	Send(*Response) error
	grpc.ServerStream
}

type deployerDeployServer struct {
	grpc.ServerStream
}

func (x *deployerDeployServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

var _Deployer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "deployment.Deployer",
	HandlerType: (*DeployerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Deploy",
			Handler:       _Deployer_Deploy_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "contracts.proto",
}

func init() { proto.RegisterFile("contracts.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 240 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x90, 0x41, 0x4b, 0x03, 0x31,
	0x10, 0x85, 0x49, 0xd7, 0xd6, 0x3a, 0x45, 0x2c, 0x43, 0x0f, 0x6b, 0x4f, 0xa5, 0xa7, 0xe2, 0x61,
	0x91, 0x7a, 0x29, 0x0a, 0x9e, 0x2c, 0x22, 0x78, 0xca, 0x3f, 0xd8, 0xee, 0x0e, 0x22, 0x6d, 0x33,
	0x6b, 0x32, 0x2d, 0xec, 0x4f, 0xf4, 0x5f, 0x49, 0x12, 0x43, 0x57, 0xa1, 0xb7, 0xf7, 0x4d, 0xe6,
	0xbd, 0xbc, 0x04, 0x6e, 0x2a, 0x36, 0x62, 0xcb, 0x4a, 0x5c, 0xd1, 0x58, 0x16, 0x46, 0xa8, 0xa9,
	0xd9, 0x71, 0xbb, 0x27, 0x23, 0xf3, 0x6f, 0x05, 0xd7, 0x2f, 0x01, 0x35, 0x7d, 0x1d, 0xc8, 0x09,
	0x4e, 0xa0, 0x2f, 0xbc, 0x25, 0x93, 0xab, 0x99, 0x5a, 0x5c, 0xe9, 0x08, 0x7e, 0xba, 0x61, 0x79,
	0xab, 0xf3, 0x5e, 0x9c, 0x06, 0xc0, 0x77, 0x18, 0x91, 0x39, 0x7e, 0x5a, 0x36, 0x3e, 0x2c, 0xcf,
	0x66, 0xd9, 0x62, 0xb4, 0xbc, 0x2b, 0x4e, 0xf9, 0xc5, 0x9f, 0xec, 0x62, 0x7d, 0x5a, 0x5e, 0x1b,
	0xb1, 0xad, 0xee, 0xda, 0xa7, 0xcf, 0x30, 0xfe, 0xbf, 0x80, 0x63, 0xc8, 0xb6, 0xd4, 0xfe, 0x76,
	0xf1, 0xd2, 0x37, 0x39, 0x96, 0xbb, 0x03, 0xa5, 0x26, 0x01, 0x1e, 0x7b, 0x2b, 0x35, 0x5f, 0xc1,
	0x50, 0x93, 0x6b, 0xd8, 0x38, 0x42, 0x84, 0x8b, 0x8a, 0x6b, 0x0a, 0xc6, 0xbe, 0x0e, 0x1a, 0x73,
	0xb8, 0xdc, 0x93, 0x73, 0xe5, 0x47, 0xf2, 0x26, 0x5c, 0xbe, 0xc2, 0x30, 0x16, 0x25, 0x8b, 0x4f,
	0x30, 0x88, 0x1a, 0x6f, 0xcf, 0x3e, 0x64, 0x3a, 0xe9, 0x1e, 0xa5, 0x4b, 0xef, 0xd5, 0x66, 0x10,
	0x7e, 0xf8, 0xe1, 0x27, 0x00, 0x00, 0xff, 0xff, 0x74, 0xc4, 0x90, 0x80, 0x74, 0x01, 0x00, 0x00,
}
