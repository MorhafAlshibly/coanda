package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetTeamsCommand struct {
	service *Service
	In      *api.GetTeamsRequest
	Out     *api.GetTeamsResponse
}

func NewGetTeamsCommand(service *Service, in *api.GetTeamsRequest) *GetTeamsCommand {
	return &GetTeamsCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTeamsCommand) Execute(ctx context.Context) error {
	limit, offset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	memberLimit, memberOffset := conversion.PaginationToLimitOffset(c.In.MemberPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	teams, err := c.service.database.GetTeams(ctx, model.GetTeamsParams{
		MemberLimit:  int64(memberLimit),
		MemberOffset: int64(memberOffset),
		Limit:        int32(limit * memberLimit),
		Offset:       int32(limit*memberLimit + offset),
	})
	if err != nil {
		return err
	}
	outs, err := unmarshalTeamsWithMembers(teams)
	if err != nil {
		return err
	}
	c.Out = &api.GetTeamsResponse{
		Success: true,
		Teams:   outs,
	}
	return nil
}
