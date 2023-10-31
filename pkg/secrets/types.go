package secrets

import "context"

type Secreter interface {
	GetSecret(ctx context.Context, name string, version *string) (string, error)
}

type MockSecreter struct {
	GetSecretFunc func(ctx context.Context, name string, version *string) (string, error)
}
