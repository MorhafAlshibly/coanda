package team

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api/pb"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateTeamScoreCommand struct {
	service *TeamService
	In      *pb.UpdateTeamScoreRequest
	Out     *pb.Team
}

func NewUpdateTeamScoreCommand(service *TeamService, in *pb.UpdateTeamScoreRequest) *UpdateTeamScoreCommand {
	return &UpdateTeamScoreCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateTeamScoreCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In.Team)
	if err != nil {
		return err
	}
	_, err = c.service.db.UpdateOne(ctx, filter, bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "score", Value: c.In.ScoreOffset},
		}},
	})
	if err != nil {
		if err.Error() == "EOF" {
			return errors.New("Team not found")
		}
		return err
	}
	c.Out, err = c.service.GetTeam(ctx, c.In.Team)
	if err != nil {
		return err
	}
	return nil
}
