package awsproviderv2

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/perkbox/cloud-access-bot/internal/settings"
)

// ResourceFinder is the New service based implementation of cloudprovider to decouple the data logic from the application code
type ResourceFinder struct {
	*S3Provider
	*DynamodbProvider
	*Validator
	settings.Settings
}

func NewAwsResourceFinder(cfg aws.Config, config settings.Settings) *ResourceFinder {
	stsClient := NewSTSClient(cfg)

	S3 := NewS3Client(cfg)
	S3.STSProvider = stsClient
	S3.Regions = config.Regions

	Dynamo := NewDynamoDBCClient(cfg)
	Dynamo.STSProvider = stsClient
	Dynamo.Regions = config.Regions

	Vali := NewValidator()

	return &ResourceFinder{
		S3Provider:       S3,
		DynamodbProvider: Dynamo,
		Validator:        Vali,
		Settings:         config,
	}
}

// ResourceFinder Returns a method to find the available resources as well as ARN metadata used to render an ARN using
// the GenerateArnForService function.
// Returns
// GetNamesHashMap Function to get resource names as a map[string]string
// bool    An easy way to see if the Aws Service has a resource finder
func (c *ResourceFinder) ResourceFinder(service string, accountName string) ([]string, bool) {
	roleArn, _ := c.Settings.GetRoleArn(accountName)

	switch service {
	case "dynamodb":
		tableNames := c.DynamodbProvider.GetDynamoTableNames(roleArn)
		return tableNames, true
	case "s3":
		bucketNames := c.S3Provider.GetBucketNames(roleArn)
		return bucketNames, true
	}
	return nil, false
}
