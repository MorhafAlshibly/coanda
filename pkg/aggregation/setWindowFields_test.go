package aggregation

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestSetWindowFields(t *testing.T) {
	documents := []bson.M{
		{"name": "a", "score": 3},
		{"name": "a", "score": 9},
		{"name": "a", "score": 6},
		{"name": "b", "score": 4},
		{"name": "b", "score": 2},
		{"name": "b", "score": 18},
	}
	c := SetWindowFieldsCommand{
		In: &SetWindowFieldsInput{
			Documents:   &documents,
			PartitionBy: "name",
			SortBy:      "score",
			Asc:         false,
			Output:      "rank",
		},
	}
	err := c.Execute(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(*c.Out.Result) != 6 {
		t.Fatal("Expected 6 documents")
	}
	if (*c.Out.Result)[0]["rank"] != 1 || (*c.Out.Result)[0]["name"] != "a" || (*c.Out.Result)[0]["score"] != 9 {
		t.Fatal("Expected rank 1, name a, score 9")
	}
	if (*c.Out.Result)[1]["rank"] != 2 || (*c.Out.Result)[1]["name"] != "a" || (*c.Out.Result)[1]["score"] != 6 {
		t.Fatal("Expected rank 2, name a, score 6")
	}
	if (*c.Out.Result)[2]["rank"] != 3 || (*c.Out.Result)[2]["name"] != "a" || (*c.Out.Result)[2]["score"] != 3 {
		t.Fatal("Expected rank 3, name a, score 3")
	}
	if (*c.Out.Result)[3]["rank"] != 1 || (*c.Out.Result)[3]["name"] != "b" || (*c.Out.Result)[3]["score"] != 18 {
		t.Fatal("Expected rank 1, name b, score 18")
	}
	if (*c.Out.Result)[4]["rank"] != 2 || (*c.Out.Result)[4]["name"] != "b" || (*c.Out.Result)[4]["score"] != 4 {
		t.Fatal("Expected rank 2, name b, score 4")
	}
	if (*c.Out.Result)[5]["rank"] != 3 || (*c.Out.Result)[5]["name"] != "b" || (*c.Out.Result)[5]["score"] != 2 {
		t.Fatal("Expected rank 3, name b, score 2")
	}
}

func TestSetWindowFieldsEmptyDocuments(t *testing.T) {
	documents := []bson.M{}
	c := SetWindowFieldsCommand{
		In: &SetWindowFieldsInput{
			Documents:   &documents,
			PartitionBy: "name",
			SortBy:      "score",
			Asc:         false,
			Output:      "rank",
		},
	}
	err := c.Execute(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(*c.Out.Result) != 0 {
		t.Fatal("Expected 0 documents")
	}
}
