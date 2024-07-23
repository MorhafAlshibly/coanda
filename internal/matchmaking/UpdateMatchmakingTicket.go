package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
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
	panic("implement me")
}
