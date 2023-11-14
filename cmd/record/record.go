package main

import (
	"context"
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
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var (
	// Flags set from command line/environment variables
	fs            = ff.NewFlagSet("record")
	service       = fs.String('s', "service", "record", "the name of the service")
	port          = fs.Uint('p', "port", 50051, "the default port to listen on")
	metricsPort   = fs.Uint('m', "metricsPort", 8081, "the port to serve metrics on")
	appConfigConn = fs.StringLong("appConfigurationConn", "", "the connection string to the app configuration service")
	// Flags set from azure app configuration
	configFs             = ff.NewFlagSet(fs.GetName())
	mongoConn            = configFs.StringLong("mongoConn", "mongodb://localhost:27017", "the connection string to the mongo database")
	mongoDatabase        = configFs.StringLong("mongoDatabase", "coanda", "the name of the mongo database")
	mongoCollection      = configFs.StringLong("mongoCollection", "record", "the name of the mongo collection")
	cacheConn            = configFs.StringLong("cacheConn", "localhost:6379", "the connection string to the cache")
	cachePassword        = configFs.StringLong("cachePassword", "", "the password to the cache")
	cacheDB              = configFs.IntLong("cacheDB", 0, "the database to use in the cache")
	cacheExpiration      = configFs.DurationLong("cacheExpiration", 30*time.Second, "the expiration time for the cache")
	minRecordNameLength  = configFs.UintLong("minRecordNameLength", 3, "the min record name length")
	maxRecordNameLength  = configFs.UintLong("maxRecordNameLength", 20, "the max record name length")
	defaultMaxPageLength = configFs.UintLong("defaultMaxPageLength", 10, "the default max page length")
	maxMaxPageLength     = configFs.UintLong("maxMaxPageLength", 100, "the max max page length")
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
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("RECORD"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		log.Fatalf("failed to parse flags: %v", err)
	}
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to get credentials: %v", err)
	}
	if *appConfigConn != "" {
		appConfig, err := flags.NewAppConfiguration(ctx, cred, *appConfigConn)
		if err != nil {
			log.Fatalf("failed to create app configuration client: %v", err)
		}
		err = appConfig.Parse(ctx, configFs, configFs.GetName())
		if err != nil {
			err = ff.Parse(configFs, os.Args[1:], ff.WithEnvVarPrefix("RECORD"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
			if err != nil {
				fmt.Printf("%s\n", ffhelp.Flags(configFs))
				log.Fatalf("failed to parse flags: %v", err)
			}
		}
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	redis := cache.NewRedisCache(*cacheConn, *cachePassword, *cacheDB, *cacheExpiration)
	db, err := database.NewMongoDatabase(ctx, database.MongoDatabaseInput{
		Connection: *mongoConn,
		Database:   *mongoDatabase,
		Collection: *mongoCollection,
		Indices:    dbIndices,
	})
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
