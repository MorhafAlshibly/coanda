package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type GetArenasCommand struct {
	service *Service
	In      *api.GetArenasRequest
	Out     *api.GetArenasResponse
}

func NewGetArenasCommand(service *Service, in *api.GetArenasRequest) *GetArenasCommand {
	return &GetArenasCommand{
		service: service,
		In:      in,
	}
}

func (c *GetArenasCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
