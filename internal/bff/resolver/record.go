package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.53

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/bff/model"
)

// CreateRecord is the resolver for the CreateRecord field.
func (r *mutationResolver) CreateRecord(ctx context.Context, input model.CreateRecordRequest) (*model.CreateRecordResponse, error) {
	resp, err := r.recordClient.CreateRecord(ctx, &api.CreateRecordRequest{
		Name:   input.Name,
		UserId: input.UserID,
		Record: input.Record,
		Data:   input.Data,
	})
	if err != nil {
		return nil, err
	}
	return &model.CreateRecordResponse{
		Success: resp.Success,
		Error:   model.CreateRecordError(resp.Error.String()),
		ID:      resp.Id,
	}, nil
}

// UpdateRecord is the resolver for the UpdateRecord field.
func (r *mutationResolver) UpdateRecord(ctx context.Context, input model.UpdateRecordRequest) (*model.UpdateRecordResponse, error) {
	var recordRequest *api.RecordRequest
	if input.Request != nil {
		recordRequest = &api.RecordRequest{
			Id: input.Request.ID,
		}
		if input.Request.NameUserID != nil {
			recordRequest = &api.RecordRequest{
				Id: input.Request.ID,
				NameUserId: &api.NameUserId{
					Name:   input.Request.NameUserID.Name,
					UserId: input.Request.NameUserID.UserID,
				},
			}
		}
	}
	resp, err := r.recordClient.UpdateRecord(ctx, &api.UpdateRecordRequest{
		Request: recordRequest,
		Record:  input.Record,
		Data:    input.Data,
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
func (r *mutationResolver) DeleteRecord(ctx context.Context, input model.RecordRequest) (*model.DeleteRecordResponse, error) {
	var nameUserId *api.NameUserId
	if input.NameUserID != nil {
		nameUserId = &api.NameUserId{
			Name:   input.NameUserID.Name,
			UserId: input.NameUserID.UserID,
		}
	}
	resp, err := r.recordClient.DeleteRecord(ctx, &api.RecordRequest{
		Id:         input.ID,
		NameUserId: nameUserId,
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
func (r *queryResolver) GetRecord(ctx context.Context, input model.RecordRequest) (*model.GetRecordResponse, error) {
	var nameUserId *api.NameUserId
	if input.NameUserID != nil {
		nameUserId = &api.NameUserId{
			Name:   input.NameUserID.Name,
			UserId: input.NameUserID.UserID,
		}
	}
	resp, err := r.recordClient.GetRecord(ctx, &api.RecordRequest{
		Id:         input.ID,
		NameUserId: nameUserId,
	})
	if err != nil {
		return nil, err
	}
	var record *model.Record
	if resp.Record != nil {
		record = &model.Record{
			ID:        resp.Record.Id,
			Name:      resp.Record.Name,
			UserID:    resp.Record.UserId,
			Record:    resp.Record.Record,
			Ranking:   resp.Record.Ranking,
			Data:      resp.Record.Data,
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
func (r *queryResolver) GetRecords(ctx context.Context, input model.GetRecordsRequest) (*model.GetRecordsResponse, error) {
	if input.Pagination == nil {
		input.Pagination = &model.Pagination{}
	}
	resp, err := r.recordClient.GetRecords(ctx, &api.GetRecordsRequest{
		Name:   input.Name,
		UserId: input.UserID,
		Pagination: &api.Pagination{
			Max:  input.Pagination.Max,
			Page: input.Pagination.Page,
		},
	})
	if err != nil {
		return nil, err
	}
	records := make([]*model.Record, len(resp.Records))
	for i, record := range resp.Records {
		records[i] = &model.Record{
			ID:        record.Id,
			Name:      record.Name,
			UserID:    record.UserId,
			Record:    record.Record,
			Ranking:   record.Ranking,
			Data:      record.Data,
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
