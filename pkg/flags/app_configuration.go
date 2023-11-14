package flags

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/appconfig/azappconfig"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/peterbourgon/ff/v4"
)

type AppConfiguration struct {
	Client *azappconfig.Client
}

func NewAppConfiguration(ctx context.Context, cred *azidentity.DefaultAzureCredential, connection string) (*AppConfiguration, error) {
	// Create a new App Configuration Client
	var client *azappconfig.Client
	var err error
	if cred == nil {
		client, err = azappconfig.NewClient(connection, cred, nil)
	} else {
		client, err = azappconfig.NewClientFromConnectionString(connection, nil)
	}
	if err != nil {
		return nil, err
	}
	return &AppConfiguration{Client: client}, nil
}

func (s *AppConfiguration) Get(ctx context.Context, key string) (string, error) {
	// Get the setting from the App Configuration service
	resp, err := s.Client.GetSetting(ctx, key, nil)
	if err != nil {
		return "", err
	}
	if resp.Value == nil {
		return "", nil
	}
	return *resp.Value, nil
}

func (s *AppConfiguration) Parse(ctx context.Context, fs *ff.FlagSet, prefix string) error {
	// Loop through flags and set them from the App Configuration service
	err := fs.WalkFlags(func(f ff.Flag) error {
		name, ok := f.GetLongName()
		if ok == false {
			return errors.New("flag does not have a long name")
		}
		// Get the value from the App Configuration service
		value, err := s.Get(ctx, prefix+name)
		if err != nil {
			return err
		}
		// If the value is not empty, set the flag
		if value != "" {
			err = f.SetValue(value)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
