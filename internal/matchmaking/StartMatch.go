package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type StartMatchCommand struct {
	service *Service
	In      *api.StartMatchRequest
	Out     *api.StartMatchResponse
}

func NewStartMatchCommand(service *Service, in *api.StartMatchRequest) *StartMatchCommand {
	return &StartMatchCommand{
		service: service,
		In:      in,
	}
}

func (c *StartMatchCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
