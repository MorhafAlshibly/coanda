package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.64

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/bff/model"
)

// CreateTeam is the resolver for the CreateTeam field.
func (r *mutationResolver) CreateTeam(ctx context.Context, input model.CreateTeamRequest) (*model.CreateTeamResponse, error) {
	resp, err := r.teamClient.CreateTeam(ctx, &api.CreateTeamRequest{
		Name:              input.Name,
		Score:             input.Score,
		FirstMemberUserId: input.FirstMemberUserID,
		Data:              input.Data,
		FirstMemberData:   input.FirstMemberData,
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

// UpdateTeam is the resolver for the UpdateTeam field.
func (r *mutationResolver) UpdateTeam(ctx context.Context, input model.UpdateTeamRequest) (*model.UpdateTeamResponse, error) {
	if input.Team == nil {
		input.Team = &model.TeamRequest{}
	}
	if input.Team.Member == nil {
		input.Team.Member = &model.TeamMemberRequest{}
	}
	resp, err := r.teamClient.UpdateTeam(ctx, &api.UpdateTeamRequest{
		Team: &api.TeamRequest{
			Id:   input.Team.ID,
			Name: input.Team.Name,
			Member: &api.TeamMemberRequest{
				Id:     input.Team.Member.ID,
				UserId: input.Team.Member.UserID,
			},
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
func (r *mutationResolver) DeleteTeam(ctx context.Context, input model.TeamRequest) (*model.TeamResponse, error) {
	if input.Member == nil {
		input.Member = &model.TeamMemberRequest{}
	}
	resp, err := r.teamClient.DeleteTeam(ctx, &api.TeamRequest{
		Id:   input.ID,
		Name: input.Name,
		Member: &api.TeamMemberRequest{
			Id:     input.Member.ID,
			UserId: input.Member.UserID,
		},
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
	if input.Team == nil {
		input.Team = &model.TeamRequest{}
	}
	if input.Team.Member == nil {
		input.Team.Member = &model.TeamMemberRequest{}
	}
	resp, err := r.teamClient.JoinTeam(ctx, &api.JoinTeamRequest{
		Team: &api.TeamRequest{
			Id:   input.Team.ID,
			Name: input.Team.Name,
			Member: &api.TeamMemberRequest{
				Id:     input.Team.Member.ID,
				UserId: input.Team.Member.UserID,
			},
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
func (r *mutationResolver) LeaveTeam(ctx context.Context, input model.TeamMemberRequest) (*model.LeaveTeamResponse, error) {
	resp, err := r.teamClient.LeaveTeam(ctx, &api.TeamMemberRequest{
		Id:     input.ID,
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
func (r *mutationResolver) UpdateTeamMember(ctx context.Context, input model.UpdateTeamMemberRequest) (*model.UpdateTeamMemberResponse, error) {
	if input.Member == nil {
		input.Member = &model.TeamMemberRequest{}
	}
	resp, err := r.teamClient.UpdateTeamMember(ctx, &api.UpdateTeamMemberRequest{
		Member: &api.TeamMemberRequest{
			Id:     input.Member.ID,
			UserId: input.Member.UserID,
		},
		Data: input.Data,
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
func (r *queryResolver) GetTeam(ctx context.Context, input model.GetTeamRequest) (*model.GetTeamResponse, error) {
	if input.Team == nil {
		input.Team = &model.TeamRequest{}
	}
	if input.Team.Member == nil {
		input.Team.Member = &model.TeamMemberRequest{}
	}
	if input.Pagination == nil {
		input.Pagination = &model.Pagination{}
	}
	resp, err := r.teamClient.GetTeam(ctx, &api.GetTeamRequest{
		Team: &api.TeamRequest{
			Id:     input.Team.ID,
			Name:   input.Team.Name,
			Member: &api.TeamMemberRequest{Id: input.Team.Member.ID, UserId: input.Team.Member.UserID},
		},
		Pagination: &api.Pagination{
			Max:  input.Pagination.Max,
			Page: input.Pagination.Page,
		},
	})
	if err != nil {
		return nil, err
	}
	var team *model.Team
	if resp.Team != nil {
		members := make([]*model.TeamMember, len(resp.Team.Members))
		for i, member := range resp.Team.Members {
			members[i] = &model.TeamMember{
				ID:        member.Id,
				UserID:    member.UserId,
				TeamID:    member.TeamId,
				Data:      member.Data,
				JoinedAt:  member.JoinedAt,
				UpdatedAt: member.UpdatedAt,
			}
		}
		team = &model.Team{
			ID:        resp.Team.Id,
			Name:      resp.Team.Name,
			Score:     resp.Team.Score,
			Ranking:   resp.Team.Ranking,
			Members:   members,
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
func (r *queryResolver) GetTeams(ctx context.Context, input model.GetTeamsRequest) (*model.GetTeamsResponse, error) {
	if input.Pagination == nil {
		input.Pagination = &model.Pagination{}
	}
	if input.MemberPagination == nil {
		input.MemberPagination = &model.Pagination{}
	}
	resp, err := r.teamClient.GetTeams(ctx, &api.GetTeamsRequest{
		Pagination: &api.Pagination{
			Max:  input.Pagination.Max,
			Page: input.Pagination.Page,
		},
		MemberPagination: &api.Pagination{
			Max:  input.MemberPagination.Max,
			Page: input.MemberPagination.Page,
		},
	})
	if err != nil {
		return nil, err
	}
	var teams []*model.Team
	for _, team := range resp.Teams {
		members := make([]*model.TeamMember, len(team.Members))
		for i, member := range team.Members {
			members[i] = &model.TeamMember{
				ID:        member.Id,
				UserID:    member.UserId,
				TeamID:    member.TeamId,
				Data:      member.Data,
				JoinedAt:  member.JoinedAt,
				UpdatedAt: member.UpdatedAt,
			}
		}
		teams = append(teams, &model.Team{
			ID:        team.Id,
			Name:      team.Name,
			Score:     team.Score,
			Ranking:   team.Ranking,
			Members:   members,
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
func (r *queryResolver) GetTeamMember(ctx context.Context, input model.TeamMemberRequest) (*model.GetTeamMemberResponse, error) {
	resp, err := r.teamClient.GetTeamMember(ctx, &api.TeamMemberRequest{
		UserId: input.UserID,
	})
	if err != nil {
		return nil, err
	}
	return &model.GetTeamMemberResponse{
		Success: resp.Success,
		Member: &model.TeamMember{
			ID:        resp.Member.Id,
			UserID:    resp.Member.UserId,
			TeamID:    resp.Member.TeamId,
			Data:      resp.Member.Data,
			JoinedAt:  resp.Member.JoinedAt,
			UpdatedAt: resp.Member.UpdatedAt,
		},
		Error: model.GetTeamMemberError(resp.Error.String()),
	}, nil
}

// SearchTeams is the resolver for the SearchTeams field.
func (r *queryResolver) SearchTeams(ctx context.Context, input model.SearchTeamsRequest) (*model.SearchTeamsResponse, error) {
	if input.Pagination == nil {
		input.Pagination = &model.Pagination{}
	}
	resp, err := r.teamClient.SearchTeams(ctx, &api.SearchTeamsRequest{
		Pagination: &api.Pagination{
			Max:  input.Pagination.Max,
			Page: input.Pagination.Page,
		},
		Query: input.Query,
	})
	if err != nil {
		return nil, err
	}
	teams := make([]*model.Team, len(resp.Teams))
	for i, team := range resp.Teams {
		members := make([]*model.TeamMember, len(team.Members))
		for j, member := range team.Members {
			members[j] = &model.TeamMember{
				ID:        member.Id,
				UserID:    member.UserId,
				TeamID:    member.TeamId,
				Data:      member.Data,
				JoinedAt:  member.JoinedAt,
				UpdatedAt: member.UpdatedAt,
			}
		}
		teams[i] = &model.Team{
			ID:        team.Id,
			Name:      team.Name,
			Score:     team.Score,
			Ranking:   team.Ranking,
			Members:   members,
			Data:      team.Data,
			CreatedAt: team.CreatedAt,
			UpdatedAt: team.UpdatedAt,
		}
	}
	return &model.SearchTeamsResponse{
		Success: resp.Success,
		Teams:   teams,
	}, nil
}
