package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type GetMatchCommand struct {
	service *Service
	In      *api.MatchRequest
	Out     *api.GetMatchResponse
}

func NewGetMatchCommand(service *Service, in *api.MatchRequest) *GetMatchCommand {
	return &GetMatchCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
