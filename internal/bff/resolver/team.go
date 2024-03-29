package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/bff/model"
	"github.com/MorhafAlshibly/coanda/pkg"
)

// CreateTeam is the resolver for the CreateTeam field.
func (r *mutationResolver) CreateTeam(ctx context.Context, input model.CreateTeamRequest) (*model.CreateTeamResponse, error) {
	data, err := pkg.MapStringAnyToMapStringString(input.Data)
	if err != nil {
		return nil, err
	}
	resp, err := r.teamClient.CreateTeam(ctx, &api.CreateTeamRequest{
		Name:                input.Name,
		Owner:               input.Owner,
		MembersWithoutOwner: input.MembersWithoutOwner,
		Score:               input.Score,
		Data:                data,
	})
	if err != nil {
		return nil, err
	}
	return &model.CreateTeamResponse{
		Success: resp.Success,
		ID:      resp.Id,
		Error:   model.CreateTeamError(resp.Error.String()),
	}, nil
}

// UpdateTeamData is the resolver for the UpdateTeamData field.
func (r *mutationResolver) UpdateTeamData(ctx context.Context, input model.UpdateTeamDataRequest) (*model.TeamResponse, error) {
	data, err := pkg.MapStringAnyToMapStringString(input.Data)
	if err != nil {
		return nil, err
	}
	resp, err := r.teamClient.UpdateTeamData(ctx, &api.UpdateTeamDataRequest{
		Team: &api.GetTeamRequest{
			Id:    input.Team.ID,
			Name:  input.Team.Name,
			Owner: input.Team.Owner,
		},
		Data: data,
	})
	if err != nil {
		return nil, err
	}
	return &model.TeamResponse{
		Success: resp.Success,
		Error:   model.TeamError(resp.Error.String()),
	}, nil
}

// UpdateTeamScore is the resolver for the UpdateTeamScore field.
func (r *mutationResolver) UpdateTeamScore(ctx context.Context, input model.UpdateTeamScoreRequest) (*model.TeamResponse, error) {
	resp, err := r.teamClient.UpdateTeamScore(ctx, &api.UpdateTeamScoreRequest{
		Team: &api.GetTeamRequest{
			Id:    input.Team.ID,
			Name:  input.Team.Name,
			Owner: input.Team.Owner,
		},
		ScoreOffset: input.ScoreOffset,
	})
	if err != nil {
		return nil, err
	}
	return &model.TeamResponse{
		Success: resp.Success,
		Error:   model.TeamError(resp.Error.String()),
	}, nil
}

// DeleteTeam is the resolver for the DeleteTeam field.
func (r *mutationResolver) DeleteTeam(ctx context.Context, input model.GetTeamRequest) (*model.TeamResponse, error) {
	resp, err := r.teamClient.DeleteTeam(ctx, &api.GetTeamRequest{
		Id:    input.ID,
		Name:  input.Name,
		Owner: input.Owner,
	})
	if err != nil {
		return nil, err
	}
	return &model.TeamResponse{
		Success: resp.Success,
		Error:   model.TeamError(resp.Error.String()),
	}, nil
}

// JoinTeam is the resolver for the JoinTeam field.
func (r *mutationResolver) JoinTeam(ctx context.Context, input model.JoinTeamRequest) (*model.JoinTeamResponse, error) {
	resp, err := r.teamClient.JoinTeam(ctx, &api.JoinTeamRequest{
		Team: &api.GetTeamRequest{
			Id:    input.Team.ID,
			Name:  input.Team.Name,
			Owner: input.Team.Owner,
		},
		UserId: input.UserID,
	})
	if err != nil {
		return nil, err
	}
	return &model.JoinTeamResponse{
		Success: resp.Success,
		Error:   model.JoinTeamError(resp.Error.String()),
	}, nil
}

// LeaveTeam is the resolver for the LeaveTeam field.
func (r *mutationResolver) LeaveTeam(ctx context.Context, input model.LeaveTeamRequest) (*model.LeaveTeamResponse, error) {
	resp, err := r.teamClient.LeaveTeam(ctx, &api.LeaveTeamRequest{
		Team: &api.GetTeamRequest{
			Id:    input.Team.ID,
			Name:  input.Team.Name,
			Owner: input.Team.Owner,
		},
		UserId: input.UserID,
	})
	if err != nil {
		return nil, err
	}
	return &model.LeaveTeamResponse{
		Success: resp.Success,
		Error:   model.LeaveTeamError(resp.Error.String()),
	}, nil
}

// GetTeam is the resolver for the GetTeam field.
func (r *queryResolver) GetTeam(ctx context.Context, input model.GetTeamRequest) (*model.GetTeamResponse, error) {
	resp, err := r.teamClient.GetTeam(ctx, &api.GetTeamRequest{
		Id:    input.ID,
		Name:  input.Name,
		Owner: input.Owner,
	})
	if err != nil {
		return nil, err
	}
	var team *model.Team
	if resp.Team != nil {
		dataMap, err := pkg.MapStringStringToMapStringAny(resp.Team.Data)
		if err != nil {
			return nil, err
		}
		team = &model.Team{
			ID:                  resp.Team.Id,
			Name:                resp.Team.Name,
			Owner:               resp.Team.Owner,
			MembersWithoutOwner: resp.Team.MembersWithoutOwner,
			Score:               resp.Team.Score,
			Rank:                resp.Team.Rank,
			Data:                dataMap,
		}
	}
	return &model.GetTeamResponse{
		Success: resp.Success,
		Team:    team,
		Error:   model.GetTeamError(resp.Error.String()),
	}, nil
}

// GetTeams is the resolver for the GetTeams field.
func (r *queryResolver) GetTeams(ctx context.Context, input model.GetTeamsRequest) (*model.GetTeamsResponse, error) {
	resp, err := r.teamClient.GetTeams(ctx, &api.GetTeamsRequest{
		Max:  input.Max,
		Page: input.Page,
	})
	if err != nil {
		return nil, err
	}
	var teams []*model.Team
	for _, team := range resp.Teams {
		dataMap, err := pkg.MapStringStringToMapStringAny(team.Data)
		if err != nil {
			return nil, err
		}
		teams = append(teams, &model.Team{
			ID:                  team.Id,
			Name:                team.Name,
			Owner:               team.Owner,
			MembersWithoutOwner: team.MembersWithoutOwner,
			Score:               team.Score,
			Rank:                team.Rank,
			Data:                dataMap,
		})
	}
	return &model.GetTeamsResponse{
		Success: resp.Success,
		Teams:   teams,
	}, nil
}

// SearchTeams is the resolver for the SearchTeams field.
func (r *queryResolver) SearchTeams(ctx context.Context, input model.SearchTeamsRequest) (*model.SearchTeamsResponse, error) {
	var pagination *api.GetTeamsRequest
	if input.Pagination != nil {
		pagination = &api.GetTeamsRequest{
			Max:  input.Pagination.Max,
			Page: input.Pagination.Page,
		}
	}
	resp, err := r.teamClient.SearchTeams(ctx, &api.SearchTeamsRequest{
		Pagination: pagination,
		Query:      input.Query,
	})
	if err != nil {
		return nil, err
	}
	var teams []*model.Team
	for _, team := range resp.Teams {
		dataMap, err := pkg.MapStringStringToMapStringAny(team.Data)
		if err != nil {
			return nil, err
		}
		teams = append(teams, &model.Team{
			ID:                  team.Id,
			Name:                team.Name,
			Owner:               team.Owner,
			MembersWithoutOwner: team.MembersWithoutOwner,
			Score:               team.Score,
			Rank:                team.Rank,
			Data:                dataMap,
		})
	}
	return &model.SearchTeamsResponse{
		Success: resp.Success,
		Teams:   teams,
	}, nil
}
