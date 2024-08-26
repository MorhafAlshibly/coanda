package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/metric"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	// Flags set from command line/environment variables
	fs                   = ff.NewFlagSet("event")
	service              = fs.String('s', "service", "event", "the name of the service")
	port                 = fs.Uint('p', "port", 50055, "the default port to listen on")
	metricPort           = fs.Uint('m', "metricPort", 8085, "the port to serve metric on")
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
		fmt.Printf("failed to parse flags: %v", err)
		return
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	redis := cache.NewRedisCache(*cacheHost, *cachePassword, *cacheDB, *cacheExpiration)
	dbConn, err := sql.Open("mysql", *dsn)
	if err != nil {
		fmt.Printf("failed to open database: %v", err)
		return
	}
	defer dbConn.Close()
	db := model.New(dbConn)
	metric, err := metric.NewPrometheusMetric(prometheus.NewRegistry(), *service, uint16(*metricPort))
	if err != nil {
		fmt.Printf("failed to create metric: %v", err)
		return
	}
	grpcServer := grpc.NewServer()
	eventService := event.NewService(
		event.WithSql(dbConn),
		event.WithDatabase(db),
		event.WithCache(redis),
		event.WithMetric(metric),
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
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
