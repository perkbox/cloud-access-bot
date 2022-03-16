package awsproviderv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSProvider struct {
	Client      SnsClientInterface
	STSProvider *STSProvider
	Regions     []string
}

type SnsClientInterface interface {
	ListTopics(ctx context.Context, params *sns.ListTopicsInput, optFns ...func(*sns.Options)) (*sns.ListTopicsOutput, error)
}

func NewSnsClient(cfg aws.Config) *SNSProvider {
	snsClient := sns.NewFromConfig(cfg)

	return &SNSProvider{
		Client: snsClient,
	}
}

func (snsp *SNSProvider) GetSNSTopics(accountRoleArn string) []string {
	p := snsp
	if accountRoleArn != "" {
		cfg, err := assumeRole(accountRoleArn, *snsp.STSProvider)
		if err != nil {
			logrus.Errorf("func:GetSNSTopics: Error assuming role %s.  AWS Error: %s", accountRoleArn, err.Error())
		}
		p = NewSnsClient(cfg)
	}

	var topicNames []string
	for _, region := range snsp.Regions {
		snsResp, err := p.Client.ListTopics(context.TODO(), &sns.ListTopicsInput{}, func(o *sns.Options) {
			o.Region = region
		})
		if err != nil {
			logrus.Errorf("func:GetSNSTopics: AWS Error: %s", err.Error())
		}

		for _, v := range snsResp.Topics {
			arnparse, err := arn.Parse(aws.ToString(v.TopicArn))
			if err != nil {
				logrus.Errorf("cant parse %s", err.Error())
			}
			topicNames = append(topicNames, arnparse.Resource)
		}
	}

	return topicNames
}
