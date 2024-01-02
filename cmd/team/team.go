package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	// Flags set from command line/environment variables
	fs                   = ff.NewFlagSet("team")
	service              = fs.String('s', "service", "team", "the name of the service")
	port                 = fs.Uint('p', "port", 50052, "the default port to listen on")
	metricsPort          = fs.Uint('m', "metricsPort", 8081, "the port to serve metrics on")
	cacheConn            = fs.StringLong("cacheConn", "localhost:6379", "the connection string to the cache")
	cachePassword        = fs.StringLong("cachePassword", "", "the password to the cache")
	cacheDB              = fs.IntLong("cacheDB", 0, "the database to use in the cache")
	cacheExpiration      = fs.DurationLong("cacheExpiration", 5*time.Second, "the expiration time for the cache")
	maxMembers           = fs.UintLong("maxMembers", 5, "the max members")
	minTeamNameLength    = fs.UintLong("minTeamNameLength", 3, "the min team name length")
	maxTeamNameLength    = fs.UintLong("maxTeamNameLength", 20, "the max team name length")
	defaultMaxPageLength = fs.UintLong("defaultMaxPageLength", 10, "the default max page length")
	maxMaxPageLength     = fs.UintLong("maxMaxPageLength", 100, "the max max page length")
	tableName            = fs.StringLong("tableName", "team", "the name of the table to use")
	region               = fs.StringLong("region", "localhost", "the region to use for the database")
	baseEndpoint         = fs.StringLong("baseEndpoint", "http://localhost:8000", "the base endpoint to use for the database")
	leaderboardIndex     = fs.StringLong("leaderboardIndex", "leaderboard", "the name of the index to create and use for leaderboard")
	readTableName        = fs.StringLong("readTableName", "teamLeaderboard", "the name of the table to use for reading")
	table                = &dynamodb.CreateTableInput{
		TableName: tableName,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: "S",
			},
			{
				AttributeName: aws.String("owner"),
				AttributeType: "N",
			},
			{
				AttributeName: aws.String("score"),
				AttributeType: "N",
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       "HASH",
			},
			{
				AttributeName: aws.String("owner"),
				KeyType:       "RANGE",
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		LocalSecondaryIndexes: []types.LocalSecondaryIndex{
			{
				IndexName: leaderboardIndex,
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("score"),
						KeyType:       "RANGE",
					},
				},
			},
		},
		StreamSpecification: &types.StreamSpecification{
			StreamEnabled:  aws.Bool(true),
			StreamViewType: "NEW_IMAGE",
		},
	}
)

func main() {
	ctx := context.TODO()
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("TEAM"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		log.Fatalf("failed to parse flags: %v", err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if *vaultConn != "" && *mongoConnSecret != "" {
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Fatalf("failed to get credentials: %v", err)
		}
		vault, err := azsecrets.NewClient(*vaultConn, cred, nil)
		if err != nil {
			log.Fatalf("failed to create vault: %v", err)
		}
		secret, err := vault.GetSecret(ctx, *mongoConnSecret, "", nil)
		if err != nil {
			log.Fatalf("failed to get mongoConn: %v", err)
		}
		*mongoConn = *secret.Value
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
	metrics, err := metrics.NewPrometheusMetrics(prometheus.NewRegistry(), "team", uint16(*metricsPort))
	if err != nil {
		log.Fatalf("failed to create metrics: %v", err)
	}
	grpcServer := grpc.NewServer()
	teamService := team.NewService(
		team.WithDatabase(db),
		team.WithCache(redis),
		team.WithMetrics(metrics),
		team.WithMaxMembers(uint8(*maxMembers)),
		team.WithMinTeamNameLength(uint8(*minTeamNameLength)),
		team.WithMaxTeamNameLength(uint8(*maxTeamNameLength)),
		team.WithDefaultMaxPageLength(uint8(*defaultMaxPageLength)),
		team.WithMaxMaxPageLength(uint8(*maxMaxPageLength)),
	)
	defer teamService.Disconnect(ctx)
	api.RegisterTeamServiceServer(grpcServer, teamService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
