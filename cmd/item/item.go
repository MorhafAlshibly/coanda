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
	"github.com/MorhafAlshibly/coanda/pkg/secrets"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	fs                   = ff.NewFlagSet("item")
	service              = fs.String('s', "service", "item", "the name of the service")
	port                 = fs.Uint('p', "port", 50051, "the default port to listen on")
	metricsPort          = fs.Uint('m', "metricsPort", 8081, "the port to serve metrics on")
	cacheConnSecret      = fs.StringLong("cacheConnSecret", "", "the name of the secret containing the cache connection string")
	cacheConn            = fs.StringLong("cacheConn", "localhost:6379", "the connection string to the cache")
	cachePasswordSecret  = fs.StringLong("cachePasswordSecret", "", "the name of the secret containing the cache password")
	cachePassword        = fs.StringLong("cachePassword", "", "the password to the cache")
	cacheDB              = fs.IntLong("cacheDB", 0, "the database to use in the cache")
	cacheExpiration      = fs.DurationLong("cacheExpiration", 30*time.Second, "the expiration time for the cache")
	defaultMaxPageLength = fs.UintLong("defaultMaxPageLength", 10, "the default max page length")
	maxMaxPageLength     = fs.UintLong("maxMaxPageLength", 100, "the max max page length")
	minTypeLength        = fs.UintLong("minTypeLength", 3, "the min type length")
	maxTypeLength        = fs.UintLong("maxTypeLength", 20, "the max type length")
)

func main() {
	ctx := context.TODO()
	gf, err := flags.GetGlobalFlags()
	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(gf.FlagSet))
		log.Fatalf("failed to parse global flags: %v", err)
	}
	err = ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("ITEM"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		log.Fatalf("failed to parse flags: %v", err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var store *storage.TableStorage
	var redis *cache.RedisCache
	if *gf.Environment == "dev" {
		redis = cache.NewRedisCache(*cacheConn, *cachePassword, *cacheDB, *cacheExpiration)
		store, err = storage.NewTableStorage(ctx, nil, *gf.TableConn, *service)
	} else {
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Fatalf("failed to create credential: %v", err)
		}
		secrets, err := secrets.NewKeyVault(cred, *gf.VaultConn)
		if err != nil {
			log.Fatalf("failed to create secrets: %v", err)
		}
		tableConnFromSecret, err := secrets.GetSecret(ctx, *gf.TableSecret, nil)
		if err != nil {
			log.Fatalf("failed to get table connection string from secret: %v", err)
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
		store, err = storage.NewTableStorage(ctx, cred, tableConnFromSecret, *service)
	}
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
