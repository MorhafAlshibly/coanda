package services

import (
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/storage"
)

type TeamService struct {
	store storage.Storer
	cache cache.Cacher
}

func NewTeamService(store storage.Storer, cache cache.Cacher) *TeamService {
	return &TeamService{
		store: store,
		cache: cache,
	}
}
