package record

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/record/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type UpdateRecordCommand struct {
	service *Service
	In      *api.UpdateRecordRequest
	Out     *api.UpdateRecordResponse
}

func NewUpdateRecordCommand(service *Service, in *api.UpdateRecordRequest) *UpdateRecordCommand {
	return &UpdateRecordCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateRecordCommand) Execute(ctx context.Context) error {
	rErr := c.service.checkForRecordRequestError(c.In.Request)
	if rErr != nil {
		c.Out = &api.UpdateRecordResponse{
			Success: false,
			Error:   conversion.Enum(*rErr, api.UpdateRecordResponse_Error_value, api.UpdateRecordResponse_NOT_FOUND),
		}
		return nil
	}
	if c.In.Record == nil && c.In.Data == nil {
		c.Out = &api.UpdateRecordResponse{
			Success: false,
			Error:   api.UpdateRecordResponse_NO_UPDATE_SPECIFIED,
		}
		return nil
	}
	var data json.RawMessage
	if c.In.Data != nil {
		var err error
		data, err = conversion.ProtobufStructToRawJson(c.In.Data)
		if err != nil {
			return err
		}
	}
	// Update the record in the store
	result, err := c.service.database.UpdateRecord(ctx, model.UpdateRecordParams{
		GetRecordParams: model.GetRecordParams{
			Id:         conversion.Uint64ToSqlNullInt64(c.In.Request.Id),
			NameUserId: convertNameUserIdToNullNameUserId(c.In.Request.NameUserId),
		},
		Record: conversion.Uint64ToSqlNullInt64(c.In.Record),
		Data:   data,
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Check if we didn't find a row
		_, err := c.service.database.GetRecord(ctx, model.GetRecordParams{
			Id:         conversion.Uint64ToSqlNullInt64(c.In.Request.Id),
			NameUserId: convertNameUserIdToNullNameUserId(c.In.Request.NameUserId),
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.Out = &api.UpdateRecordResponse{
					Success: false,
					Error:   api.UpdateRecordResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
	}
	c.Out = &api.UpdateRecordResponse{
		Success: true,
		Error:   api.UpdateRecordResponse_NONE,
	}
	return nil
}
