// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.3
// source: record.proto

package api

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

// RecordServiceClient is the client API for RecordService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RecordServiceClient interface {
	CreateRecord(ctx context.Context, in *CreateRecordRequest, opts ...grpc.CallOption) (*CreateRecordResponse, error)
	GetRecord(ctx context.Context, in *GetRecordRequest, opts ...grpc.CallOption) (*GetRecordResponse, error)
	GetRecords(ctx context.Context, in *GetRecordsRequest, opts ...grpc.CallOption) (*GetRecordsResponse, error)
	DeleteRecord(ctx context.Context, in *DeleteRecordRequest, opts ...grpc.CallOption) (*DeleteRecordResponse, error)
}

type recordServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRecordServiceClient(cc grpc.ClientConnInterface) RecordServiceClient {
	return &recordServiceClient{cc}
}

func (c *recordServiceClient) CreateRecord(ctx context.Context, in *CreateRecordRequest, opts ...grpc.CallOption) (*CreateRecordResponse, error) {
	out := new(CreateRecordResponse)
	err := c.cc.Invoke(ctx, "/api.RecordService/CreateRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recordServiceClient) GetRecord(ctx context.Context, in *GetRecordRequest, opts ...grpc.CallOption) (*GetRecordResponse, error) {
	out := new(GetRecordResponse)
	err := c.cc.Invoke(ctx, "/api.RecordService/GetRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recordServiceClient) GetRecords(ctx context.Context, in *GetRecordsRequest, opts ...grpc.CallOption) (*GetRecordsResponse, error) {
	out := new(GetRecordsResponse)
	err := c.cc.Invoke(ctx, "/api.RecordService/GetRecords", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recordServiceClient) DeleteRecord(ctx context.Context, in *DeleteRecordRequest, opts ...grpc.CallOption) (*DeleteRecordResponse, error) {
	out := new(DeleteRecordResponse)
	err := c.cc.Invoke(ctx, "/api.RecordService/DeleteRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecordServiceServer is the server API for RecordService service.
// All implementations must embed UnimplementedRecordServiceServer
// for forward compatibility
type RecordServiceServer interface {
	CreateRecord(context.Context, *CreateRecordRequest) (*CreateRecordResponse, error)
	GetRecord(context.Context, *GetRecordRequest) (*GetRecordResponse, error)
	GetRecords(context.Context, *GetRecordsRequest) (*GetRecordsResponse, error)
	DeleteRecord(context.Context, *DeleteRecordRequest) (*DeleteRecordResponse, error)
	mustEmbedUnimplementedRecordServiceServer()
}

// UnimplementedRecordServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRecordServiceServer struct {
}

func (UnimplementedRecordServiceServer) CreateRecord(context.Context, *CreateRecordRequest) (*CreateRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRecord not implemented")
}
func (UnimplementedRecordServiceServer) GetRecord(context.Context, *GetRecordRequest) (*GetRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRecord not implemented")
}
func (UnimplementedRecordServiceServer) GetRecords(context.Context, *GetRecordsRequest) (*GetRecordsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRecords not implemented")
}
func (UnimplementedRecordServiceServer) DeleteRecord(context.Context, *DeleteRecordRequest) (*DeleteRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRecord not implemented")
}
func (UnimplementedRecordServiceServer) mustEmbedUnimplementedRecordServiceServer() {}

// UnsafeRecordServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RecordServiceServer will
// result in compilation errors.
type UnsafeRecordServiceServer interface {
	mustEmbedUnimplementedRecordServiceServer()
}

func RegisterRecordServiceServer(s grpc.ServiceRegistrar, srv RecordServiceServer) {
	s.RegisterService(&RecordService_ServiceDesc, srv)
}

func _RecordService_CreateRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordServiceServer).CreateRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.RecordService/CreateRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordServiceServer).CreateRecord(ctx, req.(*CreateRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecordService_GetRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordServiceServer).GetRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.RecordService/GetRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordServiceServer).GetRecord(ctx, req.(*GetRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecordService_GetRecords_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRecordsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordServiceServer).GetRecords(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.RecordService/GetRecords",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordServiceServer).GetRecords(ctx, req.(*GetRecordsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecordService_DeleteRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordServiceServer).DeleteRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.RecordService/DeleteRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordServiceServer).DeleteRecord(ctx, req.(*DeleteRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RecordService_ServiceDesc is the grpc.ServiceDesc for RecordService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RecordService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.RecordService",
	HandlerType: (*RecordServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRecord",
			Handler:    _RecordService_CreateRecord_Handler,
		},
		{
			MethodName: "GetRecord",
			Handler:    _RecordService_GetRecord_Handler,
		},
		{
			MethodName: "GetRecords",
			Handler:    _RecordService_GetRecords_Handler,
		},
		{
			MethodName: "DeleteRecord",
			Handler:    _RecordService_DeleteRecord_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "record.proto",
}
