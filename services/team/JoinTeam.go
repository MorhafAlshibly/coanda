package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/services/team/schema"
	"github.com/bytedance/sonic"
)

type JoinTeamCommand struct {
	service *TeamService
	In      *schema.JoinTeamRequest
	Out     *schema.BoolResponse
}

func NewJoinTeamCommand(service *TeamService, in *schema.JoinTeamRequest) *JoinTeamCommand {
	return &JoinTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *JoinTeamCommand) Execute(ctx context.Context) error {
	marshalled, err := sonic.Marshal(c.In.Team)
	if err != nil {
		return err
	}
	err = c.service.queue.Enqueue(ctx, "JoinTeam", marshalled)
	if err != nil {
		return err
	}
	c.Out = &schema.BoolResponse{Value: true}
	return nil
}
