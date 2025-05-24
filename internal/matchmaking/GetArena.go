package matchmaking

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetArenaCommand struct {
	service *Service
	In      *api.ArenaRequest
	Out     *api.GetArenaResponse
}

func NewGetArenaCommand(service *Service, in *api.ArenaRequest) *GetArenaCommand {
	return &GetArenaCommand{
		service: service,
		In:      in,
	}
}

func (c *GetArenaCommand) Execute(ctx context.Context) error {
	aErr := c.service.checkForArenaRequestError(c.In)
	if aErr != nil {
		c.Out = &api.GetArenaResponse{
			Success: false,
			Error:   conversion.Enum(*aErr, api.GetArenaResponse_Error_value, api.GetArenaResponse_ARENA_ID_OR_NAME_REQUIRED),
		}
		return nil
	}
	arena, err := c.service.database.GetArena(ctx, arenaRequestToArenaParams(c.In))
	if err != nil {
		if err == sql.ErrNoRows {
			c.Out = &api.GetArenaResponse{
				Success: false,
				Error:   api.GetArenaResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	apiArena, err := unmarshalArena(arena)
	if err != nil {
		return err
	}
	c.Out = &api.GetArenaResponse{
		Success: true,
		Arena:   apiArena,
		Error:   api.GetArenaResponse_NONE,
	}
	return nil
}
