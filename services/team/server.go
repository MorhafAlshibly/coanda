package team

import (
	"context"
	"log"
	"net"
	"time"

	schema "github.com/MorhafAlshibly/coanda/services/team/schema"
	"github.com/MorhafAlshibly/coanda/services/team/service"
	"google.golang.org/grpc"
)

func Run() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	team, err := service.NewTeamService(context.TODO(), service.NewTeamServiceInput{
		DatabaseConnection: "mongodb://localhost:27017",
		DatabaseName:       "coanda",
		DatabaseCollection: "team",
		QueueConnection:    "",
		QueueName:          "",
		CacheConnection:    "localhost:6379",
		CachePassword:      "",
		CacheDB:            0,
		CacheExpiration:    1 * time.Minute,
		MaxMembers:         5,
		MinTeamNameLength:  3,
	},
	)
	if err != nil {
		log.Fatalf("failed to create team service: %v", err)
	}
	defer team.Disconnect(context.TODO())
	schema.RegisterTeamServiceServer(grpcServer, team)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
