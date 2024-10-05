package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"github.com/MorhafAlshibly/coanda/internal/bff"
)

// Mutation returns bff.MutationResolver implementation.
func (r *Resolver) Mutation() bff.MutationResolver { return &mutationResolver{r} }

// Query returns bff.QueryResolver implementation.
func (r *Resolver) Query() bff.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
