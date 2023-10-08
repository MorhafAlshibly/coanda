package tournament

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetTournamentUserCommand struct {
	service *Service
	In      *api.GetTournamentUserRequest
	Out     *api.TournamentUserResponse
}

func NewGetTournamentUserCommand(service *Service, in *api.GetTournamentUserRequest) *GetTournamentUserCommand {
	return &GetTournamentUserCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTournamentUserCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In)
	if err != nil {
		c.Out = &api.TournamentUserResponse{
			Success:        false,
			TournamentUser: nil,
			Error:          api.TournamentUserResponse_INVALID,
		}
		return nil
	}
	if c.In.TournamentIntervalUserId != nil {
		if len(c.In.TournamentIntervalUserId.Tournament) < c.service.minTournamentNameLength {
			c.Out = &api.TournamentUserResponse{
				Success:        false,
				TournamentUser: nil,
				Error:          api.TournamentUserResponse_TOURNAMENT_NAME_TOO_SHORT,
			}
			return nil
		}
	}
	// Get the item from the store
	pipelineWithMatch := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: filter},
		},
	}
	for _, stage := range pipeline {
		pipelineWithMatch = append(pipelineWithMatch, stage)
	}
	cursor, err := c.service.db.Aggregate(ctx, pipelineWithMatch)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	cursor.Next(ctx)
	tournament, err := toTournament(cursor)
	if err != nil {
		if err.Error() == "EOF" {
			c.Out = &api.TournamentUserResponse{
				Success:        false,
				TournamentUser: nil,
				Error:          api.TournamentUserResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	c.Out = &api.TournamentUserResponse{
		Success:        true,
		TournamentUser: tournament,
		Error:          api.TournamentUserResponse_NONE,
	}
	return nil
}
