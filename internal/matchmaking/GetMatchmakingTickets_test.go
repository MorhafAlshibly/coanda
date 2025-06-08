package matchmaking

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
)

func Test_GetMatchmakingTickets_NoRows_NoRowsReturned(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetMatchmakingTicketsCommand(service, &api.GetMatchmakingTicketsRequest{})
	mock.ExpectQuery("SELECT (.+) FROM matchmaking_ticket").WillReturnRows(sqlmock.NewRows(matchmakingTicketFields))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := len(c.Out.MatchmakingTickets), 0; got != want {
		t.Fatalf("Expected matchmaking tickets length to be %d, got %d", want, got)
	}
}

func Test_GetMatchmakingTickets_WithRows_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetMatchmakingTicketsCommand(service, &api.GetMatchmakingTicketsRequest{})
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_ticket_with_user_and_arena`").
		WithArgs(uint64(99), 0, 1, 0, 1).
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := len(c.Out.MatchmakingTickets), 2; got != want {
		t.Fatalf("Expected matchmaking tickets length to be %d, got %d", want, got)
	}
}
