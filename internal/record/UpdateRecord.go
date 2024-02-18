package record

import (
	"context"

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
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	// Update the record in the store
	_, err = c.service.database.UpdateRecord(ctx, model.UpdateRecordParams{
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
	c.Out = &api.UpdateRecordResponse{
		Success: true,
		Error:   api.UpdateRecordResponse_NONE,
	}
	return nil
}
