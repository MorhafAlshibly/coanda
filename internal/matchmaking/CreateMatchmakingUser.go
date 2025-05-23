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
	if c.In.ClientUserId == 0 {
		c.Out = &api.CreateMatchmakingUserResponse{
			Success: false,
			Error:   api.CreateMatchmakingUserResponse_CLIENT_USER_ID_REQUIRED,
		}
		return nil
	}
	if c.In.Data == nil {
		c.Out = &api.CreateMatchmakingUserResponse{
			Success: false,
			Error:   api.CreateMatchmakingUserResponse_DATA_REQUIRED,
		}
		return nil
	}
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	result, err := c.service.database.CreateMatchmakingUser(ctx, model.CreateMatchmakingUserParams{
		ClientUserID: c.In.ClientUserId,
		Elo:          c.In.Elo,
		Data:         data,
	})
	var matchmakingUserId int64
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if errorcode.IsDuplicateEntry(mysqlErr, "matchmaking_user", "client_user_id") {
				c.Out = &api.CreateMatchmakingUserResponse{
					Success: false,
					Error:   api.CreateMatchmakingUserResponse_ALREADY_EXISTS,
				}
				return nil
			}
		}
		return err
	}
	matchmakingUserId, err = result.LastInsertId()
	if err != nil {
		return err
	}
	c.Out = &api.CreateMatchmakingUserResponse{
		Success: true,
		Id:      conversion.ValueToPointer(uint64(matchmakingUserId)),
		Error:   api.CreateMatchmakingUserResponse_NONE,
	}
	return nil
}
