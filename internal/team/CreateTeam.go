package team

import (
	"context"
	"strings"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateTeamCommand struct {
	service *Service
	In      *api.CreateTeamRequest
	Out     *api.CreateTeamResponse
}

func NewCreateTeamCommand(service *Service, in *api.CreateTeamRequest) *CreateTeamCommand {
	return &CreateTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateTeamCommand) Execute(ctx context.Context) error {
	// Check if team name is large enough
	if len(c.In.Name) < c.service.minTeamNameLength {
		c.Out = &api.CreateTeamResponse{
			Success: false,
			Team:    nil,
			Error:   api.CreateTeamResponse_NAME_TOO_SHORT,
		}
		return nil
	}
	// Remove duplicates from members
	c.In.MembersWithoutOwner = pkg.RemoveDuplicate(c.In.MembersWithoutOwner)
	if len(c.In.MembersWithoutOwner)+1 > c.service.maxMembers {
		c.Out = &api.CreateTeamResponse{
			Success: false,
			Team:    nil,
			Error:   api.CreateTeamResponse_TOO_MANY_MEMBERS,
		}
		return nil
	}
	if c.In.Score == nil {
		c.In.Score = new(int64)
		*c.In.Score = 0
	}
	// Insert the team into the database
	id, err := c.service.db.InsertOne(ctx, bson.D{
		{Key: "name", Value: c.In.Name},
		{Key: "owner", Value: c.In.Owner},
		{Key: "membersWithoutOwner", Value: c.In.MembersWithoutOwner},
		{Key: "score", Value: c.In.Score},
		{Key: "data", Value: c.In.Data},
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			errEnum := api.CreateTeamResponse_NONE
			if strings.Contains(err.Error(), "owner") {
				errEnum = api.CreateTeamResponse_OWNER_TAKEN
			} else {
				errEnum = api.CreateTeamResponse_NAME_TAKEN
			}
			c.Out = &api.CreateTeamResponse{
				Success: false,
				Team:    nil,
				Error:   errEnum,
			}
			return nil
		}
		return err
	}
	c.Out = &api.CreateTeamResponse{
		Success: true,
		Team: &api.Team{
			Id:                  id,
			Name:                c.In.Name,
			Owner:               c.In.Owner,
			MembersWithoutOwner: c.In.MembersWithoutOwner,
			Score:               *c.In.Score,
			Data:                c.In.Data,
		},
		Error: api.CreateTeamResponse_NONE,
	}
	return nil
}
