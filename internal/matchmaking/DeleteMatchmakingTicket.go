package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
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
	result, err := c.service.database.DeleteMatchmakingTicket(ctx, params)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		c.Out = &api.DeleteMatchmakingTicketResponse{
			Success: false,
			Error:   api.DeleteMatchmakingTicketResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.DeleteMatchmakingTicketResponse{
		Success: true,
		Error:   api.DeleteMatchmakingTicketResponse_NONE,
	}
	return nil
}
