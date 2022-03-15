package awsproviderv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var _ SnsClientInterface = (*SNSMock)(nil)

type SNSMock struct{}

func (S SNSMock) ListTopics(ctx context.Context, params *sns.ListTopicsInput, optFns ...func(*sns.Options)) (*sns.ListTopicsOutput, error) {
	return &sns.ListTopicsOutput{
		Topics: []types.Topic{
			{TopicArn: aws.String("arn:::::topicA")},
		}}, nil
}
