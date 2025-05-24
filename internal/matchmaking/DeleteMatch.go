package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type DeleteMatchCommand struct {
	service *Service
	In      *api.MatchRequest
	Out     *api.DeleteMatchResponse
}

func NewDeleteMatchCommand(service *Service, in *api.MatchRequest) *DeleteMatchCommand {
	return &DeleteMatchCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteMatchCommand) Execute(ctx context.Context) error {
	mmErr := c.service.checkForMatchRequestError(c.In)
	// Check if error is found
	if mmErr != nil {
		c.Out = &api.DeleteMatchResponse{
			Success: false,
			Error:   conversion.Enum(*mmErr, api.DeleteMatchResponse_Error_value, api.DeleteMatchResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED),
		}
		return nil
	}
	params := matchRequestToMatchParams(c.In)
	tx, err := c.service.sql.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	result, err := qtx.DeleteMatch(ctx, params)
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
			c.Out = &api.DeleteMatchResponse{
				Success: false,
				Error:   api.DeleteMatchResponse_NOT_FOUND,
			}
			return nil
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.DeleteMatchResponse{
		Success: true,
		Error:   api.DeleteMatchResponse_NONE,
	}
	return nil
}
