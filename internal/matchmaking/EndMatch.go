package matchmaking

import (
	"context"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type EndMatchCommand struct {
	service *Service
	In      *api.EndMatchRequest
	Out     *api.EndMatchResponse
}

func NewEndMatchCommand(service *Service, in *api.EndMatchRequest) *EndMatchCommand {
	return &EndMatchCommand{
		service: service,
		In:      in,
	}
}

func (c *EndMatchCommand) Execute(ctx context.Context) error {
	mmErr := c.service.checkForMatchRequestError(c.In.Match)
	// Check if error is found
	if mmErr != nil {
		c.Out = &api.EndMatchResponse{
			Success: false,
			Error:   conversion.Enum(*mmErr, api.EndMatchResponse_Error_value, api.EndMatchResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED),
		}
		return nil
	}
	// Check if end time is nil
	if c.In.EndTime == nil {
		c.Out = &api.EndMatchResponse{
			Success: false,
			Error:   api.EndMatchResponse_END_TIME_REQUIRED,
		}
		return nil
	}
	if c.In.EndTime.AsTime().Before(time.Now()) {
		c.Out = &api.EndMatchResponse{
			Success: false,
			Error:   api.EndMatchResponse_INVALID_END_TIME,
		}
		return nil
	}
	params := matchRequestToMatchParams(c.In.Match)
	result, err := c.service.database.EndMatch(ctx, model.EndMatchParams{
		Match:   params,
		EndTime: c.In.EndTime.AsTime(),
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Either match wasnt found, or match already ended, or match hasn't started yet, or match end time is before start time
		match, err := c.service.database.GetMatch(ctx, model.GetMatchParams{
			Match:       params,
			TicketLimit: 1,
			UserLimit:   1,
			ArenaLimit:  1,
		})
		if err != nil {
			return err
		}
		if len(match) == 0 {
			c.Out = &api.EndMatchResponse{
				Success: false,
				Error:   api.EndMatchResponse_NOT_FOUND,
			}
			return nil
		}
		if match[0].StartedAt.Valid {
			if match[0].StartedAt.Time.After(time.Now()) {
				c.Out = &api.EndMatchResponse{
					Success: false,
					Error:   api.EndMatchResponse_HAS_NOT_STARTED,
				}
				return nil
			}
		} else {
			c.Out = &api.EndMatchResponse{
				Success: false,
				Error:   api.EndMatchResponse_HAS_NOT_STARTED,
			}
			return nil
		}
		if match[0].EndedAt.Valid {
			if match[0].EndedAt.Time.Before(time.Now()) {
				c.Out = &api.EndMatchResponse{
					Success: false,
					Error:   api.EndMatchResponse_ALREADY_ENDED,
				}
				return nil
			}
		}
		if match[0].StartedAt.Valid {
			if match[0].StartedAt.Time.After(c.In.EndTime.AsTime()) {
				c.Out = &api.EndMatchResponse{
					Success: false,
					Error:   api.EndMatchResponse_END_TIME_BEFORE_START_TIME,
				}
				return nil
			}
		}
	}
	c.Out = &api.EndMatchResponse{
		Success: true,
		Error:   api.EndMatchResponse_NONE,
	}
	return nil
}
