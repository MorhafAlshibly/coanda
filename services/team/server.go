package team

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"github.com/MorhafAlshibly/coanda/pkg/queue"
	"github.com/MorhafAlshibly/coanda/services/team/schema"
	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type TeamService struct {
	schema.UnimplementedTeamServiceServer
	db                database.Databaser
	cache             cache.Cacher
	queue             queue.Queuer
	maxMembers        int
	minTeamNameLength int
	pipeline          mongo.Pipeline
	metrics           metrics.Metrics
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
	MetricsPort        uint16
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

func Run() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	team, err := NewTeamService(context.TODO(), NewTeamServiceInput{
		DatabaseConnection: "mongodb://localhost:27017",
		DatabaseName:       "coanda",
		DatabaseCollection: "team",
		QueueConnection:    "",
		QueueName:          "",
		CacheConnection:    "localhost:6379",
		CachePassword:      "",
		CacheDB:            0,
		CacheExpiration:    1 * time.Minute,
		MetricsPort:        8081,
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
	//queue, err := queue.NewServiceBus(ctx, input.QueueConnection, input.QueueName)
	if err != nil {
		log.Fatalf("failed to create queue: %v", err)
	}
	cache := cache.NewRedisCache(input.CacheConnection, input.CachePassword, input.CacheDB, input.CacheExpiration)
	metrics, err := metrics.NewPrometheusMetrics(prometheus.NewRegistry(), "team", input.MetricsPort)
	if err != nil {
		log.Fatalf("failed to create metrics: %v", err)
	}
	return &TeamService{
		db:                db,
		cache:             cache,
		queue:             nil, //queue,
		maxMembers:        input.MaxMembers,
		minTeamNameLength: input.MinTeamNameLength,
		pipeline:          mongo.Pipeline{rankStage},
		metrics:           metrics,
	}, nil
}

func (s *TeamService) Disconnect(ctx context.Context) error {
	return s.db.Disconnect(ctx)
}

func (s *TeamService) CreateTeam(ctx context.Context, in *schema.CreateTeamRequest) (*schema.Team, error) {
	command := NewCreateTeamCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) GetTeam(ctx context.Context, in *schema.GetTeamRequest) (*schema.Team, error) {
	command := NewGetTeamCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) GetTeams(ctx context.Context, in *schema.GetTeamsRequest) (*schema.Teams, error) {
	command := NewGetTeamsCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) SearchTeams(ctx context.Context, in *schema.SearchTeamsRequest) (*schema.Teams, error) {
	command := NewSearchTeamsCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) UpdateTeamScore(ctx context.Context, in *schema.UpdateTeamScoreRequest) (*schema.Team, error) {
	command := NewUpdateTeamScoreCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) UpdateTeamData(ctx context.Context, in *schema.UpdateTeamDataRequest) (*schema.Team, error) {
	command := NewUpdateTeamDataCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) DeleteTeam(ctx context.Context, in *schema.DeleteTeamRequest) (*schema.Team, error) {
	command := NewDeleteTeamCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) JoinTeam(ctx context.Context, in *schema.JoinTeamRequest) (*schema.BoolResponse, error) {
	command := NewJoinTeamCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *TeamService) LeaveTeam(ctx context.Context, in *schema.LeaveTeamRequest) (*schema.BoolResponse, error) {
	command := NewLeaveTeamCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}
