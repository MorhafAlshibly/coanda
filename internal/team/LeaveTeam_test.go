package team

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/errorcodes"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/go-sql-driver/mysql"
)

func TestLeaveTeam(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("DELETE FROM team_member").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	c := NewLeaveTeamCommand(service, &api.LeaveTeamRequest{
		UserId: 1,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.LeaveTeamResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
}

func TestLeaveTeamNotInTeam(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("DELETE FROM team_member").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 0))
	c := NewLeaveTeamCommand(service, &api.LeaveTeamRequest{
		UserId: 1,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.LeaveTeamResponse_NOT_IN_TEAM {
		t.Fatal("Expected error to be NOT_IN_TEAM")
	}
}

func TestLeaveTeamNoUserId(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	c := NewLeaveTeamCommand(service, &api.LeaveTeamRequest{})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.LeaveTeamResponse_USER_ID_REQUIRED {
		t.Fatal("Expected error to be USER_ID_REQUIRED")
	}
}

func TestLeaveTeamMemberIsOwner(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectExec("DELETE FROM team_member").WithArgs(1).WillReturnError(&mysql.MySQLError{Number: errorcodes.MySQLErrorCodeRowIsReferenced2})
	c := NewLeaveTeamCommand(service, &api.LeaveTeamRequest{
		UserId: 1,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.LeaveTeamResponse_MEMBER_IS_OWNER {
		t.Fatal("Expected error to be MEMBER_IS_OWNER")
	}
}
