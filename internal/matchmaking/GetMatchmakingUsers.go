package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type GetMatchmakingUsersCommand struct {
	service *Service
	In      *api.GetMatchmakingUsersRequest
	Out     *api.GetMatchmakingUsersResponse
}

func NewGetMatchmakingUsersCommand(service *Service, in *api.GetMatchmakingUsersRequest) *GetMatchmakingUsersCommand {
	return &GetMatchmakingUsersCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchmakingUsersCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
