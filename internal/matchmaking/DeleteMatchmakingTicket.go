package matchmaking

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type DeleteMatchmakingTicketCommand struct {
	service *Service
	In      *api.MatchmakingTicketRequest
	Out     *api.DeleteMatchmakingTicketResponse
}

func NewDeleteMatchmakingTicketCommand(service *Service, in *api.MatchmakingTicketRequest) *DeleteMatchmakingTicketCommand {
	return &DeleteMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteMatchmakingTicketCommand) Execute(ctx context.Context) error {
	mtErr := c.service.checkForMatchmakingTicketRequestError(c.In)
	if mtErr != nil {
		c.Out = &api.DeleteMatchmakingTicketResponse{
			Success: false,
			Error:   conversion.Enum(*mtErr, api.DeleteMatchmakingTicketResponse_Error_value, api.DeleteMatchmakingTicketResponse_MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED),
		}
		return nil
	}
	params := matchmakingTicketRequestToMatchmakingTicketParams(c.In)
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	result, err := qtx.DeleteMatchmakingTicket(ctx, params)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// If no rows were affected, it means the ticket was not found or it has a match associated with it and cannot be deleted.
		ticket, err := qtx.GetMatchmakingTicket(ctx, model.GetMatchmakingTicketParams{
			MatchmakingTicket: params,
			UserLimit:         1,
			ArenaLimit:        1,
		})
		if err != nil {
			return err
		}
		if len(ticket) == 0 {
			c.Out = &api.DeleteMatchmakingTicketResponse{
				Success: false,
				Error:   api.DeleteMatchmakingTicketResponse_NOT_FOUND,
			}
			return nil
		}
		if ticket[0].MatchmakingMatchID.Valid {
			c.Out = &api.DeleteMatchmakingTicketResponse{
				Success: false,
				Error:   api.DeleteMatchmakingTicketResponse_TICKET_CURRENTLY_IN_MATCH,
			}
			return nil
		}
		// Code should not reach here, but if it does, we return an error.
		return errors.New("Ticket not found even though it has a match associated with it; this is a bug")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.DeleteMatchmakingTicketResponse{
		Success: true,
		Error:   api.DeleteMatchmakingTicketResponse_NONE,
	}
	return nil
}
