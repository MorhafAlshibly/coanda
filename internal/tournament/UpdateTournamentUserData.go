package tournament

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateTournamentUserDataCommand struct {
	service *Service
	In      *api.UpdateTournamentUserDataRequest
	Out     *api.TournamentUserResponse
}

func NewUpdateTournamentUserDataCommand(service *Service, in *api.UpdateTournamentUserDataRequest) *UpdateTournamentUserDataCommand {
	return &UpdateTournamentUserDataCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateTournamentUserDataCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In.Pagination)
	if err != nil {
		c.Out = &api.TournamentUserResponse{
			Success: false,
			Error:   api.TournamentUserResponse_INVALID,
		}
		return nil
	}
	_, err = c.service.db.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "data", Value: c.In.Data},
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
