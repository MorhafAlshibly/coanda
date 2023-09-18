package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateTeamCommand struct {
	service *Service
	In      *api.CreateTeamRequest
	Out     *api.Team
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
		return errors.New("Team name too short")
	}
	// Remove duplicates from members
	c.In.MembersWithoutOwner = pkg.RemoveDuplicate(c.In.MembersWithoutOwner)
	if len(c.In.MembersWithoutOwner)+1 > c.service.maxMembers {
		return errors.New("Too many members")
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
			fmt.Println(err.WriteErrors[0].Details)
			return errors.New("Team name already exists")
		}
		return err
	}
	c.Out = &api.Team{
		Id:                  id,
		Name:                c.In.Name,
		Owner:               c.In.Owner,
		MembersWithoutOwner: c.In.MembersWithoutOwner,
		Score:               c.In.Score,
		Data:                c.In.Data,
	}
	return nil
}
