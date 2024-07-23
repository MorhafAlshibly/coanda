package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type GetMatchesCommand struct {
	service *Service
	In      *api.GetMatchesRequest
	Out     *api.GetMatchesResponse
}

func NewGetMatchesCommand(service *Service, in *api.GetMatchesRequest) *GetMatchesCommand {
	return &GetMatchesCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchesCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
