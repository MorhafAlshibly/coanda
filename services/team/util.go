package team

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/services/team/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getFilter(input *schema.GetTeamRequest) (bson.D, error) {
	if input.Id != nil {
		id, err := primitive.ObjectIDFromHex(*input.Id)
		if err != nil {
			return nil, err
		}
		return bson.D{
			{Key: "_id", Value: id},
		}, nil
	}
	if input.Name != nil {
		return bson.D{
			{Key: "name", Value: input.Name},
		}, nil
	}
	if input.Owner != nil {
		return bson.D{
			{Key: "owner", Value: input.Owner},
		}, nil
	}
	return nil, errors.New("invalid input")
}

func toTeams(ctx context.Context, cursor *mongo.Cursor, page uint64, max uint64) (*schema.Teams, error) {
	var result []*schema.Team
	skip := (int(page) - 1) * int(max)
	for i := 0; i < skip; i++ {
		cursor.Next(ctx)
	}
	for i := 0; i < int(max); i++ {
		if !cursor.Next(ctx) {
			break
		}
		team, err := toTeam(cursor)
		if err != nil {
			return nil, err
		}
		result = append(result, team)
	}
	return &schema.Teams{Teams: result}, nil
}

func toTeam(cursor *mongo.Cursor) (*schema.Team, error) {
	var result *bson.M
	err := cursor.Decode(&result)
	if err != nil {
		return nil, err
	}
	// Convert []int64 to []uint64
	membersWithoutOwner := []uint64{}
	for _, member := range (*result)["membersWithoutOwner"].(primitive.A) {
		membersWithoutOwner = append(membersWithoutOwner, uint64(member.(int64)))
	}
	(*result)["membersWithoutOwner"] = membersWithoutOwner
	// Convert data to map[string]string
	data := (*result)["data"].(primitive.M)
	(*result)["data"] = map[string]string{}
	for key, value := range data {
		(*result)["data"].(map[string]string)[key] = value.(string)
	}
	// If rank is not given, set it to 0
	if _, ok := (*result)["rank"]; !ok {
		(*result)["rank"] = int32(0)
	}
	return &schema.Team{
		Id:                  (*result)["_id"].(primitive.ObjectID).Hex(),
		Name:                (*result)["name"].(string),
		Owner:               uint64((*result)["owner"].(int64)),
		MembersWithoutOwner: membersWithoutOwner,
		Score:               (*result)["score"].(int64),
		Rank:                int64((*result)["rank"].(int32)),
		Data:                (*result)["data"].(map[string]string),
	}, nil
}
