package awsproviderv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)

func assumeRole(accountRoleArn string, stsprovider STSProvider) (aws.Config, error) {
	cnf, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(aws.NewCredentialsCache(
			stscreds.NewAssumeRoleProvider(
				stsprovider.Client,
				accountRoleArn,
			)),
		),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return cnf, nil
}
