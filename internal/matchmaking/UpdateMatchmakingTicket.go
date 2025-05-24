package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type UpdateMatchmakingTicketCommand struct {
	service *Service
	In      *api.UpdateMatchmakingTicketRequest
	Out     *api.UpdateMatchmakingTicketResponse
}

func NewUpdateMatchmakingTicketCommand(service *Service, in *api.UpdateMatchmakingTicketRequest) *UpdateMatchmakingTicketCommand {
	return &UpdateMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateMatchmakingTicketCommand) Execute(ctx context.Context) error {
	mtErr := c.service.checkForMatchmakingTicketRequestError(c.In.MatchmakingTicket)
	// Check if error is found
	if mtErr != nil {
		c.Out = &api.UpdateMatchmakingTicketResponse{
			Success: false,
			Error:   conversion.Enum(*mtErr, api.UpdateMatchmakingTicketResponse_Error_value, api.UpdateMatchmakingTicketResponse_MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED),
		}
		return nil
	}
	// Check if data is given
	if c.In.Data == nil {
		c.Out = &api.UpdateMatchmakingTicketResponse{
			Success: false,
			Error:   api.UpdateMatchmakingTicketResponse_DATA_REQUIRED,
		}
		return nil
	}
	// Prepare data
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	params := matchmakingTicketRequestToMatchmakingTicketParams(c.In.MatchmakingTicket)
	result, err := c.service.database.UpdateMatchmakingTicket(ctx, model.UpdateMatchmakingTicketParams{
		MatchmakingTicket: params,
		Data:              data,
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
		ticket, err := c.service.database.GetMatchmakingTicket(ctx, model.GetMatchmakingTicketParams{
			MatchmakingTicket: params,
			UserLimit:         1,
			ArenaLimit:        1,
		})
		if err != nil {
			return err
		}
		if len(ticket) == 0 {
			c.Out = &api.UpdateMatchmakingTicketResponse{
				Success: false,
				Error:   api.UpdateMatchmakingTicketResponse_NOT_FOUND,
			}
			return nil
		}
	}
	c.Out = &api.UpdateMatchmakingTicketResponse{
		Success: true,
		Error:   api.UpdateMatchmakingTicketResponse_NONE,
	}
	return nil
}
