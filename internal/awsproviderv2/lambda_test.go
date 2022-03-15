package awsproviderv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

var _ LambdaClientInterface = (*LambdaMock)(nil)

type LambdaMock struct{}

func (l LambdaMock) ListFunctions(ctx context.Context, params *lambda.ListFunctionsInput, optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error) {

	return &lambda.ListFunctionsOutput{
		Functions: []types.FunctionConfiguration{
			{
				FunctionName: aws.String("FunctionA"),
			},
		},
	}, nil
}
