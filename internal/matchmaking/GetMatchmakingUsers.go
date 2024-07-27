package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetMatchmakingUsersCommand struct {
	service *Service
	In      *api.Pagination
	Out     *api.GetMatchmakingUsersResponse
}

func NewGetMatchmakingUsersCommand(service *Service, in *api.Pagination) *GetMatchmakingUsersCommand {
	return &GetMatchmakingUsersCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchmakingUsersCommand) Execute(ctx context.Context) error {
	limit, offset := conversion.PaginationToLimitOffset(c.In, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	users, err := c.service.database.GetMatchmakingUsers(ctx, model.GetMatchmakingUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return err
	}
	outs := make([]*api.MatchmakingUser, len(users))
	for i, user := range users {
		outs[i], err = unmarshalMatchmakingUser(user)
		if err != nil {
			return err
		}
	}
	c.Out = &api.GetMatchmakingUsersResponse{
		Success:          true,
		MatchmakingUsers: outs,
	}
	return nil
}
