package record

import (
	"context"
	"database/sql"
	"encoding/json"

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
			Error:   conversion.Enum(*rErr, api.UpdateRecordResponse_Error_value, api.UpdateRecordResponse_ID_OR_NAME_USER_ID_REQUIRED),
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
	if c.In.Request.NameUserId == nil {
		c.In.Request.NameUserId = &api.NameUserId{}
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	// Update the record in the store
	result, err := qtx.UpdateRecord(ctx, model.UpdateRecordParams{
		GetRecordParams: model.GetRecordParams{
			Id:     conversion.Uint64ToSqlNullInt64(c.In.Request.Id),
			Name:   conversion.StringToSqlNullString(&c.In.Request.NameUserId.Name),
			UserID: conversion.Uint64ToSqlNullInt64(&c.In.Request.NameUserId.UserId),
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
		_, err := qtx.GetRecord(ctx, model.GetRecordParams{
			Id:     conversion.Uint64ToSqlNullInt64(c.In.Request.Id),
			Name:   conversion.StringToSqlNullString(&c.In.Request.NameUserId.Name),
			UserID: conversion.Uint64ToSqlNullInt64(&c.In.Request.NameUserId.UserId),
		})
		if err != nil {
			if err == sql.ErrNoRows {
				c.Out = &api.UpdateRecordResponse{
					Success: false,
					Error:   api.UpdateRecordResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.UpdateRecordResponse{
		Success: true,
		Error:   api.UpdateRecordResponse_NONE,
	}
	return nil
}
