package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type UpdateArenaCommand struct {
	service *Service
	In      *api.UpdateArenaRequest
	Out     *api.UpdateArenaResponse
}

func NewUpdateArenaCommand(service *Service, in *api.UpdateArenaRequest) *UpdateArenaCommand {
	return &UpdateArenaCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateArenaCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
