package team

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"go.mongodb.org/mongo-driver/bson"
)

type LeaveTeamCommand struct {
	service *Service
	In      *api.LeaveTeamRequest
	Out     *api.GetTeamResponse
}

func NewLeaveTeamCommand(service *Service, in *api.LeaveTeamRequest) *LeaveTeamCommand {
	return &LeaveTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *LeaveTeamCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In.Team)
	if err != nil {
		c.Out = &api.GetTeamResponse{
			Success: false,
			Team:    nil,
			Error:   api.GetTeamResponse_INVALID,
		}
		return nil
	}
	result, err := c.service.db.UpdateOne(ctx, append(filter, bson.E{
		Key: "membersWithoutOwner", Value: c.In.UserId,
	}),
		bson.D{
			{Key: "$pull", Value: bson.D{
				{Key: "membersWithoutOwner", Value: c.In.UserId},
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
	if result.ModifiedCount == 0 {
		return errors.New("User is not a member of the team")
	}
	getTeamCommand := NewGetTeamCommand(c.service, c.In.Team)
	err = invokers.NewBasicInvoker().Invoke(ctx, getTeamCommand)
	if err != nil {
		return err
	}
	c.Out = getTeamCommand.Out
	return nil
}
