package team

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetTeamMemberCommand struct {
	service *Service
	In      *api.TeamMemberRequest
	Out     *api.GetTeamMemberResponse
}

func NewGetTeamMemberCommand(service *Service, in *api.TeamMemberRequest) *GetTeamMemberCommand {
	return &GetTeamMemberCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTeamMemberCommand) Execute(ctx context.Context) error {
	tmErr := c.service.checkForTeamMemberRequestError(c.In)
	// Check if error is found
	if tmErr != nil {
		c.Out = &api.GetTeamMemberResponse{
			Success: false,
			Error:   conversion.Enum(*tmErr, api.GetTeamMemberResponse_Error_value, api.GetTeamMemberResponse_NO_FIELD_SPECIFIED),
		}
		return nil
	}
	member, err := c.service.database.GetTeamMember(ctx, model.GetTeamMemberParams{
		ID:     conversion.Uint64ToSqlNullInt64(c.In.Id),
		UserID: conversion.Uint64ToSqlNullInt64(c.In.UserId),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			c.Out = &api.GetTeamMemberResponse{
				Success: false,
				Error:   api.GetTeamMemberResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	out, err := unmarshalTeamMember(member)
	if err != nil {
		return err
	}
	c.Out = &api.GetTeamMemberResponse{
		Success: true,
		Member:  out,
		Error:   api.GetTeamMemberResponse_NONE,
	}
	return nil
}
