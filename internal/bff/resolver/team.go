package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.43

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/bff/model"
)

// CreateTeam is the resolver for the CreateTeam field.
func (r *mutationResolver) CreateTeam(ctx context.Context, input *model.CreateTeamRequest) (*model.CreateTeamResponse, error) {
	resp, err := r.teamClient.CreateTeam(ctx, &api.CreateTeamRequest{
		Name:      input.Name,
		Owner:     input.Owner,
		Score:     input.Score,
		Data:      input.Data,
		OwnerData: input.OwnerData,
	})
	if err != nil {
		return nil, err
	}
	return &model.CreateTeamResponse{
		Success: resp.Success,
		Error:   model.CreateTeamError(resp.Error.String()),
	}, nil
}

// UpdateTeam is the resolver for the UpdateTeam field.
func (r *mutationResolver) UpdateTeam(ctx context.Context, input *model.UpdateTeamRequest) (*model.UpdateTeamResponse, error) {
	resp, err := r.teamClient.UpdateTeam(ctx, &api.UpdateTeamRequest{
		Team: &api.TeamRequest{
			Name:   input.Team.Name,
			Owner:  input.Team.Owner,
			Member: input.Team.Member,
		},
		Score:          input.Score,
		Data:           input.Data,
		IncrementScore: input.IncrementScore,
	})
	if err != nil {
		return nil, err
	}
	return &model.UpdateTeamResponse{
		Success: resp.Success,
		Error:   model.UpdateTeamError(resp.Error.String()),
	}, nil
}

// DeleteTeam is the resolver for the DeleteTeam field.
func (r *mutationResolver) DeleteTeam(ctx context.Context, input *model.TeamRequest) (*model.TeamResponse, error) {
	resp, err := r.teamClient.DeleteTeam(ctx, &api.TeamRequest{
		Name:   input.Name,
		Owner:  input.Owner,
		Member: input.Member,
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
func (r *mutationResolver) JoinTeam(ctx context.Context, input *model.JoinTeamRequest) (*model.JoinTeamResponse, error) {
	resp, err := r.teamClient.JoinTeam(ctx, &api.JoinTeamRequest{
		Team: &api.TeamRequest{
			Name:   input.Team.Name,
			Owner:  input.Team.Owner,
			Member: input.Team.Member,
		},
		UserId: input.UserID,
		Data:   input.Data,
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
func (r *mutationResolver) LeaveTeam(ctx context.Context, input *model.LeaveTeamRequest) (*model.LeaveTeamResponse, error) {
	resp, err := r.teamClient.LeaveTeam(ctx, &api.LeaveTeamRequest{
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

// UpdateTeamMember is the resolver for the UpdateTeamMember field.
func (r *mutationResolver) UpdateTeamMember(ctx context.Context, input *model.UpdateTeamMemberRequest) (*model.UpdateTeamMemberResponse, error) {
	resp, err := r.teamClient.UpdateTeamMember(ctx, &api.UpdateTeamMemberRequest{
		UserId: input.UserID,
		Data:   input.Data,
	})
	if err != nil {
		return nil, err
	}
	return &model.UpdateTeamMemberResponse{
		Success: resp.Success,
		Error:   model.UpdateTeamMemberError(resp.Error.String()),
	}, nil
}

// GetTeam is the resolver for the GetTeam field.
func (r *queryResolver) GetTeam(ctx context.Context, input *model.TeamRequest) (*model.GetTeamResponse, error) {
	resp, err := r.teamClient.GetTeam(ctx, &api.TeamRequest{
		Name:   input.Name,
		Owner:  input.Owner,
		Member: input.Member,
	})
	if err != nil {
		return nil, err
	}
	var team *model.Team
	if resp.Team != nil {
		team = &model.Team{
			Name:      resp.Team.Name,
			Owner:     resp.Team.Owner,
			Score:     resp.Team.Score,
			Ranking:   resp.Team.Ranking,
			Data:      resp.Team.Data,
			CreatedAt: resp.Team.CreatedAt,
			UpdatedAt: resp.Team.UpdatedAt,
		}
	}
	return &model.GetTeamResponse{
		Success: resp.Success,
		Team:    team,
		Error:   model.GetTeamError(resp.Error.String()),
	}, nil
}

// GetTeams is the resolver for the GetTeams field.
func (r *queryResolver) GetTeams(ctx context.Context, input *model.Pagination) (*model.GetTeamsResponse, error) {
	resp, err := r.teamClient.GetTeams(ctx, &api.Pagination{
		Max:  input.Max,
		Page: input.Page,
	})
	if err != nil {
		return nil, err
	}
	var teams []*model.Team
	for _, team := range resp.Teams {
		teams = append(teams, &model.Team{
			Name:      team.Name,
			Owner:     team.Owner,
			Score:     team.Score,
			Ranking:   team.Ranking,
			Data:      team.Data,
			CreatedAt: team.CreatedAt,
			UpdatedAt: team.UpdatedAt,
		})
	}
	return &model.GetTeamsResponse{
		Success: resp.Success,
		Teams:   teams,
	}, nil
}

// GetTeamMember is the resolver for the GetTeamMember field.
func (r *queryResolver) GetTeamMember(ctx context.Context, input *model.GetTeamMemberRequest) (*model.GetTeamMemberResponse, error) {
	resp, err := r.teamClient.GetTeamMember(ctx, &api.GetTeamMemberRequest{
		UserId: input.UserID,
	})
	if err != nil {
		return nil, err
	}
	return &model.GetTeamMemberResponse{
		Success: resp.Success,
		TeamMember: &model.TeamMember{
			Team:      resp.TeamMember.Team,
			UserID:    resp.TeamMember.UserId,
			Data:      resp.TeamMember.Data,
			JoinedAt:  resp.TeamMember.JoinedAt,
			UpdatedAt: resp.TeamMember.UpdatedAt,
		},
		Error: model.GetTeamMemberError(resp.Error.String()),
	}, nil
}

// GetTeamMembers is the resolver for the GetTeamMembers field.
func (r *queryResolver) GetTeamMembers(ctx context.Context, input *model.GetTeamMembersRequest) (*model.GetTeamMembersResponse, error) {
	resp, err := r.teamClient.GetTeamMembers(ctx, &api.GetTeamMembersRequest{
		Team: &api.TeamRequest{
			Name:   input.Team.Name,
			Owner:  input.Team.Owner,
			Member: input.Team.Member,
		},
		Pagination: &api.Pagination{
			Max:  input.Pagination.Max,
			Page: input.Pagination.Page,
		},
	})
	if err != nil {
		return nil, err
	}
	var teamMembers []*model.TeamMember
	for _, teamMember := range resp.TeamMembers {
		teamMembers = append(teamMembers, &model.TeamMember{
			Team:      teamMember.Team,
			UserID:    teamMember.UserId,
			Data:      teamMember.Data,
			JoinedAt:  teamMember.JoinedAt,
			UpdatedAt: teamMember.UpdatedAt,
		})
	}
	return &model.GetTeamMembersResponse{
		Success:     resp.Success,
		TeamMembers: teamMembers,
		Error:       model.GetTeamMembersError(resp.Error.String()),
	}, nil
}

// SearchTeams is the resolver for the SearchTeams field.
func (r *queryResolver) SearchTeams(ctx context.Context, input *model.SearchTeamsRequest) (*model.SearchTeamsResponse, error) {
	var pagination *api.Pagination
	if input.Pagination != nil {
		pagination = &api.Pagination{
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
		teams = append(teams, &model.Team{
			Name:      team.Name,
			Owner:     team.Owner,
			Score:     team.Score,
			Ranking:   team.Ranking,
			Data:      team.Data,
			CreatedAt: team.CreatedAt,
			UpdatedAt: team.UpdatedAt,
		})
	}
	return &model.SearchTeamsResponse{
		Success: resp.Success,
		Teams:   teams,
	}, nil
}
