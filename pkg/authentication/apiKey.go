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
	apiKeyHeader        string
	decodedHashedApiKey DecodedHash
}

type DecodedHash struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	salt        []byte
	hash        []byte
}

func WithApiKeyHeader(apiKeyHeader string) func(*ApiKey) {
	return func(a *ApiKey) {
		a.apiKeyHeader = apiKeyHeader
	}
}

func WithHashedApiKey(hashedApiKey string) func(*ApiKey) {
	return func(a *ApiKey) {
		decodedHashedApiKey, err := decodeHashedApiKey(hashedApiKey)
		if err != nil {
			panic(err)
		}
		a.decodedHashedApiKey = *decodedHashedApiKey
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
	encodedAndHashedApiKey := argon2.IDKey([]byte(apiKey), a.decodedHashedApiKey.salt, a.decodedHashedApiKey.iterations, a.decodedHashedApiKey.memory, a.decodedHashedApiKey.parallelism, uint32(len(a.decodedHashedApiKey.hash)))
	if subtle.ConstantTimeCompare(a.decodedHashedApiKey.hash, encodedAndHashedApiKey) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHashedApiKey(hashedApiKey string) (*DecodedHash, error) {
	values := strings.Split(hashedApiKey, "$")
	if len(values) != 6 {
		return nil, errors.New("invalid hashed api key")
	}
	var version int
	_, err := fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return nil, err
	}
	if version != argon2.Version {
		return nil, errors.New("incompatible version")
	}
	decodedHash := &DecodedHash{}
	decodedHash.salt, err = base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return nil, err
	}
	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &decodedHash.memory, &decodedHash.iterations, &decodedHash.parallelism)
	if err != nil {
		return nil, err
	}
	decodedHash.hash, err = base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return nil, err
	}
	return decodedHash, nil
}
