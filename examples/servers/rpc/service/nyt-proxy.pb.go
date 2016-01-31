// Code generated by protoc-gen-go.
// source: nyt-proxy.proto
// DO NOT EDIT!

/*
Package service is a generated protocol buffer package.

It is generated from these files:
	nyt-proxy.proto

It has these top-level messages:
	MostPopularRequest
	MostPopularResponse
	CatsRequest
	CatsResponse
*/
package service

import (
	"fmt"
	"math"

	"github.com/xhroot/gizmo/examples/nyt"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type MostPopularRequest struct {
	ResourceType   string `protobuf:"bytes,1,opt,name=resourceType" json:"resourceType,omitempty"`
	Section        string `protobuf:"bytes,2,opt,name=section" json:"section,omitempty"`
	TimePeriodDays uint32 `protobuf:"varint,3,opt,name=timePeriodDays" json:"timePeriodDays,omitempty"`
}

func (m *MostPopularRequest) Reset()                    { *m = MostPopularRequest{} }
func (m *MostPopularRequest) String() string            { return proto.CompactTextString(m) }
func (*MostPopularRequest) ProtoMessage()               {}
func (*MostPopularRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type MostPopularResponse struct {
	Result []*nyt.MostPopularResult `protobuf:"bytes,1,rep,name=result" json:"result,omitempty"`
}

func (m *MostPopularResponse) Reset()                    { *m = MostPopularResponse{} }
func (m *MostPopularResponse) String() string            { return proto.CompactTextString(m) }
func (*MostPopularResponse) ProtoMessage()               {}
func (*MostPopularResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *MostPopularResponse) GetResult() []*nyt.MostPopularResult {
	if m != nil {
		return m.Result
	}
	return nil
}

type CatsRequest struct {
}

func (m *CatsRequest) Reset()                    { *m = CatsRequest{} }
func (m *CatsRequest) String() string            { return proto.CompactTextString(m) }
func (*CatsRequest) ProtoMessage()               {}
func (*CatsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type CatsResponse struct {
	Results []*nyt.SemanticConceptArticle `protobuf:"bytes,1,rep,name=results" json:"results,omitempty"`
}

func (m *CatsResponse) Reset()                    { *m = CatsResponse{} }
func (m *CatsResponse) String() string            { return proto.CompactTextString(m) }
func (*CatsResponse) ProtoMessage()               {}
func (*CatsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *CatsResponse) GetResults() []*nyt.SemanticConceptArticle {
	if m != nil {
		return m.Results
	}
	return nil
}

func init() {
	proto.RegisterType((*MostPopularRequest)(nil), "service.MostPopularRequest")
	proto.RegisterType((*MostPopularResponse)(nil), "service.MostPopularResponse")
	proto.RegisterType((*CatsRequest)(nil), "service.CatsRequest")
	proto.RegisterType((*CatsResponse)(nil), "service.CatsResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for NYTProxyService service

type NYTProxyServiceClient interface {
	GetMostPopular(ctx context.Context, in *MostPopularRequest, opts ...grpc.CallOption) (*MostPopularResponse, error)
	GetCats(ctx context.Context, in *CatsRequest, opts ...grpc.CallOption) (*CatsResponse, error)
}

type nYTProxyServiceClient struct {
	cc *grpc.ClientConn
}

func NewNYTProxyServiceClient(cc *grpc.ClientConn) NYTProxyServiceClient {
	return &nYTProxyServiceClient{cc}
}

func (c *nYTProxyServiceClient) GetMostPopular(ctx context.Context, in *MostPopularRequest, opts ...grpc.CallOption) (*MostPopularResponse, error) {
	out := new(MostPopularResponse)
	err := grpc.Invoke(ctx, "/service.NYTProxyService/GetMostPopular", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nYTProxyServiceClient) GetCats(ctx context.Context, in *CatsRequest, opts ...grpc.CallOption) (*CatsResponse, error) {
	out := new(CatsResponse)
	err := grpc.Invoke(ctx, "/service.NYTProxyService/GetCats", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for NYTProxyService service

type NYTProxyServiceServer interface {
	GetMostPopular(context.Context, *MostPopularRequest) (*MostPopularResponse, error)
	GetCats(context.Context, *CatsRequest) (*CatsResponse, error)
}

func RegisterNYTProxyServiceServer(s *grpc.Server, srv NYTProxyServiceServer) {
	s.RegisterService(&NYTProxyService_serviceDesc, srv)
}

func _NYTProxyService_GetMostPopular_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(MostPopularRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(NYTProxyServiceServer).GetMostPopular(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _NYTProxyService_GetCats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(CatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(NYTProxyServiceServer).GetCats(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var NYTProxyService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "service.NYTProxyService",
	HandlerType: (*NYTProxyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMostPopular",
			Handler:    _NYTProxyService_GetMostPopular_Handler,
		},
		{
			MethodName: "GetCats",
			Handler:    _NYTProxyService_GetCats_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 298 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x91, 0xcf, 0x4a, 0xf3, 0x40,
	0x14, 0xc5, 0xc9, 0x57, 0x68, 0xf8, 0x6e, 0xff, 0x04, 0xc6, 0xb6, 0xc4, 0xd4, 0x85, 0x64, 0x21,
	0x5d, 0x68, 0x84, 0x0a, 0xae, 0x74, 0x21, 0x15, 0x5c, 0x88, 0x12, 0x4c, 0x37, 0x2e, 0xe3, 0x78,
	0x17, 0x81, 0x24, 0x33, 0xce, 0x4c, 0xc4, 0xbc, 0x88, 0xcf, 0xeb, 0x4d, 0x66, 0x90, 0x56, 0xdd,
	0x84, 0xe4, 0x9e, 0x93, 0xdf, 0x9c, 0x7b, 0x06, 0x82, 0xba, 0x35, 0x67, 0x52, 0x89, 0x8f, 0x36,
	0xa1, 0xa7, 0x11, 0xcc, 0xd7, 0xa8, 0xde, 0x0b, 0x8e, 0xd1, 0x9c, 0x94, 0xf3, 0x4a, 0x68, 0x23,
	0x85, 0x6c, 0xca, 0x5c, 0x59, 0x3d, 0x3a, 0xec, 0xc6, 0x1a, 0xab, 0xbc, 0x36, 0x05, 0xe7, 0xa2,
	0xe6, 0x28, 0x8d, 0x95, 0xe2, 0x0c, 0xd8, 0x03, 0xf9, 0x53, 0xeb, 0x7f, 0xc2, 0xb7, 0x06, 0xb5,
	0x61, 0x33, 0x18, 0x2b, 0xd4, 0xa2, 0x51, 0x1c, 0xb7, 0xad, 0xc4, 0xd0, 0x3b, 0xf6, 0x56, 0xff,
	0x59, 0x00, 0x74, 0x10, 0x37, 0x85, 0xa8, 0xc3, 0x7f, 0xfd, 0x60, 0x01, 0x53, 0x53, 0x54, 0x98,
	0xa2, 0x2a, 0xc4, 0xeb, 0x6d, 0xde, 0xea, 0x70, 0x40, 0xf3, 0x49, 0x7c, 0x0d, 0x07, 0x7b, 0x50,
	0x2d, 0x45, 0xad, 0x91, 0x9d, 0xc0, 0x90, 0xa8, 0x4d, 0x69, 0x88, 0x37, 0x58, 0x8d, 0xd6, 0x8b,
	0x84, 0x72, 0x25, 0xfb, 0x4e, 0x52, 0xe3, 0x09, 0x8c, 0x36, 0xb9, 0xd1, 0x2e, 0x4c, 0x7c, 0x05,
	0x63, 0xfb, 0xe9, 0x30, 0xa7, 0xe0, 0x5b, 0x8c, 0x76, 0x9c, 0x65, 0xcf, 0xc9, 0xdc, 0x7e, 0x1b,
	0xbb, 0xdf, 0x8d, 0xa2, 0xf7, 0x12, 0xd7, 0x9f, 0x1e, 0x04, 0x8f, 0xcf, 0xdb, 0xb4, 0xab, 0x2b,
	0xb3, 0x35, 0xb1, 0x7b, 0x98, 0xde, 0xa1, 0xd9, 0x39, 0x98, 0x2d, 0x13, 0x57, 0x61, 0xf2, 0xbb,
	0x8d, 0xe8, 0xe8, 0x6f, 0xd1, 0xc5, 0xb9, 0x04, 0x9f, 0x60, 0x5d, 0x42, 0x36, 0xfb, 0x36, 0xee,
	0xe4, 0x8f, 0xe6, 0x3f, 0xa6, 0xf6, 0xbf, 0x97, 0x61, 0x7f, 0x01, 0x17, 0x5f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0xb7, 0xcc, 0xf9, 0xf3, 0xce, 0x01, 0x00, 0x00,
}
