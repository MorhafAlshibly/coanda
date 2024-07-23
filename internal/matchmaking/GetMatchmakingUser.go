package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type GetMatchmakingUserCommand struct {
	service *Service
	In      *api.MatchmakingUserRequest
	Out     *api.GetMatchmakingUserResponse
}

func NewGetMatchmakingUserCommand(service *Service, in *api.MatchmakingUserRequest) *GetMatchmakingUserCommand {
	return &GetMatchmakingUserCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchmakingUserCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
