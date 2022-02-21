package awsproviderv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type STSProvider struct {
	Client STSClientInterface
}
type STSClientInterface interface {
	AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error)
}

func NewSTSClient(cfg aws.Config) *STSProvider {
	stsClient := sts.NewFromConfig(cfg)
	return &STSProvider{
		Client: stsClient,
	}
}
