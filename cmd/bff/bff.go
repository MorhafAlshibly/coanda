package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/bff"
	"github.com/MorhafAlshibly/coanda/internal/bff/resolver"
	"github.com/MorhafAlshibly/coanda/pkg/flags"
	"github.com/MorhafAlshibly/coanda/pkg/secrets"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	fs               = ff.NewFlagSet("bff")
	port             = fs.Uint('p', "port", 8080, "the default port to listen on")
	itemHostSecret   = fs.StringLong("itemHostSecret", "", "the secret containing the endpoint of the item service")
	teamHostSecret   = fs.StringLong("teamHostSecret", "", "the secret containing the endpoint of the team service")
	recordHostSecret = fs.StringLong("recordHostSecret", "", "the secret containing the endpoint of the record service")
	itemHost         = fs.StringLong("itemHost", "localhost:50051", "the endpoint of the item service")
	teamHost         = fs.StringLong("teamHost", "localhost:50052", "the endpoint of the team service")
	recordHost       = fs.StringLong("recordHost", "localhost:50053", "the endpoint of the record service")
	connOpts         = grpc.WithTransportCredentials(nil)
)

func main() {
	ctx := context.TODO()
	gf, err := flags.GetGlobalFlags()
	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(gf.FlagSet))
		log.Fatalf("failed to parse global flags: %v", err)
	}
	err = ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("BFF"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		log.Fatalf("failed to parse flags: %v", err)
	}
	if *gf.Environment == "dev" {
		connOpts = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Fatalf("failed to create credential: %v", err)
		}
		secrets, err := secrets.NewKeyVault(cred, *gf.VaultConn)
		if err != nil {
			log.Fatalf("failed to create secrets: %v", err)
		}
		*itemHost, err = secrets.GetSecret(ctx, *itemHostSecret, nil)
		if err != nil {
			log.Fatalf("failed to get item host: %v", err)
		}
		*teamHost, err = secrets.GetSecret(ctx, *teamHostSecret, nil)
		if err != nil {
			log.Fatalf("failed to get team host: %v", err)
		}
		*recordHost, err = secrets.GetSecret(ctx, *recordHostSecret, nil)
		if err != nil {
			log.Fatalf("failed to get record host: %v", err)
		}
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
	log.Fatal(http.ListenAndServe(":"+fmt.Sprintf("%d", *port), nil))
}
