package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/record"
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
	fs                   = ff.NewFlagSet("record")
	service              = fs.String('s', "service", "record", "the name of the service")
	port                 = fs.Uint('p', "port", 50053, "the default port to listen on")
	metricsPort          = fs.Uint('m', "metricsPort", 8081, "the port to serve metrics on")
	cacheConn            = fs.StringLong("cacheConn", "localhost:6379", "the connection string to the cache")
	cachePassword        = fs.StringLong("cachePassword", "", "the password to the cache")
	cacheDB              = fs.IntLong("cacheDB", 0, "the database to use in the cache")
	cacheExpiration      = fs.DurationLong("cacheExpiration", 5*time.Second, "the expiration time for the cache")
	minRecordNameLength  = fs.UintLong("minRecordNameLength", 3, "the min record name length")
	maxRecordNameLength  = fs.UintLong("maxRecordNameLength", 20, "the max record name length")
	defaultMaxPageLength = fs.UintLong("defaultMaxPageLength", 10, "the default max page length")
	maxMaxPageLength     = fs.UintLong("maxMaxPageLength", 100, "the max max page length")
	tableName            = fs.StringLong("tableName", "record", "the name of the table to use")
	region               = fs.StringLong("region", "localhost", "the region to use for the database")
	baseEndpoint         = fs.StringLong("baseEndpoint", "http://localhost:8000", "the base endpoint to use for the database")
	recordIndex          = fs.StringLong("recordIndex", "record-index", "the name of the index to create and use for sorting records")
	table                = &dynamodb.CreateTableInput{
		TableName: tableName,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: "S",
			},
			{
				AttributeName: aws.String("userId"),
				AttributeType: "N",
			},
			{
				AttributeName: aws.String("record"),
				AttributeType: "N",
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       "HASH",
			},
			{
				AttributeName: aws.String("userId"),
				KeyType:       "RANGE",
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		LocalSecondaryIndexes: []types.LocalSecondaryIndex{
			{
				IndexName: recordIndex,
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("record"),
						KeyType:       "RANGE",
					},
				},
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
	recordService := record.NewService(
		record.WithDatabase(db),
		record.WithCache(redis),
		record.WithMetrics(metrics),
		record.WithMinRecordNameLength(uint8(*minRecordNameLength)),
		record.WithMaxRecordNameLength(uint8(*maxRecordNameLength)),
		record.WithDefaultMaxPageLength(uint8(*defaultMaxPageLength)),
		record.WithMaxMaxPageLength(uint8(*maxMaxPageLength)),
	)
	grpcServer := grpc.NewServer()
	api.RegisterRecordServiceServer(grpcServer, recordService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
