package team

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api/pb"
	"go.mongodb.org/mongo-driver/bson"
)

type SearchTeamsCommand struct {
	service *TeamService
	In      *pb.SearchTeamsRequest
	Out     *pb.Teams
}

func NewSearchTeamsCommand(service *TeamService, in *pb.SearchTeamsRequest) *SearchTeamsCommand {
	return &SearchTeamsCommand{
		service: service,
		In:      in,
	}
}

func (c *SearchTeamsCommand) Execute(ctx context.Context) error {
	searchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "name", Value: bson.D{
				{Key: "$regex", Value: c.In.Query},
				{Key: "$options", Value: "i"},
			}},
		}},
	}
	if len(c.In.Query) < c.service.minTeamNameLength {
		return errors.New("Query too short")
	}
	cursor, err := c.service.db.Aggregate(ctx, append(c.service.pipeline, searchStage))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	c.Out, err = toTeams(ctx, cursor, c.In.Page, c.In.Max)
	if err != nil {
		return err
	}
	return nil
}
