package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
)

type GetTeamCommand struct {
	service *Service
	In      *api.GetTeamRequest
	Out     *api.GetTeamResponse
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
		c.Out = &api.GetTeamResponse{
			Success: false,
			Team:    nil,
			Error:   api.GetTeamResponse_INVALID,
		}
		return nil
	}
	// Get the item from the store
	matchStage := bson.D{
		{Key: "$match", Value: filter},
	}
	cursor, err := c.service.database.Aggregate(ctx, append(pipeline, matchStage))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	cursor.Next(ctx)
	team, err := toTeam(cursor)
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
	c.Out = &api.GetTeamResponse{
		Success: true,
		Team:    team,
		Error:   api.GetTeamResponse_NONE,
	}
	return nil
}
