package team

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type DeleteTeamCommand struct {
	service *Service
	In      *api.TeamRequest
	Out     *api.TeamResponse
}

func NewDeleteTeamCommand(service *Service, in *api.TeamRequest) *DeleteTeamCommand {
	return &DeleteTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteTeamCommand) Execute(ctx context.Context) error {
	field := c.service.GetTeamField(c.In)
	if field != NAME && field != OWNER && field != MEMBER {
		c.Out = &api.TeamResponse{
			Success: false,
			Error:   conversion.Enum(field, api.TeamResponse_Error_value, api.TeamResponse_NOT_FOUND),
		}
		return nil
	}
	var err error
	// Check if name or owner is provided
	if field == NAME || field == OWNER {
		err = c.service.database.DeleteTeam(ctx, model.DeleteTeamParams{
			Name: sql.NullString{
				String: *c.In.Name,
				Valid:  field == NAME,
			},
			Owner: sql.NullInt64{
				Int64: int64(*c.In.Owner),
				Valid: field == OWNER,
			}})
		// Check if member is provided
	} else if field == MEMBER {
		err = c.service.database.DeleteTeamByMember(
			ctx,
			*c.In.Member,
		)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Out = &api.TeamResponse{
				Success: false,
				Error:   api.TeamResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	c.Out = &api.TeamResponse{
		Success: true,
		Error:   api.TeamResponse_NONE,
	}
	return nil
}
