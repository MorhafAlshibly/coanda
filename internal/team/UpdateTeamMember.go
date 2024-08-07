package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type UpdateTeamMemberCommand struct {
	service *Service
	In      *api.UpdateTeamMemberRequest
	Out     *api.UpdateTeamMemberResponse
}

func NewUpdateTeamMemberCommand(service *Service, in *api.UpdateTeamMemberRequest) *UpdateTeamMemberCommand {
	return &UpdateTeamMemberCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateTeamMemberCommand) Execute(ctx context.Context) error {
	if c.In.UserId == 0 {
		c.Out = &api.UpdateTeamMemberResponse{
			Success: false,
			Error:   api.UpdateTeamMemberResponse_USER_ID_REQUIRED,
		}
		return nil
	}
	if c.In.Data == nil {
		c.Out = &api.UpdateTeamMemberResponse{
			Success: false,
			Error:   api.UpdateTeamMemberResponse_DATA_REQUIRED,
		}
		return nil
	}
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	result, err := c.service.database.UpdateTeamMember(ctx, model.UpdateTeamMemberParams{
		UserID: c.In.UserId,
		Data:   data,
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		c.Out = &api.UpdateTeamMemberResponse{
			Success: false,
			Error:   api.UpdateTeamMemberResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.UpdateTeamMemberResponse{
		Success: true,
		Error:   api.UpdateTeamMemberResponse_NONE,
	}
	return nil
}
