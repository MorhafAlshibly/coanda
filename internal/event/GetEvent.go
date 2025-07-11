package event

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetEventCommand struct {
	service *Service
	In      *api.GetEventRequest
	Out     *api.GetEventResponse
}

func NewGetEventCommand(service *Service, in *api.GetEventRequest) *GetEventCommand {
	return &GetEventCommand{
		service: service,
		In:      in,
	}
}

func (c *GetEventCommand) Execute(ctx context.Context) error {
	eErr := c.service.checkForEventRequestError(c.In.Event)
	if eErr != nil {
		c.Out = &api.GetEventResponse{
			Success: false,
			Error:   conversion.Enum(*eErr, api.GetEventResponse_Error_value, api.GetEventResponse_ID_OR_NAME_REQUIRED),
		}
		return nil
	}
	limit, offset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	event, err := qtx.GetEventWithRound(ctx, model.GetEventParams{
		ID:   conversion.Uint64ToSqlNullInt64(c.In.Event.Id),
		Name: conversion.StringToSqlNullString(c.In.Event.Name),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			c.Out = &api.GetEventResponse{
				Success: false,
				Error:   api.GetEventResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	apiEvent, err := unmarshalEventWithRound(event)
	if err != nil {
		return err
	}
	leaderboard, err := qtx.GetEventLeaderboard(ctx, model.GetEventLeaderboardParams{
		Event: model.GetEventParams{
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Event.Id),
			Name: conversion.StringToSqlNullString(c.In.Event.Name),
		},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			leaderboard = []model.EventLeaderboard{}
		} else {
			return err
		}
	}
	apiLeaderboard, err := unmarshalEventLeaderboard(leaderboard)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.GetEventResponse{
		Success:     true,
		Event:       apiEvent,
		Leaderboard: apiLeaderboard,
	}
	return nil
}
