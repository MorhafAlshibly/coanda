// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreateItem struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type GetItem struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type GetItems struct {
	Type *string `json:"type,omitempty"`
	Max  *int    `json:"max,omitempty"`
	Page *int    `json:"page,omitempty"`
}

type Item struct {
	ID   string                 `json:"id"`
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}
