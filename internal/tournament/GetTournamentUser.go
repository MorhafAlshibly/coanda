package tournament

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetTournamentUserCommand struct {
	service *Service
	In      *api.TournamentUserRequest
	Out     *api.GetTournamentUserResponse
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
		c.Out = &api.GetTournamentUserResponse{
			Success: false,
			Error:   conversion.Enum(*tErr, api.GetTournamentUserResponse_Error_value, api.GetTournamentUserResponse_NOT_FOUND),
		}
		return nil
	}
	// Get the tournament user
	result, err := c.service.Database.GetTournament(ctx, model.GetTournamentParams{
		Name:                c.In.Tournament,
		TournamentInterval:  model.TournamentTournamentInterval(c.In.Interval.String()),
		UserID:              c.In.UserId,
		TournamentStartedAt: c.service.GetTournamentStartDate(time.Now().UTC(), c.In.Interval),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Out = &api.GetTournamentUserResponse{
				Success: false,
				Error:   api.GetTournamentUserResponse_NOT_FOUND,
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
	c.Out = &api.GetTournamentUserResponse{
		Success:        true,
		TournamentUser: tournamentUser,
		Error:          api.GetTournamentUserResponse_NONE,
	}
	return nil
}
