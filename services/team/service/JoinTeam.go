package service

import (
	"context"

	"github.com/MorhafAlshibly/coanda/services/team/schema"
	"github.com/bytedance/sonic"
)

type JoinTeamCommand struct {
	Service *TeamService
	In      *schema.DeleteTeamRequest
	Out     *schema.BoolResponse
}

func (c *JoinTeamCommand) Execute(ctx context.Context) error {
	marshalled, err := sonic.Marshal(c.In.Team)
	if err != nil {
		return err
	}
	err = c.Service.Queue.Enqueue(ctx, "JoinTeam", marshalled)
	if err != nil {
		return err
	}
	c.Out = &schema.BoolResponse{Value: true}
	return nil
}
