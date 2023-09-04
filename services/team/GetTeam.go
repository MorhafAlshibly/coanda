package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/services/team/schema"
	"go.mongodb.org/mongo-driver/bson"
)

type GetTeamCommand struct {
	service *TeamService
	In      *schema.GetTeamRequest
	Out     *schema.Team
}

func NewGetTeamCommand(service *TeamService, in *schema.GetTeamRequest) *GetTeamCommand {
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
	cursor, err := c.service.Db.Aggregate(ctx, append(c.service.Pipeline, matchStage))
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
