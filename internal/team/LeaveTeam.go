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
	Out     *api.Team
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
		return err
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
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("User is not a member of the team")
	}
	c.Out, err = c.service.GetTeam(ctx, c.In.Team)
	if err != nil {
		return err
	}
	return nil
}
