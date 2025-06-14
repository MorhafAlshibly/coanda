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

func Test_GetMatches_NoRows_NoRowsReturned(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetMatchesCommand(service, &api.GetMatchesRequest{})
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WillReturnRows(sqlmock.NewRows(matchmakingTicketFields))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := len(c.Out.Matches), 0; got != want {
		t.Fatalf("Expected matches length to be %d, got %d", want, got)
	}
}

func Test_GetMatches_WithRows_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetMatchesCommand(service, &api.GetMatchesRequest{})
	mock.ExpectQuery("SELECT (.+) FROM `matchmaking_match_with_arena_and_ticket`").
		WillReturnRows(sqlmock.NewRows(matchmakingMatchWithArenaAndTicketFields).
			AddRow(
				uint64(9), nil, "PENDING", 1, 1, json.RawMessage("{}"), nil, nil, nil, time.Now(), time.Now(),
				uint64(1), "Arena1", 2, 4, 8, json.RawMessage("{}"), time.Now(), time.Now(),
				uint64(1), uint64(4), "MATCHED", 1, 1, json.RawMessage("{}"), time.Now(), time.Now(),
				uint64(4), 1200, 1, json.RawMessage("{}"), time.Now(), time.Now(),
				uint64(1), "Arena1", 2, 4, 8, 1, json.RawMessage("{}"), time.Now(), time.Now(),
			).
			AddRow(
				uint64(10), nil, "PENDING", 2, 2, json.RawMessage("{}"), nil, nil, nil, time.Now(), time.Now(),
				uint64(2), "Arena2", 3, 5, 10, json.RawMessage("{}"), time.Now(), time.Now(),
				uint64(2), uint64(5), "MATCHED", 2, 2, json.RawMessage("{}"), time.Now(), time.Now(),
				uint64(5), 1500, 2, json.RawMessage("{}"), time.Now(), time.Now(),
				uint64(2), "Arena2", 3, 5, 10, 2, json.RawMessage("{}"), time.Now(), time.Now(),
			))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := len(c.Out.Matches), 2; got != want {
		t.Fatalf("Expected matches length to be %d, got %d", want, got)
	}
	if got, want := c.Out.Matches[0].Id, uint64(9); got != want {
		t.Fatalf("Expected first match id to be %d, got %d", want, got)
	}
	if got, want := c.Out.Matches[1].Id, uint64(10); got != want {
		t.Fatalf("Expected second match id to be %d, got %d", want, got)
	}
}
