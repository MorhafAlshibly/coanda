package main

import (
	"context"
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

const service = "team"
const defaultPort = 50052
const metricsPort = 8081
const cacheConn = "localhost:6379"
const cachePassword = ""
const cacheDB = 0
const cacheExpiration = 30 * time.Second

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
	if err != nil {
		log.Fatalf("failed to create queue: %v", err)
	}
	cache := cache.NewRedisCache(cacheConn, cachePassword, cacheDB, cacheExpiration)
	metrics, err := metrics.NewPrometheusMetrics(prometheus.NewRegistry(), "team", metricsPort)
	if err != nil {
		log.Fatalf("failed to create metrics: %v", err)
	}
	grpcServer := grpc.NewServer()
	teamService := team.NewService(context.TODO(), team.NewServiceInput{
		Db:                db,
		Cache:             cache,
		Metrics:           metrics,
		MaxMembers:        5,
		MinTeamNameLength: 3,
	},
	)
	defer teamService.Disconnect(context.TODO())
	api.RegisterTeamServiceServer(grpcServer, teamService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
