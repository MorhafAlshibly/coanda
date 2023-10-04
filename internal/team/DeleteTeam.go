package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type DeleteTeamCommand struct {
	service *Service
	In      *api.GetTeamRequest
	Out     *api.TeamResponse
}

func NewDeleteTeamCommand(service *Service, in *api.GetTeamRequest) *DeleteTeamCommand {
	return &DeleteTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteTeamCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In)
	if err != nil {
		c.Out = &api.TeamResponse{
			Success: false,
			Error:   api.TeamResponse_INVALID,
		}
		return nil
	}
	result, writeErr := c.service.db.DeleteOne(ctx, filter)
	if writeErr != nil {
		return writeErr
	}
	if result.DeletedCount == 0 {
		c.Out = &api.TeamResponse{
			Success: false,
			Error:   api.TeamResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.TeamResponse{
		Success: true,
		Error:   api.TeamResponse_NONE,
	}
	return nil
}
