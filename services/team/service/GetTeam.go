package service

import (
	"context"

	"github.com/MorhafAlshibly/coanda/services/team/schema"
	"go.mongodb.org/mongo-driver/bson"
)

type GetTeamCommand struct {
	Service *TeamService
	In      *schema.GetTeamRequest
	Out     *schema.Team
}

func (c *GetTeamCommand) Execute(ctx context.Context) error {
	filter, err := GetFilter(c.In)
	if err != nil {
		return err
	}
	// Get the item from the store
	matchStage := bson.D{
		{Key: "$match", Value: filter},
	}
	cursor, err := c.Service.Db.Aggregate(ctx, append(c.Service.Pipeline, matchStage))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	cursor.Next(ctx)
	c.Out, err = ToTeam(cursor)
	if err != nil {
		return err
	}
	return nil
}
