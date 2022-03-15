package awsproviderv2

import (
	"context"

	"github.com/aws/smithy-go/middleware"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var _ dynamoClientInterface = (*DynMock)(nil)

type DynMock struct{}

func (D DynMock) ListTables(ctx context.Context, input *dynamodb.ListTablesInput, f ...func(*dynamodb.Options)) (*dynamodb.ListTablesOutput, error) {

	return &dynamodb.ListTablesOutput{
		LastEvaluatedTableName: nil,
		TableNames:             []string{"TestDynTable1", "TestDynTable2"},
		ResultMetadata:         middleware.Metadata{},
	}, nil
}
