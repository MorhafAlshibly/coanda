package tournament

import (
	"context"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
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
	// Check if tournament name is valid
	if c.In.Tournament != nil {
		if len(*c.In.Tournament) < int(c.service.minTournamentNameLength) {
			c.Out = &api.GetTournamentUsersResponse{
				Success: false,
				Error:   api.GetTournamentUsersResponse_TOURNAMENT_NAME_TOO_SHORT,
			}
			return nil
		}
		if len(*c.In.Tournament) > int(c.service.maxTournamentNameLength) {
			c.Out = &api.GetTournamentUsersResponse{
				Success: false,
				Error:   api.GetTournamentUsersResponse_TOURNAMENT_NAME_TOO_LONG,
			}
			return nil
		}
	}
	limit, offset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	result, err := c.service.Database.GetTournaments(ctx, model.GetTournamentsParams{
		Name:                conversion.StringToSqlNullString(c.In.Tournament),
		TournamentInterval:  model.TournamentTournamentInterval(c.In.Interval),
		UserID:              conversion.Uint64ToSqlNullInt64(c.In.UserId),
		TournamentStartedAt: c.service.GetTournamentStartDate(time.Now(), c.In.Interval),
		Limit:               limit,
		Offset:              offset,
	})
	if err != nil {
		return err
	}
	tournaments := make([]*api.TournamentUser, len(result))
	for i, tournament := range result {
		tournaments[i], err = unmarshalTournamentUser(&tournament)
		if err != nil {
			return err
		}
	}
	c.Out = &api.GetTournamentUsersResponse{
		Success:         true,
		TournamentUsers: tournaments,
		Error:           api.GetTournamentUsersResponse_NONE,
	}
	return nil
}
