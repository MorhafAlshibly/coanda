package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/record"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/flags"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"github.com/MorhafAlshibly/coanda/pkg/secrets"
	"github.com/peterbourgon/ff"
	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var (
	fs                   = flag.NewFlagSet("record", flag.ContinueOnError)
	service              = fs.String("service", "record", "the name of the service")
	port                 = fs.Uint("port", 50051, "the default port to listen on")
	metricsPort          = fs.Uint("metricsPort", 8081, "the port to serve metrics on")
	mongoCollection      = fs.String("mongoCollection", "record", "the name of the mongo collection")
	cacheConnSecret      = fs.String("cacheConnSecret", "", "the name of the secret containing the cache connection string")
	cacheConn            = fs.String("cacheConn", "localhost:6379", "the connection string to the cache")
	cachePasswordSecret  = fs.String("cachePasswordSecret", "", "the name of the secret containing the cache password")
	cachePassword        = fs.String("cachePassword", "", "the password to the cache")
	cacheDB              = fs.Int("cacheDB", 0, "the database to use in the cache")
	cacheExpiration      = fs.Duration("cacheExpiration", 30*time.Second, "the expiration time for the cache")
	minRecordNameLength  = fs.Uint("minRecordNameLength", 3, "the min record name length")
	maxRecordNameLength  = fs.Uint("maxRecordNameLength", 20, "the max record name length")
	defaultMaxPageLength = fs.Uint("defaultMaxPageLength", 10, "the default max page length")
	maxMaxPageLength     = fs.Uint("maxMaxPageLength", 100, "the max max page length")
	dbIndices            = []mongo.IndexModel{
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
)

func main() {
	ctx := context.TODO()
	gf, err := flags.GetGlobalFlags()
	if err != nil {
		log.Fatalf("failed to get global flags: %v", err)
	}
	err = ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("RECORD"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var db *database.MongoDatabase
	var redis *cache.RedisCache
	if *gf.Environment == "dev" {
		redis = cache.NewRedisCache(*cacheConn, *cachePassword, *cacheDB, *cacheExpiration)
		db, err = database.NewMongoDatabase(ctx, database.MongoDatabaseInput{
			Connection: *gf.MongoConn,
			Database:   *gf.MongoDatabase,
			Collection: *mongoCollection,
			Indices:    dbIndices,
		})
	} else {
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Fatalf("failed to create credential: %v", err)
		}
		secrets, err := secrets.NewKeyVault(cred, *gf.VaultConn)
		if err != nil {
			log.Fatalf("failed to create secrets: %v", err)
		}
		mongoConnFromSecret, err := secrets.GetSecret(ctx, *gf.MongoSecret, nil)
		if err != nil {
			log.Fatalf("failed to get mongo connection string from secret: %v", err)
		}
		cacheConnFromSecret, err := secrets.GetSecret(ctx, *cacheConnSecret, nil)
		if err != nil {
			log.Fatalf("failed to get cache connection string from secret: %v", err)
		}
		cachePasswordFromSecret, err := secrets.GetSecret(ctx, *cachePasswordSecret, nil)
		if err != nil {
			log.Fatalf("failed to get cache password from secret: %v", err)
		}
		redis = cache.NewRedisCache(cacheConnFromSecret, cachePasswordFromSecret, *cacheDB, *cacheExpiration)
		db, err = database.NewMongoDatabase(ctx, database.MongoDatabaseInput{
			Connection: mongoConnFromSecret,
			Database:   *gf.MongoDatabase,
			Collection: *mongoCollection,
			Indices:    dbIndices,
		})
	}
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	metrics, err := metrics.NewPrometheusMetrics(prometheus.NewRegistry(), "record", uint16(*metricsPort))
	if err != nil {
		log.Fatalf("failed to create metrics: %v", err)
	}
	grpcServer := grpc.NewServer()
	recordService := record.NewService(
		record.WithDatabase(db),
		record.WithCache(redis),
		record.WithMetrics(metrics),
		record.WithMinRecordNameLength(uint8(*minRecordNameLength)),
		record.WithMaxRecordNameLength(uint8(*maxRecordNameLength)),
		record.WithDefaultMaxPageLength(uint8(*defaultMaxPageLength)),
		record.WithMaxMaxPageLength(uint8(*maxMaxPageLength)),
	)
	defer recordService.Disconnect(ctx)
	api.RegisterRecordServiceServer(grpcServer, recordService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
