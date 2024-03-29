package main

import (
	"context"
	"fmt"
	"log"
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
	itemHost         = fs.StringLong("itemHost", "localhost:50051", "the endpoint of the item service")
	teamHost         = fs.StringLong("teamHost", "localhost:50052", "the endpoint of the team service")
	recordHost       = fs.StringLong("recordHost", "localhost:50053", "the endpoint of the record service")
	connOpts         = grpc.WithTransportCredentials(insecure.NewCredentials())
)

func main() {
	_ = context.TODO()
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("BFF"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		log.Fatalf("failed to parse flags: %v", err)
	}
	itemConn, err := grpc.Dial(*itemHost, connOpts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer itemConn.Close()
	itemClient := api.NewItemServiceClient(itemConn)
	teamConn, err := grpc.Dial(*teamHost, connOpts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer teamConn.Close()
	teamClient := api.NewTeamServiceClient(teamConn)
	recordConn, err := grpc.Dial(*recordHost, connOpts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer recordConn.Close()
	recordClient := api.NewRecordServiceClient(recordConn)
	resolver := resolver.NewResolver(&resolver.NewResolverInput{
		ItemClient:   itemClient,
		TeamClient:   teamClient,
		RecordClient: recordClient,
	})
	srv := handler.NewDefaultServer(bff.NewExecutableSchema(bff.Config{Resolvers: resolver}))
	if *enablePlayground == true {
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}
	http.Handle("/query", srv)
	log.Fatal(http.ListenAndServe(":"+fmt.Sprintf("%d", *port), nil))
}
