package team

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	schema "github.com/MorhafAlshibly/coanda/services/team/schema"
	service "github.com/MorhafAlshibly/coanda/services/team/service"
	"google.golang.org/grpc"
)

type server struct {
	schema.UnimplementedTeamServiceServer
	Team service.TeamService
}

func Run() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	schema.RegisterTeamServiceServer(grpcServer, &server{
		Team: service.NewTeamService(service.NewTeamServiceInput{
			Database:          db,
			Queue:             nil,
			Cache:             cache.NewRedisCache("localhost:6379", "", 0, 60*time.Second),
			MaxMembers:        5,
			MinTeamNameLength: 3,
		},
	})})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
