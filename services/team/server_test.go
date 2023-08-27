package team

import (
	"context"
	"testing"
)

func TestCreateTeam(t *testing.T) {
	db := &database.MockDatabase{
		InsertOneFunc: func(ctx context.Context, document interface{}) (string, error) {
			return "test", nil
		},
	}
	service := server{
		db:    db,
		queue: nil,
		cache: nil,
	}
}
