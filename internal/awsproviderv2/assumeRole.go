package awsproviderv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)

func assumeRole(accountRoleArn string, stsprovider STSProvider) (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(aws.NewCredentialsCache(
			stscreds.NewAssumeRoleProvider(
				stsprovider.Client,
				accountRoleArn,
			)),
		),
	)
}
