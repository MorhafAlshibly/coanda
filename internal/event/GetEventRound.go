package event

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetEventRoundCommand struct {
	service *Service
	In      *api.GetEventRoundRequest
	Out     *api.GetEventRoundResponse
}

func NewGetEventRoundCommand(service *Service, in *api.GetEventRoundRequest) *GetEventRoundCommand {
	return &GetEventRoundCommand{
		service: service,
		In:      in,
	}
}

func (c *GetEventRoundCommand) Execute(ctx context.Context) error {
	erErr := c.service.checkForEventRoundRequestError(c.In.Round)
	if erErr != nil {
		c.Out = &api.GetEventRoundResponse{
			Success: false,
			Error:   conversion.Enum(*erErr, api.GetEventRoundResponse_Error_value, api.GetEventRoundResponse_ID_OR_NAME_REQUIRED),
		}
		return nil
	}
	// Create blank event request if not provided
	if c.In.Round.Event == nil {
		c.In.Round.Event = &api.EventRequest{}
	}
	limit, offset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	eventRound, err := qtx.GetEventRound(ctx, model.GetEventRoundParams{
		Event: model.GetEventParams{
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Round.Event.Id),
			Name: conversion.StringToSqlNullString(c.In.Round.Event.Name),
		},
		ID:   conversion.Uint64ToSqlNullInt64(c.In.Round.Id),
		Name: conversion.StringToSqlNullString(c.In.Round.RoundName),
	}, nil)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Out = &api.GetEventRoundResponse{
				Success: false,
				Error:   api.GetEventRoundResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	apiEventRound, err := unmarshalEventRound(eventRound)
	if err != nil {
		return err
	}
	leaderboard, err := qtx.GetEventRoundLeaderboard(ctx, model.GetEventRoundLeaderboardParams{
		EventRound: model.GetEventRoundParams{
			Event: model.GetEventParams{
				ID:   conversion.Uint64ToSqlNullInt64(c.In.Round.Event.Id),
				Name: conversion.StringToSqlNullString(c.In.Round.Event.Name),
			},
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Round.Id),
			Name: conversion.StringToSqlNullString(c.In.Round.RoundName),
		},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows are found, return an empty leaderboard
			leaderboard = []model.EventRoundLeaderboard{}
		} else {
			return err
		}
	}
	apiLeaderboard, err := unmarshalEventRoundLeaderboard(leaderboard)
	if err != nil {
		return err
	}
	c.Out = &api.GetEventRoundResponse{
		Success: true,
		Round:   apiEventRound,
		Results: apiLeaderboard,
		Error:   api.GetEventRoundResponse_NONE,
	}
	return nil
}
