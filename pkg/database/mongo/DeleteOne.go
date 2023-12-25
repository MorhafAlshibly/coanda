package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type DeleteOneInput struct {
	Filter interface{}
}

type DeleteOneOutput struct {
	Result *mongo.DeleteResult
}

type DeleteOneCommand struct {
	database *MongoDatabase
	In       *DeleteOneInput
	Out      *DeleteOneOutput
}

func NewDeleteOneCommand(database *MongoDatabase, in *DeleteOneInput) *DeleteOneCommand {
	return &DeleteOneCommand{
		database: database,
		In:       in,
	}
}

func (c *DeleteOneCommand) Execute(ctx context.Context) error {
	result, err := c.database.collection.DeleteOne(ctx, c.In.Filter)
	if err != nil {
		return err
	}
	c.Out = &DeleteOneOutput{
		Result: result,
	}
	return nil
}
