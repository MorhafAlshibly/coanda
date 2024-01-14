package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type UpdateTeamCommand struct {
	service *Service
	In      *api.UpdateTeamRequest
	Out     *api.TeamResponse
}

func NewUpdateTeamCommand(service *Service, in *api.UpdateTeamRequest) *UpdateTeamCommand {
	return &UpdateTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateTeamCommand) Execute(ctx context.Context) error {
	if c.In.Team == nil {
		c.Out = &api.TeamResponse{
			Success: false,
			Error:   api.TeamResponse_NO_FIELD_SPECIFIED,
		}
		return nil
	}
	if c.In.Team.Name != nil {
		if c.In.Data != nil {

		}
	}
}
