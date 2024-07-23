package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type GetMatchmakingTicketCommand struct {
	service *Service
	In      *api.MatchmakingTicketRequest
	Out     *api.GetMatchmakingTicketResponse
}

func NewGetMatchmakingTicketCommand(service *Service, in *api.MatchmakingTicketRequest) *GetMatchmakingTicketCommand {
	return &GetMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchmakingTicketCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
