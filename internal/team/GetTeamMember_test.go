package team

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
)

var (
	teamMember = []string{"team", "user_id", "data", "joined_at", "updated_at"}
)

func TestGetTeamMemberExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM team_member").WithArgs(2).WillReturnRows(sqlmock.NewRows(teamMember).AddRow("test", 1, raw, time.Now(), time.Now()))
	c := NewGetTeamMemberCommand(service, &api.GetTeamMemberRequest{
		UserId: 2,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != true {
		t.Fatal("Expected success to be true")
	}
	if c.Out.Error != api.GetTeamMemberResponse_NONE {
		t.Fatal("Expected error to be NONE")
	}
	if c.Out.TeamMember.Team != "test" {
		t.Fatal("Expected team name to be test")
	}
	if c.Out.TeamMember.UserId != 1 {
		t.Fatal("Expected team owner to be 1")
	}
	if !reflect.DeepEqual(c.Out.TeamMember.Data, data) {
		t.Fatal("Expected team data to be empty")
	}
}

func TestGetTeamMemberNotExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := model.New(db)
	service := NewService(
		WithSql(db), WithDatabase(queries))
	mock.ExpectQuery("SELECT (.+) FROM team_member").WithArgs(2).WillReturnRows(sqlmock.NewRows(teamMember))
	c := NewGetTeamMemberCommand(service, &api.GetTeamMemberRequest{
		UserId: 2,
	})
	err = invokers.NewBasicInvoker().Invoke(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Out.Success != false {
		t.Fatal("Expected success to be false")
	}
	if c.Out.Error != api.GetTeamMemberResponse_NOT_FOUND {
		t.Fatal("Expected error to be NOT_FOUND")
	}
}
