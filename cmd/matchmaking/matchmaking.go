package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/metric"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	// Flags set from command line/environment variables
	fs                   = ff.NewFlagSet("matchmaking")
	service              = fs.String('s', "service", "matchmaking", "the name of the service")
	port                 = fs.Uint('p', "port", 50056, "the default port to listen on")
	metricPort           = fs.Uint('m', "metricPort", 8086, "the port to serve metric on")
	dsn                  = fs.StringLong("dsn", "root:password@tcp(localhost:3306)", "the data source name for the database")
	cacheHost            = fs.StringLong("cacheHost", "localhost:6379", "the connection string to the cache")
	cachePassword        = fs.StringLong("cachePassword", "", "the password to the cache")
	cacheDB              = fs.IntLong("cacheDB", 0, "the database to use in the cache")
	cacheExpiration      = fs.DurationLong("cacheExpiration", 5*time.Second, "the expiration time for the cache")
	minArenaNameLength   = fs.UintLong("minArenaNameLength", 3, "the min arena name length")
	maxArenaNameLength   = fs.UintLong("maxArenaNameLength", 20, "the max arena name length")
	expiryTimeWindow     = fs.DurationLong("expiryTimeWindow", 5*time.Second, "the expiry time window")
	lockedAtBuffer       = fs.DurationLong("lockedAtBuffer", 10*time.Second, "the locked at buffer")
	startTimeBuffer      = fs.DurationLong("startTimeBuffer", 5*time.Second, "the start time buffer")
	defaultMaxPageLength = fs.UintLong("defaultMaxPageLength", 10, "the default max page length")
	maxMaxPageLength     = fs.UintLong("maxMaxPageLength", 100, "the max max page length")
)

func main() {
	_ = context.TODO()
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("MATCHMAKING"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
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
	matchmakingService := matchmaking.NewService(
		matchmaking.WithSql(dbConn),
		matchmaking.WithDatabase(db),
		matchmaking.WithCache(redis),
		matchmaking.WithMetric(metric),
		matchmaking.WithMinArenaNameLength(uint8(*minArenaNameLength)),
		matchmaking.WithMaxArenaNameLength(uint8(*maxArenaNameLength)),
		matchmaking.WithExpiryTimeWindow(*expiryTimeWindow),
		matchmaking.WithLockedAtBuffer(*lockedAtBuffer),
		matchmaking.WithStartTimeBuffer(*startTimeBuffer),
		matchmaking.WithDefaultMaxPageLength(uint8(*defaultMaxPageLength)),
		matchmaking.WithMaxMaxPageLength(uint8(*maxMaxPageLength)),
	)
	api.RegisterMatchmakingServiceServer(grpcServer, matchmakingService)
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
