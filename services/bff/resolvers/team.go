package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"

	"github.com/MorhafAlshibly/coanda/services/bff/model"
)

func (r *mutationResolver) CreateTeam(ctx context.Context, input model.CreateTeam) (*model.QueuedTeam, error) {
	panic(fmt.Errorf("not implemented: CreateTeam - createTeam"))
}

func (r *mutationResolver) UpdateTeamData(ctx context.Context, input model.UpdateTeamData) (*model.Team, error) {
	panic(fmt.Errorf("not implemented: UpdateTeamData - updateTeamData"))
}

func (r *mutationResolver) UpdateTeamScore(ctx context.Context, input model.UpdateTeamScore) (*model.Team, error) {
	panic(fmt.Errorf("not implemented: UpdateTeamScore - updateTeamScore"))
}

func (r *mutationResolver) DeleteTeam(ctx context.Context, input model.DeleteTeam) (bool, error) {
	panic(fmt.Errorf("not implemented: DeleteTeam - deleteTeam"))
}

func (r *mutationResolver) JoinTeam(ctx context.Context, input model.JoinTeam) (bool, error) {
	panic(fmt.Errorf("not implemented: JoinTeam - joinTeam"))
}

func (r *mutationResolver) LeaveTeam(ctx context.Context, input model.LeaveTeam) (bool, error) {
	panic(fmt.Errorf("not implemented: LeaveTeam - leaveTeam"))
}

func (r *queryResolver) GetTeam(ctx context.Context, input model.GetTeam) (*model.Team, error) {
	panic(fmt.Errorf("not implemented: GetTeam - getTeam"))
}

func (r *queryResolver) GetTeams(ctx context.Context, input *model.GetTeams) ([]*model.Team, error) {
	panic(fmt.Errorf("not implemented: GetTeams - getTeams"))
}

func (r *queryResolver) SearchTeams(ctx context.Context, input model.SearchTeams) ([]*model.Team, error) {
	panic(fmt.Errorf("not implemented: SearchTeams - searchTeams"))
}