package team

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	errorcodes "github.com/MorhafAlshibly/coanda/pkg/errorcodes"
	"github.com/go-sql-driver/mysql"
)

type LeaveTeamCommand struct {
	service *Service
	In      *api.LeaveTeamRequest
	Out     *api.LeaveTeamResponse
}

func NewLeaveTeamCommand(service *Service, in *api.LeaveTeamRequest) *LeaveTeamCommand {
	return &LeaveTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *LeaveTeamCommand) Execute(ctx context.Context) error {
	// Check if user id is provided
	if c.In.UserId == 0 {
		c.Out = &api.LeaveTeamResponse{
			Success: false,
			Error:   api.LeaveTeamResponse_USER_ID_REQUIRED,
		}
		return nil
	}
	result, err := c.service.database.DeleteTeamMember(ctx, c.In.UserId)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == errorcodes.MySQLErrorCodeRowIsReferenced2 {
			c.Out = &api.LeaveTeamResponse{
				Success: false,
				Error:   api.LeaveTeamResponse_MEMBER_IS_OWNER,
			}
			return nil
		}
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		c.Out = &api.LeaveTeamResponse{
			Success: false,
			Error:   api.LeaveTeamResponse_NOT_IN_TEAM,
		}
		return nil
	}
	c.Out = &api.LeaveTeamResponse{
		Success: true,
		Error:   api.LeaveTeamResponse_NONE,
	}
	return nil
}
