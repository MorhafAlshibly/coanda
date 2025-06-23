package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/MorhafAlshibly/coanda/internal/handleMatchmaking"
	"github.com/MorhafAlshibly/coanda/internal/handleMatchmaking/model"
	lambdaFunc "github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

var (
	fs                          = ff.NewFlagSet("handleMatchmaking")
	lambda                      = fs.BoolLong("lambda", "if running as a lambda function")
	tickerInterval              = fs.DurationLong("tickerInterval", 5*time.Second, "the interval to run the handler (not for lambda)")
	dsn                         = fs.StringLong("dsn", "root:password@tcp(localhost:3306)", "the data source name for the database")
	eloWindowIncrementPerSecond = fs.UintLong("eloWindowIncrementPerSecond", 10, "the elo window increment per second elapsed since creation of the ticket")
	eloWindowMax                = fs.UintLong("eloWindowMax", 600, "the elo window max")
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
		ticker := time.NewTicker(*tickerInterval)
		defer ticker.Stop() // Always stop ticker to release resources

		for t := range ticker.C {
			fmt.Printf("Running handler at %s\n", t)
			if err := app.Handler(ctx); err != nil {
				fmt.Printf("failed to run handler: %v\n", err)
				return
			}
		}
	} else {
		// Run the lambda if not running on a cron job
		lambdaFunc.Start(app.Handler)
	}
}
