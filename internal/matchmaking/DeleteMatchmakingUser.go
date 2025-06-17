package matchmaking

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type DeleteMatchmakingUserCommand struct {
	service *Service
	In      *api.MatchmakingUserRequest
	Out     *api.DeleteMatchmakingUserResponse
}

func NewDeleteMatchmakingUserCommand(service *Service, in *api.MatchmakingUserRequest) *DeleteMatchmakingUserCommand {
	return &DeleteMatchmakingUserCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteMatchmakingUserCommand) Execute(ctx context.Context) error {
	mtErr := c.service.checkForMatchmakingUserRequestError(c.In)
	if mtErr != nil {
		c.Out = &api.DeleteMatchmakingUserResponse{
			Success: false,
			Error:   conversion.Enum(*mtErr, api.DeleteMatchmakingUserResponse_Error_value, api.DeleteMatchmakingUserResponse_MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED),
		}
		return nil
	}
	params := matchmakingUserRequestToMatchmakingUserParams(c.In)
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	result, err := qtx.DeleteMatchmakingUser(ctx, params)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// If no rows were affected, it means the User was not found or it has a match associated with it and cannot be deleted.
		user, err := qtx.GetMatchmakingUser(ctx, params, nil)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Out = &api.DeleteMatchmakingUserResponse{
					Success: false,
					Error:   api.DeleteMatchmakingUserResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
		if user.MatchmakingTicketID.Valid {
			// If the user is currently in a ticket, we cannot delete it.
			// Lets check if the ticket has a match associated with it.
			ticket, err := qtx.GetMatchmakingTicket(ctx, model.GetMatchmakingTicketParams{
				MatchmakingTicket: model.MatchmakingTicketParams{
					ID: user.MatchmakingTicketID,
				},
				UserLimit:  1,
				ArenaLimit: 1,
			})
			if err != nil {
				return err
			}
			if len(ticket) == 0 {
				return errors.New("Ticket not found even though user is in a ticket; this is a bug")
			}
			if ticket[0].MatchmakingMatchID.Valid {
				// If the ticket has a match associated with it, we cannot delete the user.
				c.Out = &api.DeleteMatchmakingUserResponse{
					Success: false,
					Error:   api.DeleteMatchmakingUserResponse_USER_CURRENTLY_IN_MATCH,
				}
				return nil
			}
			c.Out = &api.DeleteMatchmakingUserResponse{
				Success: false,
				Error:   api.DeleteMatchmakingUserResponse_USER_CURRENTLY_IN_TICKET,
			}
			return nil
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.DeleteMatchmakingUserResponse{
		Success: true,
		Error:   api.DeleteMatchmakingUserResponse_NONE,
	}
	return nil
}
