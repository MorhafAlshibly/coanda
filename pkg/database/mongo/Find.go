package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FindInput struct {
	Filter  interface{}
	Options *options.FindOptions
}

type FindOutput struct {
	Result *[]bson.M
}

type FindCommand struct {
	database *MongoDatabase
	In       *FindInput
	Out      *FindOutput
}

func NewFindCommand(database *MongoDatabase, in *FindInput) *FindCommand {
	return &FindCommand{
		database: database,
		In:       in,
	}
}

func (c *FindCommand) Execute(ctx context.Context) error {
	cursor, err := c.database.collection.Find(ctx, c.In.Filter, c.In.Options)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	var result []bson.M
	err = cursor.All(ctx, &result)
	if err != nil {
		return err
	}
	c.Out = &FindOutput{
		Result: &result,
	}
	return nil
}
