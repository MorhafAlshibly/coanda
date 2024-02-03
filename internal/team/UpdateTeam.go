package team

import (
	"context"
	"encoding/json"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/validation"
)

type UpdateTeamCommand struct {
	service *Service
	In      *api.UpdateTeamRequest
	Out     *api.UpdateTeamResponse
}

func NewUpdateTeamCommand(service *Service, in *api.UpdateTeamRequest) *UpdateTeamCommand {
	return &UpdateTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateTeamCommand) Execute(ctx context.Context) error {
	tErr := c.service.CheckForTeamRequestError(c.In.Team)
	// Check if error is found
	if tErr != nil {
		c.Out = &api.UpdateTeamResponse{
			Success: false,
			Error:   conversion.Enum(*tErr, api.UpdateTeamResponse_Error_value, api.UpdateTeamResponse_NOT_FOUND),
		}
		return nil
	}
	// Check if no update is specified
	if c.In.Score == nil && c.In.Data == nil {
		c.Out = &api.UpdateTeamResponse{
			Success: false,
			Error:   api.UpdateTeamResponse_NO_UPDATE_SPECIFIED,
		}
		return nil
	}
	// Check if score is specified without whether to increment it
	if c.In.Score != nil && c.In.IncrementScore == nil {
		c.Out = &api.UpdateTeamResponse{
			Success: false,
			Error:   api.UpdateTeamResponse_INCREMENT_SCORE_NOT_SPECIFIED,
		}
		return nil
	}
	// Prepare data
	var data json.RawMessage
	var err error
	dataExists := int64(0)
	if c.In.Data != nil {
		data, err = conversion.ProtobufStructToRawJson(c.In.Data)
		if err != nil {
			return err
		}
		dataExists = 1
	}
	// Prepare score
	incrementScore := int64(0)
	if c.In.IncrementScore != nil {
		if *c.In.IncrementScore {
			incrementScore = 1
		}
	}
	result, err := c.service.database.UpdateTeam(ctx, model.UpdateTeamParams{
		Name:           validation.ValidateAnSqlNullString(c.In.Team.Name),
		Owner:          validation.ValidateAUint64ToSqlNullInt64(c.In.Team.Owner),
		Member:         validation.ValidateAUint64ToSqlNullInt64(c.In.Team.Member),
		DataExists:     dataExists,
		Data:           data,
		Score:          validation.ValidateAnSqlNullInt64(c.In.Score),
		IncrementScore: incrementScore,
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		c.Out = &api.UpdateTeamResponse{
			Success: false,
			Error:   api.UpdateTeamResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.UpdateTeamResponse{
		Success: true,
		Error:   api.UpdateTeamResponse_NONE,
	}
	return nil

}
