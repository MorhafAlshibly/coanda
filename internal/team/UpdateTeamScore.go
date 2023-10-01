package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateTeamScoreCommand struct {
	service *Service
	In      *api.UpdateTeamScoreRequest
	Out     *api.GetTeamResponse
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
		c.Out = &api.GetTeamResponse{
			Success: false,
			Team:    nil,
			Error:   api.GetTeamResponse_INVALID,
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
			c.Out = &api.GetTeamResponse{
				Success: false,
				Team:    nil,
				Error:   api.GetTeamResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	getTeamCommand := NewGetTeamCommand(c.service, c.In.Team)
	err = invokers.NewBasicInvoker().Invoke(ctx, getTeamCommand)
	if err != nil {
		return err
	}
	c.Out = getTeamCommand.Out
	return nil
}
