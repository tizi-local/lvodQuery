// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package tizi_local_proto_lvodQuery

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// VodQueryServiceClient is the client API for VodQueryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VodQueryServiceClient interface {
	FeedQuery(ctx context.Context, in *FeedQueryReq, opts ...grpc.CallOption) (*FeedQueryResp, error)
}

type vodQueryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVodQueryServiceClient(cc grpc.ClientConnInterface) VodQueryServiceClient {
	return &vodQueryServiceClient{cc}
}

func (c *vodQueryServiceClient) FeedQuery(ctx context.Context, in *FeedQueryReq, opts ...grpc.CallOption) (*FeedQueryResp, error) {
	out := new(FeedQueryResp)
	err := c.cc.Invoke(ctx, "/tizi.local.lvodquery.VodQueryService/FeedQuery", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VodQueryServiceServer is the server API for VodQueryService service.
// All implementations must embed UnimplementedVodQueryServiceServer
// for forward compatibility
type VodQueryServiceServer interface {
	FeedQuery(context.Context, *FeedQueryReq) (*FeedQueryResp, error)
	mustEmbedUnimplementedVodQueryServiceServer()
}

// UnimplementedVodQueryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedVodQueryServiceServer struct {
}

func (UnimplementedVodQueryServiceServer) FeedQuery(context.Context, *FeedQueryReq) (*FeedQueryResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FeedQuery not implemented")
}
func (UnimplementedVodQueryServiceServer) mustEmbedUnimplementedVodQueryServiceServer() {}

// UnsafeVodQueryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VodQueryServiceServer will
// result in compilation errors.
type UnsafeVodQueryServiceServer interface {
	mustEmbedUnimplementedVodQueryServiceServer()
}

func RegisterVodQueryServiceServer(s grpc.ServiceRegistrar, srv VodQueryServiceServer) {
	s.RegisterService(&_VodQueryService_serviceDesc, srv)
}

func _VodQueryService_FeedQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FeedQueryReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VodQueryServiceServer).FeedQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tizi.local.lvodquery.VodQueryService/FeedQuery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VodQueryServiceServer).FeedQuery(ctx, req.(*FeedQueryReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _VodQueryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tizi.local.lvodquery.VodQueryService",
	HandlerType: (*VodQueryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FeedQuery",
			Handler:    _VodQueryService_FeedQuery_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "vodQuery/vodquery.proto",
}
