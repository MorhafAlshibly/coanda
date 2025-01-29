package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type DeleteAllExpiredMatchmakingTicketsCommand struct {
	service *Service
	Out     *api.DeleteAllExpiredMatchmakingTicketsResponse
}

func NewDeleteAllExpiredMatchmakingTicketsCommand(service *Service) *DeleteAllExpiredMatchmakingTicketsCommand {
	return &DeleteAllExpiredMatchmakingTicketsCommand{
		service: service,
	}
}

func (c *DeleteAllExpiredMatchmakingTicketsCommand) Execute(ctx context.Context) error {
	result, err := c.service.database.DeleteAllExpiredTickets(ctx)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	c.Out = &api.DeleteAllExpiredMatchmakingTicketsResponse{
		Success:      true,
		DeletedCount: uint64(rowsAffected),
	}
	return nil
}
