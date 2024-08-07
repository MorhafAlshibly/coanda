package tournament

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

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
	if c.In.Score != nil && c.In.IncrementScore == nil {
		c.Out = &api.UpdateTournamentUserResponse{
			Success: false,
			Error:   api.UpdateTournamentUserResponse_INCREMENT_SCORE_NOT_SPECIFIED,
		}
		return nil
	}
	data := json.RawMessage(nil)
	if c.In.Data != nil {
		var err error
		data, err = conversion.ProtobufStructToRawJson(c.In.Data)
		if err != nil {
			return err
		}
	}
	// Update the tournament user in the store
	result, err := c.service.database.UpdateTournament(ctx, model.UpdateTournamentParams{
		ID:                          conversion.Uint64ToSqlNullInt64(c.In.Tournament.Id),
		NameIntervalUserIDStartedAt: *c.service.convertTournamentIntervalUserIdToNullNameIntervalUserIDStartedAt(c.In.Tournament.TournamentIntervalUserId),
		Score:                       conversion.Int64ToSqlNullInt64(c.In.Score),
		IncrementScore:              conversion.PointerBoolToValue(c.In.IncrementScore),
		Data:                        data,
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Check if we didn't find a row
		_, err = c.service.database.GetTournament(ctx, model.GetTournamentParams{
			ID:                          conversion.Uint64ToSqlNullInt64(c.In.Tournament.Id),
			NameIntervalUserIDStartedAt: *c.service.convertTournamentIntervalUserIdToNullNameIntervalUserIDStartedAt(c.In.Tournament.TournamentIntervalUserId),
		})
		// Check if tournament user is found, if not return not found
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.Out = &api.UpdateTournamentUserResponse{
					Success: false,
					Error:   api.UpdateTournamentUserResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
	}
	c.Out = &api.UpdateTournamentUserResponse{
		Success: true,
		Error:   api.UpdateTournamentUserResponse_NONE,
	}
	return nil
}
