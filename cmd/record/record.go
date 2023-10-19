package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/record"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const service = "record"
const defaultPort = 50053
const metricsPort = 8082
const cacheConn = "localhost:6379"
const cachePassword = ""
const cacheDB = 0
const cacheExpiration = 30 * time.Second
const minRecordNameLength = 3
const defaultMaxPageLength = 10
const maxMaxPageLength = 100

var dbIndices = []mongo.IndexModel{
	{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "userId", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	},
	{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "record", Value: 1},
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
		Collection: "records",
		Indices:    dbIndices,
	})
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	cache := cache.NewRedisCache(cacheConn, cachePassword, cacheDB, cacheExpiration)
	metrics, err := metrics.NewPrometheusMetrics(prometheus.NewRegistry(), "record", metricsPort)
	if err != nil {
		log.Fatalf("failed to create metrics: %v", err)
	}
	grpcServer := grpc.NewServer()
	recordService := record.NewService(
		record.WithDatabase(db),
		record.WithCache(cache),
		record.WithMetrics(metrics),
		record.WithMinRecordNameLength(minRecordNameLength),
		record.WithDefaultMaxPageLength(defaultMaxPageLength),
		record.WithMaxMaxPageLength(maxMaxPageLength),
	)
	defer recordService.Disconnect(context.TODO())
	api.RegisterRecordServiceServer(grpcServer, recordService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
