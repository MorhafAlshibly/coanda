package mysqlTestServer

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
)

var (
	serverLock = &sync.Mutex{}
	server     *Server
)

type Server struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
	db       *sql.DB
}

type Connection struct {
	Db *sql.DB
}

func GetServer() (*Server, error) {
	if server == nil {
		serverLock.Lock()
		defer serverLock.Unlock()
		if server == nil {
			s, err := run()
			if err != nil {
				return nil, err
			}
			server = s
		}
	}
	return server, nil
}

func (s *Server) GetConnection() *Connection {
	if s.pool == nil {
		log.Fatalf("pool is nil")
	}
	db, err := sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql?multiStatements=true&parseTime=true", s.resource.GetPort("3306/tcp")))
	if err != nil {
		log.Fatalf("could not open connection: %s", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping connection: %s", err)
	}
	return &Connection{Db: db}
}

func (c *Connection) GetTransaction() *sql.Tx {
	tx, err := c.Db.Begin()
	if err != nil {
		log.Fatalf("could not begin transaction: %s", err)
	}
	return tx
}

func (s *Server) Connect(t *testing.T) *sql.Tx {
	conn := s.GetConnection()
	t.Cleanup(func() { _ = conn.Db.Close() })
	tx := conn.GetTransaction()
	t.Cleanup(func() {
		err := tx.Rollback()
		if err != nil {
			t.Fatalf("could not rollback transaction: %s", err)
		}
	})
	return tx
}

func NewServer(schemaPath string) *Server {
	var err error
	s, err := GetServer()
	if err != nil {
		log.Fatalf("could not run mysql test server: %v", err)
	}
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		s.Close()
		log.Fatalf("could not read schema file: %v", err)
	}
	conn := s.GetConnection()
	_, err = conn.Db.Exec(string(schema))
	conn.Db.Close()
	if err != nil {
		s.Close()
		log.Fatalf("could not execute schema: %v", err)
	}
	return s
}

func run() (*Server, error) {
	s := &Server{}
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not construct pool: %s", err)
	}
	s.pool = pool

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not ping Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "8.0", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %s", err)
	}
	s.resource = resource
	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err := sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql?multiStatements=true&parseTime=true", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		s.db = db
		time.Sleep(30 * time.Second)
		return db.Ping()
	}); err != nil {
		return nil, fmt.Errorf("could not connect to docker: %s", err)
	}
	return s, nil
}

func (s *Server) Close() {
	if err := s.pool.Purge(s.resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}
}
