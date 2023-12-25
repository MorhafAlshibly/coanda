package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertOneInput struct {
	Document interface{}
}

type InsertOneOutput struct {
	InsertedID primitive.ObjectID
}

type InsertOneCommand struct {
	database *MongoDatabase
	In       *InsertOneInput
	Out      *InsertOneOutput
}

func NewInsertOneCommand(database *MongoDatabase, in *InsertOneInput) *InsertOneCommand {
	return &InsertOneCommand{
		database: database,
		In:       in,
	}
}

func (c *InsertOneCommand) Execute(ctx context.Context) error {
	result, err := c.database.collection.InsertOne(ctx, c.In.Document)
	if err != nil {
		return err
	}
	c.Out = &InsertOneOutput{
		InsertedID: result.InsertedID.(primitive.ObjectID),
	}
	return nil
}
