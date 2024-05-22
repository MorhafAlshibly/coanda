package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	// Flags set from command line/environment variables
	fs                   = ff.NewFlagSet("event")
	service              = fs.String('s', "service", "event", "the name of the service")
	port                 = fs.Uint('p', "port", 50053, "the default port to listen on")
	metricsPort          = fs.Uint('m', "metricsPort", 8083, "the port to serve metrics on")
	dsn                  = fs.StringLong("dsn", "root:password@tcp(localhost:3306)", "the data source name for the database")
	cacheHost            = fs.StringLong("cacheHost", "localhost:6379", "the connection string to the cache")
	cachePassword        = fs.StringLong("cachePassword", "", "the password to the cache")
	cacheDB              = fs.IntLong("cacheDB", 0, "the database to use in the cache")
	cacheExpiration      = fs.DurationLong("cacheExpiration", 5*time.Second, "the expiration time for the cache")
	minEventNameLength   = fs.UintLong("minEventNameLength", 3, "the min event name length")
	maxEventNameLength   = fs.UintLong("maxEventNameLength", 20, "the max event name length")
	minRoundNameLength   = fs.UintLong("minRoundNameLength", 3, "the min round name length")
	maxRoundNameLength   = fs.UintLong("maxRoundNameLength", 20, "the max round name length")
	maxNumberOfRounds    = fs.UintLong("maxNumberOfRounds", 10, "the max number of rounds")
	defaultMaxPageLength = fs.UintLong("defaultMaxPageLength", 10, "the default max page length")
	maxMaxPageLength     = fs.UintLong("maxMaxPageLength", 100, "the max max page length")
)

func main() {
	_ = context.TODO()
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("EVENT"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		log.Fatalf("failed to parse flags: %v", err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	redis := cache.NewRedisCache(*cacheHost, *cachePassword, *cacheDB, *cacheExpiration)
	dbConn, err := sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer dbConn.Close()
	db := model.New(dbConn)
	metrics, err := metrics.NewPrometheusMetrics(prometheus.NewRegistry(), *service, uint16(*metricsPort))
	if err != nil {
		log.Fatalf("failed to create metrics: %v", err)
	}
	grpcServer := grpc.NewServer()
	eventService := event.NewService(
		event.WithSql(dbConn),
		event.WithDatabase(db),
		event.WithCache(redis),
		event.WithMetrics(metrics),
		event.WithMinEventNameLength(uint8(*minEventNameLength)),
		event.WithMaxEventNameLength(uint8(*maxEventNameLength)),
		event.WithMinRoundNameLength(uint8(*minRoundNameLength)),
		event.WithMaxRoundNameLength(uint8(*maxRoundNameLength)),
		event.WithMaxNumberOfRounds(uint8(*maxNumberOfRounds)),
		event.WithDefaultMaxPageLength(uint8(*defaultMaxPageLength)),
		event.WithMaxMaxPageLength(uint8(*maxMaxPageLength)),
	)
	api.RegisterEventServiceServer(grpcServer, eventService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
