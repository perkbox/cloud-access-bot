package repository

import (
	"context"
	"fmt"

	"github.com/perkbox/cloud-access-bot/internal"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamodbRepo struct {
	Client    dynamoClientInterface
	TableName string
}

func NewDynamoDBRRepo(cfg aws.Config, TableName string) *DynamodbRepo {
	dynamocClient := dynamodb.NewFromConfig(cfg)
	return &DynamodbRepo{
		Client:    dynamocClient,
		TableName: TableName,
	}
}

type dynamoClientInterface interface {
	dynamodb.QueryAPIClient
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

func (r *DynamodbRepo) UpdateApprovingUser(UserID, RequestId, approvingUser string) error {
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.TableName),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":a": &types.AttributeValueMemberS{Value: approvingUser},
		},
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: UserID},
			"SK": &types.AttributeValueMemberS{Value: RequestId},
		},
		UpdateExpression: aws.String("set ApprovingUser = :a "),
	}

	if _, err := r.Client.UpdateItem(context.TODO(), updateInput); err != nil {
		return fmt.Errorf("func:UpdateApprovingUser: error Running DynamoDB Update  %w", err)
	}
	return nil
}

func (r *DynamodbRepo) QueryAuditObjs(UserID string) ([]internal.AuditObject, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(r.TableName),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":PK": &types.AttributeValueMemberS{Value: UserID},
		},
		KeyConditionExpression: aws.String("PK = :PK"),
	}

	respitems, err := r.Client.Query(context.TODO(), input)
	if err != nil {
		return []internal.AuditObject{}, fmt.Errorf("func:QueryAuditObjs: error running DynamoDb Query response: %w", err)
	}

	records := make([]internal.AuditObject, len(respitems.Items))
	for _, item := range respitems.Items {
		var itemRecord repoAuditObject
		if err := attributevalue.UnmarshalMap(item, &itemRecord); err != nil {
			return []internal.AuditObject{}, fmt.Errorf("func:QueryAuditObjs: error unmarshalling response: %w", err)
		}
		records = append(records, itemRecord.convertFromRepoObject())
	}

	return records, nil
}

func (r *DynamodbRepo) GetAuditObj(UserID, RequestId string) (internal.AuditObject, error) {
	item, err := r.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: UserID},
			"SK": &types.AttributeValueMemberS{Value: RequestId},
		},
		TableName: aws.String(r.TableName),
	})

	if err != nil {
		return internal.AuditObject{}, fmt.Errorf("func:GetAuditObj: error Getting dynamoDb item: %w", err)
	}

	var record repoAuditObject
	if err := attributevalue.UnmarshalMap(item.Item, &record); err != nil {
		return internal.AuditObject{}, fmt.Errorf("func:GetAuditObj: error unmarshalling dynamoDb response: %w", err)
	}

	return record.convertFromRepoObject(), nil
}

func (r *DynamodbRepo) SetAuditObj(requestObj internal.AuditObject) error {
	repoObj := convertToRepoObj(requestObj)
	data, err := attributevalue.MarshalMap(repoObj)
	if err != nil {
		return fmt.Errorf("func:SetAuditObj: error marshalling response: %w", err)
	}

	if _, err := r.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      data,
		TableName: aws.String(r.TableName),
	}); err != nil {
		return fmt.Errorf("func:SetAuditObj: error Putting item into DynamoDb: %w", err)
	}

	return nil
}
