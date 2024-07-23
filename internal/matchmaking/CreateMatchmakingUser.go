package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type CreateMatchmakingUserCommand struct {
	service *Service
	In      *api.CreateMatchmakingUserRequest
	Out     *api.CreateMatchmakingUserResponse
}

func NewCreateMatchmakingUserCommand(service *Service, in *api.CreateMatchmakingUserRequest) *CreateMatchmakingUserCommand {
	return &CreateMatchmakingUserCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateMatchmakingUserCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
