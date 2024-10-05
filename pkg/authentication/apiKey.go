package authentication

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/argon2"
)

type ApiKey struct {
	apiKeyHeader string
	hashedApiKey string
}

type Argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
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

func NewApiKeyAuthentication(options ...func(*ApiKey)) *ApiKey {
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
		if apiKey == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		authorized, err := a.compareApiKeyAndHashedApiKey(apiKey)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if !authorized {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *ApiKey) compareApiKeyAndHashedApiKey(apiKey string) (bool, error) {
	params, salt, hash, err := a.decodeHashedApiKey()
	if err != nil {
		return false, err
	}
	encodedHash := argon2.IDKey([]byte(apiKey), salt, params.iterations, params.memory, params.parallelism, params.keyLength)
	if subtle.ConstantTimeCompare(hash, encodedHash) == 1 {
		return true, nil
	}
	return false, nil
}

func (a *ApiKey) decodeHashedApiKey() (*Argon2Params, []byte, []byte, error) {
	values := strings.Split(a.hashedApiKey, "$")
	if len(values) != 6 {
		return nil, nil, nil, errors.New("invalid hashed api key")
	}
	var version int
	_, err := fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible version")
	}
	salt, err := base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params := &Argon2Params{}
	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &params.memory, &params.iterations, &params.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}
	params.saltLength = uint32(len(salt))
	hash, err := base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.keyLength = uint32(len(hash))
	return params, salt, hash, nil
}
