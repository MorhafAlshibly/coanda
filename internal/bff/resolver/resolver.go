package resolver

import (
	"github.com/MorhafAlshibly/coanda/api"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	itemClient       api.ItemServiceClient
	teamClient       api.TeamServiceClient
	recordClient     api.RecordServiceClient
	tournamentClient api.TournamentServiceClient
}

type NewResolverInput struct {
	ItemClient       api.ItemServiceClient
	TeamClient       api.TeamServiceClient
	RecordClient     api.RecordServiceClient
	TournamentClient api.TournamentServiceClient
}

func NewResolver(input *NewResolverInput) *Resolver {
	return &Resolver{
		itemClient:       input.ItemClient,
		teamClient:       input.TeamClient,
		recordClient:     input.RecordClient,
		tournamentClient: input.TournamentClient,
	}
}
