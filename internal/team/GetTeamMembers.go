package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetTeamMembersCommand struct {
	service *Service
	In      *api.GetTeamMembersRequest
	Out     *api.GetTeamMembersResponse
}

func NewGetTeamMembersCommand(service *Service, in *api.GetTeamMembersRequest) *GetTeamMembersCommand {
	return &GetTeamMembersCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTeamMembersCommand) Execute(ctx context.Context) error {
	tErr := c.service.checkForTeamRequestError(c.In.Team)
	// Check if error is found
	if tErr != nil {
		c.Out = &api.GetTeamMembersResponse{
			Success: false,
			Error:   conversion.Enum(*tErr, api.GetTeamMembersResponse_Error_value, api.GetTeamMembersResponse_NO_FIELD_SPECIFIED),
		}
		return nil
	}
	// Check if team member is initialised
	if c.In.Team.Member == nil {
		c.In.Team.Member = &api.TeamMemberRequest{}
	}
	limit, offset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	members, err := c.service.database.GetTeamMembers(ctx, model.GetTeamMembersParams{
		Team: model.GetTeamParams{
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Team.Id),
			Name: conversion.StringToSqlNullString(c.In.Team.Name),
			Member: model.GetTeamMemberParams{
				ID:     conversion.Uint64ToSqlNullInt64(c.In.Team.Member.Id),
				UserID: conversion.Uint64ToSqlNullInt64(c.In.Team.Member.UserId),
			},
		},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return err
	}
	outs := make([]*api.TeamMember, len(members))
	for i, member := range members {
		outs[i], err = unmarshalTeamMember(member)
		if err != nil {
			return err
		}
	}
	c.Out = &api.GetTeamMembersResponse{
		Success: true,
		Members: outs,
		Error:   api.GetTeamMembersResponse_NONE,
	}
	return nil

}
