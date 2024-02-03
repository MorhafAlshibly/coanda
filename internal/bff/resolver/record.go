package resolver


			// This file will be automatically regenerated based on the schema, any resolver implementations
			// will be copied through when generating and any unknown code will be moved to the end.
			// Code generated by github.com/99designs/gqlgen version v0.17.43

import (
"context"
"fmt"
"io"
"strconv"
"time"
"sync"
"errors"
"bytes"
gqlparser "github.com/vektah/gqlparser/v2"
"github.com/vektah/gqlparser/v2/ast"
"github.com/99designs/gqlgen/graphql"
"github.com/99designs/gqlgen/graphql/introspection"
"github.com/MorhafAlshibly/coanda/api"
"github.com/MorhafAlshibly/coanda/internal/bff/model")


















// CreateRecord is the resolver for the CreateRecord field.
	func (r *mutationResolver) CreateRecord(ctx context.Context, input model.CreateRecordRequest) ( *model.CreateRecordResponse,  error){
		data, err := .MapStringAnyToMapStringString(input.Data)
	if err != nil {
		return nil, err
	}
	resp, err := r.recordClient.CreateRecord(ctx, &api.CreateRecordRequest{
		Name:   input.Name,
		UserId: input.UserID,
		Record: input.Record,
		Data:   data,
	})
	if err != nil {
		return nil, err
	}
	return &model.CreateRecordResponse{
		Success: resp.Success,
		ID:      resp.Id,
		Error:   model.CreateRecordError(resp.Error.String()),
	}, nil
	}

// UpdateRecord is the resolver for the UpdateRecord field.
	func (r *mutationResolver) UpdateRecord(ctx context.Context, input model.UpdateRecordRequest) ( *model.UpdateRecordResponse,  error){
		data, err := .MapStringAnyToMapStringString(input.Data)
	if err != nil {
		return nil, err
	}
	resp, err := r.recordClient.UpdateRecord(ctx, &api.UpdateRecordRequest{
		Request: &api.GetRecordRequest{
			Id: input.Request.ID,
			NameUserId: &api.NameUserId{
				Name:   input.Request.NameUserID.Name,
				UserId: input.Request.NameUserID.UserID,
			},
		},
		Record: input.Record,
		Data:   data,
	})
	if err != nil {
		return nil, err
	}
	return &model.UpdateRecordResponse{
		Success: resp.Success,
		Error:   model.UpdateRecordError(resp.Error.String()),
	}, nil
	}

// DeleteRecord is the resolver for the DeleteRecord field.
	func (r *mutationResolver) DeleteRecord(ctx context.Context, input model.GetRecordRequest) ( *model.DeleteRecordResponse,  error){
		resp, err := r.recordClient.DeleteRecord(ctx, &api.GetRecordRequest{
		Id: input.ID,
		NameUserId: &api.NameUserId{
			Name:   input.NameUserID.Name,
			UserId: input.NameUserID.UserID,
		},
	})
	if err != nil {
		return nil, err
	}
	return &model.DeleteRecordResponse{
		Success: resp.Success,
		Error:   model.DeleteRecordError(resp.Error.String()),
	}, nil
	}

// GetRecord is the resolver for the GetRecord field.
	func (r *queryResolver) GetRecord(ctx context.Context, input model.GetRecordRequest) ( *model.GetRecordResponse,  error){
		resp, err := r.recordClient.GetRecord(ctx, &api.GetRecordRequest{
		Id: input.ID,
		NameUserId: &api.NameUserId{
			Name:   input.NameUserID.Name,
			UserId: input.NameUserID.UserID,
		},
	})
	if err != nil {
		return nil, err
	}
	var record *model.Record
	if resp.Record != nil {
		dataMap, err := .MapStringStringToMapStringAny(resp.Record.Data)
		if err != nil {
			return nil, err
		}
		record = &model.Record{
			ID:        resp.Record.Id,
			Name:      resp.Record.Name,
			UserID:    resp.Record.UserId,
			Record:    resp.Record.Record,
			Rank:      resp.Record.Rank,
			Data:      dataMap,
			CreatedAt: resp.Record.CreatedAt,
			UpdatedAt: resp.Record.UpdatedAt,
		}
	}
	return &model.GetRecordResponse{
		Success: resp.Success,
		Record:  record,
		Error:   model.GetRecordError(resp.Error.String()),
	}, nil
	}

// GetRecords is the resolver for the GetRecords field.
	func (r *queryResolver) GetRecords(ctx context.Context, input model.GetRecordsRequest) ( *model.GetRecordsResponse,  error){
		resp, err := r.recordClient.GetRecords(ctx, &api.GetRecordsRequest{
		Max:  input.Max,
		Page: input.Page,
		Name: input.Name,
	})
	if err != nil {
		return nil, err
	}
	records := make([]*model.Record, len(resp.Records))
	for i, record := range resp.Records {
		dataMap, err := .MapStringStringToMapStringAny(record.Data)
		if err != nil {
			return nil, err
		}
		records[i] = &model.Record{
			ID:        record.Id,
			Name:      record.Name,
			UserID:    record.UserId,
			Record:    record.Record,
			Rank:      record.Rank,
			Data:      dataMap,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
		}
	}
	return &model.GetRecordsResponse{
		Success: resp.Success,
		Records: records,
		Error:   model.GetRecordsError(resp.Error.String()),
	}, nil
	}








