package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/MorhafAlshibly/coanda/internal/webhook"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

var (
	// Flags set from command line/environment variables
	fs      = ff.NewFlagSet("webhook")
	service = fs.String('s', "service", "webhook", "the name of the service")
	port    = fs.Uint('p', "port", 50058, "the default port to listen on")
	uri     = fs.StringLong("uri", "https://webhook.site", "the uri to send the webhook to without the trailing /")
)

func main() {
	_ = context.TODO()
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("WEBHOOK"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.PlainParser))
	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		fmt.Printf("failed to parse flags: %v", err)
		return
	}
	webhookService := webhook.NewService(webhook.WithUri(*uri))
	http.HandleFunc("/", webhookService.Handler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
