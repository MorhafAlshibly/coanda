package matchmaking

import (
	"context"
	"encoding/json"
	"testing"
	"time"

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
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_ticket_with_user_and_arena`").WillReturnRows(sqlmock.NewRows(matchmakingTicketFields))
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
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields).
			AddRow(
				uint64(10), uint64(3), "ENDED", 4, json.RawMessage("{}"),
				time.Now().Add(-2*time.Hour), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
				3, 3, 1600, 1, json.RawMessage("{}"), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
				3, "Arena3", 4, 8, 8, 0, json.RawMessage("{}"), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour),
			).
			AddRow(
				uint64(11), uint64(4), "MATCHED", 5, json.RawMessage("{}"),
				time.Now().Add(-time.Hour), time.Now(), time.Now(),
				4, 4, 1700, 1, json.RawMessage("{}"), time.Now(), time.Now(),
				4, "Arena4", 5, 10, 10, 0, json.RawMessage("{}"), time.Now(), time.Now(),
			))
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
	if got, want := c.Out.MatchmakingTickets[0].Id, uint64(10); got != want {
		t.Fatalf("Expected first matchmaking ticket id to be %d, got %d", want, got)
	}
	if got, want := c.Out.MatchmakingTickets[1].Id, uint64(11); got != want {
		t.Fatalf("Expected second matchmaking ticket id to be %d, got %d", want, got)
	}
}
