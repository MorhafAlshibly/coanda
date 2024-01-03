package tournament

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetTournamentUsersCommand struct {
	service *Service
	In      *api.GetTournamentUsersRequest
	Out     *api.GetTournamentUsersResponse
}

func NewGetTournamentUsersCommand(service *Service, in *api.GetTournamentUsersRequest) *GetTournamentUsersCommand {
	return &GetTournamentUsersCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTournamentUsersCommand) Execute(ctx context.Context) error {
	max := uint8(c.In.Max)
	if max == 0 {
		max = c.service.defaultMaxPageLength
	}
	if max > c.service.maxMaxPageLength {
		max = c.service.maxMaxPageLength
	}
	if c.In.Page == 0 {
		c.In.Page = 1
	}
	pipelineWithMatch := pipeline
	if c.In.Tournament != "" {
		if len(c.In.Tournament) < int(c.service.minTournamentNameLength) {
			c.Out = &api.GetTournamentUsersResponse{
				Success:         false,
				TournamentUsers: nil,
				Error:           api.GetTournamentUsersResponse_TOURNAMENT_NAME_TOO_SHORT,
			}
			return nil
		}
		if len(c.In.Tournament) > int(c.service.maxTournamentNameLength) {
			c.Out = &api.GetTournamentUsersResponse{
				Success:         false,
				TournamentUsers: nil,
				Error:           api.GetTournamentUsersResponse_TOURNAMENT_NAME_TOO_LONG,
			}
			return nil
		}
		pipelineWithMatch = mongo.Pipeline{
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "tournament", Value: c.In.Tournament},
				}},
			},
		}
		for _, stage := range pipeline {
			pipelineWithMatch = append(pipelineWithMatch, stage)
		}
	}
	cursor, err := c.service.database.Aggregate(ctx, pipelineWithMatch)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	tournaments, err := toTournamentUsers(ctx, cursor, c.In.Page, max)
	if err != nil {
		return err
	}
	c.Out = &api.GetTournamentUsersResponse{
		Success:         true,
		TournamentUsers: tournaments,
	}
	return nil
}
