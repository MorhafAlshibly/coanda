package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.74

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/bff/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateArena is the resolver for the CreateArena field.
func (r *mutationResolver) CreateArena(ctx context.Context, input model.CreateArenaRequest) (*model.CreateArenaResponse, error) {
	resp, err := r.matchmakingClient.CreateArena(ctx, &api.CreateArenaRequest{
		Name:                input.Name,
		MinPlayers:          input.MinPlayers,
		MaxPlayersPerTicket: input.MaxPlayersPerTicket,
		MaxPlayers:          input.MaxPlayers,
		Data:                input.Data,
	})
	if err != nil {
		return nil, err
	}
	return &model.CreateArenaResponse{
		Success: resp.Success,
		ID:      resp.Id,
		Error:   model.CreateArenaError(resp.Error.String()),
	}, nil
}

// UpdateArena is the resolver for the UpdateArena field.
func (r *mutationResolver) UpdateArena(ctx context.Context, input model.UpdateArenaRequest) (*model.UpdateArenaResponse, error) {
	if input.Arena == nil {
		input.Arena = &model.ArenaRequest{}
	}
	resp, err := r.matchmakingClient.UpdateArena(ctx, &api.UpdateArenaRequest{
		Arena: &api.ArenaRequest{
			Id:   input.Arena.ID,
			Name: input.Arena.Name,
		},
		MinPlayers:          input.MinPlayers,
		MaxPlayersPerTicket: input.MaxPlayersPerTicket,
		MaxPlayers:          input.MaxPlayers,
		Data:                input.Data,
	})
	if err != nil {
		return nil, err
	}
	return &model.UpdateArenaResponse{
		Success: resp.Success,
		Error:   model.UpdateArenaError(resp.Error.String()),
	}, nil
}

// CreateMatchmakingUser is the resolver for the CreateMatchmakingUser field.
func (r *mutationResolver) CreateMatchmakingUser(ctx context.Context, input model.CreateMatchmakingUserRequest) (*model.CreateMatchmakingUserResponse, error) {
	resp, err := r.matchmakingClient.CreateMatchmakingUser(ctx, &api.CreateMatchmakingUserRequest{
		ClientUserId: input.ClientUserID,
		Elo:          input.Elo,
		Data:         input.Data,
	})
	if err != nil {
		return nil, err
	}
	return &model.CreateMatchmakingUserResponse{
		Success: resp.Success,
		ID:      resp.Id,
		Error:   model.CreateMatchmakingUserError(resp.Error.String()),
	}, nil
}

// UpdateMatchmakingUser is the resolver for the UpdateMatchmakingUser field.
func (r *mutationResolver) UpdateMatchmakingUser(ctx context.Context, input model.UpdateMatchmakingUserRequest) (*model.UpdateMatchmakingUserResponse, error) {
	if input.MatchmakingUser == nil {
		input.MatchmakingUser = &model.MatchmakingUserRequest{}
	}
	resp, err := r.matchmakingClient.UpdateMatchmakingUser(ctx, &api.UpdateMatchmakingUserRequest{
		MatchmakingUser: &api.MatchmakingUserRequest{
			Id:           input.MatchmakingUser.ID,
			ClientUserId: input.MatchmakingUser.ClientUserID,
		},
		Data: input.Data,
	})
	if err != nil {
		return nil, err
	}
	return &model.UpdateMatchmakingUserResponse{
		Success: resp.Success,
		Error:   model.UpdateMatchmakingUserError(resp.Error.String()),
	}, nil
}

// CreateMatchmakingTicket is the resolver for the CreateMatchmakingTicket field.
func (r *mutationResolver) CreateMatchmakingTicket(ctx context.Context, input model.CreateMatchmakingTicketRequest) (*model.CreateMatchmakingTicketResponse, error) {
	matchmakingUsers := make([]*api.MatchmakingUserRequest, len(input.MatchmakingUsers))
	for i, mu := range input.MatchmakingUsers {
		if mu == nil {
			mu = &model.MatchmakingUserRequest{}
		}
		matchmakingUsers[i] = &api.MatchmakingUserRequest{
			Id:           mu.ID,
			ClientUserId: mu.ClientUserID,
		}
	}
	arenas := make([]*api.ArenaRequest, len(input.Arenas))
	for i, a := range input.Arenas {
		if a == nil {
			a = &model.ArenaRequest{}
		}
		arenas[i] = &api.ArenaRequest{
			Id:   a.ID,
			Name: a.Name,
		}
	}
	resp, err := r.matchmakingClient.CreateMatchmakingTicket(ctx, &api.CreateMatchmakingTicketRequest{
		MatchmakingUsers: matchmakingUsers,
		Arenas:           arenas,
		Data:             input.Data,
	})
	if err != nil {
		return nil, err
	}
	return &model.CreateMatchmakingTicketResponse{
		Success: resp.Success,
		ID:      resp.Id,
		Error:   model.CreateMatchmakingTicketError(resp.Error.String()),
	}, nil
}

// PollMatchmakingTicket is the resolver for the PollMatchmakingTicket field.
func (r *mutationResolver) PollMatchmakingTicket(ctx context.Context, input model.GetMatchmakingTicketRequest) (*model.PollMatchmakingTicketResponse, error) {
	if input.MatchmakingTicket == nil {
		input.MatchmakingTicket = &model.MatchmakingTicketRequest{}
	}
	if input.MatchmakingTicket.MatchmakingUser == nil {
		input.MatchmakingTicket.MatchmakingUser = &model.MatchmakingUserRequest{}
	}
	if input.UserPagination == nil {
		input.UserPagination = &model.Pagination{}
	}
	if input.ArenaPagination == nil {
		input.ArenaPagination = &model.Pagination{}
	}
	resp, err := r.matchmakingClient.PollMatchmakingTicket(ctx, &api.GetMatchmakingTicketRequest{
		MatchmakingTicket: &api.MatchmakingTicketRequest{
			Id: input.MatchmakingTicket.ID,
			MatchmakingUser: &api.MatchmakingUserRequest{
				Id:           input.MatchmakingTicket.MatchmakingUser.ID,
				ClientUserId: input.MatchmakingTicket.MatchmakingUser.ClientUserID,
			},
		},
		UserPagination: &api.Pagination{
			Page: input.UserPagination.Page,
			Max:  input.UserPagination.Max,
		},
		ArenaPagination: &api.Pagination{
			Page: input.ArenaPagination.Page,
			Max:  input.ArenaPagination.Max,
		},
	})
	if err != nil {
		return nil, err
	}
	var matchmakingTicket *model.MatchmakingTicket
	if resp.MatchmakingTicket != nil {
		matchmakingUsers := make([]*model.MatchmakingUser, len(resp.MatchmakingTicket.MatchmakingUsers))
		for i, mu := range resp.MatchmakingTicket.MatchmakingUsers {
			matchmakingUsers[i] = &model.MatchmakingUser{
				ID:           mu.Id,
				ClientUserID: mu.ClientUserId,
				Data:         mu.Data,
				Elo:          mu.Elo,
				CreatedAt:    mu.CreatedAt,
				UpdatedAt:    mu.UpdatedAt,
			}
		}
		arenas := make([]*model.Arena, len(resp.MatchmakingTicket.Arenas))
		for i, a := range resp.MatchmakingTicket.Arenas {
			arenas[i] = &model.Arena{
				ID:                  a.Id,
				Name:                a.Name,
				MinPlayers:          a.MinPlayers,
				MaxPlayersPerTicket: a.MaxPlayersPerTicket,
				MaxPlayers:          a.MaxPlayers,
				Data:                a.Data,
				CreatedAt:           a.CreatedAt,
				UpdatedAt:           a.UpdatedAt,
			}
		}
		matchmakingTicket = &model.MatchmakingTicket{
			ID:               resp.MatchmakingTicket.Id,
			MatchmakingUsers: matchmakingUsers,
			Arenas:           arenas,
			MatchID:          resp.MatchmakingTicket.MatchId,
			Status:           model.MatchmakingTicketStatus(resp.MatchmakingTicket.Status.String()),
			Data:             resp.MatchmakingTicket.Data,
			ExpiresAt:        resp.MatchmakingTicket.ExpiresAt,
			CreatedAt:        resp.MatchmakingTicket.CreatedAt,
			UpdatedAt:        resp.MatchmakingTicket.UpdatedAt,
		}
	}
	return &model.PollMatchmakingTicketResponse{
		Success:           resp.Success,
		MatchmakingTicket: matchmakingTicket,
		Error:             model.PollMatchmakingTicketError(resp.Error.String()),
	}, nil
}

// UpdateMatchmakingTicket is the resolver for the UpdateMatchmakingTicket field.
func (r *mutationResolver) UpdateMatchmakingTicket(ctx context.Context, input model.UpdateMatchmakingTicketRequest) (*model.UpdateMatchmakingTicketResponse, error) {
	if input.MatchmakingTicket == nil {
		input.MatchmakingTicket = &model.MatchmakingTicketRequest{}
	}
	if input.MatchmakingTicket.MatchmakingUser == nil {
		input.MatchmakingTicket.MatchmakingUser = &model.MatchmakingUserRequest{}
	}
	resp, err := r.matchmakingClient.UpdateMatchmakingTicket(ctx, &api.UpdateMatchmakingTicketRequest{
		MatchmakingTicket: &api.MatchmakingTicketRequest{
			Id: input.MatchmakingTicket.ID,
			MatchmakingUser: &api.MatchmakingUserRequest{
				Id:           input.MatchmakingTicket.MatchmakingUser.ID,
				ClientUserId: input.MatchmakingTicket.MatchmakingUser.ClientUserID,
			},
		},
		Data: input.Data,
	})
	if err != nil {
		return nil, err
	}
	return &model.UpdateMatchmakingTicketResponse{
		Success: resp.Success,
		Error:   model.UpdateMatchmakingTicketError(resp.Error.String()),
	}, nil
}

// ExpireMatchmakingTicket is the resolver for the ExpireMatchmakingTicket field.
func (r *mutationResolver) ExpireMatchmakingTicket(ctx context.Context, input model.MatchmakingTicketRequest) (*model.ExpireMatchmakingTicketResponse, error) {
	if input.MatchmakingUser == nil {
		input.MatchmakingUser = &model.MatchmakingUserRequest{}
	}
	resp, err := r.matchmakingClient.ExpireMatchmakingTicket(ctx, &api.MatchmakingTicketRequest{
		Id: input.ID,
		MatchmakingUser: &api.MatchmakingUserRequest{
			Id:           input.MatchmakingUser.ID,
			ClientUserId: input.MatchmakingUser.ClientUserID,
		},
	})
	if err != nil {
		return nil, err
	}
	return &model.ExpireMatchmakingTicketResponse{
		Success: resp.Success,
		Error:   model.ExpireMatchmakingTicketError(resp.Error.String()),
	}, nil
}

// DeleteMatchmakingTicket is the resolver for the DeleteMatchmakingTicket field.
func (r *mutationResolver) DeleteMatchmakingTicket(ctx context.Context, input model.MatchmakingTicketRequest) (*model.DeleteMatchmakingTicketResponse, error) {
	if input.MatchmakingUser == nil {
		input.MatchmakingUser = &model.MatchmakingUserRequest{}
	}
	resp, err := r.matchmakingClient.DeleteMatchmakingTicket(ctx, &api.MatchmakingTicketRequest{
		Id: input.ID,
		MatchmakingUser: &api.MatchmakingUserRequest{
			Id:           input.MatchmakingUser.ID,
			ClientUserId: input.MatchmakingUser.ClientUserID,
		},
	})
	if err != nil {
		return nil, err
	}
	return &model.DeleteMatchmakingTicketResponse{
		Success: resp.Success,
		Error:   model.DeleteMatchmakingTicketError(resp.Error.String()),
	}, nil
}

// DeleteAllExpiredMatchmakingTickets is the resolver for the DeleteAllExpiredMatchmakingTickets field.
func (r *mutationResolver) DeleteAllExpiredMatchmakingTickets(ctx context.Context) (*model.DeleteAllExpiredMatchmakingTicketsResponse, error) {
	resp, err := r.matchmakingClient.DeleteAllExpiredMatchmakingTickets(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return &model.DeleteAllExpiredMatchmakingTicketsResponse{
		Success: resp.Success,
	}, nil
}

// StartMatch is the resolver for the StartMatch field.
func (r *mutationResolver) StartMatch(ctx context.Context, input model.StartMatchRequest) (*model.StartMatchResponse, error) {
	if input.Match == nil {
		input.Match = &model.MatchRequest{}
	}
	if input.Match.MatchmakingTicket == nil {
		input.Match.MatchmakingTicket = &model.MatchmakingTicketRequest{}
	}
	if input.Match.MatchmakingTicket.MatchmakingUser == nil {
		input.Match.MatchmakingTicket.MatchmakingUser = &model.MatchmakingUserRequest{}
	}
	resp, err := r.matchmakingClient.StartMatch(ctx, &api.StartMatchRequest{
		Match: &api.MatchRequest{
			Id: input.Match.ID,
			MatchmakingTicket: &api.MatchmakingTicketRequest{
				Id: input.Match.MatchmakingTicket.ID,
				MatchmakingUser: &api.MatchmakingUserRequest{
					Id:           input.Match.MatchmakingTicket.MatchmakingUser.ID,
					ClientUserId: input.Match.MatchmakingTicket.MatchmakingUser.ClientUserID,
				},
			},
		},
		StartTime: input.StartTime,
	})
	if err != nil {
		return nil, err
	}
	return &model.StartMatchResponse{
		Success: resp.Success,
		Error:   model.StartMatchError(resp.Error.String()),
	}, nil
}

// EndMatch is the resolver for the EndMatch field.
func (r *mutationResolver) EndMatch(ctx context.Context, input model.EndMatchRequest) (*model.EndMatchResponse, error) {
	if input.Match == nil {
		input.Match = &model.MatchRequest{}
	}
	if input.Match.MatchmakingTicket == nil {
		input.Match.MatchmakingTicket = &model.MatchmakingTicketRequest{}
	}
	if input.Match.MatchmakingTicket.MatchmakingUser == nil {
		input.Match.MatchmakingTicket.MatchmakingUser = &model.MatchmakingUserRequest{}
	}
	resp, err := r.matchmakingClient.EndMatch(ctx, &api.EndMatchRequest{
		Match: &api.MatchRequest{
			Id: input.Match.ID,
			MatchmakingTicket: &api.MatchmakingTicketRequest{
				Id: input.Match.MatchmakingTicket.ID,
				MatchmakingUser: &api.MatchmakingUserRequest{
					Id:           input.Match.MatchmakingTicket.MatchmakingUser.ID,
					ClientUserId: input.Match.MatchmakingTicket.MatchmakingUser.ClientUserID,
				},
			},
		},
		EndTime: input.EndTime,
	})
	if err != nil {
		return nil, err
	}
	return &model.EndMatchResponse{
		Success: resp.Success,
		Error:   model.EndMatchError(resp.Error.String()),
	}, nil
}

// UpdateMatch is the resolver for the UpdateMatch field.
func (r *mutationResolver) UpdateMatch(ctx context.Context, input model.UpdateMatchRequest) (*model.UpdateMatchResponse, error) {
	if input.Match == nil {
		input.Match = &model.MatchRequest{}
	}
	if input.Match.MatchmakingTicket == nil {
		input.Match.MatchmakingTicket = &model.MatchmakingTicketRequest{}
	}
	if input.Match.MatchmakingTicket.MatchmakingUser == nil {
		input.Match.MatchmakingTicket.MatchmakingUser = &model.MatchmakingUserRequest{}
	}
	resp, err := r.matchmakingClient.UpdateMatch(ctx, &api.UpdateMatchRequest{
		Match: &api.MatchRequest{
			Id: input.Match.ID,
			MatchmakingTicket: &api.MatchmakingTicketRequest{
				Id: input.Match.MatchmakingTicket.ID,
				MatchmakingUser: &api.MatchmakingUserRequest{
					Id:           input.Match.MatchmakingTicket.MatchmakingUser.ID,
					ClientUserId: input.Match.MatchmakingTicket.MatchmakingUser.ClientUserID,
				},
			},
		},
		Data: input.Data,
	})
	if err != nil {
		return nil, err
	}
	return &model.UpdateMatchResponse{
		Success: resp.Success,
		Error:   model.UpdateMatchError(resp.Error.String()),
	}, nil
}

// SetMatchPrivateServer is the resolver for the SetMatchPrivateServer field.
func (r *mutationResolver) SetMatchPrivateServer(ctx context.Context, input model.SetMatchPrivateServerRequest) (*model.SetMatchPrivateServerResponse, error) {
	if input.Match == nil {
		input.Match = &model.MatchRequest{}
	}
	if input.Match.MatchmakingTicket == nil {
		input.Match.MatchmakingTicket = &model.MatchmakingTicketRequest{}
	}
	if input.Match.MatchmakingTicket.MatchmakingUser == nil {
		input.Match.MatchmakingTicket.MatchmakingUser = &model.MatchmakingUserRequest{}
	}
	resp, err := r.matchmakingClient.SetMatchPrivateServer(ctx, &api.SetMatchPrivateServerRequest{
		Match: &api.MatchRequest{
			Id: input.Match.ID,
			MatchmakingTicket: &api.MatchmakingTicketRequest{
				Id: input.Match.MatchmakingTicket.ID,
				MatchmakingUser: &api.MatchmakingUserRequest{
					Id:           input.Match.MatchmakingTicket.MatchmakingUser.ID,
					ClientUserId: input.Match.MatchmakingTicket.MatchmakingUser.ClientUserID,
				},
			},
		},
		PrivateServerId: input.PrivateServerID,
	})
	if err != nil {
		return nil, err
	}
	return &model.SetMatchPrivateServerResponse{
		Success:         resp.Success,
		PrivateServerID: resp.PrivateServerId,
		Error:           model.SetMatchPrivateServerError(resp.Error.String()),
	}, nil
}

// DeleteMatch is the resolver for the DeleteMatch field.
func (r *mutationResolver) DeleteMatch(ctx context.Context, input model.MatchRequest) (*model.DeleteMatchResponse, error) {
	if input.MatchmakingTicket == nil {
		input.MatchmakingTicket = &model.MatchmakingTicketRequest{}
	}
	if input.MatchmakingTicket.MatchmakingUser == nil {
		input.MatchmakingTicket.MatchmakingUser = &model.MatchmakingUserRequest{}
	}
	resp, err := r.matchmakingClient.DeleteMatch(ctx, &api.MatchRequest{
		Id: input.ID,
		MatchmakingTicket: &api.MatchmakingTicketRequest{
			Id: input.MatchmakingTicket.ID,
			MatchmakingUser: &api.MatchmakingUserRequest{
				Id:           input.MatchmakingTicket.MatchmakingUser.ID,
				ClientUserId: input.MatchmakingTicket.MatchmakingUser.ClientUserID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &model.DeleteMatchResponse{
		Success: resp.Success,
		Error:   model.DeleteMatchError(resp.Error.String()),
	}, nil
}

// GetArena is the resolver for the GetArena field.
func (r *queryResolver) GetArena(ctx context.Context, input model.ArenaRequest) (*model.GetArenaResponse, error) {
	resp, err := r.matchmakingClient.GetArena(ctx, &api.ArenaRequest{
		Id:   input.ID,
		Name: input.Name,
	})
	if err != nil {
		return nil, err
	}
	var arena *model.Arena
	if resp.Arena != nil {
		arena = &model.Arena{
			ID:                  resp.Arena.Id,
			Name:                resp.Arena.Name,
			MinPlayers:          resp.Arena.MinPlayers,
			MaxPlayersPerTicket: resp.Arena.MaxPlayersPerTicket,
			MaxPlayers:          resp.Arena.MaxPlayers,
			Data:                resp.Arena.Data,
			CreatedAt:           resp.Arena.CreatedAt,
			UpdatedAt:           resp.Arena.UpdatedAt,
		}
	}

	return &model.GetArenaResponse{
		Success: resp.Success,
		Arena:   arena,
		Error:   model.GetArenaError(resp.Error.String()),
	}, nil
}

// GetArenas is the resolver for the GetArenas field.
func (r *queryResolver) GetArenas(ctx context.Context, input model.Pagination) (*model.GetArenasResponse, error) {
	resp, err := r.matchmakingClient.GetArenas(ctx, &api.Pagination{
		Page: input.Page,
		Max:  input.Max,
	})
	if err != nil {
		return nil, err
	}
	arenas := make([]*model.Arena, len(resp.Arenas))
	for i, a := range resp.Arenas {
		arenas[i] = &model.Arena{
			ID:                  a.Id,
			Name:                a.Name,
			MinPlayers:          a.MinPlayers,
			MaxPlayersPerTicket: a.MaxPlayersPerTicket,
			MaxPlayers:          a.MaxPlayers,
			Data:                a.Data,
			CreatedAt:           a.CreatedAt,
			UpdatedAt:           a.UpdatedAt,
		}
	}
	return &model.GetArenasResponse{
		Success: resp.Success,
		Arenas:  arenas,
	}, nil
}

// GetMatchmakingUser is the resolver for the GetMatchmakingUser field.
func (r *queryResolver) GetMatchmakingUser(ctx context.Context, input model.MatchmakingUserRequest) (*model.GetMatchmakingUserResponse, error) {
	resp, err := r.matchmakingClient.GetMatchmakingUser(ctx, &api.MatchmakingUserRequest{
		Id:           input.ID,
		ClientUserId: input.ClientUserID,
	})
	if err != nil {
		return nil, err
	}
	var matchmakingUser *model.MatchmakingUser
	if resp.MatchmakingUser != nil {
		matchmakingUser = &model.MatchmakingUser{
			ID:           resp.MatchmakingUser.Id,
			ClientUserID: resp.MatchmakingUser.ClientUserId,
			Data:         resp.MatchmakingUser.Data,
			Elo:          resp.MatchmakingUser.Elo,
			CreatedAt:    resp.MatchmakingUser.CreatedAt,
			UpdatedAt:    resp.MatchmakingUser.UpdatedAt,
		}
	}
	return &model.GetMatchmakingUserResponse{
		Success:         resp.Success,
		MatchmakingUser: matchmakingUser,
		Error:           model.GetMatchmakingUserError(resp.Error.String()),
	}, nil
}

// GetMatchmakingUsers is the resolver for the GetMatchmakingUsers field.
func (r *queryResolver) GetMatchmakingUsers(ctx context.Context, input model.Pagination) (*model.GetMatchmakingUsersResponse, error) {
	resp, err := r.matchmakingClient.GetMatchmakingUsers(ctx, &api.Pagination{
		Page: input.Page,
		Max:  input.Max,
	})
	if err != nil {
		return nil, err
	}
	matchmakingUsers := make([]*model.MatchmakingUser, len(resp.MatchmakingUsers))
	for i, mu := range resp.MatchmakingUsers {
		matchmakingUsers[i] = &model.MatchmakingUser{
			ID:           mu.Id,
			ClientUserID: mu.ClientUserId,
			Data:         mu.Data,
			Elo:          mu.Elo,
			CreatedAt:    mu.CreatedAt,
			UpdatedAt:    mu.UpdatedAt,
		}
	}
	return &model.GetMatchmakingUsersResponse{
		Success:          resp.Success,
		MatchmakingUsers: matchmakingUsers,
	}, nil
}

// GetMatchmakingTicket is the resolver for the GetMatchmakingTicket field.
func (r *queryResolver) GetMatchmakingTicket(ctx context.Context, input model.GetMatchmakingTicketRequest) (*model.GetMatchmakingTicketResponse, error) {
	if input.MatchmakingTicket == nil {
		input.MatchmakingTicket = &model.MatchmakingTicketRequest{}
	}
	if input.MatchmakingTicket.MatchmakingUser == nil {
		input.MatchmakingTicket.MatchmakingUser = &model.MatchmakingUserRequest{}
	}
	if input.UserPagination == nil {
		input.UserPagination = &model.Pagination{}
	}
	if input.ArenaPagination == nil {
		input.ArenaPagination = &model.Pagination{}
	}
	resp, err := r.matchmakingClient.GetMatchmakingTicket(ctx, &api.GetMatchmakingTicketRequest{
		MatchmakingTicket: &api.MatchmakingTicketRequest{
			Id: input.MatchmakingTicket.ID,
			MatchmakingUser: &api.MatchmakingUserRequest{
				Id:           input.MatchmakingTicket.MatchmakingUser.ID,
				ClientUserId: input.MatchmakingTicket.MatchmakingUser.ClientUserID,
			},
		},
		UserPagination: &api.Pagination{
			Page: input.UserPagination.Page,
			Max:  input.UserPagination.Max,
		},
		ArenaPagination: &api.Pagination{
			Page: input.ArenaPagination.Page,
			Max:  input.ArenaPagination.Max,
		},
	})
	if err != nil {
		return nil, err
	}

	var matchmakingTicket *model.MatchmakingTicket
	if resp.MatchmakingTicket != nil {
		matchmakingUsers := make([]*model.MatchmakingUser, len(resp.MatchmakingTicket.MatchmakingUsers))
		for i, mu := range resp.MatchmakingTicket.MatchmakingUsers {
			matchmakingUsers[i] = &model.MatchmakingUser{
				ID:           mu.Id,
				ClientUserID: mu.ClientUserId,
				Data:         mu.Data,
				Elo:          mu.Elo,
				CreatedAt:    mu.CreatedAt,
				UpdatedAt:    mu.UpdatedAt,
			}
		}
		arenas := make([]*model.Arena, len(resp.MatchmakingTicket.Arenas))
		for i, a := range resp.MatchmakingTicket.Arenas {
			arenas[i] = &model.Arena{
				ID:                  a.Id,
				Name:                a.Name,
				MinPlayers:          a.MinPlayers,
				MaxPlayersPerTicket: a.MaxPlayersPerTicket,
				MaxPlayers:          a.MaxPlayers,
				Data:                a.Data,
				CreatedAt:           a.CreatedAt,
				UpdatedAt:           a.UpdatedAt,
			}
		}
		matchmakingTicket = &model.MatchmakingTicket{
			ID:               resp.MatchmakingTicket.Id,
			MatchmakingUsers: matchmakingUsers,
			Arenas:           arenas,
			MatchID:          resp.MatchmakingTicket.MatchId,
			Status:           model.MatchmakingTicketStatus(resp.MatchmakingTicket.Status.String()),
			Data:             resp.MatchmakingTicket.Data,
			ExpiresAt:        resp.MatchmakingTicket.ExpiresAt,
			CreatedAt:        resp.MatchmakingTicket.CreatedAt,
			UpdatedAt:        resp.MatchmakingTicket.UpdatedAt,
		}
	}
	return &model.GetMatchmakingTicketResponse{
		Success:           resp.Success,
		MatchmakingTicket: matchmakingTicket,
		Error:             model.GetMatchmakingTicketError(resp.Error.String()),
	}, nil
}

// GetMatchmakingTickets is the resolver for the GetMatchmakingTickets field.
func (r *queryResolver) GetMatchmakingTickets(ctx context.Context, input model.GetMatchmakingTicketsRequest) (*model.GetMatchmakingTicketsResponse, error) {
	if input.MatchmakingUser == nil {
		input.MatchmakingUser = &model.MatchmakingUserRequest{}
	}
	if input.Pagination == nil {
		input.Pagination = &model.Pagination{}
	}
	if input.UserPagination == nil {
		input.UserPagination = &model.Pagination{}
	}
	if input.ArenaPagination == nil {
		input.ArenaPagination = &model.Pagination{}
	}
	var statuses []api.MatchmakingTicket_Status
	if input.Statuses != nil {
		statuses = make([]api.MatchmakingTicket_Status, len(input.Statuses))
		for i, s := range input.Statuses {
			statuses[i] = api.MatchmakingTicket_Status(api.MatchmakingTicket_Status_value[s.String()])
		}
	}
	resp, err := r.matchmakingClient.GetMatchmakingTickets(ctx, &api.GetMatchmakingTicketsRequest{
		MatchId: input.MatchID,
		MatchmakingUser: &api.MatchmakingUserRequest{
			Id:           input.MatchmakingUser.ID,
			ClientUserId: input.MatchmakingUser.ClientUserID,
		},
		Statuses: statuses,
		Pagination: &api.Pagination{
			Page: input.Pagination.Page,
			Max:  input.Pagination.Max,
		},
		UserPagination: &api.Pagination{
			Page: input.UserPagination.Page,
			Max:  input.UserPagination.Max,
		},
		ArenaPagination: &api.Pagination{
			Page: input.ArenaPagination.Page,
			Max:  input.ArenaPagination.Max,
		},
	})
	if err != nil {
		return nil, err
	}
	matchmakingTickets := make([]*model.MatchmakingTicket, len(resp.MatchmakingTickets))
	for i, mt := range resp.MatchmakingTickets {
		matchmakingUsers := make([]*model.MatchmakingUser, len(mt.MatchmakingUsers))
		for j, mu := range mt.MatchmakingUsers {
			matchmakingUsers[j] = &model.MatchmakingUser{
				ID:           mu.Id,
				ClientUserID: mu.ClientUserId,
				Data:         mu.Data,
				Elo:          mu.Elo,
				CreatedAt:    mu.CreatedAt,
				UpdatedAt:    mu.UpdatedAt,
			}
		}
		arenas := make([]*model.Arena, len(mt.Arenas))
		for j, a := range mt.Arenas {
			arenas[j] = &model.Arena{
				ID:                  a.Id,
				Name:                a.Name,
				MinPlayers:          a.MinPlayers,
				MaxPlayersPerTicket: a.MaxPlayersPerTicket,
				MaxPlayers:          a.MaxPlayers,
				Data:                a.Data,
				CreatedAt:           a.CreatedAt,
				UpdatedAt:           a.UpdatedAt,
			}
		}
		matchmakingTickets[i] = &model.MatchmakingTicket{
			ID:               mt.Id,
			MatchmakingUsers: matchmakingUsers,
			Arenas:           arenas,
			MatchID:          mt.MatchId,
			Status:           model.MatchmakingTicketStatus(mt.Status.String()),
			Data:             mt.Data,
			ExpiresAt:        mt.ExpiresAt,
			CreatedAt:        mt.CreatedAt,
			UpdatedAt:        mt.UpdatedAt,
		}
	}
	return &model.GetMatchmakingTicketsResponse{
		Success:            resp.Success,
		MatchmakingTickets: matchmakingTickets,
	}, nil
}

// GetMatch is the resolver for the GetMatch field.
func (r *queryResolver) GetMatch(ctx context.Context, input model.GetMatchRequest) (*model.GetMatchResponse, error) {
	if input.Match == nil {
		input.Match = &model.MatchRequest{}
	}
	if input.Match.MatchmakingTicket == nil {
		input.Match.MatchmakingTicket = &model.MatchmakingTicketRequest{}
	}
	if input.Match.MatchmakingTicket.MatchmakingUser == nil {
		input.Match.MatchmakingTicket.MatchmakingUser = &model.MatchmakingUserRequest{}
	}
	if input.TicketPagination == nil {
		input.TicketPagination = &model.Pagination{}
	}
	if input.UserPagination == nil {
		input.UserPagination = &model.Pagination{}
	}
	if input.ArenaPagination == nil {
		input.ArenaPagination = &model.Pagination{}
	}
	resp, err := r.matchmakingClient.GetMatch(ctx, &api.GetMatchRequest{
		Match: &api.MatchRequest{
			Id: input.Match.ID,
			MatchmakingTicket: &api.MatchmakingTicketRequest{
				Id: input.Match.MatchmakingTicket.ID,
				MatchmakingUser: &api.MatchmakingUserRequest{
					Id:           input.Match.MatchmakingTicket.MatchmakingUser.ID,
					ClientUserId: input.Match.MatchmakingTicket.MatchmakingUser.ClientUserID,
				},
			},
		},
		TicketPagination: &api.Pagination{
			Page: input.TicketPagination.Page,
			Max:  input.TicketPagination.Max,
		},
		UserPagination: &api.Pagination{
			Page: input.UserPagination.Page,
			Max:  input.UserPagination.Max,
		},
		ArenaPagination: &api.Pagination{
			Page: input.ArenaPagination.Page,
			Max:  input.ArenaPagination.Max,
		},
	})
	if err != nil {
		return nil, err
	}
	var match *model.Match
	if resp.Match != nil {
		matchmakingTickets := make([]*model.MatchmakingTicket, len(resp.Match.Tickets))
		for i, mt := range resp.Match.Tickets {
			matchmakingUsers := make([]*model.MatchmakingUser, len(mt.MatchmakingUsers))
			for j, mu := range mt.MatchmakingUsers {
				matchmakingUsers[j] = &model.MatchmakingUser{
					ID:           mu.Id,
					ClientUserID: mu.ClientUserId,
					Data:         mu.Data,
					Elo:          mu.Elo,
					CreatedAt:    mu.CreatedAt,
					UpdatedAt:    mu.UpdatedAt,
				}
			}
			arenas := make([]*model.Arena, len(mt.Arenas))
			for j, a := range mt.Arenas {
				arenas[j] = &model.Arena{
					ID:                  a.Id,
					Name:                a.Name,
					MinPlayers:          a.MinPlayers,
					MaxPlayersPerTicket: a.MaxPlayersPerTicket,
					MaxPlayers:          a.MaxPlayers,
					Data:                a.Data,
					CreatedAt:           a.CreatedAt,
					UpdatedAt:           a.UpdatedAt,
				}
			}
			matchmakingTickets[i] = &model.MatchmakingTicket{
				ID:               mt.Id,
				MatchmakingUsers: matchmakingUsers,
				Arenas:           arenas,
				MatchID:          mt.MatchId,
				Status:           model.MatchmakingTicketStatus(mt.Status.String()),
				Data:             mt.Data,
				ExpiresAt:        mt.ExpiresAt,
				CreatedAt:        mt.CreatedAt,
				UpdatedAt:        mt.UpdatedAt,
			}
		}
		var arena *model.Arena
		if resp.Match.Arena != nil {
			arena = &model.Arena{
				ID:                  resp.Match.Arena.Id,
				Name:                resp.Match.Arena.Name,
				MinPlayers:          resp.Match.Arena.MinPlayers,
				MaxPlayersPerTicket: resp.Match.Arena.MaxPlayersPerTicket,
				MaxPlayers:          resp.Match.Arena.MaxPlayers,
				Data:                resp.Match.Arena.Data,
				CreatedAt:           resp.Match.Arena.CreatedAt,
				UpdatedAt:           resp.Match.Arena.UpdatedAt,
			}
		}
		match = &model.Match{
			ID:              resp.Match.Id,
			Arena:           arena,
			Tickets:         matchmakingTickets,
			PrivateServerID: resp.Match.PrivateServerId,
			Status:          model.MatchStatus(resp.Match.Status.String()),
			Data:            resp.Match.Data,
			LockedAt:        resp.Match.LockedAt,
			StartedAt:       resp.Match.StartedAt,
			EndedAt:         resp.Match.EndedAt,
			CreatedAt:       resp.Match.CreatedAt,
			UpdatedAt:       resp.Match.UpdatedAt,
		}
	}
	return &model.GetMatchResponse{
		Success: resp.Success,
		Match:   match,
		Error:   model.GetMatchError(resp.Error.String()),
	}, nil
}

// GetMatches is the resolver for the GetMatches field.
func (r *queryResolver) GetMatches(ctx context.Context, input model.GetMatchesRequest) (*model.GetMatchesResponse, error) {
	if input.Arena == nil {
		input.Arena = &model.ArenaRequest{}
	}
	if input.MatchmakingUser == nil {
		input.MatchmakingUser = &model.MatchmakingUserRequest{}
	}
	if input.Pagination == nil {
		input.Pagination = &model.Pagination{}
	}
	if input.TicketPagination == nil {
		input.TicketPagination = &model.Pagination{}
	}
	if input.UserPagination == nil {
		input.UserPagination = &model.Pagination{}
	}
	if input.ArenaPagination == nil {
		input.ArenaPagination = &model.Pagination{}
	}
	var statuses []api.Match_Status
	if input.Statuses != nil {
		statuses = make([]api.Match_Status, len(input.Statuses))
		for i, s := range input.Statuses {
			statuses[i] = api.Match_Status(api.Match_Status_value[s.String()])
		}
	}
	resp, err := r.matchmakingClient.GetMatches(ctx, &api.GetMatchesRequest{
		Arena: &api.ArenaRequest{
			Id:   input.Arena.ID,
			Name: input.Arena.Name,
		},
		MatchmakingUser: &api.MatchmakingUserRequest{
			Id:           input.MatchmakingUser.ID,
			ClientUserId: input.MatchmakingUser.ClientUserID,
		},
		Statuses: statuses,
		Pagination: &api.Pagination{
			Page: input.Pagination.Page,
			Max:  input.Pagination.Max,
		},
		TicketPagination: &api.Pagination{
			Page: input.TicketPagination.Page,
			Max:  input.TicketPagination.Max,
		},
		UserPagination: &api.Pagination{
			Page: input.UserPagination.Page,
			Max:  input.UserPagination.Max,
		},
		ArenaPagination: &api.Pagination{
			Page: input.ArenaPagination.Page,
			Max:  input.ArenaPagination.Max,
		},
	})
	if err != nil {
		return nil, err
	}
	matches := make([]*model.Match, len(resp.Matches))
	for i, m := range resp.Matches {
		matchmakingTickets := make([]*model.MatchmakingTicket, len(m.Tickets))
		for j, mt := range m.Tickets {
			matchmakingUsers := make([]*model.MatchmakingUser, len(mt.MatchmakingUsers))
			for k, mu := range mt.MatchmakingUsers {
				matchmakingUsers[k] = &model.MatchmakingUser{
					ID:           mu.Id,
					ClientUserID: mu.ClientUserId,
					Data:         mu.Data,
					Elo:          mu.Elo,
					CreatedAt:    mu.CreatedAt,
					UpdatedAt:    mu.UpdatedAt,
				}
			}
			arenas := make([]*model.Arena, len(mt.Arenas))
			for k, a := range mt.Arenas {
				arenas[k] = &model.Arena{
					ID:                  a.Id,
					Name:                a.Name,
					MinPlayers:          a.MinPlayers,
					MaxPlayersPerTicket: a.MaxPlayersPerTicket,
					MaxPlayers:          a.MaxPlayers,
					Data:                a.Data,
					CreatedAt:           a.CreatedAt,
					UpdatedAt:           a.UpdatedAt,
				}
			}
			matchmakingTickets[j] = &model.MatchmakingTicket{
				ID:               mt.Id,
				MatchmakingUsers: matchmakingUsers,
				Arenas:           arenas,
				MatchID:          mt.MatchId,
				Status:           model.MatchmakingTicketStatus(mt.Status.String()),
				Data:             mt.Data,
				ExpiresAt:        mt.ExpiresAt,
				CreatedAt:        mt.CreatedAt,
				UpdatedAt:        mt.UpdatedAt,
			}
		}
		var arena *model.Arena
		if m.Arena != nil {
			arena = &model.Arena{
				ID:                  m.Arena.Id,
				Name:                m.Arena.Name,
				MinPlayers:          m.Arena.MinPlayers,
				MaxPlayersPerTicket: m.Arena.MaxPlayersPerTicket,
				MaxPlayers:          m.Arena.MaxPlayers,
				Data:                m.Arena.Data,
				CreatedAt:           m.Arena.CreatedAt,
				UpdatedAt:           m.Arena.UpdatedAt,
			}
		}
		matches[i] = &model.Match{
			ID:        m.Id,
			Arena:     arena,
			Tickets:   matchmakingTickets,
			Status:    model.MatchStatus(m.Status.String()),
			Data:      m.Data,
			LockedAt:  m.LockedAt,
			StartedAt: m.StartedAt,
			EndedAt:   m.EndedAt,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}
	return &model.GetMatchesResponse{
		Success: resp.Success,
		Matches: matches,
	}, nil
}
