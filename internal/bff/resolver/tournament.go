package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/bff/model"
)

// CreateTournamentUser is the resolver for the CreateTournamentUser field.
func (r *mutationResolver) CreateTournamentUser(ctx context.Context, input model.CreateTournamentUserRequest) (*model.CreateTournamentUserResponse, error) {
	resp, err := r.tournamentClient.CreateTournamentUser(ctx, &api.CreateTournamentUserRequest{
		Tournament: input.Tournament,
		Interval:   api.TournamentInterval(api.TournamentInterval_value[input.Interval.String()]),
		UserId:     input.UserID,
		Score:      input.Score,
		Data:       input.Data,
	})
	if err != nil {
		return nil, err
	}
	return &model.CreateTournamentUserResponse{
		Success: resp.Success,
		Error:   model.CreateTournamentUserError(resp.Error.String()),
		ID:      resp.Id,
	}, nil
}

// UpdateTournamentUser is the resolver for the UpdateTournamentUser field.
func (r *mutationResolver) UpdateTournamentUser(ctx context.Context, input model.UpdateTournamentUserRequest) (*model.UpdateTournamentUserResponse, error) {
	var tournamentUserRequest *api.TournamentUserRequest
	if input.Tournament != nil {
		tournamentUserRequest = &api.TournamentUserRequest{
			Id: input.Tournament.ID,
		}
		if input.Tournament.TournamentIntervalUserID != nil {
			tournamentUserRequest = &api.TournamentUserRequest{
				Id: input.Tournament.ID,
				TournamentIntervalUserId: &api.TournamentIntervalUserId{
					Tournament: input.Tournament.TournamentIntervalUserID.Tournament,
					Interval:   api.TournamentInterval(api.TournamentInterval_value[input.Tournament.TournamentIntervalUserID.Interval.String()]),
					UserId:     input.Tournament.TournamentIntervalUserID.UserID,
				},
			}
		}
	}
	resp, err := r.tournamentClient.UpdateTournamentUser(ctx, &api.UpdateTournamentUserRequest{
		Tournament:     tournamentUserRequest,
		Score:          input.Score,
		Data:           input.Data,
		IncrementScore: input.IncrementScore,
	})
	if err != nil {
		return nil, err
	}
	return &model.UpdateTournamentUserResponse{
		Success: resp.Success,
		Error:   model.UpdateTournamentUserError(resp.Error.String()),
	}, nil
}

// DeleteTournamentUser is the resolver for the DeleteTournamentUser field.
func (r *mutationResolver) DeleteTournamentUser(ctx context.Context, input model.TournamentUserRequest) (*model.TournamentUserResponse, error) {
	var tournamentIntervalUserId *api.TournamentIntervalUserId
	if input.TournamentIntervalUserID != nil {
		tournamentIntervalUserId = &api.TournamentIntervalUserId{
			Tournament: input.TournamentIntervalUserID.Tournament,
			Interval:   api.TournamentInterval(api.TournamentInterval_value[input.TournamentIntervalUserID.Interval.String()]),
			UserId:     input.TournamentIntervalUserID.UserID,
		}
	}
	resp, err := r.tournamentClient.DeleteTournamentUser(ctx, &api.TournamentUserRequest{
		Id:                       input.ID,
		TournamentIntervalUserId: tournamentIntervalUserId,
	})
	if err != nil {
		return nil, err
	}
	return &model.TournamentUserResponse{
		Success: resp.Success,
		Error:   model.TournamentUserError(resp.Error.String()),
	}, nil
}

// GetTournamentUser is the resolver for the GetTournamentUser field.
func (r *queryResolver) GetTournamentUser(ctx context.Context, input model.TournamentUserRequest) (*model.GetTournamentUserResponse, error) {
	var tournamentIntervalUserId *api.TournamentIntervalUserId
	if input.TournamentIntervalUserID != nil {
		tournamentIntervalUserId = &api.TournamentIntervalUserId{
			Tournament: input.TournamentIntervalUserID.Tournament,
			Interval:   api.TournamentInterval(api.TournamentInterval_value[input.TournamentIntervalUserID.Interval.String()]),
			UserId:     input.TournamentIntervalUserID.UserID,
		}
	}
	resp, err := r.tournamentClient.GetTournamentUser(ctx, &api.TournamentUserRequest{
		Id:                       input.ID,
		TournamentIntervalUserId: tournamentIntervalUserId,
	})
	if err != nil {
		return nil, err
	}
	var tournamentUser *model.TournamentUser
	if resp.TournamentUser != nil {
		tournamentUser = &model.TournamentUser{
			ID:                  resp.TournamentUser.Id,
			Tournament:          resp.TournamentUser.Tournament,
			UserID:              resp.TournamentUser.UserId,
			Interval:            model.TournamentInterval(resp.TournamentUser.Interval.String()),
			Score:               resp.TournamentUser.Score,
			Ranking:             resp.TournamentUser.Ranking,
			Data:                resp.TournamentUser.Data,
			TournamentStartedAt: resp.TournamentUser.TournamentStartedAt,
			CreatedAt:           resp.TournamentUser.CreatedAt,
			UpdatedAt:           resp.TournamentUser.UpdatedAt,
		}
	}
	return &model.GetTournamentUserResponse{
		Success:        resp.Success,
		Error:          model.GetTournamentUserError(resp.Error.String()),
		TournamentUser: tournamentUser,
	}, nil
}

// GetTournamentUsers is the resolver for the GetTournamentUsers field.
func (r *queryResolver) GetTournamentUsers(ctx context.Context, input model.GetTournamentUsersRequest) (*model.GetTournamentUsersResponse, error) {
	if input.Pagination == nil {
		input.Pagination = &model.Pagination{}
	}
	resp, err := r.tournamentClient.GetTournamentUsers(ctx, &api.GetTournamentUsersRequest{
		Tournament: input.Tournament,
		Interval:   api.TournamentInterval(api.TournamentInterval_value[input.Interval.String()]),
		UserId:     input.UserID,
		Pagination: &api.Pagination{
			Max:  input.Pagination.Max,
			Page: input.Pagination.Page,
		},
	})
	if err != nil {
		return nil, err
	}
	tournamentUsers := make([]*model.TournamentUser, len(resp.TournamentUsers))
	for i, u := range resp.TournamentUsers {
		tournamentUsers[i] = &model.TournamentUser{
			ID:                  u.Id,
			Tournament:          u.Tournament,
			UserID:              u.UserId,
			Interval:            model.TournamentInterval(u.Interval.String()),
			Score:               u.Score,
			Ranking:             u.Ranking,
			Data:                u.Data,
			TournamentStartedAt: u.TournamentStartedAt,
			CreatedAt:           u.CreatedAt,
			UpdatedAt:           u.UpdatedAt,
		}
	}
	return &model.GetTournamentUsersResponse{
		Success:         resp.Success,
		Error:           model.GetTournamentUsersError(resp.Error.String()),
		TournamentUsers: tournamentUsers,
	}, nil
}
