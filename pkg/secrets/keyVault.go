package secrets

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

type KeyVault struct {
	client *azsecrets.Client
}

func NewKeyVault(cred *azidentity.DefaultAzureCredential, connection string) (*KeyVault, error) {
	client, err := azsecrets.NewClient(connection, cred, nil)
	if err != nil {
		return nil, err
	}
	return &KeyVault{client: client}, nil
}

func (k *KeyVault) GetSecret(ctx context.Context, name string, version *string) (string, error) {
	if version == nil {
		version = new(string)
		*version = ""
	}
	secret, err := k.client.GetSecret(ctx, name, *version, nil)
	if err != nil {
		return "", err
	}
	return *secret.Value, nil
}
