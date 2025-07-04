package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.76

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/bff"
	"github.com/MorhafAlshibly/coanda/internal/bff/model"
)

// Error is the resolver for the error field.
func (r *createTeamResponseResolver) Error(ctx context.Context, obj *api.CreateTeamResponse) (model.CreateTeamError, error) {
	return model.CreateTeamError(obj.Error.String()), nil
}

// Error is the resolver for the error field.
func (r *getTeamMemberResponseResolver) Error(ctx context.Context, obj *api.GetTeamMemberResponse) (model.GetTeamMemberError, error) {
	return model.GetTeamMemberError(obj.Error.String()), nil
}

// Error is the resolver for the error field.
func (r *getTeamResponseResolver) Error(ctx context.Context, obj *api.GetTeamResponse) (model.GetTeamError, error) {
	return model.GetTeamError(obj.Error.String()), nil
}

// Error is the resolver for the error field.
func (r *joinTeamResponseResolver) Error(ctx context.Context, obj *api.JoinTeamResponse) (model.JoinTeamError, error) {
	return model.JoinTeamError(obj.Error.String()), nil
}

// Error is the resolver for the error field.
func (r *leaveTeamResponseResolver) Error(ctx context.Context, obj *api.LeaveTeamResponse) (model.LeaveTeamError, error) {
	return model.LeaveTeamError(obj.Error.String()), nil
}

// CreateTeam is the resolver for the CreateTeam field.
func (r *mutationResolver) CreateTeam(ctx context.Context, input *api.CreateTeamRequest) (*api.CreateTeamResponse, error) {
	return r.teamClient.CreateTeam(ctx, input)
}

// UpdateTeam is the resolver for the UpdateTeam field.
func (r *mutationResolver) UpdateTeam(ctx context.Context, input *api.UpdateTeamRequest) (*api.UpdateTeamResponse, error) {
	return r.teamClient.UpdateTeam(ctx, input)
}

// DeleteTeam is the resolver for the DeleteTeam field.
func (r *mutationResolver) DeleteTeam(ctx context.Context, input *api.TeamRequest) (*api.TeamResponse, error) {
	return r.teamClient.DeleteTeam(ctx, input)
}

// JoinTeam is the resolver for the JoinTeam field.
func (r *mutationResolver) JoinTeam(ctx context.Context, input *api.JoinTeamRequest) (*api.JoinTeamResponse, error) {
	return r.teamClient.JoinTeam(ctx, input)
}

// LeaveTeam is the resolver for the LeaveTeam field.
func (r *mutationResolver) LeaveTeam(ctx context.Context, input *api.TeamMemberRequest) (*api.LeaveTeamResponse, error) {
	return r.teamClient.LeaveTeam(ctx, input)
}

// UpdateTeamMember is the resolver for the UpdateTeamMember field.
func (r *mutationResolver) UpdateTeamMember(ctx context.Context, input *api.UpdateTeamMemberRequest) (*api.UpdateTeamMemberResponse, error) {
	return r.teamClient.UpdateTeamMember(ctx, input)
}

// GetTeam is the resolver for the GetTeam field.
func (r *queryResolver) GetTeam(ctx context.Context, input *api.GetTeamRequest) (*api.GetTeamResponse, error) {
	return r.teamClient.GetTeam(ctx, input)
}

// GetTeams is the resolver for the GetTeams field.
func (r *queryResolver) GetTeams(ctx context.Context, input *api.GetTeamsRequest) (*api.GetTeamsResponse, error) {
	return r.teamClient.GetTeams(ctx, input)
}

// GetTeamMember is the resolver for the GetTeamMember field.
func (r *queryResolver) GetTeamMember(ctx context.Context, input *api.TeamMemberRequest) (*api.GetTeamMemberResponse, error) {
	return r.teamClient.GetTeamMember(ctx, input)
}

// SearchTeams is the resolver for the SearchTeams field.
func (r *queryResolver) SearchTeams(ctx context.Context, input *api.SearchTeamsRequest) (*api.SearchTeamsResponse, error) {
	return r.teamClient.SearchTeams(ctx, input)
}

// Error is the resolver for the error field.
func (r *searchTeamsResponseResolver) Error(ctx context.Context, obj *api.SearchTeamsResponse) (model.SearchTeamsError, error) {
	return model.SearchTeamsError(obj.Error.String()), nil
}

// Error is the resolver for the error field.
func (r *teamResponseResolver) Error(ctx context.Context, obj *api.TeamResponse) (model.TeamError, error) {
	return model.TeamError(obj.Error.String()), nil
}

// Error is the resolver for the error field.
func (r *updateTeamMemberResponseResolver) Error(ctx context.Context, obj *api.UpdateTeamMemberResponse) (model.UpdateTeamMemberError, error) {
	return model.UpdateTeamMemberError(obj.Error.String()), nil
}

// Error is the resolver for the error field.
func (r *updateTeamResponseResolver) Error(ctx context.Context, obj *api.UpdateTeamResponse) (model.UpdateTeamError, error) {
	return model.UpdateTeamError(obj.Error.String()), nil
}

// CreateTeamResponse returns bff.CreateTeamResponseResolver implementation.
func (r *Resolver) CreateTeamResponse() bff.CreateTeamResponseResolver {
	return &createTeamResponseResolver{r}
}

// GetTeamMemberResponse returns bff.GetTeamMemberResponseResolver implementation.
func (r *Resolver) GetTeamMemberResponse() bff.GetTeamMemberResponseResolver {
	return &getTeamMemberResponseResolver{r}
}

// GetTeamResponse returns bff.GetTeamResponseResolver implementation.
func (r *Resolver) GetTeamResponse() bff.GetTeamResponseResolver { return &getTeamResponseResolver{r} }

// JoinTeamResponse returns bff.JoinTeamResponseResolver implementation.
func (r *Resolver) JoinTeamResponse() bff.JoinTeamResponseResolver {
	return &joinTeamResponseResolver{r}
}

// LeaveTeamResponse returns bff.LeaveTeamResponseResolver implementation.
func (r *Resolver) LeaveTeamResponse() bff.LeaveTeamResponseResolver {
	return &leaveTeamResponseResolver{r}
}

// SearchTeamsResponse returns bff.SearchTeamsResponseResolver implementation.
func (r *Resolver) SearchTeamsResponse() bff.SearchTeamsResponseResolver {
	return &searchTeamsResponseResolver{r}
}

// TeamResponse returns bff.TeamResponseResolver implementation.
func (r *Resolver) TeamResponse() bff.TeamResponseResolver { return &teamResponseResolver{r} }

// UpdateTeamMemberResponse returns bff.UpdateTeamMemberResponseResolver implementation.
func (r *Resolver) UpdateTeamMemberResponse() bff.UpdateTeamMemberResponseResolver {
	return &updateTeamMemberResponseResolver{r}
}

// UpdateTeamResponse returns bff.UpdateTeamResponseResolver implementation.
func (r *Resolver) UpdateTeamResponse() bff.UpdateTeamResponseResolver {
	return &updateTeamResponseResolver{r}
}

type createTeamResponseResolver struct{ *Resolver }
type getTeamMemberResponseResolver struct{ *Resolver }
type getTeamResponseResolver struct{ *Resolver }
type joinTeamResponseResolver struct{ *Resolver }
type leaveTeamResponseResolver struct{ *Resolver }
type searchTeamsResponseResolver struct{ *Resolver }
type teamResponseResolver struct{ *Resolver }
type updateTeamMemberResponseResolver struct{ *Resolver }
type updateTeamResponseResolver struct{ *Resolver }
