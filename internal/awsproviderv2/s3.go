package awsproviderv2

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Provider struct {
	Client      S3ClientInterface
	STSProvider *STSProvider
	Regions     []string
}

func NewS3Client(cfg aws.Config) *S3Provider {
	s3Client := s3.NewFromConfig(cfg)

	return &S3Provider{
		Client: s3Client,
	}
}

type S3ClientInterface interface {
	// ListBuckets Required as there is currently no ListBucketsAPIClient
	ListBuckets(ctx context.Context,
		params *s3.ListBucketsInput,
		optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
}

func (s3p *S3Provider) GetBucketNames(accountRoleArn string) []string {
	p := s3p
	if accountRoleArn != "" {
		cfg, err := assumeRole(accountRoleArn, *s3p.STSProvider)
		if err != nil {
			logrus.Errorf("Error assuming role %s.  AWS Error: %s", accountRoleArn, err.Error())
		}
		p = NewS3Client(cfg)
	}

	var bucketNames []string
	for _, region := range s3p.Regions {
		result, err := p.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{}, func(o *s3.Options) {
			o.Region = region
		})
		if err != nil {
			logrus.Errorf("Error fetching S3 buckets from region %s with role %s", region, accountRoleArn)
			continue
		}

		for _, b := range result.Buckets {
			bucketNames = append(bucketNames, aws.ToString(b.Name))
		}
	}

	return bucketNames
}
