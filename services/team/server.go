package team

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/queue"
	schema "github.com/MorhafAlshibly/coanda/services/team/schema"
	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct {
	schema.UnimplementedTeamServiceServer
	queue queue.Queuer
	store mongo.Collection
	cache cache.Cacher
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
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://foo:bar@localhost:27017"))
	defer mongoClient.Disconnect(context.Background())
	database := mongoClient.Database("coanda")
	store := database.Collection("teams")
	schema.RegisterTeamServiceServer(grpcServer, &server{
		queue: queue,
		store: *store,
		cache: cache.NewRedisCache("localhost:6379", "", 0, 60*time.Second),
	})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) CreateTeam(ctx context.Context, req *schema.CreateTeamRequest) (*schema.BoolResponse, error) {
	marshalled, err := sonic.Marshal(req)
	if err != nil {
		return nil, err
	}
	err = s.queue.Enqueue(ctx, "CreateTeam", marshalled)
	if err != nil {
		return nil, err
	}
	return &schema.BoolResponse{Value: true}, nil
}

func (s *server) GetTeam(ctx context.Context, req *schema.GetTeamRequest) (*schema.Team, error) {
	data, err := s.cache.Get(ctx, *req.Name)
	// If the item is not in the cache, get it from the store, else marshal it to output
	if err == nil {
		var out *schema.Team
		err = sonic.Unmarshal([]byte(data), &out)
		if err != nil {
			return nil, err
		}
		return out, nil
	}
	// Get the item from the store
	var result bson.M
	err = s.store.FindOne(
		ctx,
		bson.D{{Key: "name", Value: req.Name}},
		options.FindOne().SetSort(bson.D{{Key: "name", Value: 1}}),
	).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	out := bsonToTeam(&result)
	// Marshal the final output to a string to be cached
	marshalled, err := sonic.Marshal(out)
	if err != nil {
		return nil, err
	}
	// Add the item to the cache in a separate thread
	go s.cache.Add(ctx, *req.Name, string(marshalled))
	return out, nil
}

func (s *server) GetTeams(ctx context.Context, req *schema.GetTeamsRequest) (*schema.Teams, error) {
	var result []*schema.Team
	cursor, err := s.store.Find(
		ctx,
		bson.D{},
		options.Find().SetSort(bson.D{{Key: "rank", Value: 1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	skip := (int(req.Page) - 1) * int(req.Max)
	for i := 0; i < skip; i++ {
		cursor.Next(ctx)
	}
	for i := 0; i < int(req.Max); i++ {
		if !cursor.Next(ctx) {
			break
		}
		var team bson.M
		err = cursor.Decode(&team)
		if err != nil {
			return nil, err
		}
		result = append(result, bsonToTeam(&team))
	}
	return &schema.Teams{Teams: result}, nil
}

func (s *server) SearchTeams(ctx context.Context, req *schema.SearchTeamsRequest) (*schema.Teams, error) {
	return nil, nil
}

func (s *server) UpdateTeamScore(ctx context.Context, req *schema.UpdateTeamScoreRequest) (*schema.BoolResponse, error) {
	marshalled, err := sonic.Marshal(req)
	if err != nil {
		return nil, err
	}
	err = s.queue.Enqueue(ctx, "UpdateTeamScore", marshalled)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *server) UpdateTeamData(ctx context.Context, req *schema.UpdateTeamDataRequest) (*schema.Team, error) {
	return nil, nil
}

func bsonToTeam(result *bson.M) *schema.Team {
	return &schema.Team{
		Name:                (*result)["name"].(string),
		Owner:               (*result)["owner"].(string),
		MembersWithoutOwner: (*result)["membersWithoutOwner"].([]string),
		Score:               (*result)["score"].(int32),
		Rank:                (*result)["rank"].(int32),
		Data:                (*result)["data"].(map[string]string),
	}
}
