package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type DeleteArenaCommand struct {
	service *Service
	In      *api.ArenaRequest
	Out     *api.ArenaResponse
}

func NewDeleteArenaCommand(service *Service, in *api.ArenaRequest) *DeleteArenaCommand {
	return &DeleteArenaCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteArenaCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
