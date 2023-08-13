// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type CreateItem struct {
	// The type of the item. Can be any string.
	Type string `json:"type"`
	// The data associated with the item. Can be any JSON object.
	Data map[string]interface{} `json:"data"`
	// The timestamp of when the item will expire.
	Expire *time.Time `json:"expire,omitempty"`
}

type CreateTeam struct {
	Name  string                 `json:"name"`
	Score *int                   `json:"score,omitempty"`
	Data  map[string]interface{} `json:"data,omitempty"`
}

type DeleteTeam struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type GetItem struct {
	// The unique identifier for the item you want to retrieve.
	ID string `json:"id"`
	// The type of the item you want to retrieve.
	Type string `json:"type"`
}

type GetItems struct {
	// The type of the items you want to retrieve.
	Type *string `json:"type,omitempty"`
	// The maximum number of items to retrieve.
	Max *int `json:"max,omitempty"`
	// The page number of the items specified by the max parameter.
	Page *int `json:"page,omitempty"`
}

type GetTeam struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type GetTeams struct {
	Max  *int `json:"max,omitempty"`
	Page *int `json:"page,omitempty"`
}

// An item is a generic object that can be created, read, updated, and deleted.
type Item struct {
	// The unique identifier for the item.
	ID string `json:"id"`
	// The type of the item. Used in partitioning.
	Type string `json:"type"`
	// The data associated with the item. Can be any JSON object.
	Data map[string]interface{} `json:"data"`
	// The timestamp of when the item will expire.
	Expire *time.Time `json:"expire,omitempty"`
}

type JoinTeam struct {
	ID     *string `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	UserID string  `json:"userId"`
}

type LeaveTeam struct {
	ID     *string `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	UserID string  `json:"userId"`
}

type QueuedTeam struct {
	Name    string                 `json:"name"`
	Members []*string              `json:"members"`
	Score   int                    `json:"score"`
	Rank    int                    `json:"rank"`
	Data    map[string]interface{} `json:"data"`
}

type SearchTeams struct {
	Name string `json:"name"`
}

type Team struct {
	ID      string                 `json:"id"`
	Name    string                 `json:"name"`
	Members []*string              `json:"members"`
	Score   int                    `json:"score"`
	Rank    int                    `json:"rank"`
	Data    map[string]interface{} `json:"data"`
}

type UpdateTeamData struct {
	ID   *string                `json:"id,omitempty"`
	Name *string                `json:"name,omitempty"`
	Data map[string]interface{} `json:"data"`
}

type UpdateTeamScore struct {
	ID    *string `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`
	Score int     `json:"score"`
}
