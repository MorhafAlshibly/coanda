package matchmaking

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/go-sql-driver/mysql"
)

type CreateArenaCommand struct {
	service *Service
	In      *api.CreateArenaRequest
	Out     *api.CreateArenaResponse
}

func NewCreateArenaCommand(service *Service, in *api.CreateArenaRequest) *CreateArenaCommand {
	return &CreateArenaCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateArenaCommand) Execute(ctx context.Context) error {
	if len(c.In.Name) < int(c.service.minArenaNameLength) {
		c.Out = &api.CreateArenaResponse{
			Success: false,
			Error:   api.CreateArenaResponse_NAME_TOO_SHORT,
		}
		return nil
	}
	if len(c.In.Name) > int(c.service.maxArenaNameLength) {
		c.Out = &api.CreateArenaResponse{
			Success: false,
			Error:   api.CreateArenaResponse_NAME_TOO_LONG,
		}
		return nil
	}
	if c.In.MinPlayers == 0 {
		c.Out = &api.CreateArenaResponse{
			Success: false,
			Error:   api.CreateArenaResponse_MIN_PLAYERS_REQUIRED,
		}
		return nil
	}
	if c.In.MaxPlayersPerTicket == 0 {
		c.Out = &api.CreateArenaResponse{
			Success: false,
			Error:   api.CreateArenaResponse_MAX_PLAYERS_PER_TICKET_REQUIRED,
		}
		return nil
	}
	if c.In.MaxPlayers == 0 {
		c.Out = &api.CreateArenaResponse{
			Success: false,
			Error:   api.CreateArenaResponse_MAX_PLAYERS_REQUIRED,
		}
		return nil
	}
	if c.In.MinPlayers > c.In.MaxPlayers {
		c.Out = &api.CreateArenaResponse{
			Success: false,
			Error:   api.CreateArenaResponse_MIN_PLAYERS_CANNOT_BE_GREATER_THAN_MAX_PLAYERS,
		}
		return nil
	}
	if c.In.MaxPlayersPerTicket > c.In.MaxPlayers {
		c.Out = &api.CreateArenaResponse{
			Success: false,
			Error:   api.CreateArenaResponse_MAX_PLAYERS_PER_TICKET_CANNOT_BE_GREATER_THAN_MAX_PLAYERS,
		}
		return nil
	}
	if c.In.Data == nil {
		c.Out = &api.CreateArenaResponse{
			Success: false,
			Error:   api.CreateArenaResponse_DATA_REQUIRED,
		}
		return nil
	}
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	result, err := c.service.database.CreateArena(ctx, model.CreateArenaParams{
		Name:                c.In.Name,
		MinPlayers:          uint32(c.In.MinPlayers),
		MaxPlayersPerTicket: uint32(c.In.MaxPlayersPerTicket),
		MaxPlayers:          uint32(c.In.MaxPlayers),
		Data:                data,
	})
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if errorcode.IsDuplicateEntry(mysqlErr, "matchmaking_arena", "name") {
				c.Out = &api.CreateArenaResponse{
					Success: false,
					Error:   api.CreateArenaResponse_ALREADY_EXISTS,
				}
				return nil
			}
		}
		return err
	}
	arenaId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	c.Out = &api.CreateArenaResponse{
		Success: true,
		Id:      conversion.ValueToPointer(uint64(arenaId)),
		Error:   api.CreateArenaResponse_NONE,
	}
	return nil
}
