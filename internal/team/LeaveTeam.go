package team

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
)

type LeaveTeamCommand struct {
	service *Service
	In      *api.LeaveTeamRequest
	Out     *api.LeaveTeamResponse
}

func NewLeaveTeamCommand(service *Service, in *api.LeaveTeamRequest) *LeaveTeamCommand {
	return &LeaveTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *LeaveTeamCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In.Team)
	if err != nil || c.In.UserId == 0 {
		c.Out = &api.LeaveTeamResponse{
			Success: false,
			Error:   api.LeaveTeamResponse_INVALID,
		}
		return nil
	}
	if c.In.UserId == 0 {
		return errors.New("UserId is required")
	}
	result, writeErr := c.service.db.UpdateOne(ctx, filter,
		bson.D{
			{Key: "$pull", Value: bson.D{
				{Key: "membersWithoutOwner", Value: c.In.UserId},
			}},
		})
	if writeErr != nil {
		return writeErr
	}
	if result.MatchedCount == 0 {
		c.Out = &api.LeaveTeamResponse{
			Success: false,
			Error:   api.LeaveTeamResponse_NOT_FOUND,
		}
		return nil
	}
	if result.ModifiedCount == 0 {
		c.Out = &api.LeaveTeamResponse{
			Success: false,
			Error:   api.LeaveTeamResponse_NOT_MEMBER,
		}
		return nil
	}
	c.Out = &api.LeaveTeamResponse{
		Success: true,
		Error:   api.LeaveTeamResponse_NONE,
	}
	return nil
}
