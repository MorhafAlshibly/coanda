package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateTeamScoreCommand struct {
	service *Service
	In      *api.UpdateTeamScoreRequest
	Out     *api.TeamResponse
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
		c.Out = &api.TeamResponse{
			Success: false,
			Error:   api.TeamResponse_INVALID,
		}
		return nil
	}
	_, err = c.service.db.UpdateOne(ctx, filter, bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "score", Value: c.In.ScoreOffset},
		}},
	})
	if err != nil {
		if err.Error() == "EOF" {
			c.Out = &api.TeamResponse{
				Success: false,
				Error:   api.TeamResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	c.Out = &api.TeamResponse{
		Success: true,
		Error:   api.TeamResponse_NONE,
	}
	return nil
}
