package matchmaking

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
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
	muErr := c.service.checkForMatchmakingUserRequestError(c.In.MatchmakingUser)
	// Check if error is found
	if muErr != nil {
		c.Out = &api.UpdateMatchmakingUserResponse{
			Success: false,
			Error:   conversion.Enum(*muErr, api.UpdateMatchmakingUserResponse_Error_value, api.UpdateMatchmakingUserResponse_MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED),
		}
		return nil
	}
	// Check if data is given
	if c.In.Data == nil {
		c.Out = &api.UpdateMatchmakingUserResponse{
			Success: false,
			Error:   api.UpdateMatchmakingUserResponse_DATA_REQUIRED,
		}
		return nil
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := model.New(tx)
	// Prepare data
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	result, err := qtx.UpdateMatchmakingUser(ctx, model.UpdateMatchmakingUserParams{
		MatchmakingUser: matchmakingUserRequestToMatchmakingUserParams(c.In.MatchmakingUser),
		Data:            data,
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Check if we didn't find a row
		_, err = qtx.GetMatchmakingUser(ctx, matchmakingUserRequestToMatchmakingUserParams(c.In.MatchmakingUser))
		if err != nil {
			if err == sql.ErrNoRows {
				// If we didn't find a row
				c.Out = &api.UpdateMatchmakingUserResponse{
					Success: false,
					Error:   api.UpdateMatchmakingUserResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.UpdateMatchmakingUserResponse{
		Success: true,
		Error:   api.UpdateMatchmakingUserResponse_NONE,
	}
	return nil
}
