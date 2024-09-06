package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// The webhhook contains data in the url, get the remaining url after the first / as a string
		webhookUriData := r.URL.EscapedPath()
		webhookUri := fmt.Sprintf("%s%s", *uri, webhookUriData)
		response, err := http.Post(webhookUri, "application/json", r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer response.Body.Close()
		// Return the exact response from the webhook
		w.WriteHeader(response.StatusCode)
		// Return the exact response body from the webhook
		body, err := io.ReadAll(response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, err := w.Write(body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	})
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
