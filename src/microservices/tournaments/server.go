package main

import (
	"context"
	"log"
	"net"

	schema "github.com/MorhafAlshibly/coanda/src/microservices/tournaments/schema"
	"google.golang.org/grpc"
)

type server struct {
	schema.UnimplementedTournamentsServer
}

func (s *server) CreateTournament(ctx context.Context, req *schema.CreateTournamentRequest) (*schema.CreateTournamentResponse, error) {
	return &schema.CreateTournamentResponse{Id: 3}, nil
}

func (s *server) GetTournament(ctx context.Context, req *schema.GetTournamentRequest) (*schema.GetTournamentResponse, error) {
	return &schema.GetTournamentResponse{Id: 3, Name: "test", Game: "test"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	schema.RegisterTournamentsServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
