package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/metric"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	// Flags set from command line/environment variables
	fs                   = ff.NewFlagSet("team")
	service              = fs.String('s', "service", "team", "the name of the service")
	port                 = fs.Uint('p', "port", 50053, "the default port to listen on")
	metricPort           = fs.Uint('m', "metricPort", 8083, "the port to serve metric on")
	dsn                  = fs.StringLong("dsn", "root:password@tcp(localhost:3306)", "the data source name for the database")
	cacheHost            = fs.StringLong("cacheHost", "localhost:6379", "the connection string to the cache")
	cachePassword        = fs.StringLong("cachePassword", "", "the password to the cache")
	cacheDB              = fs.IntLong("cacheDB", 0, "the database to use in the cache")
	cacheExpiration      = fs.DurationLong("cacheExpiration", 5*time.Second, "the expiration time for the cache")
	maxMembers           = fs.UintLong("maxMembers", 5, "the max members")
	minTeamNameLength    = fs.UintLong("minTeamNameLength", 3, "the min team name length")
	maxTeamNameLength    = fs.UintLong("maxTeamNameLength", 20, "the max team name length")
	defaultMaxPageLength = fs.UintLong("defaultMaxPageLength", 10, "the default max page length")
	maxMaxPageLength     = fs.UintLong("maxMaxPageLength", 100, "the max max page length")
)

func main() {
	_ = context.TODO()
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("TEAM"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
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
	teamService := team.NewService(
		team.WithSql(dbConn),
		team.WithDatabase(db),
		team.WithCache(redis),
		team.WithMetric(metric),
		team.WithMaxMembers(uint8(*maxMembers)),
		team.WithMinTeamNameLength(uint8(*minTeamNameLength)),
		team.WithMaxTeamNameLength(uint8(*maxTeamNameLength)),
		team.WithDefaultMaxPageLength(uint8(*defaultMaxPageLength)),
		team.WithMaxMaxPageLength(uint8(*maxMaxPageLength)),
	)
	api.RegisterTeamServiceServer(grpcServer, teamService)
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
