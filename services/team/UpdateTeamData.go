package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/services/team/schema"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateTeamDataCommand struct {
	service *TeamService
	In      *schema.UpdateTeamDataRequest
	Out     *schema.Team
}

func NewUpdateTeamDataCommand(service *TeamService, in *schema.UpdateTeamDataRequest) *UpdateTeamDataCommand {
	return &UpdateTeamDataCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateTeamDataCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In.Team)
	if err != nil {
		return err
	}
	_, err = c.service.Db.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "data", Value: c.In.Data},
		}},
	})
	if err != nil {
		return err
	}
	c.Out, err = c.service.GetTeam(ctx, c.In.Team)
	if err != nil {
		return err
	}
	return nil
}
