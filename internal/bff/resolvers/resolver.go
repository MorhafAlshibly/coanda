package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import "github.com/MorhafAlshibly/coanda/internal/bff/item"

type Resolver struct {
	// An item service is injected into the resolver to provide functionality
	ItemService *item.Service
}
