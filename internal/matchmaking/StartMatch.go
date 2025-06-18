package matchmaking

import (
	"context"
	"errors"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type StartMatchCommand struct {
	service *Service
	In      *api.StartMatchRequest
	Out     *api.StartMatchResponse
}

func NewStartMatchCommand(service *Service, in *api.StartMatchRequest) *StartMatchCommand {
	return &StartMatchCommand{
		service: service,
		In:      in,
	}
}

func (c *StartMatchCommand) Execute(ctx context.Context) error {
	mmErr := c.service.checkForMatchRequestError(c.In.Match)
	// Check if error is found
	if mmErr != nil {
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   conversion.Enum(*mmErr, api.StartMatchResponse_Error_value, api.StartMatchResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED),
		}
		return nil
	}
	// Check if start time is nil
	if c.In.StartTime == nil {
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   api.StartMatchResponse_START_TIME_REQUIRED,
		}
		return nil
	}
	if c.In.StartTime.AsTime().Before(time.Now()) {
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   api.StartMatchResponse_INVALID_START_TIME,
		}
		return nil
	}
	if c.In.StartTime.AsTime().Before(time.Now().Add(c.service.startTimeBuffer)) {
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   api.StartMatchResponse_START_TIME_TOO_SOON,
		}
		return nil
	}
	lockTime := c.In.StartTime.AsTime().Add(-c.service.lockedAtBuffer)
	// Check if lock time is before now
	if lockTime.Before(time.Now()) {
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   api.StartMatchResponse_START_TIME_TOO_SOON,
		}
		return nil
	}
	// Get the match
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	params := matchRequestToMatchParams(c.In.Match)
	// We dont need to lock the match here, as we are only checking if it exists and if it has enough players to start.
	// This data is not designed to change once it is valid.
	// The StartMatch request will only start the match if started_at is not set.
	match, err := qtx.GetMatch(ctx, model.GetMatchParams{
		Match:       params,
		TicketLimit: 1,
		UserLimit:   1,
		ArenaLimit:  1,
	})
	if err != nil {
		return err
	}
	if len(match) == 0 {
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   api.StartMatchResponse_NOT_FOUND,
		}
		return nil
	}
	// Check if match has enough players
	if match[0].UserCount < uint64(match[0].ArenaMinPlayers) {
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   api.StartMatchResponse_NOT_ENOUGH_PLAYERS_TO_START,
		}
		return nil
	}
	// Check if match already has start time
	if match[0].StartedAt.Valid {
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   api.StartMatchResponse_ALREADY_HAS_START_TIME,
		}
		return nil
	}
	// Check if match has a private server
	if !match[0].PrivateServerID.Valid {
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   api.StartMatchResponse_PRIVATE_SERVER_NOT_SET,
		}
		return nil
	}
	// Start match
	result, err := qtx.StartMatch(ctx, model.StartMatchParams{
		Match:     params,
		StartTime: c.In.StartTime.AsTime(),
		LockTime:  lockTime,
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("race condition occured however it has been safely handled. retry the request")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.StartMatchResponse{
		Success: true,
		Error:   api.StartMatchResponse_NONE,
	}
	return nil
}
