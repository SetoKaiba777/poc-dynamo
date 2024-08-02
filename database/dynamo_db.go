package database

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient struct {
	tableName   string
	awsRegion   string
	awsEndpoint string
	client      *dynamodb.Client
}

func NewDynamoDBClient(tableName, awsRegion, awsEndpoint string) *DynamoDBClient {
	dbClient := DynamoDBClient{
		tableName:   tableName,
		awsRegion:   awsRegion,
		awsEndpoint: awsEndpoint,
	}
	dbClient.client = dbClient.loadDynamoDBClient()
	return &dbClient
}

func (c *DynamoDBClient) loadDynamoDBClient() *dynamodb.Client {
	awsconfig, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(_, _ string, _ ...interface{}) (aws.Endpoint, error) {
		if c.awsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           c.awsEndpoint,
				SigningRegion: c.awsRegion,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})), awsConfig.WithRegion(c.awsRegion))

	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(awsconfig, func(opt *dynamodb.Options) {
		opt.Region = awsconfig.Region
	})
}

func (c *DynamoDBClient) GetItemByID(primaryKey string, value string) (map[string]interface{}, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]types.AttributeValue{
			primaryKey: &types.AttributeValueMemberS{Value: value},
		},
	}

	result, err := c.client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("no item found with the specified ID")
	}
	out, err := convertToReadableMap(result.Item)
	if err != nil {
		return nil, fmt.Errorf("conversion failed: %v", err)
	}
	return out, nil
}

func convertToReadableMap(item map[string]types.AttributeValue) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for k, v := range item {
		switch v := v.(type) {
		case *types.AttributeValueMemberS:
			result[k] = v.Value
		case *types.AttributeValueMemberN:
			result[k] = v.Value
		case *types.AttributeValueMemberBOOL:
			result[k] = v.Value
		case *types.AttributeValueMemberB:
			result[k] = v.Value
		case *types.AttributeValueMemberSS:
			result[k] = v.Value
		case *types.AttributeValueMemberNS:
			result[k] = v.Value
		case *types.AttributeValueMemberBS:
			result[k] = v.Value
		case *types.AttributeValueMemberM:
			nested, err := convertToReadableMap(v.Value)
			if err != nil {
				return nil, err
			}
			result[k] = nested
		case *types.AttributeValueMemberL:
			var list []interface{}
			for _, item := range v.Value {
				convertedItem, err := convertToReadableMap(map[string]types.AttributeValue{"": item})
				if err != nil {
					return nil, err
				}
				list = append(list, convertedItem[""])
			}
			result[k] = list
		default:
			return nil, fmt.Errorf("unsupported attribute value type: %T", v)
		}
	}
	return result, nil
}
