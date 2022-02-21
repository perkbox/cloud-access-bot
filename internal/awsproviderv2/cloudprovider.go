package awsproviderv2

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/perkbox/cloud-access-bot/internal/settings"
	"github.com/sirupsen/logrus"
)

// ResourceFinder is the New service based implementation of cloudprovider to decouple the data logic from the application code
type ResourceFinder struct {
	*S3Provider
	*DynamodbProvider
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

	return &ResourceFinder{
		S3Provider:       S3,
		DynamodbProvider: Dynamo,
		Settings:         config,
	}
}

// ResourceFinder Returns a method to find the available resources as well as ARN metadata used to render an ARN using
// the GenerateArnForService function.
// Returns
// GetNamesHashMap Function to get resource names as a map[string]string
// bool    An easy way to see if the Aws Service has a resource finder
func (c *ResourceFinder) ResourceFinder(service string, accountName string) ([]string, bool) {
	var (
		roleArn string
		err     error
	)
	if accountName != "" {
		roleArn, err = c.Settings.GetRoleArn(accountName)
		if err != nil {
			logrus.Errorf("error fetching role for account %s Err: %s", accountName, err.Error())
			return nil, false
		}
	}

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

func (v *ResourceFinder) ValidateResourcesFormat(resources []string) []string {
	var failedResources []string
	for _, r := range resources {
		if !arn.IsARN(r) && r != "" {
			failedResources = append(failedResources, r)
		}
	}

	return failedResources
}
