package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type GetArenaCommand struct {
	service *Service
	In      *api.ArenaRequest
	Out     *api.GetArenaResponse
}

func NewGetArenaCommand(service *Service, in *api.ArenaRequest) *GetArenaCommand {
	return &GetArenaCommand{
		service: service,
		In:      in,
	}
}

func (c *GetArenaCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
