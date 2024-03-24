package tournament

import (
	"context"
	"errors"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcodes"
	"github.com/MorhafAlshibly/coanda/pkg/tournament"
	"github.com/go-sql-driver/mysql"
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
	// Check if data is provided
	if c.In.Data == nil {
		c.Out = &api.CreateTournamentUserResponse{
			Success: false,
			Error:   api.CreateTournamentUserResponse_DATA_REQUIRED,
		}
		return nil
	}
	raw, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	score := int64(0)
	if c.In.Score != nil {
		score = *c.In.Score
	}
	_, err = c.service.database.CreateTournament(ctx, model.CreateTournamentParams{
		Name:               c.In.Tournament,
		TournamentInterval: model.TournamentTournamentInterval(c.In.Interval.String()),
		UserID:             c.In.UserId,
		Score:              score,
		Data:               raw,
		TournamentStartedAt: tournament.GetStartTime(time.Now().UTC(), c.In.Interval, tournament.WipeTimes{
			DailyTournamentMinute:   c.service.dailyTournamentMinute,
			WeeklyTournamentMinute:  c.service.weeklyTournamentMinute,
			WeeklyTournamentDay:     c.service.weeklyTournamentDay,
			MonthlyTournamentMinute: c.service.monthlyTournamentMinute,
			MonthlyTournamentDay:    c.service.monthlyTournamentDay,
		}),
	})
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) && mysqlError.Number == errorcodes.MySQLErrorCodeDuplicateEntry {
			c.Out = &api.CreateTournamentUserResponse{
				Success: false,
				Error:   api.CreateTournamentUserResponse_ALREADY_EXISTS,
			}
			return nil
		}
		return err
	}
	c.Out = &api.CreateTournamentUserResponse{
		Success: true,
		Error:   api.CreateTournamentUserResponse_NONE,
	}
	return nil
}
