package tournament

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type DeleteTournamentUserCommand struct {
	service *Service
	In      *api.GetTournamentUserRequest
	Out     *api.DeleteTournamentUserResponse
}

func NewDeleteTournamentUserCommand(service *Service, in *api.GetTournamentUserRequest) *DeleteTournamentUserCommand {
	return &DeleteTournamentUserCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteTournamentUserCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In)
	if err != nil {
		c.Out = &api.DeleteTournamentUserResponse{
			Success: false,
			Error:   api.DeleteTournamentUserResponse_INVALID,
		}
		return nil
	}
	if c.In.TournamentIntervalUserId != nil {
		if len(c.In.TournamentIntervalUserId.Tournament) < int(c.service.minTournamentNameLength) {
			c.Out = &api.DeleteTournamentUserResponse{
				Success: false,
				Error:   api.DeleteTournamentUserResponse_TOURNAMENT_NAME_TOO_SHORT,
			}
			return nil
		}
		if len(c.In.TournamentIntervalUserId.Tournament) > int(c.service.maxTournamentNameLength) {
			c.Out = &api.DeleteTournamentUserResponse{
				Success: false,
				Error:   api.DeleteTournamentUserResponse_TOURNAMENT_NAME_TOO_LONG,
			}
			return nil
		}
	}
	result, writeErr := c.service.database.DeleteOne(ctx, filter)
	if writeErr != nil {
		return writeErr
	}
	if result.DeletedCount == 0 {
		c.Out = &api.DeleteTournamentUserResponse{
			Success: false,
			Error:   api.DeleteTournamentUserResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.DeleteTournamentUserResponse{
		Success: true,
		Error:   api.DeleteTournamentUserResponse_NONE,
	}
	return nil
}
