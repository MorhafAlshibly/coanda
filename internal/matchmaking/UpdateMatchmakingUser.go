package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type UpdateMatchmakingUserCommand struct {
	service *Service
	In      *api.UpdateMatchmakingUserRequest
	Out     *api.UpdateMatchmakingUserResponse
}

func NewUpdateMatchmakingUserCommand(service *Service, in *api.UpdateMatchmakingUserRequest) *UpdateMatchmakingUserCommand {
	return &UpdateMatchmakingUserCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateMatchmakingUserCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
