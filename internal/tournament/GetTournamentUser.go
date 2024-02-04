package tournament

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetTournamentUserCommand struct {
	service *Service
	In      *api.TournamentUserRequest
	Out     *api.TournamentUserResponse
}

func NewGetTournamentUserCommand(service *Service, in *api.TournamentUserRequest) *GetTournamentUserCommand {
	return &GetTournamentUserCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTournamentUserCommand) Execute(ctx context.Context) error {
	// Check for errors
	tErr := c.service.checkForTournamentUserRequestError(c.In)
	if tErr != nil {
		c.Out = &api.TournamentUserResponse{
			Success: false,
			Error:   conversion.Enum(*tErr, api.TournamentUserResponse_Error_value, api.TournamentUserResponse_NOT_FOUND),
		}
		return nil
	}
	// Get the tournament user
	result, err := c.service.database.GetTournament(ctx, model.GetTournamentParams{
		Name:               c.In.Tournament,
		TournamentInterval: model.TournamentTournamentInterval(c.In.Interval.String()),
		UserID:             c.In.UserId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Out = &api.TournamentUserResponse{
				Success: false,
				Error:   api.TournamentUserResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	// Unmarshal the tournament user
	tournamentUser, err := unmarshalTournamentUser(&result)
	if err != nil {
		return err
	}
	c.Out = &api.TournamentUserResponse{
		Success:        true,
		TournamentUser: tournamentUser,
		Error:          api.TournamentUserResponse_NONE,
	}
	return nil
}
