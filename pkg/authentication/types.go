package authentication

import "net/http"

type Authenticator interface {
	Middleware(http.Handler) http.Handler
}
