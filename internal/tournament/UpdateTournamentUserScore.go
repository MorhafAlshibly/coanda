package tournament

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateTournamentUserScoreCommand struct {
	service *Service
	In      *api.UpdateTournamentUserScoreRequest
	Out     *api.TournamentUserResponse
}

func NewUpdateTournamentUserScoreCommand(service *Service, in *api.UpdateTournamentUserScoreRequest) *UpdateTournamentUserScoreCommand {
	return &UpdateTournamentUserScoreCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateTournamentUserScoreCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In.Pagination)

	if err != nil {
		c.Out = &api.TournamentUserResponse{
			Success: false,
			Error:   api.TournamentUserResponse_INVALID,
		}
		return nil
	}
	_, err = c.service.db.UpdateOne(ctx, filter, bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "score", Value: c.In.ScoreOffset},
		}},
	})
	if err != nil {
		if err.Error() == "EOF" {
			c.Out = &api.TournamentUserResponse{
				Success: false,
				Error:   api.TournamentUserResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	c.Out = &api.TournamentUserResponse{
		Success: true,
		Error:   api.TournamentUserResponse_NONE,
	}
	return nil
}
