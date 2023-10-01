package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type DeleteTeamCommand struct {
	service *Service
	In      *api.DeleteTeamRequest
	Out     *api.DeleteTeamResponse
}

func NewDeleteTeamCommand(service *Service, in *api.DeleteTeamRequest) *DeleteTeamCommand {
	return &DeleteTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteTeamCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In.Team)
	if err != nil {
		c.Out = &api.DeleteTeamResponse{
			Success: false,
			Error:   api.DeleteTeamResponse_INVALID,
		}
		return nil
	}
	result, writeErr := c.service.db.DeleteOne(ctx, filter)
	if writeErr != nil {
		return writeErr
	}
	if result.DeletedCount == 0 {
		c.Out = &api.DeleteTeamResponse{
			Success: false,
			Error:   api.DeleteTeamResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.DeleteTeamResponse{
		Success: true,
		Error:   api.DeleteTeamResponse_NONE,
	}
	return nil
}
