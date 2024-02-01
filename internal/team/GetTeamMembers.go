package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
)

type GetTeamMembersCommand struct {
	service *Service
	In      *api.TeamRequest
	Out     *api.GetTeamMembersResponse
}

func NewGetTeamMembersCommand(service *Service, in *api.TeamRequest) *GetTeamMembersCommand {
	return &GetTeamMembersCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTeamMembersCommand) Execute(ctx context.Context) error {
	field := c.service.GetTeamField(c.In)
	var members []model.TeamMember
	var err error
	// Check if name or owner is provided
	if field == NAME {
		members, err = c.service.database.GetTeamMembers(ctx, model.GetTeamMembersParams{
			Team:   *c.In.Name,
			Limit:  int32(c.service.maxMembers),
			Offset: 0,
		})
		// Check if owner is provided
	} else if field == OWNER {
		members, err = c.service.database.GetTeamMembersByOwner(ctx, model.GetTeamMembersByOwnerParams{
			Owner:  *c.In.Owner,
			Limit:  int32(c.service.maxMembers),
			Offset: 0,
		})
		// Check if member is provided
	} else if field == MEMBER {
		members, err = c.service.database.GetTeamMembersByMember(
			ctx,
			model.GetTeamMembersByMemberParams{
				UserID: *c.In.Member,
				Limit:  int32(c.service.maxMembers),
				Offset: 0,
			},
		)
	} else {
		c.Out = &api.GetTeamMembersResponse{
			Success: false,
			Error:   api.GetTeamMembersResponse_NOT_FOUND,
		}
		return nil
	}
	if err != nil {
		return err
	}
	outs := make([]*api.TeamMember, len(members))
	for i, member := range members {
		outs[i], err = UnmarshalTeamMember(member)
		if err != nil {
			return err
		}
	}
	c.Out = &api.GetTeamMembersResponse{
		Success:     true,
		TeamMembers: outs,
		Error:       api.GetTeamMembersResponse_NONE,
	}
	return nil

}
