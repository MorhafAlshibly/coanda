package tournament

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
)

func TestDeleteTournamentUserTournamentNameTooShort(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteTournamentUserCommand(service, &api.TournamentUserRequest{
		Tournament: "t",
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.TournamentUserResponse_TOURNAMENT_NAME_TOO_SHORT {
		t.Fatal("Expected error to be TOURNAMENT_NAME_TOO_SHORT")
	}
}

func TestDeleteTournamentUserTournamentNameTooLong(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries), WithMaxTournamentNameLength(5))
	c := NewDeleteTournamentUserCommand(service, &api.TournamentUserRequest{
		Tournament: "aaaaaaa",
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.TournamentUserResponse_TOURNAMENT_NAME_TOO_LONG {
		t.Fatal("Expected error to be TOURNAMENT_NAME_TOO_LONG")
	}
}

func TestDeleteTournamentUserNoUserId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewDeleteTournamentUserCommand(service, &api.TournamentUserRequest{
		Tournament: "test",
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.TournamentUserResponse_USER_ID_REQUIRED {
		t.Fatal("Expected error to be USER_ID_REQUIRED")
	}
}

func TestDeleteTeamSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("DELETE FROM tournament").WithArgs("test", "DAILY", 1, time.Now().Truncate(time.Hour*24).UTC()).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewDeleteTournamentUserCommand(service, &api.TournamentUserRequest{
		Tournament: "test",
		Interval:   api.TournamentInterval_DAILY,
		UserId:     1,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
}

func TestDeleteTeamNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("DELETE FROM tournament").WithArgs("test", "DAILY", 1, time.Now().Truncate(time.Hour*24).UTC()).WillReturnResult(sqlmock.NewResult(0, 0))
	c := NewDeleteTournamentUserCommand(service, &api.TournamentUserRequest{
		Tournament: "test",
		Interval:   api.TournamentInterval_DAILY,
		UserId:     1,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.TournamentUserResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}
