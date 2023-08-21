package queue

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

// ServiceBus is a wrapper around the Azure Service Bus client
type ServiceBus struct {
	Client *azservicebus.Client
	Sender *azservicebus.Sender
}

// NewServiceBus creates a new ServiceBus client
func NewServiceBus(ctx context.Context, connection string, queueName string) (*ServiceBus, error) {
	client, err := azservicebus.NewClientFromConnectionString(connection, nil)
	if err != nil {
		return nil, err
	}
	sender, err := client.NewSender(queueName, nil)
	return &ServiceBus{Client: client, Sender: sender}, nil
}

// Enqueue adds a message to the queue
func (s *ServiceBus) Enqueue(ctx context.Context, data []byte) error {
	msg := &azservicebus.Message{
		Body: data,
	}
	return s.Sender.SendMessage(ctx, msg, nil)
}
