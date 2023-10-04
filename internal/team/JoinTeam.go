package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
)

type JoinTeamCommand struct {
	service *Service
	In      *api.JoinTeamRequest
	Out     *api.JoinTeamResponse
}

func NewJoinTeamCommand(service *Service, in *api.JoinTeamRequest) *JoinTeamCommand {
	return &JoinTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *JoinTeamCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In.Team)
	if err != nil || c.In.UserId == 0 {
		c.Out = &api.JoinTeamResponse{
			Success: false,
			Error:   api.JoinTeamResponse_INVALID,
		}
		return nil
	}
	result, writeErr := c.service.db.UpdateOne(ctx, append(filter, bson.E{
		Key: "$expr", Value: bson.D{
			{Key: "$lt", Value: bson.A{
				bson.D{
					{Key: "$size", Value: "$membersWithoutOwner"},
				},
				c.service.maxMembers - 1,
			}},
		}},
	), bson.D{
		{Key: "$addToSet", Value: bson.D{
			{Key: "membersWithoutOwner", Value: c.In.UserId},
		}},
	})
	if writeErr != nil {
		return writeErr
	}
	if result.MatchedCount == 0 {
		c.Out = &api.JoinTeamResponse{
			Success: false,
			Error:   api.JoinTeamResponse_NOT_FOUND_OR_TEAM_FULL,
		}
		return nil
	}
	if result.ModifiedCount == 0 {
		c.Out = &api.JoinTeamResponse{
			Success: false,
			Error:   api.JoinTeamResponse_ALREADY_MEMBER,
		}
		return nil
	}
	c.Out = &api.JoinTeamResponse{
		Success: true,
		Error:   api.JoinTeamResponse_NONE,
	}
	return nil
}
