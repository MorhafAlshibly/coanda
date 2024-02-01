package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcodes"
	"github.com/MorhafAlshibly/coanda/pkg/file"
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	dbName  = "testdb"
	address = "localhost"
	port    = 3306
)

func TestMain(m *testing.M) {
	engine := sqle.NewDefault(
		memory.NewDBProvider(
			func() *memory.Database {
				db := memory.NewDatabase(dbName)
				db.EnablePrimaryKeyIndexes()
				return db
			}(),
		))
	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", address, port),
	}
	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}
	go func() {
		if err = s.Start(); err != nil {
			panic(err)
		}
	}()
	// Parse the schema file
	root := file.FindMyRootDir()
	schema, err := os.ReadFile(fmt.Sprintf("%smigration/team.sql", root))
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("mysql", fmt.Sprintf("root@tcp(%s:%d)/%s?multiStatements=true&parseTime=true", address, port, dbName))
	if err != nil {
		panic(err)
	}
	res, err := db.Exec(string(schema))
	if err != nil {
		panic(err)
	}
	fmt.Println("Applied schema:", res)
	os.Exit(m.Run())
}

func TestCreateTeam(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root@tcp(%s:%d)/%s?multiStatements=true&parseTime=true", address, port, dbName))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := New(db)
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	result, err := queries.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team",
		Owner: 1,
		Score: 0,
		Data:  raw,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}
}

func TestCreateTeamNameTaken(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root@tcp(%s:%d)/%s?multiStatements=true&parseTime=true", address, port, dbName))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := New(db)
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	_, err = queries.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team",
		Owner: 1,
		Score: 0,
		Data:  raw,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = queries.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team",
		Owner: 1,
		Score: 0,
		Data:  raw,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		t.Fatal("expected mysql error")
	}
	if mysqlErr.Number != errorcodes.MySQLErrorCodeDuplicateEntry {
		t.Fatal("expected duplicate entry error")
	}
}

func TestCreateTeamOwnerTaken(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root@tcp(%s:%d)/%s?multiStatements=true&parseTime=true", address, port, dbName))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := New(db)
	data, err := conversion.MapToProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := conversion.ProtobufStructToRawJson(data)
	if err != nil {
		t.Fatal(err)
	}
	_, err = queries.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team",
		Owner: 1,
		Score: 0,
		Data:  raw,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = queries.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team2",
		Owner: 1,
		Score: 0,
		Data:  raw,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		t.Fatal("expected mysql error")
	}
	if mysqlErr.Number != errorcodes.MySQLErrorCodeDuplicateEntry {
		t.Fatal("expected duplicate entry error")
	}
}

func TestCreateTeamMember(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root@tcp(%s:%d)/%s?multiStatements=true&parseTime=true", address, port, dbName))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := New(db)
	raw, err := conversion.ProtobufStructToRawJson(&structpb.Struct{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = queries.CreateTeam(context.Background(), CreateTeamParams{
		Name:  "team",
		Owner: 1,
		Score: 0,
		Data:  raw,
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := queries.CreateTeamMember(context.Background(), CreateTeamMemberParams{
		Team:       "team",
		UserID:     1,
		Data:       raw,
		MaxMembers: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}
}
