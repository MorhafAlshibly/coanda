package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type UpdateOneInput struct {
	Filter interface{}
	Update interface{}
}

type UpdateOneOutput struct {
	Result *mongo.UpdateResult
}

type UpdateOneCommand struct {
	database *MongoDatabase
	In       *UpdateOneInput
	Out      *UpdateOneOutput
}

func NewUpdateOneCommand(database *MongoDatabase, in *UpdateOneInput) *UpdateOneCommand {
	return &UpdateOneCommand{
		database: database,
		In:       in,
	}
}

func (c *UpdateOneCommand) Execute(ctx context.Context) error {
	result, err := c.database.collection.UpdateOne(ctx, c.In.Filter, c.In.Update)
	if err != nil {
		return err
	}
	c.Out = &UpdateOneOutput{
		Result: result,
	}
	return nil
}
