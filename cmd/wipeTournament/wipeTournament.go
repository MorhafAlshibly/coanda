package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/MorhafAlshibly/coanda/internal/wipeTournament"
	"github.com/MorhafAlshibly/coanda/internal/wipeTournament/model"
	lambdaFunc "github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/peterbourgon/ff/v4"
	"github.com/robfig/cron/v3"
)

var (
	fs                      = ff.NewFlagSet("wipeTournament")
	lambda                  = fs.BoolLong("lambda", "if running as a lambda function")
	cronSchedule            = fs.StringLong("cronSchedule", "* * * * *", "the cron schedule to run the handler")
	dsn                     = fs.StringLong("dsn", "root:password@tcp(localhost:3306)", "the data source name for the database")
	dailyTournamentMinute   = fs.UintLong("dailyTournamentMinute", 0, "the minute of the day to start the daily tournament")
	weeklyTournamentMinute  = fs.UintLong("weeklyTournamentMinute", 0, "the minute of the week to start the weekly tournament")
	weeklyTournamentDay     = fs.UintLong("weeklyTournamentDay", 0, "the day of the week to start the weekly tournament")
	monthlyTournamentMinute = fs.UintLong("monthlyTournamentMinute", 0, "the minute of the month to start the monthly tournament")
	monthlyTournamentDay    = fs.UintLong("monthlyTournamentDay", 1, "the day of the month to start the monthly tournament")
)

func main() {
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("WIPE_TOURNAMENT"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}
	dbConn, err := sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer dbConn.Close()
	db := model.New(dbConn)
	// Create the app
	app := wipeTournament.NewApp(
		wipeTournament.WithSql(dbConn),
		wipeTournament.WithDatabase(db),
		wipeTournament.WithDailyTournamentMinute(uint16(*dailyTournamentMinute)),
		wipeTournament.WithWeeklyTournamentMinute(uint16(*weeklyTournamentMinute)),
		wipeTournament.WithWeeklyTournamentDay(time.Weekday(*weeklyTournamentDay)),
		wipeTournament.WithMonthlyTournamentMinute(uint16(*monthlyTournamentMinute)),
		wipeTournament.WithMonthlyTournamentDay(uint8(*monthlyTournamentDay)),
	)
	if !*lambda {
		// Run the handler on a cron job
		c := cron.New()
		c.AddFunc(*cronSchedule, func() {
			if err := app.Handler(context.Background()); err != nil {
				log.Fatalf("failed to run handler: %v", err)
			}
		})
		c.Start()
		return
	}
	// Run the lambda if not running on a cron job
	lambdaFunc.Start(app.Handler)
}
