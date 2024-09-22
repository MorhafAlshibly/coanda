package authentication

import "net/http"

type ApiKey struct {
	apiKeyHeader string
	hashedApiKey string
}

func WithApiKeyHeader(apiKeyHeader string) func(*ApiKey) {
	return func(a *ApiKey) {
		a.apiKeyHeader = apiKeyHeader
	}
}

func WithHashedApiKey(hashedApiKey string) func(*ApiKey) {
	return func(a *ApiKey) {
		a.hashedApiKey = hashedApiKey
	}
}

func NewAuthentication(options ...func(*ApiKey)) *ApiKey {
	a := &ApiKey{
		apiKeyHeader: "X-API-KEY",
	}
	for _, option := range options {
		option(a)
	}
	return a
}

func (a *ApiKey) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(a.apiKeyHeader)
		if apiKey != a.hashedApiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
