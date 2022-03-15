package awsproviderv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var _ SqsClientInterface = (*SQSMock)(nil)

type SQSMock struct{}

func (S SQSMock) ListQueues(ctx context.Context, params *sqs.ListQueuesInput, optFns ...func(*sqs.Options)) (*sqs.ListQueuesOutput, error) {

	return &sqs.ListQueuesOutput{
		QueueUrls: []string{"https://sqs.eu-west-1.amazonaws.com/123456789/qeueA"},
	}, nil
}
