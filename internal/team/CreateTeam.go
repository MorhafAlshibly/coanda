package team

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	errorcode "github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/go-sql-driver/mysql"
)

type CreateTeamCommand struct {
	service *Service
	In      *api.CreateTeamRequest
	Out     *api.CreateTeamResponse
}

func NewCreateTeamCommand(service *Service, in *api.CreateTeamRequest) *CreateTeamCommand {
	return &CreateTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateTeamCommand) Execute(ctx context.Context) error {
	// Check if team name is correct length
	if len(c.In.Name) < int(c.service.minTeamNameLength) {
		c.Out = &api.CreateTeamResponse{
			Success: false,
			Error:   api.CreateTeamResponse_NAME_TOO_SHORT,
		}
		return nil
	}
	if len(c.In.Name) > int(c.service.maxTeamNameLength) {
		c.Out = &api.CreateTeamResponse{
			Success: false,
			Error:   api.CreateTeamResponse_NAME_TOO_LONG,
		}
		return nil
	}
	if c.In.FirstMemberUserId == 0 {
		c.Out = &api.CreateTeamResponse{
			Success: false,
			Error:   api.CreateTeamResponse_FIRST_MEMBER_USER_ID_REQUIRED,
		}
		return nil
	}
	// Check if data is provided
	if c.In.Data == nil {
		c.Out = &api.CreateTeamResponse{
			Success: false,
			Error:   api.CreateTeamResponse_DATA_REQUIRED,
		}
		return nil
	}
	// Check if first member data is provided
	if c.In.FirstMemberData == nil {
		c.Out = &api.CreateTeamResponse{
			Success: false,
			Error:   api.CreateTeamResponse_FIRST_MEMBER_DATA_REQUIRED,
		}
		return nil
	}
	// If score is not provided, set it to 0
	if c.In.Score == nil {
		c.In.Score = new(int64)
		*c.In.Score = 0
	}
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	firstMemberData, err := conversion.ProtobufStructToRawJson(c.In.FirstMemberData)
	if err != nil {
		return err
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	// Create the team
	result, err := qtx.CreateTeam(ctx, model.CreateTeamParams{
		Name:  c.In.Name,
		Score: *c.In.Score,
		Data:  data,
	})
	// If the team already exists, return appropriate error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == errorcode.MySQLErrorCodeDuplicateEntry {
				c.Out = &api.CreateTeamResponse{
					Success: false,
					Error:   api.CreateTeamResponse_NAME_TAKEN,
				}
				return nil
			}
		}
		return err
	}
	teamID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	// Create the first member
	_, err = qtx.CreateTeamMember(ctx, model.CreateTeamMemberParams{
		UserID:       c.In.FirstMemberUserId,
		TeamID:       uint64(teamID),
		MemberNumber: 1,
		Data:         firstMemberData,
	})
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == errorcode.MySQLErrorCodeDuplicateEntry {
			c.Out = &api.CreateTeamResponse{
				Success: false,
				Error:   api.CreateTeamResponse_FIRST_MEMBER_ALREADY_IN_A_TEAM,
			}
			return nil
		}
		return err
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.CreateTeamResponse{
		Success: true,
		Id:      conversion.ValueToPointer(uint64(teamID)),
		Error:   api.CreateTeamResponse_NONE,
	}
	return nil
}
