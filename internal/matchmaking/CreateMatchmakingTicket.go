package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type CreateMatchmakingTicketCommand struct {
	service *Service
	In      *api.CreateMatchmakingTicketRequest
	Out     *api.CreateMatchmakingTicketResponse
}

func NewCreateMatchmakingTicketCommand(service *Service, in *api.CreateMatchmakingTicketRequest) *CreateMatchmakingTicketCommand {
	return &CreateMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateMatchmakingTicketCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
