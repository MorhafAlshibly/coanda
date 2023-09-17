package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api/pb"
)

type DeleteTeamCommand struct {
	service *TeamService
	In      *pb.DeleteTeamRequest
	Out     *pb.Team
}

func NewDeleteTeamCommand(service *TeamService, in *pb.DeleteTeamRequest) *DeleteTeamCommand {
	return &DeleteTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteTeamCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In.Team)
	if err != nil {
		return err
	}
	c.Out, err = c.service.GetTeam(ctx, c.In.Team)
	if err != nil {
		return err
	}
	_, err = c.service.db.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
