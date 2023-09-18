package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
)

type GetTeamCommand struct {
	service *Service
	In      *api.GetTeamRequest
	Out     *api.Team
}

func NewGetTeamCommand(service *Service, in *api.GetTeamRequest) *GetTeamCommand {
	return &GetTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTeamCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In)
	if err != nil {
		return err
	}
	// Get the item from the store
	matchStage := bson.D{
		{Key: "$match", Value: filter},
	}
	cursor, err := c.service.db.Aggregate(ctx, append(pipeline, matchStage))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	cursor.Next(ctx)
	c.Out, err = toTeam(cursor)
	if err != nil {
		return err
	}
	return nil
}
