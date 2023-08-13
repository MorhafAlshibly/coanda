package bff

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
	"github.com/MorhafAlshibly/coanda/services/bff/graph"
	"github.com/MorhafAlshibly/coanda/services/bff/resolvers"
	"github.com/MorhafAlshibly/coanda/services/bff/services"
)

const defaultPort = "8080"
const tableConn = "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;"

func Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create the item store
	itemStore, err := storage.NewTableStorage(context.Background(), tableConn, "items")
	if err != nil {
		log.Fatal(err)
	}
	// Create the item cache
	itemCache := cache.NewRedisCache("localhost:6379", "", 0, time.Second*120)
	// Create the item service
	itemService := services.NewItemService(itemStore, itemCache)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolvers.Resolver{ItemService: itemService}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
