package handleMatchmaking

import (
	"context"
	"database/sql"
	"log"

	"github.com/MorhafAlshibly/coanda/internal/handleMatchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type App struct {
	sql                *sql.DB
	database           *model.Queries
	eloWindowIncrement uint16
	eloWindowMax       uint16
	limit              int32
}

func WithSql(sql *sql.DB) func(*App) {
	return func(input *App) {
		input.sql = sql
	}
}

func WithDatabase(database *model.Queries) func(*App) {
	return func(input *App) {
		input.database = database
	}
}

func WithEloWindowIncrement(eloWindowIncrement uint16) func(*App) {
	return func(input *App) {
		input.eloWindowIncrement = eloWindowIncrement
	}
}

func WithEloWindowMax(eloWindowMax uint16) func(*App) {
	return func(input *App) {
		input.eloWindowMax = eloWindowMax
	}
}

func WithLimit(limit int32) func(*App) {
	return func(input *App) {
		input.limit = limit
	}
}

func NewApp(options ...func(*App)) *App {
	app := &App{
		eloWindowIncrement: 50,
		eloWindowMax:       200,
	}
	for _, option := range options {
		option(app)
	}
	return app
}

func (a *App) Handler(ctx context.Context) error {
	err := a.createNewMatches(ctx)
	if err != nil {
		log.Fatalf("failed to create new matches: %v", err)
		return err
	}
	err = a.incrementTicketEloWindow(ctx)
	if err != nil {
		log.Fatalf("failed to increment ticket elo window: %v", err)
		return err
	}
	err = a.matchmakeTickets(ctx)
	if err != nil {
		log.Fatalf("failed to matchmake tickets: %v", err)
		return err
	}
	return nil
}

func (a *App) createNewMatches(ctx context.Context) error {
	limit := a.limit
	offset := int32(0)
	// Loop until we have created all matches
	for {
		// Get all tickets that have been aged
		agedTickets, err := a.database.GetAgedMatchmakingTickets(ctx, model.GetAgedMatchmakingTicketsParams{
			EloWindowMax: uint32(a.eloWindowMax),
			Limit:        limit,
			Offset:       offset,
		})
		if err != nil {
			return err
		}
		for _, ticket := range agedTickets {
			tx, err := a.sql.BeginTx(ctx, nil)
			if err != nil {
				continue
			}
			defer tx.Rollback()
			qtx := a.database.WithTx(tx)
			// Get most popular arena on the ticket
			arena, err := qtx.GetMostPopularArenaOnTicket(ctx, ticket.ID)
			if err != nil {
				continue
			}
			// Create a new match
			matchResult, err := qtx.CreateMatch(ctx, arena.ID)
			if err != nil {
				continue
			}
			// Get the match ID
			matchID, err := matchResult.LastInsertId()
			if err != nil {
				continue
			}
			// Update the ticket with the match ID
			addMatchIDResult, err := qtx.AddMatchIDToTicket(ctx, model.AddMatchIDToTicketParams{
				ID:                 ticket.ID,
				MatchmakingMatchID: conversion.Int64ToSqlNullInt64(&matchID),
			})
			if err != nil {
				continue
			}
			rowsAffected, err := addMatchIDResult.RowsAffected()
			if err != nil {
				continue
			}
			if rowsAffected == 0 {
				// If the ticket was not updated, then it already has a match ID
				// We need to rollback the transaction and the match creation
				tx.Rollback()
				continue
			}
			// Commit the transaction
			err = tx.Commit()
			if err != nil {
				continue
			}
		}
		// If we have less than the limit of tickets, we have created all matches
		if int32(len(agedTickets)) < limit {
			break
		}
		offset += limit
	}
	return nil
}

func (a *App) incrementTicketEloWindow(ctx context.Context) error {
	_, err := a.database.IncrementEloWindow(ctx, model.IncrementEloWindowParams{
		EloWindowIncrement: uint32(a.eloWindowIncrement),
		EloWindowMax:       uint32(a.eloWindowMax),
	})
	return err
}

func (a *App) matchmakeTickets(ctx context.Context) error {
	limit := a.limit
	offset := int32(0)
	for {
		// Get all tickets that have not been aged
		tickets, err := a.database.GetNonAgedMatchmakingTickets(ctx, model.GetNonAgedMatchmakingTicketsParams{
			EloWindowMax: uint32(a.eloWindowMax),
			Limit:        limit,
			Offset:       offset,
		})
		if err != nil {
			return err
		}
		for _, ticket := range tickets {
			tx, err := a.sql.BeginTx(ctx, nil)
			if err != nil {
				continue
			}
			defer tx.Rollback()
			qtx := a.database.WithTx(tx)
			closestMatch, err := qtx.GetClosestMatch(ctx, ticket.ID)
			if err != nil {
				continue
			}
			// Update the ticket with the match ID
			_, err = qtx.AddMatchIDToTicket(ctx, model.AddMatchIDToTicketParams{
				ID:                 ticket.ID,
				MatchmakingMatchID: conversion.Uint64ToSqlNullInt64(&closestMatch.ID),
			})
		}
		if int32(len(tickets)) < limit {
			break
		}
		offset += limit
	}
	return nil
}
