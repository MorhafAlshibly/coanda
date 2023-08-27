package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/queue"
	"github.com/MorhafAlshibly/coanda/services/team/schema"
	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TeamService struct {
	db                database.Databaser
	cache             cache.Cacher
	queue             queue.Queuer
	maxMembers        int
	minTeamNameLength int
}

type NewTeamServiceInput struct {
	DatabaseConnection string
	DatabaseName       string
	DatabaseCollection string
	QueueConnection    string
	QueueName          string
	CacheConnection    string
	CachePassword      string
	CacheDB            int
	CacheExpiration    time.Duration
	MaxMembers         int
	MinTeamNameLength  int
}

// Pipeline
var rankStage = bson.D{
	{Key: "$setWindowFields", Value: bson.D{
		{Key: "sortBy", Value: bson.D{
			{Key: "score", Value: -1},
		}},
		{Key: "output", Value: bson.D{
			{Key: "rank", Value: bson.D{
				{Key: "$rank", Value: bson.D{}},
			}},
		}},
	}},
}

var dbIndices = []mongo.IndexModel{
	{
		Keys: bson.D{
			{Key: "name", Value: "text"},
		},
	},
	{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	},
	{
		Keys: bson.D{
			{Key: "owner", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	},
	{
		Keys: bson.D{
			{Key: "score", Value: -1},
		},
	},
}

func NewTeamService(ctx context.Context, input NewTeamServiceInput) (*TeamService, error) {
	db, err := database.NewMongoDatabase(ctx, database.MongoDatabaseInput{
		Connection: "mongodb://localhost:27017",
		Database:   "coanda",
		Collection: "teams",
		Indices:    dbIndices,
	})
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	defer db.Disconnect(ctx)
	queue, err := queue.NewServiceBus(ctx, input.QueueConnection, input.QueueName)
	if err != nil {
		log.Fatalf("failed to create queue: %v", err)
	}
	cache := cache.NewRedisCache(input.CacheConnection, input.CachePassword, input.CacheDB, input.CacheExpiration)
	return &TeamService{
		db:                db,
		cache:             cache,
		queue:             queue,
		maxMembers:        input.MaxMembers,
		minTeamNameLength: input.MinTeamNameLength,
	}, nil
}

func (s *TeamService) Disconnect(ctx context.Context) error {
	return s.db.Disconnect(ctx)
}

func (s *server) CreateTeam(ctx context.Context, req *schema.CreateTeamRequest) (*schema.Team, error) {
	// Check if team name is large enough
	if len(req.Name) < 3 {
		return nil, errors.New("team name too short")
	}
	// Remove duplicates from members
	req.MembersWithoutOwner = removeDuplicate(req.MembersWithoutOwner)
	if len(req.MembersWithoutOwner)+1 > s.maxMembers {
		return nil, errors.New("too many members")
	}
	// Check if score is given
	if req.Score == nil {
		req.Score = new(int64)
		*req.Score = 0
	}
	// Insert the team into the database
	id, err := s.db.InsertOne(ctx, bson.D{
		{Key: "name", Value: req.Name},
		{Key: "owner", Value: req.Owner},
		{Key: "membersWithoutOwner", Value: req.MembersWithoutOwner},
		{Key: "score", Value: *req.Score},
		{Key: "data", Value: req.Data},
	})
	if err != nil {
		return nil, err
	}
	return &schema.Team{
		Id:                  id,
		Name:                req.Name,
		Owner:               req.Owner,
		MembersWithoutOwner: req.MembersWithoutOwner,
		Score:               *req.Score,
		Data:                req.Data,
	}, nil
}

func (s *server) GetTeam(ctx context.Context, req *schema.GetTeamRequest) (*schema.Team, error) {
	filter, err := getFilter(req)
	if err != nil {
		return nil, err
	}
	data, err := s.cache.Get(ctx, fmt.Sprintf("%v", filter[0].Value))
	// If the item is not in the cache, get it from the store, else marshal it to output
	if err == nil {
		var out *schema.Team
		err = sonic.Unmarshal([]byte(data), &out)
		if err != nil {
			return nil, err
		}
		return out, nil
	}
	// Get the item from the store
	matchStage := bson.D{
		{Key: "$match", Value: filter},
	}
	cursor, err := s.db.Aggregate(ctx, mongo.Pipeline{rankStage, matchStage})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.Next(ctx)
	out, err := toTeam(cursor)
	if err != nil {
		return nil, err
	}
	// Cache the item
	go func() {
		err = cacheTeam(context.Background(), &s.cache, out)
		if err != nil {
			log.Println(err)
		}
	}()
	return out, nil
}

func (s *server) GetTeams(ctx context.Context, req *schema.GetTeamsRequest) (*schema.Teams, error) {
	// Create encoded cache key
	encodedKey := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(int(req.Max)) + "|" + strconv.Itoa(int(req.Page))))
	data, err := s.cache.Get(ctx, encodedKey)
	// If the item is not in the cache, get it from the store, else marshal it to output
	if err == nil {
		var out *schema.Teams
		err = sonic.Unmarshal([]byte(data), &out)
		if err != nil {
			return nil, err
		}
		return out, nil
	}
	cursor, err := s.db.Aggregate(ctx, mongo.Pipeline{rankStage})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	result, err := toTeams(ctx, cursor, req.Page, req.Max)
	if err != nil {
		return nil, err
	}
	// Cache the item
	go func() {
		marshalled, err := sonic.Marshal(&schema.Teams{Teams: result})
		if err != nil {
			log.Println(err)
		}
		err = s.cache.Add(context.Background(), encodedKey, string(marshalled))
		if err != nil {
			log.Println(err)
		}
	}()
	return &schema.Teams{Teams: result}, nil
}

func (s *server) SearchTeams(ctx context.Context, req *schema.SearchTeamsRequest) (*schema.Teams, error) {
	searchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "name", Value: bson.D{
				{Key: "$regex", Value: req.Query},
				{Key: "$options", Value: "i"},
			}},
		}},
	}
	if len(req.Query) < 3 {
		return nil, errors.New("query too short")
	}
	cursor, err := s.db.Aggregate(ctx, mongo.Pipeline{rankStage, searchStage})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	result, err := toTeams(ctx, cursor, req.Page, req.Max)
	if err != nil {
		return nil, err
	}
	return &schema.Teams{Teams: result}, nil
}

func (s *server) UpdateTeamScore(ctx context.Context, req *schema.UpdateTeamScoreRequest) (*schema.Team, error) {
	filter, err := getFilter(req.Team)
	if err != nil {
		return nil, err
	}
	var out *schema.Team
	_, err = s.db.UpdateOne(ctx, filter, bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "score", Value: req.ScoreOffset},
		}},
	})
	if err != nil {
		return nil, err
	}
	out, err = s.GetTeam(ctx, req.Team)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *server) UpdateTeamData(ctx context.Context, req *schema.UpdateTeamDataRequest) (*schema.Team, error) {
	filter, err := getFilter(req.Team)
	if err != nil {
		return nil, err
	}
	var out *schema.Team
	_, err = s.db.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "data", Value: req.Data},
		}},
	})
	if err != nil {
		return nil, err
	}
	out, err = s.GetTeam(ctx, req.Team)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *server) DeleteTeam(ctx context.Context, req *schema.DeleteTeamRequest) (*schema.Team, error) {
	filter, err := getFilter(req.Team)
	if err != nil {
		return nil, err
	}
	_, err = s.db.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *server) JoinTeam(ctx context.Context, req *schema.JoinTeamRequest) (*schema.BoolResponse, error) {
	marshalled, err := sonic.Marshal(req.Team)
	if err != nil {
		return nil, err
	}
	err = s.queue.Enqueue(ctx, "JoinTeam", marshalled)
	if err != nil {
		return nil, err
	}
	return &schema.BoolResponse{Value: true}, nil
}

func (s *server) LeaveTeam(ctx context.Context, req *schema.LeaveTeamRequest) (*schema.BoolResponse, error) {
	marshalled, err := sonic.Marshal(req.Team)
	if err != nil {
		return nil, err
	}
	err = s.queue.Enqueue(ctx, "LeaveTeam", marshalled)
	if err != nil {
		return nil, err
	}
	return &schema.BoolResponse{Value: true}, nil
}

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

func toTeams(ctx context.Context, cursor *mongo.Cursor, page uint64, max uint64) ([]*schema.Team, error) {
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
	return result, nil
}

type TeamResult interface {
	Decode(v interface{}) error
}

func toTeam(cursor TeamResult) (*schema.Team, error) {
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

func cacheTeam(ctx context.Context, cache *cache.Cacher, team *schema.Team) error {
	marshalled, err := sonic.Marshal(team)
	if err != nil {
		return err
	}
	err = (*cache).Add(ctx, team.Id, string(marshalled))
	if err != nil {
		return err
	}
	err = (*cache).Add(ctx, team.Name, string(marshalled))
	if err != nil {
		return err
	}
	err = (*cache).Add(ctx, strconv.Itoa(int(team.Owner)), string(marshalled))
	if err != nil {
		return err
	}
	return nil
}

func removeDuplicate[T string | int | uint64](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
