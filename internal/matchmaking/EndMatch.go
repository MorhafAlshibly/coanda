package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type EndMatchCommand struct {
	service *Service
	In      *api.EndMatchRequest
	Out     *api.EndMatchResponse
}

func NewEndMatchCommand(service *Service, in *api.EndMatchRequest) *EndMatchCommand {
	return &EndMatchCommand{
		service: service,
		In:      in,
	}
}

func (c *EndMatchCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
