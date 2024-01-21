package team

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	errorcodes "github.com/MorhafAlshibly/coanda/pkg/errorCodes"
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
	if c.In.Owner == 0 {
		c.Out = &api.CreateTeamResponse{
			Success: false,
			Error:   api.CreateTeamResponse_OWNER_REQUIRED,
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
	// If score is not provided, set it to 0
	if c.In.Score == nil {
		c.In.Score = new(int64)
		*c.In.Score = 0
	}
	// Insert the team into the database
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	tx, err := c.service.sql.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	// Create the team owner
	_, err = qtx.CreateTeamOwner(ctx, model.CreateTeamOwnerParams{
		Team:   c.In.Name,
		UserID: c.In.Owner,
	})
	// If the team owner already exists, return appropriate error
	var mysqlErr *mysql.MySQLError
	if err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == errorcodes.MySQLErrorCodeDuplicateEntry {
			c.Out = &api.CreateTeamResponse{
				Success: false,
				Error:   api.CreateTeamResponse_OWNER_OWNS_ANOTHER_TEAM,
			}
			return nil
		}
		return err
	}
	// Create the team
	_, err = qtx.CreateTeam(ctx, model.CreateTeamParams{
		Name:  c.In.Name,
		Owner: c.In.Owner,
		Score: *c.In.Score,
		Data:  data,
	})
	// If the team already exists, return appropriate error
	if err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == errorcodes.MySQLErrorCodeDuplicateEntry {
			c.Out = &api.CreateTeamResponse{
				Success: false,
				Error:   api.CreateTeamResponse_NAME_TAKEN,
			}
			return nil
		}
		return err
	}
	// Create the owner as a member of the team
	_, err = qtx.CreateTeamMember(ctx, model.CreateTeamMemberParams{
		Team:       c.In.Name,
		UserID:     c.In.Owner,
		Data:       data,
		MaxMembers: int64(c.service.maxMembers),
	})
	// If the team member already exists, return appropriate error
	if err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == errorcodes.MySQLErrorCodeDuplicateEntry {
			c.Out = &api.CreateTeamResponse{
				Success: false,
				Error:   api.CreateTeamResponse_OWNER_ALREADY_IN_TEAM,
			}
			return nil
		}
		return err
	}
	c.Out = &api.CreateTeamResponse{
		Success: true,
		Error:   api.CreateTeamResponse_NONE,
	}
	return nil
}
