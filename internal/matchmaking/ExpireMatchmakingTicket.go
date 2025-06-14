package matchmaking

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type ExpireMatchmakingTicketCommand struct {
	service *Service
	In      *api.MatchmakingTicketRequest
	Out     *api.ExpireMatchmakingTicketResponse
}

func NewExpireMatchmakingTicketCommand(service *Service, in *api.MatchmakingTicketRequest) *ExpireMatchmakingTicketCommand {
	return &ExpireMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *ExpireMatchmakingTicketCommand) Execute(ctx context.Context) error {
	mtErr := c.service.checkForMatchmakingTicketRequestError(c.In)
	if mtErr != nil {
		c.Out = &api.ExpireMatchmakingTicketResponse{
			Success: false,
			Error:   conversion.Enum(*mtErr, api.ExpireMatchmakingTicketResponse_Error_value, api.ExpireMatchmakingTicketResponse_MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED),
		}
		return nil
	}
	params := matchmakingTicketRequestToMatchmakingTicketParams(c.In)
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := model.New(tx)
	result, err := qtx.ExpireMatchmakingTicket(ctx, params)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Check if we didn't find a row
		ticket, err := qtx.GetMatchmakingTicket(ctx, model.GetMatchmakingTicketParams{
			MatchmakingTicket: params,
			UserLimit:         1,
			ArenaLimit:        1,
		})
		// Check if ticket is found
		if err != nil {
			return err
		}
		if len(ticket) == 0 {
			c.Out = &api.ExpireMatchmakingTicketResponse{
				Success: false,
				Error:   api.ExpireMatchmakingTicketResponse_NOT_FOUND,
			}
			return nil
		}
		if ticket[0].Status == "EXPIRED" {
			c.Out = &api.ExpireMatchmakingTicketResponse{
				Success: false,
				Error:   api.ExpireMatchmakingTicketResponse_ALREADY_EXPIRED,
			}
			return nil
		} else if ticket[0].Status == "MATCHED" {
			c.Out = &api.ExpireMatchmakingTicketResponse{
				Success: false,
				Error:   api.ExpireMatchmakingTicketResponse_ALREADY_MATCHED,
			}
			return nil
		} else if ticket[0].Status == "ENDED" {
			c.Out = &api.ExpireMatchmakingTicketResponse{
				Success: false,
				Error:   api.ExpireMatchmakingTicketResponse_ALREADY_ENDED,
			}
			return nil
		} else {
			// Unexpected error
			return errors.New("could not expire ticket")
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.ExpireMatchmakingTicketResponse{
		Success: true,
		Error:   api.ExpireMatchmakingTicketResponse_NONE,
	}
	return nil
}
