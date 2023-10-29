package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var service = flag.String("service", "team", "the name of the service")
var defaultPort = flag.Uint("defaultPort", 50052, "the default port to listen on")
var metricsPort = flag.Uint("metricsPort", 8081, "the port to serve metrics on")
var cacheConn = flag.String("cacheConn", "localhost:6379", "the connection string to the cache")
var cachePassword = flag.String("cachePassword", "", "the password to the cache")
var cacheDB = flag.Int("cacheDB", 0, "the database to use in the cache")
var cacheExpiration = flag.Duration("cacheExpiration", 30*time.Second, "the expiration time for the cache")
var maxMembers = flag.Uint("maxMembers", 5, "the max members")
var minTeamNameLength = flag.Uint("minTeamNameLength", 3, "the min team name length")
var maxTeamNameLength = flag.Uint("maxTeamNameLength", 20, "the max team name length")
var defaultMaxPageLength = flag.Uint("defaultMaxPageLength", 10, "the default max page length")
var maxMaxPageLength = flag.Uint("maxMaxPageLength", 100, "the max max page length")

var dbIndices = []mongo.IndexModel{
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

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", defaultPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	db, err := database.NewMongoDatabase(context.TODO(), database.MongoDatabaseInput{
		Connection: "mongodb://localhost:27017",
		Database:   "coanda",
		Collection: "teams",
		Indices:    dbIndices,
	})
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	cache := cache.NewRedisCache(*cacheConn, *cachePassword, *cacheDB, *cacheExpiration)
	metrics, err := metrics.NewPrometheusMetrics(prometheus.NewRegistry(), "team", uint16(*metricsPort))
	if err != nil {
		log.Fatalf("failed to create metrics: %v", err)
	}
	grpcServer := grpc.NewServer()
	teamService := team.NewService(
		team.WithDatabase(db),
		team.WithCache(cache),
		team.WithMetrics(metrics),
		team.WithMaxMembers(uint8(*maxMembers)),
		team.WithMinTeamNameLength(uint8(*minTeamNameLength)),
		team.WithMaxTeamNameLength(uint8(*maxTeamNameLength)),
		team.WithDefaultMaxPageLength(uint8(*defaultMaxPageLength)),
		team.WithMaxMaxPageLength(uint8(*maxMaxPageLength)),
	)
	defer teamService.Disconnect(context.TODO())
	api.RegisterTeamServiceServer(grpcServer, teamService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
