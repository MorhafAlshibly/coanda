package team

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/queue"
	schema "github.com/MorhafAlshibly/coanda/services/team/schema"
	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Pipeline
var rankStage = bson.D{
	{Key: "$setWindowFields", Value: bson.D{
		{Key: "sortBy", Value: bson.D{
			{Key: "score", Value: -1},
		}},
		{Key: "output", Value: bson.D{
			{Key: "rank", Value: bson.D{
				{Key: "$rank", Value: bson.D{}},
			}},
		}},
	}},
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
	_, err = store.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "name", Value: "text"},
				{Key: "owner", Value: "text"},
			},
		},
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "owner", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "score", Value: -1},
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to create index: %v", err)
	}
	schema.RegisterTeamServiceServer(grpcServer, &server{
		queue: queue,
		store: *store,
		cache: cache.NewRedisCache("localhost:6379", "", 0, 60*time.Second),
	})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) CreateTeam(ctx context.Context, req *schema.CreateTeamRequest) (*schema.Team, error) {
	result, err := s.store.InsertOne(ctx, bson.D{
		{Key: "name", Value: req.Name},
		{Key: "owner", Value: req.Owner},
		{Key: "membersWithoutOwner", Value: req.MembersWithoutOwner},
		{Key: "score", Value: req.Score},
		{Key: "data", Value: req.Data},
	})
	if err != nil {
		return nil, err
	}
	return &schema.Team{
		Id:                  result.InsertedID.(primitive.ObjectID).String(),
		Name:                req.Name,
		Owner:               req.Owner,
		MembersWithoutOwner: req.MembersWithoutOwner,
		Score:               *req.Score,
		Data:                req.Data,
	}, nil
}

func (s *server) GetTeam(ctx context.Context, req *schema.GetTeamRequest) (*schema.Team, error) {
	filter, err := getFilter(req)
	if err != nil {
		return nil, err
	}
	data, err := s.cache.Get(ctx, filter[0].Value.(string))
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
	var out *schema.Team
	matchStage := bson.D{
		{Key: "$match", Value: filter},
	}
	cursor, err := s.store.Aggregate(ctx, mongo.Pipeline{rankStage, matchStage})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.Next(ctx)
	err = cursor.Decode(out)
	if err != nil {
		return nil, err
	}
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
	cursor, err := s.store.Aggregate(ctx, mongo.Pipeline{rankStage})
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
	var result []*schema.Team
	cursor, err := s.store.Find(ctx, bson.D{
		{Key: "$text", Value: bson.D{
			{Key: "$search", Value: req.Query},
		}},
	})
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

func (s *server) UpdateTeamScore(ctx context.Context, req *schema.UpdateTeamScoreRequest) (*schema.Team, error) {
	filter, err := getFilter(req.Team)
	if err != nil {
		return nil, err
	}
	var out *schema.Team
	err = s.store.FindOneAndUpdate(ctx, filter, bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "score", Value: req.ScoreOffset},
		}},
	}).Decode(out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *server) UpdateTeamData(ctx context.Context, req *schema.UpdateTeamDataRequest) (*schema.Team, error) {
	filter, err := getFilter(req.Team)
	if err != nil {
		return nil, err
	}
	var out *schema.Team
	err = s.store.FindOneAndUpdate(ctx, filter, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "data", Value: req.Data},
		}},
	}).Decode(out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *server) DeleteTeam(ctx context.Context, req *schema.DeleteTeamRequest) (*schema.Team, error) {
	filter, err := getFilter(req.Team)
	if err != nil {
		return nil, err
	}
	var out *schema.Team
	err = s.store.FindOneAndDelete(ctx, filter).Decode(out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *server) JoinTeam(ctx context.Context, req *schema.JoinTeamRequest) (*schema.BoolResponse, error) {
	marshalled, err := sonic.Marshal(req.Team)
	if err != nil {
		return nil, err
	}
	err = s.queue.Enqueue(ctx, "JoinTeam", marshalled)
	if err != nil {
		return nil, err
	}
	return &schema.BoolResponse{Value: true}, nil
}

func (s *server) LeaveTeam(ctx context.Context, req *schema.LeaveTeamRequest) (*schema.BoolResponse, error) {
	marshalled, err := sonic.Marshal(req.Team)
	if err != nil {
		return nil, err
	}
	err = s.queue.Enqueue(ctx, "LeaveTeam", marshalled)
	if err != nil {
		return nil, err
	}
	return &schema.BoolResponse{Value: true}, nil
}

func getFilter(input *schema.GetTeamRequest) (bson.D, error) {
	if input.Id != nil {
		return bson.D{
			{Key: "_id", Value: input.Id},
		}, nil
	}
	if input.Name != nil {
		return bson.D{
			{Key: "name", Value: input.Name},
		}, nil
	}
	return nil, errors.New("invalid input")
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
