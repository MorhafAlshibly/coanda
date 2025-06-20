package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MorhafAlshibly/coanda/internal/handleMatchmaking"
	"github.com/MorhafAlshibly/coanda/internal/handleMatchmaking/model"
	lambdaFunc "github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/robfig/cron/v3"
)

var (
	fs                          = ff.NewFlagSet("handleMatchmaking")
	lambda                      = fs.BoolLong("lambda", "if running as a lambda function")
	cronSchedule                = fs.StringLong("cronSchedule", "*/5 * * * * *", "the cron schedule to run the handler (not for lambda)")
	dsn                         = fs.StringLong("dsn", "root:password@tcp(localhost:3306)", "the data source name for the database")
	eloWindowIncrementPerSecond = fs.UintLong("eloWindowIncrementPerSecond", 10, "the elo window increment per second elapsed since creation of the ticket")
	eloWindowMax                = fs.UintLong("eloWindowMax", 200, "the elo window max")
	limit                       = fs.UintLong("limit", 100, "the limit of tickets handled each loop, tweak this based on performance")
)

func main() {
	ctx := context.TODO()
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("HANDLE_MATCHMAKING"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		fmt.Printf("failed to parse flags: %v", err)
		return
	}
	dbConn, err := sql.Open("mysql", *dsn)
	if err != nil {
		fmt.Printf("failed to open database: %v", err)
		return
	}
	defer dbConn.Close()
	db := model.New(dbConn)
	// Create the app
	app := handleMatchmaking.NewApp(
		handleMatchmaking.WithSql(dbConn),
		handleMatchmaking.WithDatabase(db),
		handleMatchmaking.WithEloWindowIncrementPerSecond(uint16(*eloWindowIncrementPerSecond)),
		handleMatchmaking.WithEloWindowMax(uint16(*eloWindowMax)),
		handleMatchmaking.WithLimit(int32(*limit)),
	)
	if !*lambda {
		// Run the handler on a cron job
		c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))
		c.AddFunc(*cronSchedule, func() {
			if err := app.Handler(ctx); err != nil {
				fmt.Printf("failed to run handler: %v", err)
				return
			}
		})
		c.Start()
		// Wait for a signal to stop the cron job
		sig := make(chan os.Signal, 1)                    // Fix: Use a buffered channel
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM) // Fix: Replace os.Kill with syscall.SIGTERM
		<-sig
	} else {
		// Run the lambda if not running on a cron job
		lambdaFunc.Start(app.Handler)
	}
}
