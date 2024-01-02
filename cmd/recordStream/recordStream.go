package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
// Flags set from command line/environment variables
// fs              = ff.NewFlagSet("record")
// tableName       = fs.StringLong("tableName", "record", "the name of the table to use")
// indexName       = fs.StringLong("indexName", "leaderboard", "the name of the index to create and use for leaderboard")
// targetTableName = fs.StringLong("targetTableName", "leaderboardRankings", "the name of the table to use for leaderboard rankings")
)

// Lambda function to update a table record in DynamoDB
func main() {
	lambda.Start(handler)
}

// Handler for the lambda function
func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	payload := req.Body
	log.Println("payload", payload)

	var record api.Record
	err := json.Unmarshal([]byte(payload), &record)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}
	log.Println("record", record)
	return events.APIGatewayV2HTTPResponse{}, nil
}
