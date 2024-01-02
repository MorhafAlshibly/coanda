package dynamoTable

import (
	"context"

	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoTable struct {
	client        *dynamodb.Client
	cache         cache.Cacher
	tableName     string
	readTableName string
}

type DynamoTableInput struct {
	Options          *dynamodb.Options
	CreateTableInput *dynamodb.CreateTableInput
	Cache            cache.Cacher
	ReadTableName    *string
}

func NewDynamoTable(ctx context.Context, input *DynamoTableInput) (*DynamoTable, error) {
	client := dynamodb.New(*input.Options)
	_, err := client.CreateTable(ctx, input.CreateTableInput)
	if err != nil {
		_, err = client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: input.CreateTableInput.TableName,
		})
		if err != nil {
			return nil, err
		}
	}
	var readTableName string
	if input.ReadTableName != nil {
		readTableName = *input.ReadTableName
	} else {
		readTableName = *input.CreateTableInput.TableName
	}
	return &DynamoTable{
		client:        client,
		tableName:     *input.CreateTableInput.TableName,
		readTableName: readTableName,
		cache:         input.Cache,
	}, nil
}

func (d *DynamoTable) PutItem(ctx context.Context, input *PutItemInput) error {
	command := &PutItemCommand{
		table: d,
		In:    input,
	}
	invoker := invokers.NewBasicInvoker()
	return invoker.Invoke(ctx, command)
}

func (d *DynamoTable) GetItem(ctx context.Context, input *GetItemInput) (map[string]any, error) {
	command := &GetItemCommand{
		table: d,
		In:    input,
	}
	invoker := invokers.NewCacheInvoker(d.cache)
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (d *DynamoTable) DeleteItem(ctx context.Context, input *DeleteItemInput) error {
	command := &DeleteItemCommand{
		table: d,
		In:    input,
	}
	invoker := invokers.NewBasicInvoker()
	return invoker.Invoke(ctx, command)
}

func (d *DynamoTable) UpdateItem(ctx context.Context, input *UpdateItemInput) error {
	command := &UpdateItemCommand{
		table: d,
		In:    input,
	}
	invoker := invokers.NewBasicInvoker()
	return invoker.Invoke(ctx, command)
}

func (d *DynamoTable) Query(ctx context.Context, input *QueryInput) ([]map[string]any, error) {
	command := &QueryCommand{
		table: d,
		In:    input,
	}
	invoker := invokers.NewCacheInvoker(d.cache)
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (d *DynamoTable) Scan(ctx context.Context, input *ScanInput) ([]map[string]any, error) {
	command := &ScanCommand{
		table: d,
		In:    input,
	}
	invoker := invokers.NewCacheInvoker(d.cache)
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func NullifyEmptyMapAttributeValue(mapStringAny map[string]any, mapAttrivuteValue map[string]types.AttributeValue) map[string]types.AttributeValue {
	if mapStringAny == nil {
		mapAttrivuteValue = nil
	}
	return mapAttrivuteValue
}

func NullifyEmptyStringMapAttributeValue(mapStringString map[string]string, mapAttrivuteValue map[string]types.AttributeValue) map[string]types.AttributeValue {
	if mapStringString == nil {
		mapAttrivuteValue = nil
	}
	return mapAttrivuteValue
}

func NullifyEmptyString(str *string) *string {
	if *str == "" {
		str = nil
	}
	return str
}
