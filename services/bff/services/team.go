package services

import (
	"github.com/MorhafAlshibly/coanda/libs/cache"
	"github.com/MorhafAlshibly/coanda/libs/storage"
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
