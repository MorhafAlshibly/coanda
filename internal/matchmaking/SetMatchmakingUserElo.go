package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type SetMatchmakingUserEloCommand struct {
	service *Service
	In      *api.SetMatchmakingUserEloRequest
	Out     *api.SetMatchmakingUserEloResponse
}

func NewSetMatchmakingUserEloCommand(service *Service, in *api.SetMatchmakingUserEloRequest) *SetMatchmakingUserEloCommand {
	return &SetMatchmakingUserEloCommand{
		service: service,
		In:      in,
	}
}

func (c *SetMatchmakingUserEloCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
