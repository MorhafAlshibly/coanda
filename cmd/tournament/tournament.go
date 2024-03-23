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
	"github.com/MorhafAlshibly/coanda/internal/tournament"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	// Flags set from command line/environment variables
	fs                      = ff.NewFlagSet("tournament")
	service                 = fs.String('s', "service", "tournament", "the name of the service")
	port                    = fs.Uint('p', "port", 50054, "the default port to listen on")
	metricsPort             = fs.Uint('m', "metricsPort", 8084, "the port to serve metrics on")
	dsn                     = fs.StringLong("dsn", "root:password@tcp(localhost:3306)", "the data source name for the database")
	cacheHost               = fs.StringLong("cacheHost", "localhost:6379", "the connection string to the cache")
	cachePassword           = fs.StringLong("cachePassword", "", "the password to the cache")
	cacheDB                 = fs.IntLong("cacheDB", 0, "the database to use in the cache")
	cacheExpiration         = fs.DurationLong("cacheExpiration", 5*time.Second, "the expiration time for the cache")
	dailyTournamentMinute   = fs.UintLong("dailyTournamentMinute", 0, "the minute of the day to start the daily tournament")
	weeklyTournamentMinute  = fs.UintLong("weeklyTournamentMinute", 0, "the minute of the week to start the weekly tournament")
	weeklyTournamentDay     = fs.UintLong("weeklyTournamentDay", 0, "the day of the week to start the weekly tournament")
	monthlyTournamentMinute = fs.UintLong("monthlyTournamentMinute", 0, "the minute of the month to start the monthly tournament")
	monthlyTournamentDay    = fs.UintLong("monthlyTournamentDay", 1, "the day of the month to start the monthly tournament")
	minTournamentNameLength = fs.UintLong("minTournamentNameLength", 3, "the min tournament name length")
	maxTournamentNameLength = fs.UintLong("maxTournamentNameLength", 20, "the max tournament name length")
	defaultMaxPageLength    = fs.UintLong("defaultMaxPageLength", 10, "the default max page length")
	maxMaxPageLength        = fs.UintLong("maxMaxPageLength", 100, "the max max page length")
)

func main() {
	_ = context.TODO()
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("TOURNAMENT"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
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
	tournamentService := tournament.NewService(
		tournament.WithSql(dbConn),
		tournament.WithDatabase(db),
		tournament.WithCache(redis),
		tournament.WithMetrics(metrics),
		tournament.WithMinTournamentNameLength(uint8(*minTournamentNameLength)),
		tournament.WithMaxTournamentNameLength(uint8(*maxTournamentNameLength)),
		tournament.WithDailyTournamentMinute(uint16(*dailyTournamentMinute)),
		tournament.WithWeeklyTournamentMinute(uint16(*weeklyTournamentMinute)),
		tournament.WithWeeklyTournamentDay(time.Weekday(*weeklyTournamentDay)),
		tournament.WithMonthlyTournamentMinute(uint16(*monthlyTournamentMinute)),
		tournament.WithMonthlyTournamentDay(uint8(*monthlyTournamentDay)),
		tournament.WithDefaultMaxPageLength(uint8(*defaultMaxPageLength)),
		tournament.WithMaxMaxPageLength(uint8(*maxMaxPageLength)),
	)
	api.RegisterTournamentServiceServer(grpcServer, tournamentService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
