package resolvers

import (
	"github.com/MorhafAlshibly/coanda/src/bff/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ItemService *services.ItemService
}
