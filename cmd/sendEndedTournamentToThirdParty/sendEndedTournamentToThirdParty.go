package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MorhafAlshibly/coanda/internal/sendEndedTournamentToThirdParty"
	"github.com/MorhafAlshibly/coanda/internal/sendEndedTournamentToThirdParty/model"
	lambdaFunc "github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/robfig/cron/v3"
)

var (
	fs                      = ff.NewFlagSet("sendEndedTournamentToThirdParty")
	lambda                  = fs.BoolLong("lambda", "if running as a lambda function")
	cronSchedule            = fs.StringLong("cronSchedule", "*/5 * * * *", "the cron schedule to run the handler (not for lambda)")
	dsn                     = fs.StringLong("dsn", "root:password@tcp(localhost:3306)", "the data source name for the database")
	dailyTournamentMinute   = fs.UintLong("dailyTournamentMinute", 0, "the minute of the day to start the daily tournament")
	weeklyTournamentMinute  = fs.UintLong("weeklyTournamentMinute", 0, "the minute of the week to start the weekly tournament")
	weeklyTournamentDay     = fs.UintLong("weeklyTournamentDay", 0, "the day of the week to start the weekly tournament")
	monthlyTournamentMinute = fs.UintLong("monthlyTournamentMinute", 0, "the minute of the month to start the monthly tournament")
	monthlyTournamentDay    = fs.UintLong("monthlyTournamentDay", 1, "the day of the month to start the monthly tournament")
	thirdPartyUri           = fs.StringLong("thirdPartyUri", "http://localhost:8080", "the uri of the third party to send the tournament to")
	apiKeyHeader            = fs.StringLong("apiKeyHeader", "x-api-key", "the header to send the api key in")
	apiKey                  = fs.StringLong("apiKey", "password", "the api key to send to the third party")
	topLimit                = fs.UintLong("topLimit", 100, "the top limit of tournament leaderboard to send to the third party")
	limit                   = fs.UintLong("limit", 100, "the limit of tournaments per loop, tweak this based on performance")
)

func main() {
	ctx := context.TODO()
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("SEND_ENDED_TOURNAMENT_TO_THIRD_PARTY"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
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
	app := sendEndedTournamentToThirdParty.NewApp(
		sendEndedTournamentToThirdParty.WithSql(dbConn),
		sendEndedTournamentToThirdParty.WithDatabase(db),
		sendEndedTournamentToThirdParty.WithDailyTournamentMinute(uint16(*dailyTournamentMinute)),
		sendEndedTournamentToThirdParty.WithWeeklyTournamentMinute(uint16(*weeklyTournamentMinute)),
		sendEndedTournamentToThirdParty.WithWeeklyTournamentDay(time.Weekday(*weeklyTournamentDay)),
		sendEndedTournamentToThirdParty.WithMonthlyTournamentMinute(uint16(*monthlyTournamentMinute)),
		sendEndedTournamentToThirdParty.WithMonthlyTournamentDay(uint8(*monthlyTournamentDay)),
		sendEndedTournamentToThirdParty.WithThirdPartyUri(*thirdPartyUri),
		sendEndedTournamentToThirdParty.WithApiKeyHeader(*apiKeyHeader),
		sendEndedTournamentToThirdParty.WithApiKey(*apiKey),
		sendEndedTournamentToThirdParty.WithTopLimit(int32(*topLimit)),
		sendEndedTournamentToThirdParty.WithLimit(int32(*limit)),
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
