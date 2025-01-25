package matchmaking

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/go-sql-driver/mysql"
	"open-match.dev/open-match/pkg/pb"
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
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	// Check if user has an active ticket
	ticketClient, err := c.service.queryServiceClient.QueryTicketIds(ctx, &pb.QueryTicketIdsRequest{
		Pool: &pb.Pool{
			Name: "default",
			TagPresentFilters: []*pb.TagPresentFilter{
				{
					Tag: fmt.Sprintf("User_%d", c.In.ClientUserId),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = ticketClient.Recv()
	if err != io.EOF {
		c.Out = &api.CreateMatchmakingUserResponse{
			Success: false,
			Error:   api.CreateMatchmakingUserResponse_ALREADY_EXISTS,
		}
		return nil
	}
	result, err := c.service.database.CreateMatchmakingUser(ctx, model.CreateMatchmakingUserParams{
		ClientUserID: c.In.ClientUserId,
		Elo:          c.In.Elo,
		Data:         data,
	})
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if errorcode.IsDuplicateEntry(mysqlErr, "matchmaking_user", "client_user_id") {
				result, err = qtx.UpdateMatchmakingUserByClientUserId(ctx, model.UpdateMatchmakingUserByClientUserIdParams{
					ClientUserID: c.In.ClientUserId,
					Elo:          c.In.Elo,
					Data:         data,
				})
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			return err
		}
	}
	matchmakingUserId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	err = tx.Commit()
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
