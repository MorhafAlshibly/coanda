package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetArenasCommand struct {
	service *Service
	In      *api.Pagination
	Out     *api.GetArenasResponse
}

func NewGetArenasCommand(service *Service, in *api.Pagination) *GetArenasCommand {
	return &GetArenasCommand{
		service: service,
		In:      in,
	}
}

func (c *GetArenasCommand) Execute(ctx context.Context) error {
	limit, offset := conversion.PaginationToLimitOffset(c.In, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	arenas, err := c.service.database.GetArenas(ctx, model.GetArenasParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return err
	}
	outs := make([]*api.Arena, len(arenas))
	for i, arena := range arenas {
		outs[i], err = unmarshalArena(arena)
		if err != nil {
			return err
		}
	}
	c.Out = &api.GetArenasResponse{
		Success: true,
		Arenas:  outs,
	}
	return nil
}
