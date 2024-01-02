package dynamoTable

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/MorhafAlshibly/coanda/pkg/cache"
// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/credentials"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
// )

// func DynamoCreateTable() (*DynamoTable, error) {
// 	table, err := NewDynamoTable(context.Background(), &DynamoTableInput{
// 		Config: &aws.Config{
// 			Region:      aws.String("localhost"),
// 			Endpoint:    aws.String("http://localhost:8000"),
// 			Credentials: credentials.NewStaticCredentials("test", "test", "test"),
// 		},
// 		CreateTableInput: &dynamodb.CreateTableInput{
// 			TableName: aws.String("test"),
// 			AttributeDefinitions: []*dynamodb.AttributeDefinition{
// 				{
// 					AttributeName: aws.String("id"),
// 					AttributeType: aws.String("S"),
// 				},
// 				{
// 					AttributeName: aws.String("name"),
// 					AttributeType: aws.String("S"),
// 				},
// 			},
// 			KeySchema: []*dynamodb.KeySchemaElement{
// 				{
// 					AttributeName: aws.String("id"),
// 					KeyType:       aws.String("HASH"),
// 				},
// 				{
// 					AttributeName: aws.String("name"),
// 					KeyType:       aws.String("RANGE"),
// 				},
// 			},
// 			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
// 				ReadCapacityUnits:  aws.Int64(5),
// 				WriteCapacityUnits: aws.Int64(5),
// 			},
// 		},
// 		Cache: cache.NewRedisCache("localhost:6379", "", 0, time.Second*30),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return table, nil
// }

// func TestDynamoPutItem(t *testing.T) {
// 	table, err := DynamoCreateTable()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = table.PutItem(context.Background(), &PutItemInput{
// 		Item: map[string]*dynamodb.AttributeValue{
// 			"id": {
// 				S: aws.String("1"),
// 			},
// 			"name": {
// 				S: aws.String("test"),
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestDynamoGetItem(t *testing.T) {
// 	table, err := DynamoCreateTable()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = table.PutItem(context.Background(), &PutItemInput{
// 		Item: map[string]*dynamodb.AttributeValue{
// 			"id": {
// 				S: aws.String("1"),
// 			},
// 			"name": {
// 				S: aws.String("test"),
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	item, err := table.GetItem(context.Background(), &GetItemInput{
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"id": {
// 				S: aws.String("1"),
// 			},
// 			"name": {
// 				S: aws.String("test"),
// 			},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if item["id"].S == nil {
// 		t.Fatal("id is nil")
// 	}
// 	if item["name"].S == nil {
// 		t.Fatal("name is nil")
// 	}
// 	if *item["id"].S != "1" {
// 		t.Fatal("id is not 1")
// 	}
// 	if *item["name"].S != "test" {
// 		t.Fatal("name is not test")
// 	}
// }
