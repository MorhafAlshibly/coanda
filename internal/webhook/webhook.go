package webhook

import (
	"fmt"
	"io"
	"net/http"
)

type Service struct {
	uri string
}

func WithUri(uri string) func(*Service) {
	return func(input *Service) {
		input.uri = uri
	}
}

func NewService(options ...func(*Service)) *Service {
	service := Service{}
	for _, option := range options {
		option(&service)
	}
	return &service
}

func (s *Service) Handler(w http.ResponseWriter, r *http.Request) {
	// The webhhook contains data in the url, get the remaining url after the first / as a string
	webhookUriData := r.URL.EscapedPath()
	webhookUri := fmt.Sprintf("%s%s", s.uri, webhookUriData)
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
}
