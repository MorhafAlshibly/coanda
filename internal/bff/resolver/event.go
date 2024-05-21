package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/bff/model"
)

// CreateEvent is the resolver for the CreateEvent field.
func (r *mutationResolver) CreateEvent(ctx context.Context, input model.CreateEventRequest) (*model.CreateEventResponse, error) {
	var rounds []*api.CreateEventRound
	for _, round := range input.Rounds {
		rounds = append(rounds, &api.CreateEventRound{
			Name:    round.Name,
			Data:    round.Data,
			Scoring: round.Scoring,
			EndedAt: round.EndedAt,
		})
	}
	resp, err := r.eventClient.CreateEvent(ctx, &api.CreateEventRequest{
		Name:      input.Name,
		Data:      input.Data,
		StartedAt: input.StartedAt,
		Rounds:    rounds,
	})
	if err != nil {
		return nil, err
	}
	return &model.CreateEventResponse{
		Success: resp.Success,
		Error:   model.CreateEventError(resp.Error.String()),
		ID:      resp.Id,
	}, nil
}
