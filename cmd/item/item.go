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
	"github.com/MorhafAlshibly/coanda/internal/item"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/flags"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	// Flags set from command line/environment variables
	fs            = ff.NewFlagSet("item")
	service       = fs.String('s', "service", "item", "the name of the service")
	port          = fs.Uint('p', "port", 50051, "the default port to listen on")
	metricsPort   = fs.Uint('m', "metricsPort", 8081, "the port to serve metrics on")
	appConfigConn = fs.StringLong("appConfigurationConn", "", "the connection string to the app configuration service")
	// Flags set from azure app configuration
	configFs             = ff.NewFlagSet(fs.GetName())
	tableConn            = configFs.StringLong("tableConn", "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", "the connection string to the table storage")
	cacheConn            = configFs.StringLong("cacheConn", "localhost:6379", "the connection string to the cache")
	cachePassword        = configFs.StringLong("cachePassword", "", "the password to the cache")
	cacheDB              = configFs.IntLong("cacheDB", 0, "the database to use in the cache")
	cacheExpiration      = configFs.DurationLong("cacheExpiration", 30*time.Second, "the expiration time for the cache")
	defaultMaxPageLength = configFs.UintLong("defaultMaxPageLength", 10, "the default max page length")
	maxMaxPageLength     = configFs.UintLong("maxMaxPageLength", 100, "the max max page length")
	minTypeLength        = configFs.UintLong("minTypeLength", 3, "the min type length")
	maxTypeLength        = configFs.UintLong("maxTypeLength", 20, "the max type length")
)

func main() {
	ctx := context.TODO()
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("ITEM"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
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
			err = ff.Parse(configFs, os.Args[1:], ff.WithEnvVarPrefix("ITEM"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
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
	store, err := storage.NewTableStorage(ctx, cred, *tableConn, *service)
	if err != nil {
		log.Fatalf("failed to create store: %v", err)
	}
	metrics, err := metrics.NewPrometheusMetrics(prometheus.NewRegistry(), *service, uint16(*metricsPort))
	if err != nil {
		log.Fatalf("failed to create metrics: %v", err)
	}
	itemService := item.NewService(
		item.WithStore(store),
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
