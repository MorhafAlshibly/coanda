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

func Test_GetArenas_NoRows_NoRowsReturned(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetArenasCommand(service, &api.Pagination{})
	mock.ExpectQuery("SELECT (.+) FROM matchmaking_arena").WillReturnRows(sqlmock.NewRows(matchmakingArenaFields))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := len(c.Out.Arenas), 0; got != want {
		t.Fatalf("Expected arenas length to be %d, got %d", want, got)
	}
}

func Test_GetArenas_WithRows_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetArenasCommand(service, &api.Pagination{})
	mock.ExpectQuery("SELECT (.+) FROM matchmaking_arena").
		WillReturnRows(sqlmock.NewRows(matchmakingArenaFields).
			AddRow(uint64(1), "Arena 1", 2, 4, 8, json.RawMessage("{}"), time.Now(), time.Now()).
			AddRow(uint64(2), "Arena 2", 3, 5, 10, json.RawMessage("{}"), time.Now(), time.Now()))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := len(c.Out.Arenas), 2; got != want {
		t.Fatalf("Expected arenas length to be %d, got %d", want, got)
	}
	if got, want := c.Out.Arenas[0].Id, uint64(1); got != want {
		t.Fatalf("Expected first arena id to be %d, got %d", want, got)
	}
	if got, want := c.Out.Arenas[0].Name, "Arena 1"; got != want {
		t.Fatalf("Expected first arena name to be %s, got %s", want, got)
	}
	if got, want := c.Out.Arenas[1].Id, uint64(2); got != want {
		t.Fatalf("Expected second arena id to be %d, got %d", want, got)
	}
	if got, want := c.Out.Arenas[1].Name, "Arena 2"; got != want {
		t.Fatalf("Expected second arena name to be %s, got %s", want, got)
	}
}
