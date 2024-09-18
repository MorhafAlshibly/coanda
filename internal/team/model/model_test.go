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
		Owner: 1,
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
		Owner: 2,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team1",
		Owner: 3,
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
	if errorcode.IsDuplicateEntry(mysqlErr, "team", "name") {
		t.Fatalf("expected duplicate entry error, got %d", mysqlErr.Number)
	}
}

func Test_CreateTeam_TeamOwnerExists_TeamNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team2",
		Owner: 4,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team3",
		Owner: 4,
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
	if !errorcode.IsDuplicateEntry(mysqlErr, "team", "owner") {
		t.Fatalf("expected duplicate entry error, got %d", mysqlErr.Number)
	}
}

func Test_CreateTeamOwner_TeamOwner_TeamOwnerCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team4",
		Owner: 5,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err := q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team4",
		UserID:       5,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	result, err = q.CreateTeamOwner(context.Background(), CreateTeamOwnerParams{
		Team:   "team4",
		UserID: 5,
	})
	if err != nil {
		t.Fatalf("could not create team owner: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("could not get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}
}

func Test_CreateTeamOwner_TeamOwnerExists_TeamOwnerNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team5",
		Owner: 6,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team5",
		UserID:       6,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamOwner(context.Background(), CreateTeamOwnerParams{
		Team:   "team5",
		UserID: 6,
	})
	if err != nil {
		t.Fatalf("could not create team owner: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team6",
		Owner: 7,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamOwner(context.Background(), CreateTeamOwnerParams{
		Team:   "team6",
		UserID: 6,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if !errorcode.IsDuplicateEntry(mysqlErr, "team_owner", "user_id") {
		t.Fatalf("expected duplicate entry error, got %d", mysqlErr.Number)
	}
}

func Test_CreateTeamOwner_TeamAlreadyHasOwner_TeamOwnerNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team7",
		Owner: 8,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team7",
		UserID:       8,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team7",
		UserID:       9,
		MemberNumber: 2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamOwner(context.Background(), CreateTeamOwnerParams{
		Team:   "team7",
		UserID: 8,
	})
	if err != nil {
		t.Fatalf("could not create team owner: %v", err)
	}
	_, err = q.CreateTeamOwner(context.Background(), CreateTeamOwnerParams{
		Team:   "team7",
		UserID: 9,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		t.Fatalf("expected mysql error, got %v", err)
	}
	if errorcode.IsDuplicateEntry(mysqlErr, "team_owner", "team") {
		t.Fatalf("expected duplicate entry error, got %d", mysqlErr.Number)
	}
}

func Test_CreateTeamOwner_TeamDoesNotExist_TeamOwnerNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeamOwner(context.Background(), CreateTeamOwnerParams{
		Team:   "team8",
		UserID: 10,
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

func Test_CreateTeamMember_TeamMember_TeamMemberCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team9",
		Owner: 11,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err := q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team9",
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
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team10",
		Owner: 12,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team10",
		UserID:       12,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team10",
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
	if errorcode.IsDuplicateEntry(mysqlErr, "team_member", "user_id") {
		t.Fatalf("expected duplicate entry error, got %d", mysqlErr.Number)
	}
}

func Test_CreateTeamMember_TeamMemberNumberExists_TeamMemberNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team11",
		Owner: 13,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team11",
		UserID:       13,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team11",
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
	if errorcode.IsDuplicateEntry(mysqlErr, "team_member", "member_number") {
		t.Fatalf("expected duplicate entry error, got %d", mysqlErr.Number)
	}
}

func Test_CreateTeamMember_TeamDoesNotExist_TeamMemberNotCreated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team12",
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
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team13",
		Owner: 16,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team14",
		Owner: 17,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team13",
		UserID:       18,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team14",
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

func Test_DeleteTeamMember_TeamMemberExists_TeamMemberDeleted(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team15",
		Owner: 19,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team15",
		UserID:       20,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	result, err := q.DeleteTeamMember(context.Background(), 20)
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

func Test_DeleteTeamMember_TeamMemberDoesNotExist_TeamMemberNotDeleted(t *testing.T) {
	q := New(db)
	result, err := q.DeleteTeamMember(context.Background(), 21)
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

func Test_GetHighestMemberNumber_TeamNoMembers_ReturnNil(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team16",
		Owner: 22,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.GetHighestMemberNumber(context.Background(), "team16")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetHighestMemberNumber_TeamHasMembers_ReturnHighestMemberNumber(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team17",
		Owner: 23,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team17",
		UserID:       24,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team17",
		UserID:       25,
		MemberNumber: 2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	memberNumber, err := q.GetHighestMemberNumber(context.Background(), "team17")
	if err != nil {
		t.Fatalf("could not get highest member number: %v", err)
	}
	if memberNumber != 2 {
		t.Fatalf("expected 2, got %d", memberNumber)
	}
}

func Test_GetHighestMemberNumber_TeamDoesNotExist_ReturnError(t *testing.T) {
	q := New(db)
	_, err := q.GetHighestMemberNumber(context.Background(), "team18")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetTeamMember_TeamMemberExists_ReturnTeamMember(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team19",
		Owner: 26,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team19",
		UserID:       27,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMember, err := q.GetTeamMember(context.Background(), 27)
	if err != nil {
		t.Fatalf("could not get team member: %v", err)
	}
	if teamMember.Team != "team19" {
		t.Fatalf("expected team19, got %s", teamMember.Team)
	}
	if teamMember.UserID != 27 {
		t.Fatalf("expected 27, got %d", teamMember.UserID)
	}
	if teamMember.MemberNumber != 1 {
		t.Fatalf("expected 1, got %d", teamMember.MemberNumber)
	}
}

func Test_GetTeamMember_TeamMemberDoesNotExist_ReturnNil(t *testing.T) {
	q := New(db)
	_, err := q.GetTeamMember(context.Background(), 28)
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
		Owner: 29,
		Score: 1,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team21",
		Owner: 30,
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
		Owner: 32,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "teamwithspecialwordinthemiddletest",
		Owner: 33,
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
		Owner: 34,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "teamwithspecialwordatend",
		Owner: 35,
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
		Owner: 36,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "specialwordatstartteam",
		Owner: 37,
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
		Owner: 38,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "specialwordteam",
		Owner: 39,
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
		Owner: 40,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "teamwithspecialwordinthemiddlecaseinsensitivetest",
		Owner: 41,
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

func Test_UpdateTeamMember_TeamMemberExists_TeamMemberUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team28",
		Owner: 42,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team28",
		UserID:       43,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	result, err := q.UpdateTeamMember(context.Background(), UpdateTeamMemberParams{
		UserID: 43,
		Data:   json.RawMessage(`{"test": "test"}`),
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
	team, err := q.GetTeamMember(context.Background(), 43)
	if err != nil {
		t.Fatalf("could not get team member: %v", err)
	}
	if !reflect.DeepEqual(team.Data, json.RawMessage(`{"test": "test"}`)) {
		t.Fatalf("expected {\"test\": \"test\"}, got %s", team.Data)
	}
}

func Test_UpdateTeamMember_TeamMemberDoesNotExist_TeamMemberNotUpdated(t *testing.T) {
	q := New(db)
	result, err := q.UpdateTeamMember(context.Background(), UpdateTeamMemberParams{
		UserID: 44,
		Data:   json.RawMessage(`{"test": "test"}`),
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
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team29",
		Owner: 45,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Name: sql.NullString{String: "team29", Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if team.Name != "team29" {
		t.Fatalf("expected team29, got %s", team.Name)
	}
	if team.Owner != 45 {
		t.Fatalf("expected 45, got %d", team.Owner)
	}
}

func Test_GetTeam_ByOwnerId_Team(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team30",
		Owner: 46,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Owner: sql.NullInt64{Int64: 46, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if team.Name != "team30" {
		t.Fatalf("expected team30, got %s", team.Name)
	}
	if team.Owner != 46 {
		t.Fatalf("expected 46, got %d", team.Owner)
	}
}

func Test_GetTeam_ByMemberId_Team(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team31",
		Owner: 47,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team31",
		UserID:       48,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	team, err := q.GetTeam(context.Background(), GetTeamParams{
		Member: sql.NullInt64{Int64: 48, Valid: true},
	})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if team.Name != "team31" {
		t.Fatalf("expected team31, got %s", team.Name)
	}
	if team.Owner != 47 {
		t.Fatalf("expected 47, got %d", team.Owner)
	}
}

func Test_GetTeam_TeamDoesNotExist_ReturnError(t *testing.T) {
	q := New(db)
	_, err := q.GetTeam(context.Background(), GetTeamParams{
		Name: sql.NullString{String: "team32", Valid: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected no rows error, got %v", err)
	}
}

func Test_GetTeamMembers_ByTeamName_TeamMembers(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team33",
		Owner: 49,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team33",
		UserID:       50,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team33",
		UserID:       51,
		MemberNumber: 2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMembers, err := q.GetTeamMembers(context.Background(), GetTeamMembersParams{
		Team:   GetTeamParams{Name: sql.NullString{String: "team33", Valid: true}},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get team members: %v", err)
	}
	if len(teamMembers) != 2 {
		t.Fatalf("expected 2 team members, got %d", len(teamMembers))
	}
}

func Test_GetTeamMembers_ByTeamOwner_TeamMembers(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team34",
		Owner: 52,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team34",
		UserID:       53,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team34",
		UserID:       54,
		MemberNumber: 2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMembers, err := q.GetTeamMembers(context.Background(), GetTeamMembersParams{
		Team:   GetTeamParams{Owner: sql.NullInt64{Int64: 52, Valid: true}},
		Limit:  2,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("could not get team members: %v", err)
	}
	if len(teamMembers) != 2 {
		t.Fatalf("expected 2 team members, got %d", len(teamMembers))
	}
}

func Test_GetTeamMembers_ByTeamMember_TeamMembers(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team35",
		Owner: 55,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team35",
		UserID:       56,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team35",
		UserID:       57,
		MemberNumber: 2,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	teamMembers, err := q.GetTeamMembers(context.Background(), GetTeamMembersParams{
		Team:   GetTeamParams{Member: sql.NullInt64{Int64: 56, Valid: true}},
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

func Test_GetTeamMembers_TeamDoesNotExist_ReturnError(t *testing.T) {
	q := New(db)
	teamMembers, err := q.GetTeamMembers(context.Background(), GetTeamMembersParams{
		Team: GetTeamParams{Name: sql.NullString{String: "team36", Valid: true}},
	})
	if err != nil {
		t.Fatalf("could not get team members: %v", err)
	}
	if len(teamMembers) != 0 {
		t.Fatalf("expected 0 team members, got %d", len(teamMembers))
	}
}

func Test_GetTeamMembers_TeamMemberDoesNotExist_ReturnError(t *testing.T) {
	q := New(db)
	teamMembers, err := q.GetTeamMembers(context.Background(), GetTeamMembersParams{
		Team: GetTeamParams{Member: sql.NullInt64{Int64: 58, Valid: true}},
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
		Owner: 59,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err := q.DeleteTeam(context.Background(), GetTeamParams{Name: sql.NullString{String: "team37", Valid: true}})
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
	_, err = q.GetTeam(context.Background(), GetTeamParams{Name: sql.NullString{String: "team37", Valid: true}})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_DeleteTeam_ByOwner_TeamDeleted(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team38",
		Owner: 60,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamOwner(context.Background(), CreateTeamOwnerParams{
		Team:   "team38",
		UserID: 61,
	})
	result, err := q.DeleteTeam(context.Background(), GetTeamParams{Owner: sql.NullInt64{Int64: 60, Valid: true}})
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
	_, err = q.GetTeam(context.Background(), GetTeamParams{Name: sql.NullString{String: "team38", Valid: true}})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_DeleteTeam_ByMember_TeamDeleted(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team39",
		Owner: 62,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team39",
		UserID:       63,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	result, err := q.DeleteTeam(context.Background(), GetTeamParams{Member: sql.NullInt64{Int64: 63, Valid: true}})
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
	_, err = q.GetTeam(context.Background(), GetTeamParams{Name: sql.NullString{String: "team39", Valid: true}})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	// check if team member was deleted
	_, err = q.GetTeamMember(context.Background(), 63)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_DeleteTeam_TeamDoesNotExist_TeamNotDeleted(t *testing.T) {
	q := New(db)
	result, err := q.DeleteTeam(context.Background(), GetTeamParams{Name: sql.NullString{String: "team40", Valid: true}})
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
		Owner: 64,
		Score: 0,
		Data:  json.RawMessage(`{"test": "test"}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err := q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team: GetTeamParams{Name: sql.NullString{String: "team41", Valid: true}},
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
	team, err := q.GetTeam(context.Background(), GetTeamParams{Name: sql.NullString{String: "team41", Valid: true}})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if !reflect.DeepEqual(team.Data, json.RawMessage(`{"test": "test2"}`)) {
		t.Fatalf("expected {\"test\": \"test2\"}, got %s", team.Data)
	}
}

func Test_UpdateTeam_UpdateDataByTeamOwner_TeamUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team42",
		Owner: 65,
		Score: 0,
		Data:  json.RawMessage(`{"test": "test"}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err := q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team: GetTeamParams{Owner: sql.NullInt64{Int64: 65, Valid: true}},
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
	team, err := q.GetTeam(context.Background(), GetTeamParams{Name: sql.NullString{String: "team42", Valid: true}})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if !reflect.DeepEqual(team.Data, json.RawMessage(`{"test": "test2"}`)) {
		t.Fatalf("expected {\"test\": \"test2\"}, got %s", team.Data)
	}
}

func Test_UpdateTeam_UpdateDataByTeamMember_TeamUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team43",
		Owner: 66,
		Score: 0,
		Data:  json.RawMessage(`{"test": "test"}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team43",
		UserID:       67,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	result, err := q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team: GetTeamParams{Member: sql.NullInt64{Int64: 67, Valid: true}},
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
	team, err := q.GetTeam(context.Background(), GetTeamParams{Name: sql.NullString{String: "team43", Valid: true}})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if !reflect.DeepEqual(team.Data, json.RawMessage(`{"test": "test2"}`)) {
		t.Fatalf("expected {\"test\": \"test2\"}, got %s", team.Data)
	}
}

func Test_UpdateTeam_UpdateScoreByTeamName_TeamUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team44",
		Owner: 68,
		Score: 0,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err := q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team:  GetTeamParams{Name: sql.NullString{String: "team44", Valid: true}},
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
	team, err := q.GetTeam(context.Background(), GetTeamParams{Name: sql.NullString{String: "team44", Valid: true}})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if team.Score != 1 {
		t.Fatalf("expected 1, got %d", team.Score)
	}
}

func Test_UpdateTeam_IncrementScoreByTeamName_TeamUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team45",
		Owner: 69,
		Score: 10,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err := q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team:           GetTeamParams{Name: sql.NullString{String: "team45", Valid: true}},
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
	team, err := q.GetTeam(context.Background(), GetTeamParams{Name: sql.NullString{String: "team45", Valid: true}})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if team.Score != 20 {
		t.Fatalf("expected 20, got %d", team.Score)
	}
}

func Test_UpdateTeam_DecrementScoreByTeamName_TeamUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team46",
		Owner: 70,
		Score: 10,
		Data:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	result, err := q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team:           GetTeamParams{Name: sql.NullString{String: "team46", Valid: true}},
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
	team, err := q.GetTeam(context.Background(), GetTeamParams{Name: sql.NullString{String: "team46", Valid: true}})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if team.Score != 0 {
		t.Fatalf("expected 0, got %d", team.Score)
	}
}

func Test_UpdateTeam_UpdateDataAndScoreByTeamMember_TeamUpdated(t *testing.T) {
	q := New(db)
	_, err := q.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team47",
		Owner: 71,
		Score: 0,
		Data:  json.RawMessage(`{"test": "test"}`),
	})
	if err != nil {
		t.Fatalf("could not create team: %v", err)
	}
	_, err = q.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:         "team47",
		UserID:       72,
		MemberNumber: 1,
		Data:         json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("could not create team member: %v", err)
	}
	result, err := q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team:  GetTeamParams{Member: sql.NullInt64{Int64: 72, Valid: true}},
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
	team, err := q.GetTeam(context.Background(), GetTeamParams{Name: sql.NullString{String: "team47", Valid: true}})
	if err != nil {
		t.Fatalf("could not get team: %v", err)
	}
	if !reflect.DeepEqual(team.Data, json.RawMessage(`{"test": "test2"}`)) {
		t.Fatalf("expected {\"test\": \"test2\"}, got %s", team.Data)
	}
	if team.Score != 1 {
		t.Fatalf("expected 1, got %d", team.Score)
	}
}

func Test_UpdateTeam_TeamDoesNotExist_TeamNotUpdated(t *testing.T) {
	q := New(db)
	result, err := q.UpdateTeam(context.Background(), UpdateTeamParams{
		Team: GetTeamParams{Name: sql.NullString{String: "team48", Valid: true}},
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
