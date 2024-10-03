package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/bff"
	"github.com/MorhafAlshibly/coanda/internal/bff/resolver"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// Flags set from command line/environment variables
	fs               = ff.NewFlagSet("bff")
	port             = fs.Uint('p', "port", 8080, "the default port to listen on")
	enablePlayground = fs.BoolLong("enablePlayground", "enable the graphql playground")
	apiKeyHeader     = fs.StringLong("apiKeyHeader", "X-API-KEY", "the header key for the api key")
	hashedApiKey     = fs.StringLong("hashedApiKey", "", "the hashed api key")
	itemHost         = fs.StringLong("itemHost", "localhost:50051", "the endpoint of the item service")
	recordHost       = fs.StringLong("recordHost", "localhost:50052", "the endpoint of the record service")
	teamHost         = fs.StringLong("teamHost", "localhost:50053", "the endpoint of the team service")
	tournamentHost   = fs.StringLong("tournamentHost", "localhost:50054", "the endpoint of the tournament service")
	eventHost        = fs.StringLong("eventHost", "localhost:50055", "the endpoint of the event service")
	matchmakingHost  = fs.StringLong("matchmakingHost", "localhost:50056", "the endpoint of the matchmaking service")
	taskHost         = fs.StringLong("taskHost", "localhost:50057", "the endpoint of the task service")
	webhookHost      = fs.StringLong("webhookHost", "localhost:50058", "the endpoint of the webhook service")
	connOpts         = grpc.WithTransportCredentials(insecure.NewCredentials())
)

func main() {
	_ = context.TODO()
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("BFF"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		fmt.Printf("failed to parse flags: %v", err)
		return
	}
	itemConn, err := grpc.Dial(*itemHost, connOpts)
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		return
	}
	defer itemConn.Close()
	itemClient := api.NewItemServiceClient(itemConn)
	teamConn, err := grpc.Dial(*teamHost, connOpts)
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		return
	}
	defer teamConn.Close()
	teamClient := api.NewTeamServiceClient(teamConn)
	recordConn, err := grpc.Dial(*recordHost, connOpts)
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		return
	}
	defer recordConn.Close()
	recordClient := api.NewRecordServiceClient(recordConn)
	tournamentConn, err := grpc.Dial(*tournamentHost, connOpts)
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		return
	}
	defer tournamentConn.Close()
	tournamentClient := api.NewTournamentServiceClient(tournamentConn)
	eventConn, err := grpc.Dial(*eventHost, connOpts)
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		return
	}
	defer eventConn.Close()
	eventClient := api.NewEventServiceClient(eventConn)
	matchmakingConn, err := grpc.Dial(*matchmakingHost, connOpts)
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		return
	}
	defer matchmakingConn.Close()
	matchmakingClient := api.NewMatchmakingServiceClient(matchmakingConn)
	taskConn, err := grpc.Dial(*taskHost, connOpts)
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		return
	}
	defer taskConn.Close()
	taskClient := api.NewTaskServiceClient(taskConn)
	webhookConn, err := grpc.Dial(*webhookHost, connOpts)
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		return
	}
	defer webhookConn.Close()
	webhookClient := api.NewWebhookServiceClient(webhookConn)
	resolver := resolver.NewResolver(&resolver.NewResolverInput{
		ItemClient:        itemClient,
		TeamClient:        teamClient,
		RecordClient:      recordClient,
		TournamentClient:  tournamentClient,
		EventClient:       eventClient,
		MatchmakingClient: matchmakingClient,
		TaskClient:        taskClient,
		WebhookClient:     webhookClient,
	})
	srv := handler.NewDefaultServer(bff.NewExecutableSchema(bff.Config{Resolvers: resolver}))
	if *enablePlayground {
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}
	http.Handle("/query", srv)
	err = http.ListenAndServe(":"+fmt.Sprintf("%d", *port), nil)
	if err != nil {
		fmt.Printf("failed to listen and serve: %v", err)
		return
	}
}
