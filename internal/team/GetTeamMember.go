package team

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
)

type GetTeamMemberCommand struct {
	service *Service
	In      *api.GetTeamMemberRequest
	Out     *api.GetTeamMemberResponse
}

func NewGetTeamMemberCommand(service *Service, in *api.GetTeamMemberRequest) *GetTeamMemberCommand {
	return &GetTeamMemberCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTeamMemberCommand) Execute(ctx context.Context) error {
	if c.In.UserId == 0 {
		c.Out = &api.GetTeamMemberResponse{
			Success: false,
			Error:   api.GetTeamMemberResponse_USER_ID_REQUIRED,
		}
		return nil
	}
	member, err := c.service.database.GetTeamMember(ctx, c.In.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Out = &api.GetTeamMemberResponse{
				Success: false,
				Error:   api.GetTeamMemberResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	out, err := UnmarshalTeamMember(member)
	if err != nil {
		return err
	}
	c.Out = &api.GetTeamMemberResponse{
		Success:    true,
		TeamMember: out,
		Error:      api.GetTeamMemberResponse_NONE,
	}
	return nil
}
