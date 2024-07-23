package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type GetMatchmakingTicketsCommand struct {
	service *Service
	In      *api.GetMatchmakingTicketsRequest
	Out     *api.GetMatchmakingTicketsResponse
}

func NewGetMatchmakingTicketsCommand(service *Service, in *api.GetMatchmakingTicketsRequest) *GetMatchmakingTicketsCommand {
	return &GetMatchmakingTicketsCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchmakingTicketsCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
