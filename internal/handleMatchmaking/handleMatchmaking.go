package handleMatchmaking

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MorhafAlshibly/coanda/internal/handleMatchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type App struct {
	sql                         *sql.DB
	database                    *model.Queries
	eloWindowIncrementPerSecond uint16
	eloWindowMax                uint16
	limit                       int32
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

func WithEloWindowIncrementPerSecond(eloWindowIncrementPerSecond uint16) func(*App) {
	return func(input *App) {
		input.eloWindowIncrementPerSecond = eloWindowIncrementPerSecond
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
		eloWindowIncrementPerSecond: 10,
		eloWindowMax:                600,
		limit:                       100,
	}
	for _, option := range options {
		option(app)
	}
	return app
}

func (a *App) Handler(ctx context.Context) error {
	// Two parts of the background job:
	// 1. Create new matches for aged tickets
	// 2. Matchmake tickets that have not been aged yet
	err := a.createNewMatches(ctx)
	if err != nil {
		fmt.Printf("failed to create new matches: %v", err)
		return err
	}
	err = a.matchmakeTickets(ctx)
	if err != nil {
		fmt.Printf("failed to matchmake tickets: %v", err)
		return err
	}
	return nil
}

func (a *App) createNewMatch(ctx context.Context, ticketID uint64) error {
	tx, err := a.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := a.database.WithTx(tx)
	// Lock the ticket for update
	_, err = qtx.LockTicketForUpdate(ctx, ticketID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If we didn't find a row, it means the ticket was deleted or does not exist
			return fmt.Errorf("ticket with ID %d not found", ticketID)
		}
		return err
	}
	// Get most popular arena on the ticket
	arena, err := qtx.GetMostPopularArenaOnTicket(ctx, ticketID)
	if err != nil {
		return err
	}
	// Create a new match
	matchResult, err := qtx.CreateMatch(ctx, arena.ID)
	if err != nil {
		return err
	}
	// Get the match ID
	matchID, err := matchResult.LastInsertId()
	if err != nil {
		return err
	}
	// Update the ticket with the match ID
	addMatchIDResult, err := qtx.AddMatchIDToTicket(ctx, model.AddMatchIDToTicketParams{
		ID:                 ticketID,
		MatchmakingMatchID: conversion.Int64ToSqlNullInt64(&matchID),
	})
	if err != nil {
		return err
	}
	rowsAffected, err := addMatchIDResult.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// If the ticket was not updated, then it already has a match ID
		// We need to rollback the transaction and the match creation
		return err
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
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
			EloWindowIncrementPerSecond: int64(a.eloWindowIncrementPerSecond),
			EloWindowMax:                int64(a.eloWindowMax),
			Limit:                       limit,
			Offset:                      offset,
		})
		if err != nil {
			return err
		}
		for _, ticket := range agedTickets {
			fmt.Println("Creating new match for ticket ID:", ticket.ID)
			err := a.createNewMatch(ctx, ticket.ID)
			if err != nil {
				fmt.Printf("failed to create new match for ticket ID %d: %v\n", ticket.ID, err)
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

func (a *App) matchmakeTicket(ctx context.Context, ticketID uint64) error {
	tx, err := a.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := a.database.WithTx(tx)
	// Lock the ticket for update
	ticket, err := qtx.LockTicketForUpdate(ctx, ticketID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If we didn't find a row, it means the ticket was deleted or does not exist
			return fmt.Errorf("ticket with ID %d not found", ticketID)
		}
		return err
	}
	// Gets the closest match that has enough capacity and locks it for update
	closestMatch, err := qtx.GetClosestMatch(ctx, model.GetClosestMatchParams{
		TicketID: conversion.Uint64ToSqlNullInt64(&ticketID),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			// If we didn't find a row, it means there are no matches available
			return fmt.Errorf("No matches available for ticket ID %d", ticketID)
		}
		return err
	}
	// Check if the closest match's elo difference is within the acceptable range
	// Elo window increment per second multiplied by time now minus the ticket's created at time
	eloWindow := int64(a.eloWindowIncrementPerSecond) * int64(time.Since(ticket.CreatedAt).Seconds())
	if closestMatch.EloDifference > eloWindow {
		// If the elo difference is greater than the elo window, we cannot matchmake this ticket
		return fmt.Errorf("Elo difference %d is greater than the elo window %d for ticket ID %d", closestMatch.EloDifference, eloWindow, ticketID)
	}
	if closestMatch.LockedAt.Valid {
		if closestMatch.LockedAt.Time.Before(time.Now()) {
			// Match is locked
			return fmt.Errorf("Match with ID %d is locked and cannot be used for ticket ID %d", closestMatch.MatchID, ticketID)
		}
	}
	// Update the ticket with the match ID
	_, err = qtx.AddMatchIDToTicket(ctx, model.AddMatchIDToTicketParams{
		ID:                 ticketID,
		MatchmakingMatchID: conversion.Uint64ToSqlNullInt64(&closestMatch.MatchID),
	})
	if err != nil {
		return err
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	fmt.Println("Matchmaking successful for ticket ID:", ticketID, "with match ID:", closestMatch.MatchID)
	return nil
}

func (a *App) matchmakeTickets(ctx context.Context) error {
	limit := a.limit
	offset := int32(0)
	for {
		// Get all tickets that have not been aged
		tickets, err := a.database.GetNonAgedMatchmakingTickets(ctx, model.GetNonAgedMatchmakingTicketsParams{
			EloWindowIncrementPerSecond: int64(a.eloWindowIncrementPerSecond),
			EloWindowMax:                int64(a.eloWindowMax),
			Limit:                       limit,
			Offset:                      offset,
		})
		if err != nil {
			return err
		}
		for _, ticket := range tickets {
			fmt.Println("Finding closest match for ticket ID:", ticket.ID)
			err := a.matchmakeTicket(ctx, ticket.ID)
			if err != nil {
				fmt.Printf("failed to matchmake ticket ID %d: %v\n", ticket.ID, err)
				continue
			}
		}
		if int32(len(tickets)) < limit {
			break
		}
		offset += limit
	}
	return nil
}
