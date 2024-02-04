package tournament

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type UpdateTournamentUserCommand struct {
	service *Service
	In      *api.UpdateTournamentUserRequest
	Out     *api.UpdateTournamentUserResponse
}

func NewUpdateTournamentUserCommand(service *Service, in *api.UpdateTournamentUserRequest) *UpdateTournamentUserCommand {
	return &UpdateTournamentUserCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateTournamentUserCommand) Execute(ctx context.Context) error {
	tErr := c.service.checkForTournamentUserRequestError(c.In.Tournament)
	if tErr != nil {
		c.Out = &api.UpdateTournamentUserResponse{
			Success: false,
			Error:   conversion.Enum(*tErr, api.UpdateTournamentUserResponse_Error_value, api.UpdateTournamentUserResponse_NOT_FOUND),
		}
		return nil
	}
	if c.In.Score == nil && c.In.Data == nil {
		c.Out = &api.UpdateTournamentUserResponse{
			Success: false,
			Error:   api.UpdateTournamentUserResponse_NO_UPDATE_SPECIFIED,
		}
		return nil
	}
	// Update the tournament user in the store
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if c.In.Score != nil {
		result, err := c.service.database.UpdateTournamentScore(ctx, model.UpdateTournamentScoreParams{
			Name:               c.In.Tournament.Tournament,
			TournamentInterval: model.TournamentTournamentInterval(c.In.Tournament.Interval.String()),
			UserID:             c.In.Tournament.UserId,
			Score:              *c.In.Score,
		})
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			c.Out = &api.UpdateTournamentUserResponse{
				Success: false,
				Error:   api.UpdateTournamentUserResponse_NOT_FOUND,
			}
			return nil
		}
	}
	if c.In.Data != nil {
		data, err := conversion.ProtobufStructToRawJson(c.In.Data)
		if err != nil {
			return err
		}
		result, err := c.service.database.UpdateTournamentData(ctx, model.UpdateTournamentDataParams{
			Name:               c.In.Tournament.Tournament,
			TournamentInterval: model.TournamentTournamentInterval(c.In.Tournament.Interval.String()),
			UserID:             c.In.Tournament.UserId,
			Data:               data,
		})
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			c.Out = &api.UpdateTournamentUserResponse{
				Success: false,
				Error:   api.UpdateTournamentUserResponse_NOT_FOUND,
			}
			return nil
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.UpdateTournamentUserResponse{
		Success: true,
		Error:   api.UpdateTournamentUserResponse_NONE,
	}
	return nil
}
