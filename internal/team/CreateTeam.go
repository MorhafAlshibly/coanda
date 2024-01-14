package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/database/sqlc"
	"github.com/MorhafAlshibly/coanda/pkg/validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	// Check if team name is correct length
	if len(c.In.Name) < int(c.service.minTeamNameLength) {
		c.Out = &api.CreateTeamResponse{
			Success: false,
			Error:   api.CreateTeamResponse_NAME_TOO_SHORT,
		}
		return nil
	}
	if len(c.In.Name) > int(c.service.maxTeamNameLength) {
		c.Out = &api.CreateTeamResponse{
			Success: false,
			Error:   api.CreateTeamResponse_NAME_TOO_LONG,
		}
		return nil
	}
	if c.In.Owner == 0 {
		c.Out = &api.CreateTeamResponse{
			Success: false,
			Error:   api.CreateTeamResponse_OWNER_REQUIRED,
		}
		return nil
	}
	// If score is not provided, set it to 0
	if c.In.Score == nil {
		c.In.Score = new(int64)
		*c.In.Score = 0
	}
	// Remove duplicates from members
	c.In.MembersWithoutOwner = conversion.ArrayToSet(c.In.MembersWithoutOwner)
	// Check if owner is in members
	if validation.CheckArrayContains(c.In.MembersWithoutOwner, c.In.Owner) {
		c.In.MembersWithoutOwner = conversion.ArrayRemove(c.In.MembersWithoutOwner, c.In.Owner)
	}
	// Check if team is too big
	if len(c.In.MembersWithoutOwner)+1 > int(c.service.maxMembers) {
		c.Out = &api.CreateTeamResponse{
			Success: false,
			Error:   api.CreateTeamResponse_TOO_MANY_MEMBERS,
		}
		return nil
	}
	// Insert the team into the database
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	tx, err := c.service.sql.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	_, err = qtx.CreateTeam(ctx, sqlc.CreateTeamParams{
		Name:  c.In.Name,
		Owner: c.In.Owner,
		Score: *c.In.Score,
		Data:  data,
	})
	if err != nil {
		return err
	}
	// Create team members for the team
	if err != nil {
		if mongo.IsDuplicateKeyError(writeErr) {
			errEnum := api.CreateTeamResponse_NONE
			findName, err := c.service.database.Find(ctx, bson.D{
				{Key: "name", Value: c.In.Name}}, &options.FindOptions{
				Projection: bson.D{
					{Key: "_id", Value: 1},
				},
			})
			if err != nil {
				return err
			}
			if findName.Next(ctx) {
				errEnum = api.CreateTeamResponse_NAME_TAKEN
			} else {
				errEnum = api.CreateTeamResponse_OWNER_TAKEN
			}
			c.Out = &api.CreateTeamResponse{
				Success: false,
				Error:   errEnum,
			}
			return nil
		}
		return writeErr
	}
	c.Out = &api.CreateTeamResponse{
		Success: true,
		Id:      id.Hex(),
		Error:   api.CreateTeamResponse_NONE,
	}
	return nil
}
