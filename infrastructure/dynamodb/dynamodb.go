package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DB struct {
	client *dynamodb.Client
}

func NewDB(ctx context.Context, endpoint, region string) (*DB, error) {
	resolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID {
			return aws.Endpoint{
				URL: endpoint,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithEndpointResolver(resolver),
	)
	if err != nil {
		return nil, err
	}
	return &DB{dynamodb.NewFromConfig(cfg)}, nil
}

func (db DB) PutItem(ctx context.Context, table string, item map[string]types.AttributeValue) error {
	_, err := db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &table,
		Item:      item,
	})
	if err != nil {
		return err
	}
	return nil
}

func (db DB) GetItem(ctx context.Context, table, partitionKey, sortKey string) (map[string]types.AttributeValue, error) {
	stmt := fmt.Sprintf("SELECT * FROM %s WHERE partition_key=? AND sort_key=?", table)
	params := []types.AttributeValue{
		&types.AttributeValueMemberS{Value: partitionKey},
		&types.AttributeValueMemberS{Value: sortKey},
	}
	input := &dynamodb.ExecuteStatementInput{
		Statement:  &stmt,
		Parameters: params,
	}
	res, err := db.client.ExecuteStatement(ctx, input)
	if err != nil {
		return nil, err
	}
	if len(res.Items) == 0 {
		return nil, nil
	}
	return res.Items[0], nil
}

func (db DB) GetItemsWithSortKeyRange(ctx context.Context, table, partitionKey, since, until string) ([]map[string]types.AttributeValue, error) {
	stmt := fmt.Sprintf("SELECT * FROM %s WHERE partition_key=? AND sort_key BETWEEN ? AND ?", table)
	params := []types.AttributeValue{
		&types.AttributeValueMemberS{Value: partitionKey},
		&types.AttributeValueMemberS{Value: since},
		&types.AttributeValueMemberS{Value: until},
	}
	input := &dynamodb.ExecuteStatementInput{
		Statement:  &stmt,
		Parameters: params,
	}
	res, err := db.client.ExecuteStatement(ctx, input)
	if err != nil {
		return nil, err
	}
	if len(res.Items) == 0 {
		return nil, nil
	}
	return res.Items, nil
}
