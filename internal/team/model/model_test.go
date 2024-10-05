package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/MorhafAlshibly/coanda/pkg/mysqlTestServer"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func TestMain(m *testing.M) {
	server, err := mysqlTestServer.GetServer()
	if err != nil {
		log.Fatalf("could not run mysql test server: %v", err)
	}
	defer server.Close()
	db = server.Db
	schema, err := os.ReadFile("../../../migration/team.sql")
	if err != nil {
		log.Fatalf("could not read schema file: %v", err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalf("could not execute schema: %v", err)
	}

	m.Run()
}

func Test_CreateTeam_Team_TeamCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateTeam_TeamNameExists_TeamNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team1",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team1",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if !errorcode.IsDuplicateEntry(mysqlErr, "team", "team_name_idx") {
		t.Fatalf("expected duplicate entry error, got %d", mysqlErr.Number)
	}
}

func Test_CreateTeamMember_TeamMember_TeamMemberCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team9",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	result, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       11,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateTeamMember_TeamMemberExists_TeamMemberNotCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team10",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       12,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       12,
		MemberNumber: 2,
		Data:         json.RawMessage(`{}`),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if !errorcode.IsDuplicateEntry(mysqlErr, "team_member", "team_member_user_id_idx") {
		t.Fatalf("expected duplicate entry error, got %d", mysqlErr.Number)
	}
}

func Test_CreateTeamMember_TeamMemberNumberExists_TeamMemberNotCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team11",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       13,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       14,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if !errorcode.IsDuplicateEntry(mysqlErr, "team_member", "team_member_team_id_member_number_idx") {
		t.Fatalf("expected duplicate entry error, got %d", mysqlErr.Number)
	}
}

func Test_CreateTeamMember_TeamDoesNotExist_TeamMemberNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       999999,
		UserID:       15,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if mysqlErr.Number != errorcode.MySQLErrorCodeNoReferencedRow2 {
		t.Fatalf("expected foreign key constraint error, got %d", mysqlErr.Number)
	}
}

func Test_CreateTeamMember_TeamMemberInAnotherTeam_TeamMemberNotCreated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team13",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId1, err := result.LastInsertId()
	result, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team14",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId2, err := result.LastInsertId()
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId1),
		UserID:       18,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId2),
		UserID:       18,
		MemberNumber: 2,
		Data:         json.RawMessage(`{}`),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if mysqlErr.Number != errorcode.MySQLErrorCodeDuplicateEntry {
		t.Fatalf("expected duplicate key constraint error, got %d", mysqlErr.Number)
	}
}

func Test_DeleteTeamMember_ByUserId_TeamMemberDeleted(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team15",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       20,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	result, err = q.DeleteTeamMember(context.Background(), GetTeamMemberParams{
		UserID: sql.NullInt64{Int64: 20, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete team member: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_DeleteTeamMember_ById_TeamMemberDeleted(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team150",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	result, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       200,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMemberId, err := result.LastInsertId()
	result, err = q.DeleteTeamMember(context.Background(), GetTeamMemberParams{
		ID: sql.NullInt64{Int64: teamMemberId, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete team member: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_DeleteTeamMember_ByUserIdTeamMemberDoesNotExist_TeamMemberNotDeleted(t *testing.T) {
	q := New(db)
	result, err := q.DeleteTeamMember(context.Background(), GetTeamMemberParams{
		UserID: sql.NullInt64{Int64: 99999, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete team member: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_DeleteTeamMember_ByIdTeamMemberDoesNotExist_TeamMemberNotDeleted(t *testing.T) {
	q := New(db)
	result, err := q.DeleteTeamMember(context.Background(), GetTeamMemberParams{
		ID: sql.NullInt64{Int64: 99999, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete team member: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_GetFirstOpenMemberNumber_TeamNoMembers_ReturnOne(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team16",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamID, err := result.LastInsertId()
	member_number, err := q.GetFirstOpenMemberNumber(context.Background(), uint64(teamID))
	if err != nil {
		t.Fatalf("could not get first open member number: %v", err)
	}
	if member_number != 1 {
		t.Fatalf("expected 1, got %d", member_number)
	}
}

func Test_GetFirstOpenMemberNumber_TeamHasMembers_ReturnNextMemberNumber(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team17",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamID, err := result.LastInsertId()
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamID),
		UserID:       24,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	member_number, err := q.GetFirstOpenMemberNumber(context.Background(), uint64(teamID))
	if err != nil {
		t.Fatalf("could not get first open member number: %v", err)
	}
	if member_number != 2 {
		t.Fatalf("expected 2, got %d", member_number)
	}
}

func Test_GetFirstOpenMemberNumber_TeamWithGapInMemberNumbers_ReturnGapMemberNumber(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team18",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamID, err := result.LastInsertId()
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamID),
		UserID:       25,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamID),
		UserID:       26,
		MemberNumber: 3,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	member_number, err := q.GetFirstOpenMemberNumber(context.Background(), uint64(teamID))
	if err != nil {
		t.Fatalf("could not get first open member number: %v", err)
	}
	if member_number != 2 {
		t.Fatalf("expected 2, got %d", member_number)
	}
}

func Test_GetFirstOpenMemberNumber_TeamDoesNotExist_ReturnError(t *testing.T) {
	q := New(db)
	_, err := q.GetFirstOpenMemberNumber(context.Background(), 9999999)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetTeamMember_ByUserId_ReturnTeamMember(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team19",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamID, err := result.LastInsertId()
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamID),
		UserID:       27,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMember, err := q.GetTeamMember(context.Background(), GetTeamMemberParams{
		UserID: sql.NullInt64{Int64: 27, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get team member: %v", err)
	}
	if teamMember.TeamID != uint64(teamID) {
		t.Fatalf("expected %d, got %d", teamID, teamMember.TeamID)
	}
	if teamMember.UserID != 27 {
		t.Fatalf("expected 27, got %d", teamMember.UserID)
	}
	if teamMember.MemberNumber != 1 {
		t.Fatalf("expected 1, got %d", teamMember.MemberNumber)
	}
}

func Test_GetTeamMember_ById_ReturnTeamMember(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team190",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamID, err := result.LastInsertId()
	result, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamID),
		UserID:       270,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMemberId, err := result.LastInsertId()
	teamMember, err := q.GetTeamMember(context.Background(), GetTeamMemberParams{
		ID: sql.NullInt64{Int64: teamMemberId, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get team member: %v", err)
	}
	if teamMember.TeamID != uint64(teamID) {
		t.Fatalf("expected %d, got %d", teamID, teamMember.TeamID)
	}
	if teamMember.UserID != 270 {
		t.Fatalf("expected 270, got %d", teamMember.UserID)
	}
	if teamMember.MemberNumber != 1 {
		t.Fatalf("expected 1, got %d", teamMember.MemberNumber)
	}
}

func Test_GetTeamMember_ByUserIdTeamMemberDoesNotExist_ReturnNil(t *testing.T) {
	q := New(db)
	_, err := q.GetTeamMember(context.Background(), GetTeamMemberParams{
		UserID: sql.NullInt64{Int64: 9999999, Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetTeamMember_ByIdTeamMemberDoesNotExist_ReturnNil(t *testing.T) {
	q := New(db)
	_, err := q.GetTeamMember(context.Background(), GetTeamMemberParams{
		ID: sql.NullInt64{Int64: 9999999, Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetTeams_TwoTeams_ReturnTeams(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team20",
		Score: 1,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team21",
		Score: 1,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teams, err := q.GetTeams(context.Background(), GetTeamsParams{
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get teams: %v", err)
	}
	if len(teams) != 2 {
		t.Fatalf("expected 2 teams, got %d", len(teams))
	}
}

func Test_SearchTeams_QueryForSpecialWord_TeamsWithSpecialWordInMiddle(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team23",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "teamwithspecialwordinthemiddletest",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teams, err := q.SearchTeams(context.Background(), SearchTeamsParams{
		Query:  "specialwordinthemiddle",
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not search teams: %v", err)
	}
	if len(teams) != 1 {
		t.Fatalf("expected 1 team, got %d", len(teams))
	}
	if teams[0].Name != "teamwithspecialwordinthemiddletest" {
		t.Fatalf("expected teamwithspecialwordinthemiddletest, got %s", teams[0].Name)
	}
}

func Test_SearchTeams_QueryForSpecialWord_TeamsWithSpecialWordAtEnd(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team24",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "teamwithspecialwordatend",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teams, err := q.SearchTeams(context.Background(), SearchTeamsParams{
		Query:  "specialwordatend",
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not search teams: %v", err)
	}
	if len(teams) != 1 {
		t.Fatalf("expected 1 team, got %d", len(teams))
	}
	if teams[0].Name != "teamwithspecialwordatend" {
		t.Fatalf("expected teamwithspecialwordatend, got %s", teams[0].Name)
	}
}

func Test_SearchTeams_QueryForSpecialWord_TeamsWithSpecialWordAtStart(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team25",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "specialwordatstartteam",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teams, err := q.SearchTeams(context.Background(), SearchTeamsParams{
		Query:  "specialwordatstart",
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not search teams: %v", err)
	}
	if len(teams) != 1 {
		t.Fatalf("expected 1 team, got %d", len(teams))
	}
	if teams[0].Name != "specialwordatstartteam" {
		t.Fatalf("expected specialwordatstartteam, got %s", teams[0].Name)
	}
}

func Test_SearchTeams_QueryForSpecialWord_TeamsWithSpecialWordOnly(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team26",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "specialwordteam",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teams, err := q.SearchTeams(context.Background(), SearchTeamsParams{
		Query:  "specialwordteam",
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not search teams: %v", err)
	}
	if len(teams) != 1 {
		t.Fatalf("expected 1 team, got %d", len(teams))
	}
	if teams[0].Name != "specialwordteam" {
		t.Fatalf("expected specialwordteam, got %s", teams[0].Name)
	}
}

func Test_SearchTeams_QueryForSpecialWord_TeamsWithSpecialWordInMiddleCaseInsensitive(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team27",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "teamwithspecialwordinthemiddlecaseinsensitivetest",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teams, err := q.SearchTeams(context.Background(), SearchTeamsParams{
		Query:  "SpecialWordInTheMiddleCaseInsensitive",
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not search teams: %v", err)
	}
	if len(teams) != 1 {
		t.Fatalf("expected 1 team, got %d", len(teams))
	}
	if teams[0].Name != "teamwithspecialwordinthemiddlecaseinsensitivetest" {
		t.Fatalf("expected teamwithspecialwordinthemiddlecaseinsensitivetest, got %s", teams[0].Name)
	}
}

func Test_SearchTeams_QueryForSpecialWordNoTeams_ReturnEmpty(t *testing.T) {
	q := New(db)
	teams, err := q.SearchTeams(context.Background(), SearchTeamsParams{
		Query:  "specialwordnoteams",
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not search teams: %v", err)
	}
	if len(teams) != 0 {
		t.Fatalf("expected 0 teams, got %d", len(teams))
	}
}

func Test_UpdateTeamMember_ByUserId_TeamMemberUpdated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team28",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       43,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	result, err = q.UpdateTeamMember(context.Background(), UpdateTeamMemberParams{
		TeamMember: GetTeamMemberParams{
			UserID: sql.NullInt64{Int64: 43, Valid: true},
		},
		Data: json.RawMessage(`{"test": "test"}`),
	})
	if err != nil {
		t.Fatalf("could not update team member: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	team, err := q.GetTeamMember(context.Background(), GetTeamMemberParams{
		UserID: sql.NullInt64{Int64: 43, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get team member: %v", err)
	}
	if !reflect.DeepEqual(team.Data, json.RawMessage(`{"test": "test"}`)) {
		t.Fatalf("expected {\"test\": \"test\"}, got %s", team.Data)
	}
}

func Test_UpdateTeamMember_ById_TeamMemberUpdated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team280",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	result, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       430,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMemberId, err := result.LastInsertId()
	result, err = q.UpdateTeamMember(context.Background(), UpdateTeamMemberParams{
		TeamMember: GetTeamMemberParams{
			ID: sql.NullInt64{Int64: teamMemberId, Valid: true},
		},
		Data: json.RawMessage(`{"test": "test"}`),
	})
	if err != nil {
		t.Fatalf("could not update team member: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	team, err := q.GetTeamMember(context.Background(), GetTeamMemberParams{
		ID: sql.NullInt64{Int64: teamMemberId, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get team member: %v", err)
	}
	if !reflect.DeepEqual(team.Data, json.RawMessage(`{"test": "test"}`)) {
		t.Fatalf("expected {\"test\": \"test\"}, got %s", team.Data)
	}
}

func Test_UpdateTeamMember_ByUserIdTeamMemberDoesNotExist_TeamMemberNotUpdated(t *testing.T) {
	q := New(db)
	result, err := q.UpdateTeamMember(context.Background(), UpdateTeamMemberParams{
		TeamMember: GetTeamMemberParams{
			UserID: sql.NullInt64{Int64: 9999999, Valid: true},
		},
		Data: json.RawMessage(`{"test": "test"}`),
	})
	if err != nil {
		t.Fatalf("could not update team member: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_UpdateTeamMember_ByIdTeamMemberDoesNotExist_TeamMemberNotUpdated(t *testing.T) {
	q := New(db)
	result, err := q.UpdateTeamMember(context.Background(), UpdateTeamMemberParams{
		TeamMember: GetTeamMemberParams{
			ID: sql.NullInt64{Int64: 9999999, Valid: true},
		},
		Data: json.RawMessage(`{"test": "test"}`),
	})
	if err != nil {
		t.Fatalf("could not update team member: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_GetTeam_ByName_Team(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team29",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team29", Valid: true},
		},
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if team[0].ID != uint64(teamId) {
		t.Fatalf("expected %d, got %d", team[0].ID, teamId)
	}
	if team[0].Name != "team29" {
		t.Fatalf("expected team29, got %s", team[0].Name)
	}
}

func Test_GetTeam_ById_Team(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team30",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			ID: sql.NullInt64{Int64: teamId, Valid: true},
		},
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if team[0].ID != uint64(teamId) {
		t.Fatalf("expected %d, got %d", team[0].ID, teamId)
	}
	if team[0].Name != "team30" {
		t.Fatalf("expected team30, got %s", team[0].Name)
	}
}

func Test_GetTeam_ByMemberId_Team(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team31",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	result, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       48,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMemberId, err := result.LastInsertId()
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			Member: GetTeamMemberParams{
				ID: sql.NullInt64{Int64: teamMemberId, Valid: true},
			},
		},
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if team[0].ID != uint64(teamId) {
		t.Fatalf("expected %d, got %d", team[0].ID, teamId)
	}
	if team[0].Name != "team31" {
		t.Fatalf("expected team31, got %s", team[0].Name)
	}
}

func Test_GetTeam_ByMemberUserId_Team(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team310",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	result, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       480,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			Member: GetTeamMemberParams{
				UserID: sql.NullInt64{Int64: 480, Valid: true},
			},
		},
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if team[0].ID != uint64(teamId) {
		t.Fatalf("expected %d, got %d", team[0].ID, teamId)
	}
	if team[0].Name != "team310" {
		t.Fatalf("expected team310, got %s", team[0].Name)
	}
}

func Test_GetTeam_ByNameTeamDoesNotExist_ReturnError(t *testing.T) {
	q := New(db)
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team32", Valid: true},
		},
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if len(team) != 0 {
		t.Fatalf("expected 0 teams, got %d", len(team))
	}
}

func Test_GetTeamMembers_ByTeamName_TeamMembers(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team33",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       50,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       51,
		MemberNumber: 2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMembers, err := q.GetTeamMembers(context.Background(), GetTeamMembersParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team33", Valid: true},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get team members: %v", err)
	}
	if len(teamMembers) != 2 {
		t.Fatalf("expected 2 team members, got %d", len(teamMembers))
	}
	if teamMembers[0].UserID != 50 {
		t.Fatalf("expected 50, got %d", teamMembers[0].UserID)
	}
	if teamMembers[1].UserID != 51 {
		t.Fatalf("expected 51, got %d", teamMembers[1].UserID)
	}
}

func Test_GetTeamMembers_ByTeamId_TeamMembers(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team34",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       53,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       54,
		MemberNumber: 2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMembers, err := q.GetTeamMembers(context.Background(), GetTeamMembersParams{
		Team: TeamParams{
			ID: sql.NullInt64{Int64: teamId, Valid: true},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get team members: %v", err)
	}
	if len(teamMembers) != 2 {
		t.Fatalf("expected 2 team members, got %d", len(teamMembers))
	}
	if teamMembers[0].UserID != 53 {
		t.Fatalf("expected 53, got %d", teamMembers[0].UserID)
	}
	if teamMembers[1].UserID != 54 {
		t.Fatalf("expected 54, got %d", teamMembers[1].UserID)
	}
}

func Test_GetTeamMembers_ByTeamMemberId_TeamMembers(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team35",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	result, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       56,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMemberId, err := result.LastInsertId()
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       57,
		MemberNumber: 2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMembers, err := q.GetTeamMembers(context.Background(), GetTeamMembersParams{
		Team: TeamParams{
			Member: GetTeamMemberParams{
				ID: sql.NullInt64{Int64: teamMemberId, Valid: true},
			},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get team members: %v", err)
	}
	if len(teamMembers) != 2 {
		t.Fatalf("expected 2 team members, got %d", len(teamMembers))
	}
	if teamMembers[0].UserID != 56 {
		t.Fatalf("expected 56, got %d", teamMembers[0].UserID)
	}
	if teamMembers[1].UserID != 57 {
		t.Fatalf("expected 57, got %d", teamMembers[1].UserID)
	}
}

func Test_GetTeamMembers_ByTeamMemberUserId_TeamMembers(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team350",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	result, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       560,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       570,
		MemberNumber: 2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMembers, err := q.GetTeamMembers(context.Background(), GetTeamMembersParams{
		Team: TeamParams{
			Member: GetTeamMemberParams{
				UserID: sql.NullInt64{Int64: 560, Valid: true},
			},
		},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get team members: %v", err)
	}
	if len(teamMembers) != 2 {
		t.Fatalf("expected 2 team members, got %d", len(teamMembers))
	}
	if teamMembers[0].UserID != 560 {
		t.Fatalf("expected 56, got %d", teamMembers[0].UserID)
	}
	if teamMembers[1].UserID != 570 {
		t.Fatalf("expected 57, got %d", teamMembers[1].UserID)
	}
}

func Test_GetTeamMembers_ByNameTeamDoesNotExist_ReturnNoMembers(t *testing.T) {
	q := New(db)
	teamMembers, err := q.GetTeamMembers(context.Background(), GetTeamMembersParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team36", Valid: true},
		},
	})
	if err != nil {
		t.Fatalf("could not get team members: %v", err)
	}
	if len(teamMembers) != 0 {
		t.Fatalf("expected 0 team members, got %d", len(teamMembers))
	}
}

func Test_GetTeamMembers_TeamMemberDoesNotExist_ReturnNoMembers(t *testing.T) {
	q := New(db)
	teamMembers, err := q.GetTeamMembers(context.Background(), GetTeamMembersParams{
		Team: TeamParams{
			Member: GetTeamMemberParams{
				ID: sql.NullInt64{Int64: 9999999, Valid: true},
			},
		},
	})
	if err != nil {
		t.Fatalf("could not get team members: %v", err)
	}
	if len(teamMembers) != 0 {
		t.Fatalf("expected 0 team members, got %d", len(teamMembers))
	}
}

func Test_DeleteTeam_ByName_TeamDeleted(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team37",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err := q.DeleteTeam(context.Background(), TeamParams{Name: sql.NullString{String: "team37", Valid: true}})
	if err != nil {
		t.Fatalf("could not delete team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	// check if team was deleted
	_, err = q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team37", Valid: true},
		},
		Limit:  1,
		Offset: 0,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_DeleteTeam_ById_TeamDeleted(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team38",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	result, err = q.DeleteTeam(context.Background(), TeamParams{
		ID: sql.NullInt64{Int64: teamId, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not delete team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	// check if team was deleted
	_, err = q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team38", Valid: true},
		},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_DeleteTeam_ByMemberId_TeamDeleted(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team39",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	result, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       63,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMemberId, err := result.LastInsertId()
	result, err = q.DeleteTeam(context.Background(), TeamParams{
		Member: GetTeamMemberParams{
			ID: sql.NullInt64{Int64: teamMemberId, Valid: true},
		},
	})
	if err != nil {
		t.Fatalf("could not delete team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	// check if team was deleted
	_, err = q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team39", Valid: true},
		},
		Limit:  1,
		Offset: 0,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	// check if team member was deleted
	_, err = q.GetTeamMember(context.Background(), GetTeamMemberParams{
		ID: sql.NullInt64{Int64: teamMemberId, Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_DeleteTeam_ByMemberUserId_TeamDeleted(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team390",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	result, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       630,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMemberId, err := result.LastInsertId()
	result, err = q.DeleteTeam(context.Background(), TeamParams{
		Member: GetTeamMemberParams{
			UserID: sql.NullInt64{Int64: 630, Valid: true},
		},
	})
	if err != nil {
		t.Fatalf("could not delete team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	// check if team was deleted
	_, err = q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team390", Valid: true},
		},
		Limit:  1,
		Offset: 0,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	// check if team member was deleted
	_, err = q.GetTeamMember(context.Background(), GetTeamMemberParams{
		ID: sql.NullInt64{Int64: teamMemberId, Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_DeleteTeam_TeamDoesNotExist_TeamNotDeleted(t *testing.T) {
	q := New(db)
	result, err := q.DeleteTeam(context.Background(), TeamParams{Name: sql.NullString{String: "team40", Valid: true}})
	if err != nil {
		t.Fatalf("could not delete team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func Test_UpdateTeam_UpdateDataByTeamName_TeamUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team41",
		Score: 0,
		Data:  json.RawMessage(`{"test": "test"}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err := q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team: TeamParams{Name: sql.NullString{String: "team41", Valid: true}},
		Data: json.RawMessage(`{"test": "test2"}`),
	})
	if err != nil {
		t.Fatalf("could not update team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team41", Valid: true},
		},
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if !reflect.DeepEqual(team[0].Data, json.RawMessage(`{"test": "test2"}`)) {
		t.Fatalf("expected {\"test\": \"test2\"}, got %s", team[0].Data)
	}
}

func Test_UpdateTeam_UpdateDataByTeamId_TeamUpdated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team42",
		Score: 0,
		Data:  json.RawMessage(`{"test": "test"}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	result, err = q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team: TeamParams{
			ID: sql.NullInt64{Int64: teamId, Valid: true},
		},
		Data: json.RawMessage(`{"test": "test2"}`),
	})
	if err != nil {
		t.Fatalf("could not update team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team42", Valid: true},
		},
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if !reflect.DeepEqual(team[0].Data, json.RawMessage(`{"test": "test2"}`)) {
		t.Fatalf("expected {\"test\": \"test2\"}, got %s", team[0].Data)
	}
}

func Test_UpdateTeam_UpdateDataByTeamMemberId_TeamUpdated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team43",
		Score: 0,
		Data:  json.RawMessage(`{"test": "test"}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       67,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMemberId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	result, err = q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team: TeamParams{
			Member: GetTeamMemberParams{
				ID: sql.NullInt64{Int64: teamMemberId, Valid: true},
			},
		},
		Data: json.RawMessage(`{"test": "test2"}`),
	})
	if err != nil {
		t.Fatalf("could not update team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team43", Valid: true},
		},
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if !reflect.DeepEqual(team[0].Data, json.RawMessage(`{"test": "test2"}`)) {
		t.Fatalf("expected {\"test\": \"test2\"}, got %s", team[0].Data)
	}
}

func Test_UpdateTeam_UpdateDataByTeamMemberUserId_TeamUpdated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team430",
		Score: 0,
		Data:  json.RawMessage(`{"test": "test"}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       670,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	result, err = q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team: TeamParams{
			Member: GetTeamMemberParams{
				UserID: sql.NullInt64{Int64: 670, Valid: true},
			},
		},
		Data: json.RawMessage(`{"test": "test2"}`),
	})
	if err != nil {
		t.Fatalf("could not update team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team430", Valid: true},
		},
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if !reflect.DeepEqual(team[0].Data, json.RawMessage(`{"test": "test2"}`)) {
		t.Fatalf("expected {\"test\": \"test2\"}, got %s", team[0].Data)
	}
}

func Test_UpdateTeam_UpdateScoreByTeamName_TeamUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team44",
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err := q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team44", Valid: true},
		},
		Score: sql.NullInt64{Int64: 1, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not update team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team44", Valid: true},
		},
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if team[0].Score != 1 {
		t.Fatalf("expected 1, got %d", team[0].Score)
	}
}

func Test_UpdateTeam_IncrementScoreByTeamName_TeamUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team45",
		Score: 10,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err := q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team:           TeamParams{Name: sql.NullString{String: "team45", Valid: true}},
		Score:          sql.NullInt64{Int64: 10, Valid: true},
		IncrementScore: true,
	})
	if err != nil {
		t.Fatalf("could not update team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Team:   TeamParams{Name: sql.NullString{String: "team45", Valid: true}},
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if team[0].Score != 20 {
		t.Fatalf("expected 20, got %d", team[0].Score)
	}
}

func Test_UpdateTeam_DecrementScoreByTeamName_TeamUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team46",
		Score: 10,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err := q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team:           TeamParams{Name: sql.NullString{String: "team46", Valid: true}},
		Score:          sql.NullInt64{Int64: -10, Valid: true},
		IncrementScore: true,
	})
	if err != nil {
		t.Fatalf("could not update team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Team:   TeamParams{Name: sql.NullString{String: "team46", Valid: true}},
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if team[0].Score != 0 {
		t.Fatalf("expected 0, got %d", team[0].Score)
	}
}

func Test_UpdateTeam_UpdateDataAndScoreByTeamMemberId_TeamUpdated(t *testing.T) {
	q := New(db)
	result, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team47",
		Score: 0,
		Data:  json.RawMessage(`{"test": "test"}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	teamId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		TeamID:       uint64(teamId),
		UserID:       72,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMemberId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	result, err = q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team: TeamParams{
			Member: GetTeamMemberParams{
				ID: sql.NullInt64{Int64: teamMemberId, Valid: true},
			},
		},
		Data:  json.RawMessage(`{"test": "test2"}`),
		Score: sql.NullInt64{Int64: 1, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not update team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Team: TeamParams{
			Name: sql.NullString{String: "team47", Valid: true},
		},
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if !reflect.DeepEqual(team[0].Data, json.RawMessage(`{"test": "test2"}`)) {
		t.Fatalf("expected {\"test\": \"test2\"}, got %s", team[0].Data)
	}
	if team[0].Score != 1 {
		t.Fatalf("expected 1, got %d", team[0].Score)
	}
}

func Test_UpdateTeam_TeamDoesNotExist_TeamNotUpdated(t *testing.T) {
	q := New(db)
	result, err := q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team: TeamParams{Name: sql.NullString{String: "team48", Valid: true}},
		Data: json.RawMessage(`{"test": "test"}`),
	})
	if err != nil {
		t.Fatalf("could not update team: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}
