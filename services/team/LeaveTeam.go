package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/services/team/schema"
	"github.com/bytedance/sonic"
)

type LeaveTeamCommand struct {
	service *TeamService
	In      *schema.LeaveTeamRequest
	Out     *schema.BoolResponse
}

func NewLeaveTeamCommand(service *TeamService, in *schema.LeaveTeamRequest) *LeaveTeamCommand {
	return &LeaveTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *LeaveTeamCommand) Execute(ctx context.Context) error {
	marshalled, err := sonic.Marshal(c.In.Team)
	if err != nil {
		return err
	}
	err = c.service.Queue.Enqueue(ctx, "LeaveTeam", marshalled)
	if err != nil {
		return err
	}
	c.Out = &schema.BoolResponse{Value: true}
	return nil
}
