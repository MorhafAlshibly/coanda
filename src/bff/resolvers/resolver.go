package resolvers

import (
	"github.com/MorhafAlshibly/coanda/src/bff/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	// An item service is injected into the resolver to provide functionality
	ItemService *services.ItemService
}
