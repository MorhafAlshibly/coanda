package mysqlTestServer

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

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
	Db       *sql.DB
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

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err := sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql?multiStatements=true&parseTime=true", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(5)
		s.Db = db
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
