package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"

	"github.com/MorhafAlshibly/coanda/internal/bff/model"
)

// CreateArena is the resolver for the CreateArena field.
func (r *mutationResolver) CreateArena(ctx context.Context, input model.CreateArenaRequest) (*model.CreateArenaResponse, error) {
	panic(fmt.Errorf("not implemented: CreateArena - CreateArena"))
}

// UpdateArena is the resolver for the UpdateArena field.
func (r *mutationResolver) UpdateArena(ctx context.Context, input model.UpdateArenaRequest) (*model.UpdateArenaResponse, error) {
	panic(fmt.Errorf("not implemented: UpdateArena - UpdateArena"))
}

// CreateMatchmakingUser is the resolver for the CreateMatchmakingUser field.
func (r *mutationResolver) CreateMatchmakingUser(ctx context.Context, input model.CreateMatchmakingUserRequest) (*model.CreateMatchmakingUserResponse, error) {
	panic(fmt.Errorf("not implemented: CreateMatchmakingUser - CreateMatchmakingUser"))
}

// UpdateMatchmakingUser is the resolver for the UpdateMatchmakingUser field.
func (r *mutationResolver) UpdateMatchmakingUser(ctx context.Context, input model.UpdateMatchmakingUserRequest) (*model.UpdateMatchmakingUserResponse, error) {
	panic(fmt.Errorf("not implemented: UpdateMatchmakingUser - UpdateMatchmakingUser"))
}

// SetMatchmakingUserElo is the resolver for the SetMatchmakingUserElo field.
func (r *mutationResolver) SetMatchmakingUserElo(ctx context.Context, input model.SetMatchmakingUserEloRequest) (*model.SetMatchmakingUserEloResponse, error) {
	panic(fmt.Errorf("not implemented: SetMatchmakingUserElo - SetMatchmakingUserElo"))
}

// CreateMatchmakingTicket is the resolver for the CreateMatchmakingTicket field.
func (r *mutationResolver) CreateMatchmakingTicket(ctx context.Context, input model.CreateMatchmakingTicketRequest) (*model.CreateMatchmakingTicketResponse, error) {
	panic(fmt.Errorf("not implemented: CreateMatchmakingTicket - CreateMatchmakingTicket"))
}

// PollMatchmakingTicket is the resolver for the PollMatchmakingTicket field.
func (r *mutationResolver) PollMatchmakingTicket(ctx context.Context, input model.MatchmakingTicketRequest) (*model.MatchmakingTicketResponse, error) {
	panic(fmt.Errorf("not implemented: PollMatchmakingTicket - PollMatchmakingTicket"))
}

// UpdateMatchmakingTicket is the resolver for the UpdateMatchmakingTicket field.
func (r *mutationResolver) UpdateMatchmakingTicket(ctx context.Context, input model.UpdateMatchmakingTicketRequest) (*model.UpdateMatchmakingTicketResponse, error) {
	panic(fmt.Errorf("not implemented: UpdateMatchmakingTicket - UpdateMatchmakingTicket"))
}

// ExpireMatchmakingTicket is the resolver for the ExpireMatchmakingTicket field.
func (r *mutationResolver) ExpireMatchmakingTicket(ctx context.Context, input model.MatchmakingTicketRequest) (*model.ExpireMatchmakingTicketResponse, error) {
	panic(fmt.Errorf("not implemented: ExpireMatchmakingTicket - ExpireMatchmakingTicket"))
}

// StartMatch is the resolver for the StartMatch field.
func (r *mutationResolver) StartMatch(ctx context.Context, input model.StartMatchRequest) (*model.StartMatchResponse, error) {
	panic(fmt.Errorf("not implemented: StartMatch - StartMatch"))
}

// EndMatch is the resolver for the EndMatch field.
func (r *mutationResolver) EndMatch(ctx context.Context, input model.EndMatchRequest) (*model.EndMatchResponse, error) {
	panic(fmt.Errorf("not implemented: EndMatch - EndMatch"))
}

// UpdateMatch is the resolver for the UpdateMatch field.
func (r *mutationResolver) UpdateMatch(ctx context.Context, input model.UpdateMatchRequest) (*model.UpdateMatchResponse, error) {
	panic(fmt.Errorf("not implemented: UpdateMatch - UpdateMatch"))
}

// GetArena is the resolver for the GetArena field.
func (r *queryResolver) GetArena(ctx context.Context, input model.ArenaRequest) (*model.GetArenaResponse, error) {
	panic(fmt.Errorf("not implemented: GetArena - GetArena"))
}

// GetArenas is the resolver for the GetArenas field.
func (r *queryResolver) GetArenas(ctx context.Context, input model.Pagination) (*model.GetArenasResponse, error) {
	panic(fmt.Errorf("not implemented: GetArenas - GetArenas"))
}

// GetMatchmakingUser is the resolver for the GetMatchmakingUser field.
func (r *queryResolver) GetMatchmakingUser(ctx context.Context, input model.MatchmakingUserRequest) (*model.GetMatchmakingUserResponse, error) {
	panic(fmt.Errorf("not implemented: GetMatchmakingUser - GetMatchmakingUser"))
}

// GetMatchmakingUsers is the resolver for the GetMatchmakingUsers field.
func (r *queryResolver) GetMatchmakingUsers(ctx context.Context, input model.Pagination) (*model.GetMatchmakingUsersResponse, error) {
	panic(fmt.Errorf("not implemented: GetMatchmakingUsers - GetMatchmakingUsers"))
}

// GetMatchmakingTicket is the resolver for the GetMatchmakingTicket field.
func (r *queryResolver) GetMatchmakingTicket(ctx context.Context, input model.MatchmakingTicketRequest) (*model.GetMatchmakingTicketResponse, error) {
	panic(fmt.Errorf("not implemented: GetMatchmakingTicket - GetMatchmakingTicket"))
}

// GetMatchmakingTickets is the resolver for the GetMatchmakingTickets field.
func (r *queryResolver) GetMatchmakingTickets(ctx context.Context, input model.GetMatchmakingTicketsRequest) (*model.GetMatchmakingTicketsResponse, error) {
	panic(fmt.Errorf("not implemented: GetMatchmakingTickets - GetMatchmakingTickets"))
}

// GetMatch is the resolver for the GetMatch field.
func (r *queryResolver) GetMatch(ctx context.Context, input model.MatchRequest) (*model.GetMatchResponse, error) {
	panic(fmt.Errorf("not implemented: GetMatch - GetMatch"))
}

// GetMatches is the resolver for the GetMatches field.
func (r *queryResolver) GetMatches(ctx context.Context, input model.GetMatchesRequest) (*model.GetMatchesResponse, error) {
	panic(fmt.Errorf("not implemented: GetMatches - GetMatches"))
}
