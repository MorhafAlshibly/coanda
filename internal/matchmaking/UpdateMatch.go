package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type UpdateMatchCommand struct {
	service *Service
	In      *api.UpdateMatchRequest
	Out     *api.UpdateMatchResponse
}

func NewUpdateMatchCommand(service *Service, in *api.UpdateMatchRequest) *UpdateMatchCommand {
	return &UpdateMatchCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateMatchCommand) Execute(ctx context.Context) error {
	mmErr := c.service.checkForMatchRequestError(c.In.Match)
	// Check if error is found
	if mmErr != nil {
		c.Out = &api.UpdateMatchResponse{
			Success: false,
			Error:   conversion.Enum(*mmErr, api.UpdateMatchResponse_Error_value, api.UpdateMatchResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED),
		}
		return nil
	}
	// Check if data is given
	if c.In.Data == nil {
		c.Out = &api.UpdateMatchResponse{
			Success: false,
			Error:   api.UpdateMatchResponse_DATA_REQUIRED,
		}
		return nil
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := model.New(tx)
	// Prepare data
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	params := matchRequestToMatchParams(c.In.Match)
	result, err := qtx.UpdateMatch(ctx, model.UpdateMatchParams{
		Match: params,
		Data:  data,
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
		match, err := qtx.GetMatch(ctx, model.GetMatchParams{
			Match:       params,
			TicketLimit: 1,
			UserLimit:   1,
			ArenaLimit:  1,
		})
		if err != nil {
			return err
		}
		if len(match) == 0 {
			c.Out = &api.UpdateMatchResponse{
				Success: false,
				Error:   api.UpdateMatchResponse_NOT_FOUND,
			}
			return nil
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.UpdateMatchResponse{
		Success: true,
		Error:   api.UpdateMatchResponse_NONE,
	}
	return nil
}
