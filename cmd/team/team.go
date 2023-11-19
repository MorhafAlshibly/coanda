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
	fs                   = ff.NewFlagSet("team")
	service              = fs.String('s', "service", "team", "the name of the service")
	port                 = fs.Uint('p', "port", 50052, "the default port to listen on")
	metricsPort          = fs.Uint('m', "metricsPort", 8081, "the port to serve metrics on")
	vaultConn            = fs.StringLong("vaultConn", "", "the connection string to the vault (optional)")
	mongoConnSecret      = fs.StringLong("mongoConnSecret", "", "the name of the secret containing the mongo connection string (optional, requires vaultConn)")
	mongoConn            = fs.StringLong("mongoConn", "mongodb://localhost:27017", "the connection string to the mongo database")
	mongoDatabase        = fs.StringLong("mongoDatabase", "coanda", "the name of the mongo database")
	mongoCollection      = fs.StringLong("mongoCollection", "team", "the name of the mongo collection")
	cacheConn            = fs.StringLong("cacheConn", "localhost:6379", "the connection string to the cache")
	cachePassword        = fs.StringLong("cachePassword", "", "the password to the cache")
	cacheDB              = fs.IntLong("cacheDB", 0, "the database to use in the cache")
	cacheExpiration      = fs.DurationLong("cacheExpiration", 5*time.Second, "the expiration time for the cache")
	maxMembers           = fs.UintLong("maxMembers", 5, "the max members")
	minTeamNameLength    = fs.UintLong("minTeamNameLength", 3, "the min team name length")
	maxTeamNameLength    = fs.UintLong("maxTeamNameLength", 20, "the max team name length")
	defaultMaxPageLength = fs.UintLong("defaultMaxPageLength", 10, "the default max page length")
	maxMaxPageLength     = fs.UintLong("maxMaxPageLength", 100, "the max max page length")
	dbIndices            = []mongo.IndexModel{
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
