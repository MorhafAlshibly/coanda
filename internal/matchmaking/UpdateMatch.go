package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type UpdateMatchCommand struct {
	service *Service
	In      *api.UpdateMatchRequest
	Out     *api.UpdateMatchResponse
}

func NewUpdateMatchCommand(service *Service, in *api.UpdateMatchRequest) *UpdateMatchCommand {
	return &UpdateMatchCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateMatchCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
