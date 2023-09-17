package team

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api/pb"
	"go.mongodb.org/mongo-driver/bson"
)

type JoinTeamCommand struct {
	service *TeamService
	In      *pb.JoinTeamRequest
	Out     *pb.Team
}

func NewJoinTeamCommand(service *TeamService, in *pb.JoinTeamRequest) *JoinTeamCommand {
	return &JoinTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *JoinTeamCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In.Team)
	if err != nil {
		return err
	}
	result, err := c.service.db.UpdateOne(ctx, append(append(filter, bson.E{
		Key: "$expr", Value: bson.D{
			{Key: "$lt", Value: bson.A{
				bson.D{
					{Key: "$size", Value: "$membersWithoutOwner"},
				},
				c.service.maxMembers,
			}},
		}},
	), bson.E{
		Key: "membersWithoutOwner", Value: bson.D{
			{Key: "$ne", Value: c.In.UserId},
		},
	},
	), bson.D{
		{Key: "$addToSet", Value: bson.D{
			{Key: "membersWithoutOwner", Value: c.In.UserId},
		}},
	})
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("Team is full")
	}
	c.Out, err = c.service.GetTeam(ctx, c.In.Team)
	if err != nil {
		return err
	}
	return nil
}
