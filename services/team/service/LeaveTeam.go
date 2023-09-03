package service

import (
	"context"

	"github.com/MorhafAlshibly/coanda/services/team/schema"
	"github.com/bytedance/sonic"
)

type LeaveTeamCommand struct {
	Service *TeamService
	In      *schema.LeaveTeamRequest
	Out     *schema.BoolResponse
}

func (c *LeaveTeamCommand) Execute(ctx context.Context) error {
	marshalled, err := sonic.Marshal(c.In.Team)
	if err != nil {
		return err
	}
	err = c.Service.Queue.Enqueue(ctx, "LeaveTeam", marshalled)
	if err != nil {
		return err
	}
	c.Out = &schema.BoolResponse{Value: true}
	return nil
}
