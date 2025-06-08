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

func Test_GetMatchmakingUsers_NoRows_NoRowsReturned(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetMatchmakingUsersCommand(service, &api.Pagination{})
	mock.ExpectQuery("SELECT (.+) FROM matchmaking_user").
		WillReturnRows(sqlmock.NewRows(matchmakingUserFields))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := len(c.Out.MatchmakingUsers), 0; got != want {
		t.Fatalf("Expected matchmaking users length to be %d, got %d", want, got)
	}
}

func Test_GetMatchmakingUsers_WithRows_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewGetMatchmakingUsersCommand(service, &api.Pagination{})
	mock.ExpectQuery("SELECT (.+) FROM matchmaking_user").
		WillReturnRows(sqlmock.NewRows(matchmakingUserFields).
			AddRow(uint64(1), nil, 100, 0, json.RawMessage("{}"), time.Now(), time.Now()).
			AddRow(uint64(2), nil, 200, 1, json.RawMessage("{}"), time.Now(), time.Now()))
	err = invoker.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.Out.Success, true; got != want {
		t.Fatalf("Expected success to be %v, got %v", want, got)
	}
	if got, want := len(c.Out.MatchmakingUsers), 2; got != want {
		t.Fatalf("Expected matchmaking users length to be %d, got %d", want, got)
	}
	if got, want := c.Out.MatchmakingUsers[0].Id, uint64(1); got != want {
		t.Fatalf("Expected matchmaking user ID to be %d, got %d", want, got)
	}
	if got, want := c.Out.MatchmakingUsers[1].Id, uint64(2); got != want {
		t.Fatalf("Expected matchmaking user ID to be %d, got %d", want, got)
	}
}
