package awsproviderv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
)

var _ S3ClientInterface = (*S3Mock)(nil)

type S3Mock struct{}

func (s S3Mock) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return &s3.ListBucketsOutput{Buckets: []types.Bucket{
		{
			Name: aws.String("BucketA"),
		},
		{
			Name: aws.String("BucketB"),
		},
	}}, nil
}
