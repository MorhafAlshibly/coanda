package team

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/queue"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
	schema "github.com/MorhafAlshibly/coanda/services/team/schema"
	"github.com/bytedance/sonic"
	"google.golang.org/grpc"
)

type server struct {
	schema.UnimplementedTeamServiceServer
	Queue queue.Queuer
	Store storage.Storer
	Cache cache.Cacher
}

func (s *server) CreateTeam(ctx context.Context, req *schema.CreateTeamRequest) (*schema.BoolResponse, error) {
	marshalled, err := sonic.Marshal(req)
	if err != nil {
		return nil, err
	}
	err = s.Queue.Enqueue(ctx, marshalled)
	if err != nil {
		return nil, err
	}
	return &schema.BoolResponse{Value: true}, nil
}

func (s *server) GetTeam(ctx context.Context, req *schema.GetTeamRequest) (*schema.GetTeamResponse, error) {

}

func Run() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	queue, err := queue.NewServiceBus(context.Background(), "", "")
	if err != nil {
		log.Fatalf("failed to create queue: %v", err)
	}
	store, err := storage.NewTableStorage(context.Background(), "", "")
	if err != nil {
		log.Fatalf("failed to create store: %v", err)
	}
	schema.RegisterTeamServiceServer(grpcServer, &server{
		Queue: queue,
		Store: store,
		Cache: cache.NewRedisCache("localhost:6379", "", 0, 60*time.Second),
	})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
