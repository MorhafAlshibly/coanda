package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/item"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var service = flag.String("service", "item", "the name of the service")
var tableConn = flag.String("tableConn", "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", "the connection string to the table storage")
var defaultPort = flag.Uint("defaultPort", 50051, "the default port to listen on")
var metricsPort = flag.Uint("metricsPort", 8081, "the port to serve metrics on")
var cacheConn = flag.String("cacheConn", "localhost:6379", "the connection string to the cache")
var cachePassword = flag.String("cachePassword", "", "the password to the cache")
var cacheDB = flag.Int("cacheDB", 0, "the database to use in the cache")
var cacheExpiration = flag.Duration("cacheExpiration", 30*time.Second, "the expiration time for the cache")
var defaultMaxPageLength = flag.Uint("defaultMaxPageLength", 10, "the default max page length")
var maxMaxPageLength = flag.Uint("maxMaxPageLength", 100, "the max max page length")
var minTypeLength = flag.Uint("minTypeLength", 3, "the min type length")
var maxTypeLength = flag.Uint("maxTypeLength", 20, "the max type length")

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", defaultPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	store, err := storage.NewTableStorage(context.TODO(), *tableConn, *service)
	if err != nil {
		log.Fatalf("failed to create store: %v", err)
	}
	cache := cache.NewRedisCache(*cacheConn, *cachePassword, *cacheDB, *cacheExpiration)
	metrics, err := metrics.NewPrometheusMetrics(prometheus.NewRegistry(), *service, uint16(*metricsPort))
	if err != nil {
		log.Fatalf("failed to create metrics: %v", err)
	}
	itemService := item.NewService(
		item.WithStore(store),
		item.WithCache(cache),
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
