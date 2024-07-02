package event

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/event/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
)

type Service struct {
	api.UnimplementedEventServiceServer
	sql                  *sql.DB
	database             *model.Queries
	cache                cache.Cacher
	metrics              metrics.Metrics
	minEventNameLength   uint8
	maxEventNameLength   uint8
	minRoundNameLength   uint8
	maxRoundNameLength   uint8
	maxNumberOfRounds    uint8
	defaultMaxPageLength uint8
	maxMaxPageLength     uint8
}

func WithSql(sql *sql.DB) func(*Service) {
	return func(input *Service) {
		input.sql = sql
	}
}

func WithDatabase(database *model.Queries) func(*Service) {
	return func(input *Service) {
		input.database = database
	}
}

func WithCache(cache cache.Cacher) func(*Service) {
	return func(input *Service) {
		input.cache = cache
	}
}

func WithMetrics(metrics metrics.Metrics) func(*Service) {
	return func(input *Service) {
		input.metrics = metrics
	}
}

func WithMinEventNameLength(minEventNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.minEventNameLength = minEventNameLength
	}
}

func WithMaxEventNameLength(maxEventNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.maxEventNameLength = maxEventNameLength
	}
}

func WithMinRoundNameLength(minRoundNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.minRoundNameLength = minRoundNameLength
	}
}

func WithMaxRoundNameLength(maxRoundNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.maxRoundNameLength = maxRoundNameLength
	}
}

func WithMaxNumberOfRounds(maxNumberOfRounds uint8) func(*Service) {
	return func(input *Service) {
		input.maxNumberOfRounds = maxNumberOfRounds
	}
}

func WithDefaultMaxPageLength(defaultMaxPageLength uint8) func(*Service) {
	return func(input *Service) {
		input.defaultMaxPageLength = defaultMaxPageLength
	}
}

func WithMaxMaxPageLength(maxMaxPageLength uint8) func(*Service) {
	return func(input *Service) {
		input.maxMaxPageLength = maxMaxPageLength
	}
}

func NewService(opts ...func(*Service)) *Service {
	service := &Service{
		minEventNameLength:   3,
		maxEventNameLength:   20,
		minRoundNameLength:   3,
		maxRoundNameLength:   20,
		maxNumberOfRounds:    10,
		defaultMaxPageLength: 10,
		maxMaxPageLength:     100,
	}
	for _, opt := range opts {
		opt(service)
	}
	return service
}

func (s *Service) CreateEvent(ctx context.Context, in *api.CreateEventRequest) (*api.CreateEventResponse, error) {
	command := NewCreateEventCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetEvent(ctx context.Context, in *api.GetEventRequest) (*api.GetEventResponse, error) {
	command := NewGetEventCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateEvent(ctx context.Context, in *api.UpdateEventRequest) (*api.UpdateEventResponse, error) {
	command := NewUpdateEventCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) DeleteEvent(ctx context.Context, in *api.EventRequest) (*api.EventResponse, error) {
	command := NewDeleteEventCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) CreateEventRound(ctx context.Context, in *api.CreateEventRoundRequest) (*api.CreateEventRoundResponse, error) {
	command := NewCreateEventRoundCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetEventRound(ctx context.Context, in *api.GetEventRoundRequest) (*api.GetEventRoundResponse, error) {
	command := NewGetEventRoundCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateEventRound(ctx context.Context, in *api.UpdateEventRoundRequest) (*api.UpdateEventRoundResponse, error) {
	command := NewUpdateEventRoundCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetEventUser(ctx context.Context, in *api.GetEventUserRequest) (*api.GetEventUserResponse, error) {
	command := NewGetEventUserCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateEventUser(ctx context.Context, in *api.UpdateEventUserRequest) (*api.UpdateEventUserResponse, error) {
	command := NewUpdateEventUserCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) DeleteEventUser(ctx context.Context, in *api.EventUserRequest) (*api.EventUserResponse, error) {
	command := NewDeleteEventUserCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) AddEventResult(ctx context.Context, in *api.AddEventResultRequest) (*api.AddEventResultResponse, error) {
	command := NewAddEventResultCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) RemoveEventResult(ctx context.Context, in *api.EventRoundUserRequest) (*api.RemoveEventResultResponse, error) {
	command := NewRemoveEventResultCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

// Utility functions

func UnmarshalEventWithRound(event []model.EventWithRound) (*api.Event, error) {
	if len(event) == 0 {
		return nil, nil
	}
	data, err := conversion.RawJsonToProtobufStruct(event[0].Data)
	if err != nil {
		return nil, err
	}
	eventWithRound := api.Event{
		Id:               event[0].ID,
		Name:             event[0].Name,
		CurrentRoundId:   event[0].CurrentRoundID,
		CurrentRoundName: event[0].CurrentRoundName,
		Data:             data,
		StartedAt:        conversion.TimeToTimestamppb(&event[0].StartedAt),
		CreatedAt:        conversion.TimeToTimestamppb(&event[0].CreatedAt),
		UpdatedAt:        conversion.TimeToTimestamppb(&event[0].UpdatedAt),
	}
	rounds := make([]*api.EventRound, 0, len(event))
	for _, round := range event {
		roundId := conversion.SqlNullInt64ToInt64(round.RoundID)
		if roundId == nil {
			return nil, errors.New("round id is null")
		}
		roundName := conversion.SqlNullStringToString(round.RoundName)
		if roundName == nil {
			return nil, errors.New("round name is null")
		}
		roundScoring, err := conversion.RawJsonToMap(round.RoundScoring)
		if err != nil {
			return nil, err
		}
		// Check if we have scoring
		if _, ok := roundScoring["scoring"]; !ok {
			return nil, errors.New("round scoring is missing")
		}
		// Convert round scoring to uint64 array
		scoringField := roundScoring["scoring"].([]interface{})
		scoringArray := make([]uint64, 0, len(scoringField))
		for _, score := range scoringField {
			scoringArray = append(scoringArray, uint64(score.(float64)))
		}
		if len(scoringArray) == 0 {
			return nil, errors.New("round scoring is empty")
		}
		roundData, err := conversion.RawJsonToProtobufStruct(round.RoundData)
		if err != nil {
			return nil, err
		}
		endedAt := conversion.SqlNullTimeToTimestamp(round.RoundEndedAt)
		if endedAt == nil {
			return nil, errors.New("round ended at is null")
		}
		createdAt := conversion.SqlNullTimeToTimestamp(round.RoundCreatedAt)
		if createdAt == nil {
			return nil, errors.New("round created at is null")
		}
		updatedAt := conversion.SqlNullTimeToTimestamp(round.RoundUpdatedAt)
		if updatedAt == nil {
			return nil, errors.New("round updated at is null")
		}
		rounds = append(rounds, &api.EventRound{
			Id:        uint64(*roundId),
			Name:      *roundName,
			Scoring:   scoringArray,
			Data:      roundData,
			EndedAt:   endedAt,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}
	eventWithRound.Rounds = rounds
	return &eventWithRound, nil
}

func UnmarshalEventLeaderboard(leaderboard []model.EventLeaderboard) ([]*api.EventUser, error) {
	eventUsers := make([]*api.EventUser, 0, len(leaderboard))
	for _, eventUser := range leaderboard {
		data, err := conversion.RawJsonToProtobufStruct(eventUser.Data)
		if err != nil {
			return nil, err
		}
		eventUsers = append(eventUsers, &api.EventUser{
			Id:        eventUser.ID,
			EventId:   eventUser.EventID,
			UserId:    eventUser.UserID,
			Score:     eventUser.Score,
			Ranking:   eventUser.Ranking,
			Data:      data,
			CreatedAt: conversion.TimeToTimestamppb(&eventUser.CreatedAt),
			UpdatedAt: conversion.TimeToTimestamppb(&eventUser.UpdatedAt),
		})
	}
	return eventUsers, nil
}

func UnmarshalEventRound(eventRound model.EventRound) (*api.EventRound, error) {
	data, err := conversion.RawJsonToProtobufStruct(eventRound.Data)
	if err != nil {
		return nil, err
	}
	// Convert scoring to uint64 array
	scoringField, err := conversion.RawJsonToMap(eventRound.Scoring)
	if err != nil {
		return nil, err
	}
	scoringArray := make([]uint64, 0, len(scoringField))
	for _, score := range scoringField["scoring"].([]interface{}) {
		scoringArray = append(scoringArray, uint64(score.(float64)))
	}
	endedAt := conversion.TimeToTimestamppb(&eventRound.EndedAt)
	createdAt := conversion.TimeToTimestamppb(&eventRound.CreatedAt)
	updatedAt := conversion.TimeToTimestamppb(&eventRound.UpdatedAt)
	return &api.EventRound{
		Id:        eventRound.ID,
		EventId:   eventRound.EventID,
		Name:      eventRound.Name,
		Scoring:   scoringArray,
		Data:      data,
		EndedAt:   endedAt,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func UnmarshalEventRoundLeaderboard(leaderboard []model.EventRoundLeaderboard) ([]*api.EventRoundUser, error) {
	eventUsers := make([]*api.EventRoundUser, 0, len(leaderboard))
	for _, eventUser := range leaderboard {
		data, err := conversion.RawJsonToProtobufStruct(eventUser.Data)
		if err != nil {
			return nil, err
		}
		eventUsers = append(eventUsers, &api.EventRoundUser{
			Id:           eventUser.ID,
			EventUserId:  eventUser.EventUserID,
			EventRoundId: eventUser.EventRoundID,
			Result:       eventUser.Result,
			Ranking:      eventUser.Ranking,
			Data:         data,
			CreatedAt:    conversion.TimeToTimestamppb(&eventUser.CreatedAt),
			UpdatedAt:    conversion.TimeToTimestamppb(&eventUser.UpdatedAt),
		})
	}
	return eventUsers, nil
}

func UnmarshalEventUser(eventUser model.EventLeaderboard) (*api.EventUser, error) {
	data, err := conversion.RawJsonToProtobufStruct(eventUser.Data)
	if err != nil {
		return nil, err
	}
	createdAt := conversion.TimeToTimestamppb(&eventUser.CreatedAt)
	updatedAt := conversion.TimeToTimestamppb(&eventUser.UpdatedAt)
	return &api.EventUser{
		Id:        eventUser.ID,
		EventId:   eventUser.EventID,
		UserId:    eventUser.UserID,
		Score:     eventUser.Score,
		Ranking:   eventUser.Ranking,
		Data:      data,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// Enum for errors
type EventRequestError string

const (
	NAME_TOO_SHORT      EventRequestError = "NAME_TOO_SHORT"
	NAME_TOO_LONG       EventRequestError = "NAME_TOO_LONG"
	ID_OR_NAME_REQUIRED EventRequestError = "ID_OR_NAME_REQUIRED"
	// Event round errors
	EVENT_ROUND_OR_ID_REQUIRED EventRequestError = "EVENT_ROUND_OR_ID_REQUIRED"
	// Event user errors
	USER_ID_REQUIRED          EventRequestError = "USER_ID_REQUIRED"
	EVENT_USER_OR_ID_REQUIRED EventRequestError = "EVENT_USER_OR_ID_REQUIRED"
	// Event round user errors
	EVENT_ROUND_USER_OR_ID_REQUIRED EventRequestError = "EVENT_ROUND_USER_OR_ID_REQUIRED"
)

func (s *Service) checkForEventRequestError(request *api.EventRequest) *EventRequestError {
	if request == nil {
		return conversion.ValueToPointer(ID_OR_NAME_REQUIRED)
	}
	if request.Id != nil {
		return nil
	}
	if request.Name == nil {
		return conversion.ValueToPointer(ID_OR_NAME_REQUIRED)
	}
	if len(*request.Name) < int(s.minEventNameLength) {
		return conversion.ValueToPointer(NAME_TOO_SHORT)
	}
	if len(*request.Name) > int(s.maxEventNameLength) {
		return conversion.ValueToPointer(NAME_TOO_LONG)
	}
	return nil
}

func (s *Service) checkForEventRoundRequestError(request *api.EventRoundRequest) *EventRequestError {
	if request == nil {
		return conversion.ValueToPointer(EVENT_ROUND_OR_ID_REQUIRED)
	}
	if request.Id != nil {
		return nil
	}
	if request.RoundName != nil {
		if len(*request.RoundName) < int(s.minRoundNameLength) {
			return conversion.ValueToPointer(NAME_TOO_SHORT)
		}
		if len(*request.RoundName) > int(s.maxRoundNameLength) {
			return conversion.ValueToPointer(NAME_TOO_LONG)
		}
	}
	return s.checkForEventRequestError(request.Event)
}

func (s *Service) checkForEventUserRequestError(request *api.EventUserRequest) *EventRequestError {
	if request == nil {
		return conversion.ValueToPointer(EVENT_USER_OR_ID_REQUIRED)
	}
	if request.Id != nil {
		return nil
	}
	if request.UserId == nil {
		return conversion.ValueToPointer(USER_ID_REQUIRED)
	}
	return s.checkForEventRequestError(request.Event)
}

func (s *Service) checkForEventRoundUserRequestError(request *api.EventRoundUserRequest) *EventRequestError {
	if request == nil {
		return conversion.ValueToPointer(EVENT_ROUND_USER_OR_ID_REQUIRED)
	}
	if request.Id != nil {
		return nil
	}
	if request.Round != nil {
		if len(*request.Round) < int(s.minRoundNameLength) {
			return conversion.ValueToPointer(NAME_TOO_SHORT)
		}
		if len(*request.Round) > int(s.maxRoundNameLength) {
			return conversion.ValueToPointer(NAME_TOO_LONG)
		}
	}
	return s.checkForEventUserRequestError(request.User)
}
