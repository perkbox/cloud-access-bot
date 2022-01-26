package awsproviderv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var _ STSClientInterface = (*STSMock)(nil)

type STSMock struct{}

func (S STSMock) AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error) {
	panic("implement me")
}
