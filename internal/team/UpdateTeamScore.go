package team

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateTeamScoreCommand struct {
	service *Service
	In      *api.UpdateTeamScoreRequest
	Out     *api.Team
}

func NewUpdateTeamScoreCommand(service *Service, in *api.UpdateTeamScoreRequest) *UpdateTeamScoreCommand {
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
