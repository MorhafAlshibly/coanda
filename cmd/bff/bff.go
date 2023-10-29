package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/bff"
	"github.com/MorhafAlshibly/coanda/internal/bff/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var defaultPort = flag.Uint("port", 8080, "default port to listen on")
var itemHost = flag.String("itemHost", "localhost:50051", "the endpoint of the item service")
var teamHost = flag.String("teamHost", "localhost:50052", "the endpoint of the team service")
var recordHost = flag.String("recordHost", "localhost:50053", "the endpoint of the record service")

var connOpts = grpc.WithTransportCredentials(insecure.NewCredentials())

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", *defaultPort)
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

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
