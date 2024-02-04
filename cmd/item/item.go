package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/item"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/database/dynamoTable"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	// Flags set from command line/environment variables
	fs                   = ff.NewFlagSet("item")
	service              = fs.String('s', "service", "item", "the name of the service")
	port                 = fs.Uint('p', "port", 50051, "the default port to listen on")
	metricsPort          = fs.Uint('m', "metricsPort", 8081, "the port to serve metrics on")
	cacheConn            = fs.StringLong("cacheConn", "localhost:6379", "the connection string to the cache")
	cachePassword        = fs.StringLong("cachePassword", "", "the password to the cache")
	cacheDB              = fs.IntLong("cacheDB", 0, "the database to use in the cache")
	cacheExpiration      = fs.DurationLong("cacheExpiration", 5*time.Second, "the expiration time for the cache")
	defaultMaxPageLength = fs.UintLong("defaultMaxPageLength", 10, "the default max page length")
	maxMaxPageLength     = fs.UintLong("maxMaxPageLength", 100, "the max max page length")
	minTypeLength        = fs.UintLong("minTypeLength", 3, "the min type length")
	maxTypeLength        = fs.UintLong("maxTypeLength", 20, "the max type length")
	tableName            = fs.StringLong("tableName", "item", "the name of the table to use")
	region               = fs.StringLong("region", "localhost", "the region to use for the database")
	baseEndpoint         = fs.StringLong("baseEndpoint", "http://localhost:8000", "the base endpoint to use for the database")
	table                = &dynamodb.CreateTableInput{
		TableName: tableName,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: "S",
			},
			{
				AttributeName: aws.String("type"),
				AttributeType: "S",
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("type"),
				KeyType:       "HASH",
			},
			{
				AttributeName: aws.String("id"),
				KeyType:       "RANGE",
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}
)

func main() {
	ctx := context.TODO()
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("ITEM"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		log.Fatalf("failed to parse flags: %v", err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	redis := cache.NewRedisCache(*cacheConn, *cachePassword, *cacheDB, *cacheExpiration)
	db, err := dynamoTable.NewDynamoTable(ctx, &dynamoTable.DynamoTableInput{
		Options: &dynamodb.Options{
			Region:       *region,
			BaseEndpoint: baseEndpoint,
			Credentials:  credentials.NewStaticCredentialsProvider("test", "test", "test"),
		},
		CreateTableInput: table,
		Cache:            redis,
	})
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	metrics, err := metrics.NewPrometheusMetrics(prometheus.NewRegistry(), *service, uint16(*metricsPort))
	if err != nil {
		log.Fatalf("failed to create metrics: %v", err)
	}
	itemService := item.NewService(
		item.WithDatabase(db),
		item.WithCache(redis),
		item.WithMetrics(metrics),
		item.WithDefaultMaxPageLength(uint8(*defaultMaxPageLength)),
		item.WithMaxMaxPageLength(uint8(*maxMaxPageLength)),
		item.WithMinTypeLength(uint8(*minTypeLength)),
		item.WithMaxTypeLength(uint8(*maxTypeLength)),
	)
	grpcServer := grpc.NewServer()
	api.RegisterItemServiceServer(grpcServer, itemService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
