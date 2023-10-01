package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateTeamDataCommand struct {
	service *Service
	In      *api.UpdateTeamDataRequest
	Out     *api.GetTeamResponse
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
		c.Out = &api.GetTeamResponse{
			Success: false,
			Team:    nil,
			Error:   api.GetTeamResponse_INVALID,
		}
		return nil
	}
	_, err = c.service.db.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "data", Value: c.In.Data},
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
