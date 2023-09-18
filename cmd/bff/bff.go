package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/ysugimoto/grpc-graphql-gateway/runtime"
)

const defaultPort = 8080
const tableConn = "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;"

func main() {
	mux := runtime.NewServeMux()

	if err := api.RegisterItemServiceGraphql(mux); err != nil {
		log.Fatalln(err)
	}
	http.Handle("/graphql", mux)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", defaultPort), nil))
}
