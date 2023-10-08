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
	if c.In.Max == nil {
		c.In.Max = new(uint64)
		*c.In.Max = c.service.defaultMaxPageLength
	}
	if c.In.Page == nil {
		c.In.Page = new(uint64)
		*c.In.Page = 1
	}
	pipelineWithMatch := pipeline
	if c.In.Tournament != nil {
		if len(*c.In.Tournament) < c.service.minTournamentNameLength {
			c.Out = &api.GetTournamentUsersResponse{
				Success:         false,
				TournamentUsers: nil,
				Error:           api.GetTournamentUsersResponse_TOURNAMENT_NAME_TOO_SHORT,
			}
			return nil
		}
		pipelineWithMatch = mongo.Pipeline{
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "tournament", Value: *c.In.Tournament},
				}},
			},
		}
		for _, stage := range pipeline {
			pipelineWithMatch = append(pipelineWithMatch, stage)
		}
	}
	cursor, err := c.service.db.Aggregate(ctx, pipelineWithMatch)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	tournaments, err := toTournaments(ctx, cursor, *c.In.Page, *c.In.Max)
	if err != nil {
		return err
	}
	c.Out = &api.GetTournamentUsersResponse{
		Success:         true,
		TournamentUsers: tournaments,
	}
	return nil
}
