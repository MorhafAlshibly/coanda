package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/queue"
	"github.com/MorhafAlshibly/coanda/services/team/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TeamService struct {
	schema.UnimplementedTeamServiceServer
	Db                database.Databaser
	Cache             cache.Cacher
	Queue             queue.Queuer
	MaxMembers        int
	MinTeamNameLength int
	Pipeline          mongo.Pipeline
}

type NewTeamServiceInput struct {
	DatabaseConnection string
	DatabaseName       string
	DatabaseCollection string
	QueueConnection    string
	QueueName          string
	CacheConnection    string
	CachePassword      string
	CacheDB            int
	CacheExpiration    time.Duration
	MaxMembers         int
	MinTeamNameLength  int
}

var (
	dbIndices = []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "name", Value: "text"},
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
	}
	rankStage = bson.D{
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
)

func NewTeamService(ctx context.Context, input NewTeamServiceInput) (*TeamService, error) {
	db, err := database.NewMongoDatabase(ctx, database.MongoDatabaseInput{
		Connection: "mongodb://localhost:27017",
		Database:   "coanda",
		Collection: "teams",
		Indices:    dbIndices,
	})
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	queue, err := queue.NewServiceBus(ctx, input.QueueConnection, input.QueueName)
	if err != nil {
		log.Fatalf("failed to create queue: %v", err)
	}
	cache := cache.NewRedisCache(input.CacheConnection, input.CachePassword, input.CacheDB, input.CacheExpiration)
	return &TeamService{
		Db:                db,
		Cache:             cache,
		Queue:             queue,
		MaxMembers:        input.MaxMembers,
		MinTeamNameLength: input.MinTeamNameLength,
		Pipeline:          mongo.Pipeline{rankStage},
	}, nil
}

func (s *TeamService) Disconnect(ctx context.Context) error {
	return s.Db.Disconnect(ctx)
}

func (s *TeamService) CreateTeam(ctx context.Context, in *schema.CreateTeamRequest) (*schema.Team, error) {
	command := &CreateTeamCommand{
		Service: s,
		In:      in,
	}
	invoker := &invokers.BasicInvoker{}
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) GetTeam(ctx context.Context, in *schema.GetTeamRequest) (*schema.Team, error) {
	command := &GetTeamCommand{
		Service: s,
		In:      in,
	}
	invoker := invokers.NewCacheInvoker(s.Cache)
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) GetTeams(ctx context.Context, in *schema.GetTeamsRequest) (*schema.Teams, error) {
	command := &GetTeamsCommand{
		Service: s,
		In:      in,
	}
	invoker := invokers.NewCacheInvoker(s.Cache)
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) SearchTeams(ctx context.Context, in *schema.SearchTeamsRequest) (*schema.Teams, error) {
	command := &SearchTeamsCommand{
		Service: s,
		In:      in,
	}
	invoker := invokers.NewCacheInvoker(s.Cache)
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) UpdateTeamScore(ctx context.Context, in *schema.UpdateTeamScoreRequest) (*schema.Team, error) {
	command := &UpdateTeamScoreCommand{
		Service: s,
		In:      in,
	}
	invoker := &invokers.BasicInvoker{}
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) UpdateTeamData(ctx context.Context, in *schema.UpdateTeamDataRequest) (*schema.Team, error) {
	command := &UpdateTeamDataCommand{
		Service: s,
		In:      in,
	}
	invoker := &invokers.BasicInvoker{}
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) DeleteTeam(ctx context.Context, in *schema.DeleteTeamRequest) (*schema.Team, error) {
	command := &DeleteTeamCommand{
		Service: s,
		In:      in,
	}
	invoker := &invokers.BasicInvoker{}
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) JoinTeam(ctx context.Context, in *schema.JoinTeamRequest) (*schema.BoolResponse, error) {
	command := &JoinTeamCommand{
		Service: s,
		In:      in,
	}
	invoker := &invokers.BasicInvoker{}
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) LeaveTeam(ctx context.Context, in *schema.LeaveTeamRequest) (*schema.BoolResponse, error) {
	command := &LeaveTeamCommand{
		Service: s,
		In:      in,
	}
	invoker := &invokers.BasicInvoker{}
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func GetFilter(input *schema.GetTeamRequest) (bson.D, error) {
	if input.Id != nil {
		id, err := primitive.ObjectIDFromHex(*input.Id)
		if err != nil {
			return nil, err
		}
		return bson.D{
			{Key: "_id", Value: id},
		}, nil
	}
	if input.Name != nil {
		return bson.D{
			{Key: "name", Value: input.Name},
		}, nil
	}
	if input.Owner != nil {
		return bson.D{
			{Key: "owner", Value: input.Owner},
		}, nil
	}
	return nil, errors.New("invalid input")
}

func ToTeams(ctx context.Context, cursor *mongo.Cursor, page uint64, max uint64) (*schema.Teams, error) {
	var result []*schema.Team
	skip := (int(page) - 1) * int(max)
	for i := 0; i < skip; i++ {
		cursor.Next(ctx)
	}
	for i := 0; i < int(max); i++ {
		if !cursor.Next(ctx) {
			break
		}
		team, err := ToTeam(cursor)
		if err != nil {
			return nil, err
		}
		result = append(result, team)
	}
	return &schema.Teams{Teams: result}, nil
}

func ToTeam(cursor *mongo.Cursor) (*schema.Team, error) {
	var result *bson.M
	err := cursor.Decode(&result)
	if err != nil {
		return nil, err
	}
	// Convert []int64 to []uint64
	membersWithoutOwner := []uint64{}
	for _, member := range (*result)["membersWithoutOwner"].(primitive.A) {
		membersWithoutOwner = append(membersWithoutOwner, uint64(member.(int64)))
	}
	(*result)["membersWithoutOwner"] = membersWithoutOwner
	// Convert data to map[string]string
	data := (*result)["data"].(primitive.M)
	(*result)["data"] = map[string]string{}
	for key, value := range data {
		(*result)["data"].(map[string]string)[key] = value.(string)
	}
	// If rank is not given, set it to 0
	if _, ok := (*result)["rank"]; !ok {
		(*result)["rank"] = int32(0)
	}
	return &schema.Team{
		Id:                  (*result)["_id"].(primitive.ObjectID).Hex(),
		Name:                (*result)["name"].(string),
		Owner:               uint64((*result)["owner"].(int64)),
		MembersWithoutOwner: membersWithoutOwner,
		Score:               (*result)["score"].(int64),
		Rank:                int64((*result)["rank"].(int32)),
		Data:                (*result)["data"].(map[string]string),
	}, nil
}
