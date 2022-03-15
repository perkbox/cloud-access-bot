package awsproviderv2

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSProvider struct {
	Client      SqsClientInterface
	STSProvider *STSProvider
	Regions     []string
}

type SqsClientInterface interface {
	ListQueues(ctx context.Context, params *sqs.ListQueuesInput, optFns ...func(*sqs.Options)) (*sqs.ListQueuesOutput, error)
}

func NewSqsClient(cfg aws.Config) *SQSProvider {
	sqsClient := sqs.NewFromConfig(cfg)

	return &SQSProvider{
		Client: sqsClient,
	}
}

func (sqsp *SQSProvider) GetSQSQueues(accountRoleArn string) []string {
	p := sqsp
	if accountRoleArn != "" {
		cfg, err := assumeRole(accountRoleArn, *sqsp.STSProvider)
		if err != nil {
			logrus.Errorf("func:GetLambdaFunctions: Error assuming role %s.  AWS Error: %s", accountRoleArn, err.Error())
		}
		p = NewSqsClient(cfg)
	}

	var queueNames []string
	for _, region := range p.Regions {
		sqsResp, err := p.Client.ListQueues(context.TODO(), &sqs.ListQueuesInput{}, func(o *sqs.Options) {
			o.Region = region
		})
		if err != nil {
			logrus.Errorf("func:GetSQSQueues: Error assuming role %s.  AWS Error: %s", accountRoleArn, err.Error())
		}

		for _, k := range sqsResp.QueueUrls {
			split := strings.Split(k, "/")
			queueNames = append(queueNames, split[4])
		}
	}

	return queueNames
}
