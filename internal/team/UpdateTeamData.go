package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateTeamDataCommand struct {
	service *Service
	In      *api.UpdateTeamDataRequest
	Out     *api.Team
}

func NewUpdateTeamDataCommand(service *Service, in *api.UpdateTeamDataRequest) *UpdateTeamDataCommand {
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
	_, err = c.service.db.UpdateOne(ctx, filter, bson.D{
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
