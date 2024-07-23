package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type ExpireMatchmakingTicketCommand struct {
	service *Service
	In      *api.MatchmakingTicketRequest
	Out     *api.MatchmakingTicketResponse
}

func NewExpireMatchmakingTicketCommand(service *Service, in *api.MatchmakingTicketRequest) *ExpireMatchmakingTicketCommand {
	return &ExpireMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *ExpireMatchmakingTicketCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
