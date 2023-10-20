package main

import (
	"context"
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

const service = "item"
const tableConn = "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;"
const defaultPort = 50051
const metricsPort = 8081
const cacheConn = "localhost:6379"
const cachePassword = ""
const cacheDB = 0
const cacheExpiration = 30 * time.Second
const defaultMaxPageLength = 10
const maxMaxPageLength = 100
const minTypeLength = 3
const maxTypeLength = 20

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", defaultPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	store, err := storage.NewTableStorage(context.TODO(), tableConn, service)
	if err != nil {
		log.Fatalf("failed to create store: %v", err)
	}
	cache := cache.NewRedisCache(cacheConn, cachePassword, cacheDB, cacheExpiration)
	metrics, err := metrics.NewPrometheusMetrics(prometheus.NewRegistry(), service, metricsPort)
	if err != nil {
		log.Fatalf("failed to create metrics: %v", err)
	}
	itemService := item.NewService(
		item.WithStore(store),
		item.WithCache(cache),
		item.WithMetrics(metrics),
		item.WithDefaultMaxPageLength(defaultMaxPageLength),
		item.WithMaxMaxPageLength(maxMaxPageLength),
		item.WithMinTypeLength(minTypeLength),
		item.WithMaxTypeLength(maxTypeLength),
	)
	grpcServer := grpc.NewServer()
	api.RegisterItemServiceServer(grpcServer, itemService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
