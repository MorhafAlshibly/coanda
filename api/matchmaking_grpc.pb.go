// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.3
// source: matchmaking.proto

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

// MatchmakingServiceClient is the client API for MatchmakingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MatchmakingServiceClient interface {
	CreateArena(ctx context.Context, in *CreateArenaRequest, opts ...grpc.CallOption) (*CreateArenaResponse, error)
	GetArena(ctx context.Context, in *ArenaRequest, opts ...grpc.CallOption) (*GetArenaResponse, error)
	GetArenas(ctx context.Context, in *GetArenasRequest, opts ...grpc.CallOption) (*GetArenasResponse, error)
	UpdateArena(ctx context.Context, in *UpdateArenaRequest, opts ...grpc.CallOption) (*UpdateArenaResponse, error)
	DeleteArena(ctx context.Context, in *ArenaRequest, opts ...grpc.CallOption) (*ArenaResponse, error)
	CreateMatchmakingUser(ctx context.Context, in *CreateMatchmakingUserRequest, opts ...grpc.CallOption) (*CreateMatchmakingUserResponse, error)
	GetMatchmakingUser(ctx context.Context, in *MatchmakingUserRequest, opts ...grpc.CallOption) (*GetMatchmakingUserResponse, error)
	GetMatchmakingUsers(ctx context.Context, in *GetMatchmakingUsersRequest, opts ...grpc.CallOption) (*GetMatchmakingUsersResponse, error)
	UpdateMatchmakingUser(ctx context.Context, in *UpdateMatchmakingUserRequest, opts ...grpc.CallOption) (*UpdateMatchmakingUserResponse, error)
	SetMatchmakingUserElo(ctx context.Context, in *SetMatchmakingUserEloRequest, opts ...grpc.CallOption) (*SetMatchmakingUserEloResponse, error)
	CreateMatchmakingTicket(ctx context.Context, in *CreateMatchmakingTicketRequest, opts ...grpc.CallOption) (*CreateMatchmakingTicketResponse, error)
	PollMatchmakingTicket(ctx context.Context, in *MatchmakingTicketRequest, opts ...grpc.CallOption) (*MatchmakingTicketResponse, error)
	GetMatchmakingTicket(ctx context.Context, in *MatchmakingTicketRequest, opts ...grpc.CallOption) (*GetMatchmakingTicketResponse, error)
	GetMatchmakingTickets(ctx context.Context, in *GetMatchmakingTicketsRequest, opts ...grpc.CallOption) (*GetMatchmakingTicketsResponse, error)
	UpdateMatchmakingTicket(ctx context.Context, in *UpdateMatchmakingTicketRequest, opts ...grpc.CallOption) (*UpdateMatchmakingTicketResponse, error)
	ExpireMatchmakingTicket(ctx context.Context, in *MatchmakingTicketRequest, opts ...grpc.CallOption) (*ExpireMatchmakingTicketResponse, error)
	StartMatch(ctx context.Context, in *StartMatchRequest, opts ...grpc.CallOption) (*StartMatchResponse, error)
	EndMatch(ctx context.Context, in *EndMatchRequest, opts ...grpc.CallOption) (*EndMatchResponse, error)
	GetMatch(ctx context.Context, in *MatchRequest, opts ...grpc.CallOption) (*GetMatchResponse, error)
	GetMatches(ctx context.Context, in *GetMatchesRequest, opts ...grpc.CallOption) (*GetMatchesResponse, error)
	UpdateMatch(ctx context.Context, in *UpdateMatchRequest, opts ...grpc.CallOption) (*UpdateMatchResponse, error)
}

type matchmakingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMatchmakingServiceClient(cc grpc.ClientConnInterface) MatchmakingServiceClient {
	return &matchmakingServiceClient{cc}
}

func (c *matchmakingServiceClient) CreateArena(ctx context.Context, in *CreateArenaRequest, opts ...grpc.CallOption) (*CreateArenaResponse, error) {
	out := new(CreateArenaResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/CreateArena", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) GetArena(ctx context.Context, in *ArenaRequest, opts ...grpc.CallOption) (*GetArenaResponse, error) {
	out := new(GetArenaResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/GetArena", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) GetArenas(ctx context.Context, in *GetArenasRequest, opts ...grpc.CallOption) (*GetArenasResponse, error) {
	out := new(GetArenasResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/GetArenas", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) UpdateArena(ctx context.Context, in *UpdateArenaRequest, opts ...grpc.CallOption) (*UpdateArenaResponse, error) {
	out := new(UpdateArenaResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/UpdateArena", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) DeleteArena(ctx context.Context, in *ArenaRequest, opts ...grpc.CallOption) (*ArenaResponse, error) {
	out := new(ArenaResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/DeleteArena", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) CreateMatchmakingUser(ctx context.Context, in *CreateMatchmakingUserRequest, opts ...grpc.CallOption) (*CreateMatchmakingUserResponse, error) {
	out := new(CreateMatchmakingUserResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/CreateMatchmakingUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) GetMatchmakingUser(ctx context.Context, in *MatchmakingUserRequest, opts ...grpc.CallOption) (*GetMatchmakingUserResponse, error) {
	out := new(GetMatchmakingUserResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/GetMatchmakingUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) GetMatchmakingUsers(ctx context.Context, in *GetMatchmakingUsersRequest, opts ...grpc.CallOption) (*GetMatchmakingUsersResponse, error) {
	out := new(GetMatchmakingUsersResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/GetMatchmakingUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) UpdateMatchmakingUser(ctx context.Context, in *UpdateMatchmakingUserRequest, opts ...grpc.CallOption) (*UpdateMatchmakingUserResponse, error) {
	out := new(UpdateMatchmakingUserResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/UpdateMatchmakingUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) SetMatchmakingUserElo(ctx context.Context, in *SetMatchmakingUserEloRequest, opts ...grpc.CallOption) (*SetMatchmakingUserEloResponse, error) {
	out := new(SetMatchmakingUserEloResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/SetMatchmakingUserElo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) CreateMatchmakingTicket(ctx context.Context, in *CreateMatchmakingTicketRequest, opts ...grpc.CallOption) (*CreateMatchmakingTicketResponse, error) {
	out := new(CreateMatchmakingTicketResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/CreateMatchmakingTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) PollMatchmakingTicket(ctx context.Context, in *MatchmakingTicketRequest, opts ...grpc.CallOption) (*MatchmakingTicketResponse, error) {
	out := new(MatchmakingTicketResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/PollMatchmakingTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) GetMatchmakingTicket(ctx context.Context, in *MatchmakingTicketRequest, opts ...grpc.CallOption) (*GetMatchmakingTicketResponse, error) {
	out := new(GetMatchmakingTicketResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/GetMatchmakingTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) GetMatchmakingTickets(ctx context.Context, in *GetMatchmakingTicketsRequest, opts ...grpc.CallOption) (*GetMatchmakingTicketsResponse, error) {
	out := new(GetMatchmakingTicketsResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/GetMatchmakingTickets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) UpdateMatchmakingTicket(ctx context.Context, in *UpdateMatchmakingTicketRequest, opts ...grpc.CallOption) (*UpdateMatchmakingTicketResponse, error) {
	out := new(UpdateMatchmakingTicketResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/UpdateMatchmakingTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) ExpireMatchmakingTicket(ctx context.Context, in *MatchmakingTicketRequest, opts ...grpc.CallOption) (*ExpireMatchmakingTicketResponse, error) {
	out := new(ExpireMatchmakingTicketResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/ExpireMatchmakingTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) StartMatch(ctx context.Context, in *StartMatchRequest, opts ...grpc.CallOption) (*StartMatchResponse, error) {
	out := new(StartMatchResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/StartMatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) EndMatch(ctx context.Context, in *EndMatchRequest, opts ...grpc.CallOption) (*EndMatchResponse, error) {
	out := new(EndMatchResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/EndMatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) GetMatch(ctx context.Context, in *MatchRequest, opts ...grpc.CallOption) (*GetMatchResponse, error) {
	out := new(GetMatchResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/GetMatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) GetMatches(ctx context.Context, in *GetMatchesRequest, opts ...grpc.CallOption) (*GetMatchesResponse, error) {
	out := new(GetMatchesResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/GetMatches", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchmakingServiceClient) UpdateMatch(ctx context.Context, in *UpdateMatchRequest, opts ...grpc.CallOption) (*UpdateMatchResponse, error) {
	out := new(UpdateMatchResponse)
	err := c.cc.Invoke(ctx, "/MatchmakingService/UpdateMatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MatchmakingServiceServer is the server API for MatchmakingService service.
// All implementations must embed UnimplementedMatchmakingServiceServer
// for forward compatibility
type MatchmakingServiceServer interface {
	CreateArena(context.Context, *CreateArenaRequest) (*CreateArenaResponse, error)
	GetArena(context.Context, *ArenaRequest) (*GetArenaResponse, error)
	GetArenas(context.Context, *GetArenasRequest) (*GetArenasResponse, error)
	UpdateArena(context.Context, *UpdateArenaRequest) (*UpdateArenaResponse, error)
	DeleteArena(context.Context, *ArenaRequest) (*ArenaResponse, error)
	CreateMatchmakingUser(context.Context, *CreateMatchmakingUserRequest) (*CreateMatchmakingUserResponse, error)
	GetMatchmakingUser(context.Context, *MatchmakingUserRequest) (*GetMatchmakingUserResponse, error)
	GetMatchmakingUsers(context.Context, *GetMatchmakingUsersRequest) (*GetMatchmakingUsersResponse, error)
	UpdateMatchmakingUser(context.Context, *UpdateMatchmakingUserRequest) (*UpdateMatchmakingUserResponse, error)
	SetMatchmakingUserElo(context.Context, *SetMatchmakingUserEloRequest) (*SetMatchmakingUserEloResponse, error)
	CreateMatchmakingTicket(context.Context, *CreateMatchmakingTicketRequest) (*CreateMatchmakingTicketResponse, error)
	PollMatchmakingTicket(context.Context, *MatchmakingTicketRequest) (*MatchmakingTicketResponse, error)
	GetMatchmakingTicket(context.Context, *MatchmakingTicketRequest) (*GetMatchmakingTicketResponse, error)
	GetMatchmakingTickets(context.Context, *GetMatchmakingTicketsRequest) (*GetMatchmakingTicketsResponse, error)
	UpdateMatchmakingTicket(context.Context, *UpdateMatchmakingTicketRequest) (*UpdateMatchmakingTicketResponse, error)
	ExpireMatchmakingTicket(context.Context, *MatchmakingTicketRequest) (*ExpireMatchmakingTicketResponse, error)
	StartMatch(context.Context, *StartMatchRequest) (*StartMatchResponse, error)
	EndMatch(context.Context, *EndMatchRequest) (*EndMatchResponse, error)
	GetMatch(context.Context, *MatchRequest) (*GetMatchResponse, error)
	GetMatches(context.Context, *GetMatchesRequest) (*GetMatchesResponse, error)
	UpdateMatch(context.Context, *UpdateMatchRequest) (*UpdateMatchResponse, error)
	mustEmbedUnimplementedMatchmakingServiceServer()
}

// UnimplementedMatchmakingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMatchmakingServiceServer struct {
}

func (UnimplementedMatchmakingServiceServer) CreateArena(context.Context, *CreateArenaRequest) (*CreateArenaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateArena not implemented")
}
func (UnimplementedMatchmakingServiceServer) GetArena(context.Context, *ArenaRequest) (*GetArenaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArena not implemented")
}
func (UnimplementedMatchmakingServiceServer) GetArenas(context.Context, *GetArenasRequest) (*GetArenasResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArenas not implemented")
}
func (UnimplementedMatchmakingServiceServer) UpdateArena(context.Context, *UpdateArenaRequest) (*UpdateArenaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateArena not implemented")
}
func (UnimplementedMatchmakingServiceServer) DeleteArena(context.Context, *ArenaRequest) (*ArenaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteArena not implemented")
}
func (UnimplementedMatchmakingServiceServer) CreateMatchmakingUser(context.Context, *CreateMatchmakingUserRequest) (*CreateMatchmakingUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMatchmakingUser not implemented")
}
func (UnimplementedMatchmakingServiceServer) GetMatchmakingUser(context.Context, *MatchmakingUserRequest) (*GetMatchmakingUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMatchmakingUser not implemented")
}
func (UnimplementedMatchmakingServiceServer) GetMatchmakingUsers(context.Context, *GetMatchmakingUsersRequest) (*GetMatchmakingUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMatchmakingUsers not implemented")
}
func (UnimplementedMatchmakingServiceServer) UpdateMatchmakingUser(context.Context, *UpdateMatchmakingUserRequest) (*UpdateMatchmakingUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMatchmakingUser not implemented")
}
func (UnimplementedMatchmakingServiceServer) SetMatchmakingUserElo(context.Context, *SetMatchmakingUserEloRequest) (*SetMatchmakingUserEloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetMatchmakingUserElo not implemented")
}
func (UnimplementedMatchmakingServiceServer) CreateMatchmakingTicket(context.Context, *CreateMatchmakingTicketRequest) (*CreateMatchmakingTicketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMatchmakingTicket not implemented")
}
func (UnimplementedMatchmakingServiceServer) PollMatchmakingTicket(context.Context, *MatchmakingTicketRequest) (*MatchmakingTicketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PollMatchmakingTicket not implemented")
}
func (UnimplementedMatchmakingServiceServer) GetMatchmakingTicket(context.Context, *MatchmakingTicketRequest) (*GetMatchmakingTicketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMatchmakingTicket not implemented")
}
func (UnimplementedMatchmakingServiceServer) GetMatchmakingTickets(context.Context, *GetMatchmakingTicketsRequest) (*GetMatchmakingTicketsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMatchmakingTickets not implemented")
}
func (UnimplementedMatchmakingServiceServer) UpdateMatchmakingTicket(context.Context, *UpdateMatchmakingTicketRequest) (*UpdateMatchmakingTicketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMatchmakingTicket not implemented")
}
func (UnimplementedMatchmakingServiceServer) ExpireMatchmakingTicket(context.Context, *MatchmakingTicketRequest) (*ExpireMatchmakingTicketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExpireMatchmakingTicket not implemented")
}
func (UnimplementedMatchmakingServiceServer) StartMatch(context.Context, *StartMatchRequest) (*StartMatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartMatch not implemented")
}
func (UnimplementedMatchmakingServiceServer) EndMatch(context.Context, *EndMatchRequest) (*EndMatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EndMatch not implemented")
}
func (UnimplementedMatchmakingServiceServer) GetMatch(context.Context, *MatchRequest) (*GetMatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMatch not implemented")
}
func (UnimplementedMatchmakingServiceServer) GetMatches(context.Context, *GetMatchesRequest) (*GetMatchesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMatches not implemented")
}
func (UnimplementedMatchmakingServiceServer) UpdateMatch(context.Context, *UpdateMatchRequest) (*UpdateMatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMatch not implemented")
}
func (UnimplementedMatchmakingServiceServer) mustEmbedUnimplementedMatchmakingServiceServer() {}

// UnsafeMatchmakingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MatchmakingServiceServer will
// result in compilation errors.
type UnsafeMatchmakingServiceServer interface {
	mustEmbedUnimplementedMatchmakingServiceServer()
}

func RegisterMatchmakingServiceServer(s grpc.ServiceRegistrar, srv MatchmakingServiceServer) {
	s.RegisterService(&MatchmakingService_ServiceDesc, srv)
}

func _MatchmakingService_CreateArena_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateArenaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).CreateArena(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/CreateArena",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).CreateArena(ctx, req.(*CreateArenaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_GetArena_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArenaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).GetArena(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/GetArena",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).GetArena(ctx, req.(*ArenaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_GetArenas_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetArenasRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).GetArenas(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/GetArenas",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).GetArenas(ctx, req.(*GetArenasRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_UpdateArena_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateArenaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).UpdateArena(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/UpdateArena",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).UpdateArena(ctx, req.(*UpdateArenaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_DeleteArena_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArenaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).DeleteArena(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/DeleteArena",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).DeleteArena(ctx, req.(*ArenaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_CreateMatchmakingUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMatchmakingUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).CreateMatchmakingUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/CreateMatchmakingUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).CreateMatchmakingUser(ctx, req.(*CreateMatchmakingUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_GetMatchmakingUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchmakingUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).GetMatchmakingUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/GetMatchmakingUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).GetMatchmakingUser(ctx, req.(*MatchmakingUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_GetMatchmakingUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMatchmakingUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).GetMatchmakingUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/GetMatchmakingUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).GetMatchmakingUsers(ctx, req.(*GetMatchmakingUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_UpdateMatchmakingUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMatchmakingUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).UpdateMatchmakingUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/UpdateMatchmakingUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).UpdateMatchmakingUser(ctx, req.(*UpdateMatchmakingUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_SetMatchmakingUserElo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetMatchmakingUserEloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).SetMatchmakingUserElo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/SetMatchmakingUserElo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).SetMatchmakingUserElo(ctx, req.(*SetMatchmakingUserEloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_CreateMatchmakingTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMatchmakingTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).CreateMatchmakingTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/CreateMatchmakingTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).CreateMatchmakingTicket(ctx, req.(*CreateMatchmakingTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_PollMatchmakingTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchmakingTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).PollMatchmakingTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/PollMatchmakingTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).PollMatchmakingTicket(ctx, req.(*MatchmakingTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_GetMatchmakingTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchmakingTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).GetMatchmakingTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/GetMatchmakingTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).GetMatchmakingTicket(ctx, req.(*MatchmakingTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_GetMatchmakingTickets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMatchmakingTicketsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).GetMatchmakingTickets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/GetMatchmakingTickets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).GetMatchmakingTickets(ctx, req.(*GetMatchmakingTicketsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_UpdateMatchmakingTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMatchmakingTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).UpdateMatchmakingTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/UpdateMatchmakingTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).UpdateMatchmakingTicket(ctx, req.(*UpdateMatchmakingTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_ExpireMatchmakingTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchmakingTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).ExpireMatchmakingTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/ExpireMatchmakingTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).ExpireMatchmakingTicket(ctx, req.(*MatchmakingTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_StartMatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartMatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).StartMatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/StartMatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).StartMatch(ctx, req.(*StartMatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_EndMatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndMatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).EndMatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/EndMatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).EndMatch(ctx, req.(*EndMatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_GetMatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).GetMatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/GetMatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).GetMatch(ctx, req.(*MatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_GetMatches_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMatchesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).GetMatches(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/GetMatches",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).GetMatches(ctx, req.(*GetMatchesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchmakingService_UpdateMatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchmakingServiceServer).UpdateMatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchmakingService/UpdateMatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchmakingServiceServer).UpdateMatch(ctx, req.(*UpdateMatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MatchmakingService_ServiceDesc is the grpc.ServiceDesc for MatchmakingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MatchmakingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "MatchmakingService",
	HandlerType: (*MatchmakingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateArena",
			Handler:    _MatchmakingService_CreateArena_Handler,
		},
		{
			MethodName: "GetArena",
			Handler:    _MatchmakingService_GetArena_Handler,
		},
		{
			MethodName: "GetArenas",
			Handler:    _MatchmakingService_GetArenas_Handler,
		},
		{
			MethodName: "UpdateArena",
			Handler:    _MatchmakingService_UpdateArena_Handler,
		},
		{
			MethodName: "DeleteArena",
			Handler:    _MatchmakingService_DeleteArena_Handler,
		},
		{
			MethodName: "CreateMatchmakingUser",
			Handler:    _MatchmakingService_CreateMatchmakingUser_Handler,
		},
		{
			MethodName: "GetMatchmakingUser",
			Handler:    _MatchmakingService_GetMatchmakingUser_Handler,
		},
		{
			MethodName: "GetMatchmakingUsers",
			Handler:    _MatchmakingService_GetMatchmakingUsers_Handler,
		},
		{
			MethodName: "UpdateMatchmakingUser",
			Handler:    _MatchmakingService_UpdateMatchmakingUser_Handler,
		},
		{
			MethodName: "SetMatchmakingUserElo",
			Handler:    _MatchmakingService_SetMatchmakingUserElo_Handler,
		},
		{
			MethodName: "CreateMatchmakingTicket",
			Handler:    _MatchmakingService_CreateMatchmakingTicket_Handler,
		},
		{
			MethodName: "PollMatchmakingTicket",
			Handler:    _MatchmakingService_PollMatchmakingTicket_Handler,
		},
		{
			MethodName: "GetMatchmakingTicket",
			Handler:    _MatchmakingService_GetMatchmakingTicket_Handler,
		},
		{
			MethodName: "GetMatchmakingTickets",
			Handler:    _MatchmakingService_GetMatchmakingTickets_Handler,
		},
		{
			MethodName: "UpdateMatchmakingTicket",
			Handler:    _MatchmakingService_UpdateMatchmakingTicket_Handler,
		},
		{
			MethodName: "ExpireMatchmakingTicket",
			Handler:    _MatchmakingService_ExpireMatchmakingTicket_Handler,
		},
		{
			MethodName: "StartMatch",
			Handler:    _MatchmakingService_StartMatch_Handler,
		},
		{
			MethodName: "EndMatch",
			Handler:    _MatchmakingService_EndMatch_Handler,
		},
		{
			MethodName: "GetMatch",
			Handler:    _MatchmakingService_GetMatch_Handler,
		},
		{
			MethodName: "GetMatches",
			Handler:    _MatchmakingService_GetMatches_Handler,
		},
		{
			MethodName: "UpdateMatch",
			Handler:    _MatchmakingService_UpdateMatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "matchmaking.proto",
}
