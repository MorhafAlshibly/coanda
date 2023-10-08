package tournament

import (
	"context"

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
	if len(c.In.Tournament) < c.service.minTournamentNameLength {
		c.Out = &api.CreateTournamentUserResponse{
			Success: false,
			Error:   api.CreateTournamentUserResponse_TOURNAMENT_NAME_TOO_SHORT,
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
	// If score is not provided, set it to 0
	if c.In.Score == nil {
		c.In.Score = new(int64)
		*c.In.Score = 0
	}
	// Insert the tournament into the database
	id, writeErr := c.service.db.InsertOne(ctx, bson.D{
		{Key: "tournament", Value: c.In.Tournament},
		{Key: "interval", Value: c.In.Interval},
		{Key: "userId", Value: c.In.UserId},
		{Key: "score", Value: c.In.Score},
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
