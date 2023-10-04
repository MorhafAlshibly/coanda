package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateTeamDataCommand struct {
	service *Service
	In      *api.UpdateTeamDataRequest
	Out     *api.TeamResponse
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
		c.Out = &api.TeamResponse{
			Success: false,
			Error:   api.TeamResponse_INVALID,
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
