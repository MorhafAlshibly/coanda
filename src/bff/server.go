package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/MorhafAlshibly/coanda/src/bff/graph"
	"github.com/MorhafAlshibly/coanda/src/bff/resolvers"
	"github.com/MorhafAlshibly/coanda/src/bff/services"
	"github.com/MorhafAlshibly/coanda/src/bff/storage"
)

const defaultPort = "8080"
const tableConn = "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	itemStore, err := storage.NewTableStorage(context.Background(), tableConn, "items")
	if err != nil {
		log.Fatal(err)
	}
	itemCache := storage.NewRedisCache("localhost:6379", "", 0, time.Second*120)
	itemService := services.NewItemService(itemStore, itemCache)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolvers.Resolver{ItemService: itemService}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
