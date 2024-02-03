package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcodes"
	"github.com/MorhafAlshibly/coanda/pkg/file"
	"github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	user     = "root"
	password = "secret"
	address  = "localhost"
	port     = 3306
)

var db *sql.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "8.0", []string{"MYSQL_ROOT_PASSWORD=" + password})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql?multiStatements=true&parseTime=true", user, password, address, port))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	// apply schema
	root := file.FindMyRootDir()
	schema, err := os.ReadFile(fmt.Sprintf("%smigration/team.sql", root))
	if err != nil {
		panic(err)
	}
	res, err := db.Exec(string(schema))
	if err != nil {
		panic(err)
	}
	fmt.Println("Applied schema: ", res)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	db.Close()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestCreateTeam(t *testing.T) {
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
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql?multiStatements=true&parseTime=true", user, password, address, port))
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
