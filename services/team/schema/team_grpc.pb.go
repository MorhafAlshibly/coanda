// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.3
// source: team.proto

package schema

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

// TeamServiceClient is the client API for TeamService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TeamServiceClient interface {
	GetTeam(ctx context.Context, in *GetTeamRequest, opts ...grpc.CallOption) (*Team, error)
	GetTeams(ctx context.Context, in *GetTeamsRequest, opts ...grpc.CallOption) (*Teams, error)
	SearchTeams(ctx context.Context, in *SearchTeamsRequest, opts ...grpc.CallOption) (*Teams, error)
	CreateTeam(ctx context.Context, in *CreateTeamRequest, opts ...grpc.CallOption) (*BoolResponse, error)
	UpdateTeamData(ctx context.Context, in *UpdateTeamDataRequest, opts ...grpc.CallOption) (*Team, error)
	UpdateTeamScore(ctx context.Context, in *UpdateTeamScoreRequest, opts ...grpc.CallOption) (*Team, error)
	DeleteTeam(ctx context.Context, in *DeleteTeamRequest, opts ...grpc.CallOption) (*BoolResponse, error)
	JoinTeam(ctx context.Context, in *JoinTeamRequest, opts ...grpc.CallOption) (*BoolResponse, error)
	LeaveTeam(ctx context.Context, in *LeaveTeamRequest, opts ...grpc.CallOption) (*BoolResponse, error)
}

type teamServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTeamServiceClient(cc grpc.ClientConnInterface) TeamServiceClient {
	return &teamServiceClient{cc}
}

func (c *teamServiceClient) GetTeam(ctx context.Context, in *GetTeamRequest, opts ...grpc.CallOption) (*Team, error) {
	out := new(Team)
	err := c.cc.Invoke(ctx, "/schema.TeamService/GetTeam", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServiceClient) GetTeams(ctx context.Context, in *GetTeamsRequest, opts ...grpc.CallOption) (*Teams, error) {
	out := new(Teams)
	err := c.cc.Invoke(ctx, "/schema.TeamService/GetTeams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServiceClient) SearchTeams(ctx context.Context, in *SearchTeamsRequest, opts ...grpc.CallOption) (*Teams, error) {
	out := new(Teams)
	err := c.cc.Invoke(ctx, "/schema.TeamService/SearchTeams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServiceClient) CreateTeam(ctx context.Context, in *CreateTeamRequest, opts ...grpc.CallOption) (*BoolResponse, error) {
	out := new(BoolResponse)
	err := c.cc.Invoke(ctx, "/schema.TeamService/CreateTeam", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServiceClient) UpdateTeamData(ctx context.Context, in *UpdateTeamDataRequest, opts ...grpc.CallOption) (*Team, error) {
	out := new(Team)
	err := c.cc.Invoke(ctx, "/schema.TeamService/UpdateTeamData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServiceClient) UpdateTeamScore(ctx context.Context, in *UpdateTeamScoreRequest, opts ...grpc.CallOption) (*Team, error) {
	out := new(Team)
	err := c.cc.Invoke(ctx, "/schema.TeamService/UpdateTeamScore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServiceClient) DeleteTeam(ctx context.Context, in *DeleteTeamRequest, opts ...grpc.CallOption) (*BoolResponse, error) {
	out := new(BoolResponse)
	err := c.cc.Invoke(ctx, "/schema.TeamService/DeleteTeam", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServiceClient) JoinTeam(ctx context.Context, in *JoinTeamRequest, opts ...grpc.CallOption) (*BoolResponse, error) {
	out := new(BoolResponse)
	err := c.cc.Invoke(ctx, "/schema.TeamService/JoinTeam", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServiceClient) LeaveTeam(ctx context.Context, in *LeaveTeamRequest, opts ...grpc.CallOption) (*BoolResponse, error) {
	out := new(BoolResponse)
	err := c.cc.Invoke(ctx, "/schema.TeamService/LeaveTeam", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TeamServiceServer is the server API for TeamService service.
// All implementations must embed UnimplementedTeamServiceServer
// for forward compatibility
type TeamServiceServer interface {
	GetTeam(context.Context, *GetTeamRequest) (*Team, error)
	GetTeams(context.Context, *GetTeamsRequest) (*Teams, error)
	SearchTeams(context.Context, *SearchTeamsRequest) (*Teams, error)
	CreateTeam(context.Context, *CreateTeamRequest) (*BoolResponse, error)
	UpdateTeamData(context.Context, *UpdateTeamDataRequest) (*Team, error)
	UpdateTeamScore(context.Context, *UpdateTeamScoreRequest) (*Team, error)
	DeleteTeam(context.Context, *DeleteTeamRequest) (*BoolResponse, error)
	JoinTeam(context.Context, *JoinTeamRequest) (*BoolResponse, error)
	LeaveTeam(context.Context, *LeaveTeamRequest) (*BoolResponse, error)
	mustEmbedUnimplementedTeamServiceServer()
}

// UnimplementedTeamServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTeamServiceServer struct {
}

func (UnimplementedTeamServiceServer) GetTeam(context.Context, *GetTeamRequest) (*Team, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTeam not implemented")
}
func (UnimplementedTeamServiceServer) GetTeams(context.Context, *GetTeamsRequest) (*Teams, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTeams not implemented")
}
func (UnimplementedTeamServiceServer) SearchTeams(context.Context, *SearchTeamsRequest) (*Teams, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchTeams not implemented")
}
func (UnimplementedTeamServiceServer) CreateTeam(context.Context, *CreateTeamRequest) (*BoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTeam not implemented")
}
func (UnimplementedTeamServiceServer) UpdateTeamData(context.Context, *UpdateTeamDataRequest) (*Team, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTeamData not implemented")
}
func (UnimplementedTeamServiceServer) UpdateTeamScore(context.Context, *UpdateTeamScoreRequest) (*Team, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTeamScore not implemented")
}
func (UnimplementedTeamServiceServer) DeleteTeam(context.Context, *DeleteTeamRequest) (*BoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTeam not implemented")
}
func (UnimplementedTeamServiceServer) JoinTeam(context.Context, *JoinTeamRequest) (*BoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinTeam not implemented")
}
func (UnimplementedTeamServiceServer) LeaveTeam(context.Context, *LeaveTeamRequest) (*BoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveTeam not implemented")
}
func (UnimplementedTeamServiceServer) mustEmbedUnimplementedTeamServiceServer() {}

// UnsafeTeamServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TeamServiceServer will
// result in compilation errors.
type UnsafeTeamServiceServer interface {
	mustEmbedUnimplementedTeamServiceServer()
}

func RegisterTeamServiceServer(s grpc.ServiceRegistrar, srv TeamServiceServer) {
	s.RegisterService(&TeamService_ServiceDesc, srv)
}

func _TeamService_GetTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServiceServer).GetTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.TeamService/GetTeam",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServiceServer).GetTeam(ctx, req.(*GetTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamService_GetTeams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTeamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServiceServer).GetTeams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.TeamService/GetTeams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServiceServer).GetTeams(ctx, req.(*GetTeamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamService_SearchTeams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchTeamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServiceServer).SearchTeams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.TeamService/SearchTeams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServiceServer).SearchTeams(ctx, req.(*SearchTeamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamService_CreateTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServiceServer).CreateTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.TeamService/CreateTeam",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServiceServer).CreateTeam(ctx, req.(*CreateTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamService_UpdateTeamData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTeamDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServiceServer).UpdateTeamData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.TeamService/UpdateTeamData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServiceServer).UpdateTeamData(ctx, req.(*UpdateTeamDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamService_UpdateTeamScore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTeamScoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServiceServer).UpdateTeamScore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.TeamService/UpdateTeamScore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServiceServer).UpdateTeamScore(ctx, req.(*UpdateTeamScoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamService_DeleteTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServiceServer).DeleteTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.TeamService/DeleteTeam",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServiceServer).DeleteTeam(ctx, req.(*DeleteTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamService_JoinTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServiceServer).JoinTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.TeamService/JoinTeam",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServiceServer).JoinTeam(ctx, req.(*JoinTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamService_LeaveTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServiceServer).LeaveTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.TeamService/LeaveTeam",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServiceServer).LeaveTeam(ctx, req.(*LeaveTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TeamService_ServiceDesc is the grpc.ServiceDesc for TeamService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TeamService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "schema.TeamService",
	HandlerType: (*TeamServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTeam",
			Handler:    _TeamService_GetTeam_Handler,
		},
		{
			MethodName: "GetTeams",
			Handler:    _TeamService_GetTeams_Handler,
		},
		{
			MethodName: "SearchTeams",
			Handler:    _TeamService_SearchTeams_Handler,
		},
		{
			MethodName: "CreateTeam",
			Handler:    _TeamService_CreateTeam_Handler,
		},
		{
			MethodName: "UpdateTeamData",
			Handler:    _TeamService_UpdateTeamData_Handler,
		},
		{
			MethodName: "UpdateTeamScore",
			Handler:    _TeamService_UpdateTeamScore_Handler,
		},
		{
			MethodName: "DeleteTeam",
			Handler:    _TeamService_DeleteTeam_Handler,
		},
		{
			MethodName: "JoinTeam",
			Handler:    _TeamService_JoinTeam_Handler,
		},
		{
			MethodName: "LeaveTeam",
			Handler:    _TeamService_LeaveTeam_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "team.proto",
}
