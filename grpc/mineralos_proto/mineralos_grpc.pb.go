// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package mineralos_proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MineralosClient is the client API for Mineralos service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MineralosClient interface {
	// Sends a greeting
	ReportStatus(ctx context.Context, in *Payload, opts ...grpc.CallOption) (*ServerReply, error)
}

type mineralosClient struct {
	cc grpc.ClientConnInterface
}

func NewMineralosClient(cc grpc.ClientConnInterface) MineralosClient {
	return &mineralosClient{cc}
}

func (c *mineralosClient) ReportStatus(ctx context.Context, in *Payload, opts ...grpc.CallOption) (*ServerReply, error) {
	out := new(ServerReply)
	err := c.cc.Invoke(ctx, "/mineralos_proto.Mineralos/ReportStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MineralosServer is the server API for Mineralos service.
// All implementations must embed UnimplementedMineralosServer
// for forward compatibility
type MineralosServer interface {
	// Sends a greeting
	ReportStatus(context.Context, *Payload) (*ServerReply, error)
	mustEmbedUnimplementedMineralosServer()
}

// UnimplementedMineralosServer must be embedded to have forward compatible implementations.
type UnimplementedMineralosServer struct {
}

func (UnimplementedMineralosServer) ReportStatus(context.Context, *Payload) (*ServerReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportStatus not implemented")
}
func (UnimplementedMineralosServer) mustEmbedUnimplementedMineralosServer() {}

// UnsafeMineralosServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MineralosServer will
// result in compilation errors.
type UnsafeMineralosServer interface {
	mustEmbedUnimplementedMineralosServer()
}

func RegisterMineralosServer(s grpc.ServiceRegistrar, srv MineralosServer) {
	s.RegisterService(&Mineralos_ServiceDesc, srv)
}

func _Mineralos_ReportStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Payload)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MineralosServer).ReportStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mineralos_proto.Mineralos/ReportStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MineralosServer).ReportStatus(ctx, req.(*Payload))
	}
	return interceptor(ctx, in, info, handler)
}

// Mineralos_ServiceDesc is the grpc.ServiceDesc for Mineralos service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Mineralos_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mineralos_proto.Mineralos",
	HandlerType: (*MineralosServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReportStatus",
			Handler:    _Mineralos_ReportStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mineralos_proto/mineralos.proto",
}