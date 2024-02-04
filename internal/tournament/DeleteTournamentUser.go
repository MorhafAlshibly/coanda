package tournament

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type DeleteTournamentUserCommand struct {
	service *Service
	In      *api.TournamentUserRequest
	Out     *api.TournamentUserResponse
}

func NewDeleteTournamentUserCommand(service *Service, in *api.TournamentUserRequest) *DeleteTournamentUserCommand {
	return &DeleteTournamentUserCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteTournamentUserCommand) Execute(ctx context.Context) error {
	// Check for errors
	tErr := c.service.checkForTournamentUserRequestError(c.In)
	if tErr != nil {
		c.Out = &api.TournamentUserResponse{
			Success: false,
			Error:   conversion.Enum(*tErr, api.TournamentUserResponse_Error_value, api.TournamentUserResponse_NOT_FOUND),
		}
		return nil
	}
	// Delete the tournament user
	result, err := c.service.database.DeleteTournament(ctx, model.DeleteTournamentParams{
		Name:               c.In.Tournament,
		TournamentInterval: model.TournamentTournamentInterval(c.In.Interval.String()),
		UserID:             c.In.UserId,
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// Check if the tournament user was deleted
	if rowsAffected == 0 {
		c.Out = &api.TournamentUserResponse{
			Success: false,
			Error:   api.TournamentUserResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.TournamentUserResponse{
		Success: true,
		Error:   api.TournamentUserResponse_NONE,
	}
	return nil
}
