package tournament

import (
	"context"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateTournamentUserCommand struct {
	service *Service
	In      *api.CreateTournamentUserRequest
	Out     *api.CreateTournamentUserResponse
}

func NewCreateTournamentUserCommand(service *Service, in *api.CreateTournamentUserRequest) *CreateTournamentUserCommand {
	return &CreateTournamentUserCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateTournamentUserCommand) Execute(ctx context.Context) error {
	// Check if tournament name is large enough
	if len(c.In.Tournament) < int(c.service.minTournamentNameLength) {
		c.Out = &api.CreateTournamentUserResponse{
			Success: false,
			Error:   api.CreateTournamentUserResponse_TOURNAMENT_NAME_TOO_SHORT,
		}
		return nil
	}
	// Check if tournament name is small enough
	if len(c.In.Tournament) > int(c.service.maxTournamentNameLength) {
		c.Out = &api.CreateTournamentUserResponse{
			Success: false,
			Error:   api.CreateTournamentUserResponse_TOURNAMENT_NAME_TOO_LONG,
		}
		return nil
	}
	// Check if user id is valid
	if c.In.UserId == 0 {
		c.Out = &api.CreateTournamentUserResponse{
			Success: false,
			Error:   api.CreateTournamentUserResponse_USER_ID_REQUIRED,
		}
		return nil
	}
	// Insert the tournament into the database
	id, writeErr := c.service.database.InsertOne(ctx, bson.D{
		{Key: "tournament", Value: c.In.Tournament},
		{Key: "interval", Value: c.In.Interval},
		{Key: "userId", Value: c.In.UserId},
		{Key: "score", Value: c.In.Score},
		{Key: "tournamentStartDate", Value: c.service.getTournamentStartDate(time.Now(), c.In.Interval)},
		{Key: "data", Value: c.In.Data},
	})
	if writeErr != nil {
		if mongo.IsDuplicateKeyError(writeErr) {
			c.Out = &api.CreateTournamentUserResponse{
				Success: false,
				Error:   api.CreateTournamentUserResponse_ALREADY_EXISTS,
			}
			return nil
		}
		return writeErr
	}
	c.Out = &api.CreateTournamentUserResponse{
		Success: true,
		Id:      id.Hex(),
		Error:   api.CreateTournamentUserResponse_NONE,
	}
	return nil
}
