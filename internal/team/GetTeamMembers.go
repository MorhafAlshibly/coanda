package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/validation"
)

type GetTeamMembersCommand struct {
	service *Service
	In      *api.GetTeamMembersRequest
	Out     *api.GetTeamMembersResponse
}

func NewGetTeamMembersCommand(service *Service, in *api.GetTeamMembersRequest) *GetTeamMembersCommand {
	return &GetTeamMembersCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTeamMembersCommand) Execute(ctx context.Context) error {
	tErr := c.service.CheckForTeamRequestError(c.In.Team)
	// Check if error is found
	if tErr != nil {
		c.Out = &api.GetTeamMembersResponse{
			Success: false,
			Error:   conversion.Enum(*tErr, api.GetTeamMembersResponse_Error_value, api.GetTeamMembersResponse_NOT_FOUND),
		}
		return nil
	}
	limit, offset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	members, err := c.service.database.GetTeamMembers(ctx, model.GetTeamMembersParams{
		Name:   validation.ValidateAnSqlNullString(c.In.Team.Name),
		Owner:  validation.ValidateAUint64ToSqlNullInt64(c.In.Team.Owner),
		Member: validation.ValidateAUint64ToSqlNullInt64(c.In.Team.Member),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return err
	}
	outs := make([]*api.TeamMember, len(members))
	for i, member := range members {
		outs[i], err = UnmarshalTeamMember(member)
		if err != nil {
			return err
		}
	}
	c.Out = &api.GetTeamMembersResponse{
		Success:     true,
		TeamMembers: outs,
		Error:       api.GetTeamMembersResponse_NONE,
	}
	return nil

}
