package policy

import (
	"context"
	"fmt"
	"net/url"

	localconfig "github.com/perkbox/cloud-access-bot/internal/settings"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type IAMProvider struct {
	Client      IAMClientInterface
	Settings    localconfig.Settings
	STSProvider *sts.Client
}

type IAMClientInterface interface {
	iam.ListRolePoliciesAPIClient
	iam.GetRoleAPIClient
	PutRolePolicy(ctx context.Context, params *iam.PutRolePolicyInput, optFns ...func(*iam.Options)) (*iam.PutRolePolicyOutput, error)
	DeleteRolePolicy(ctx context.Context, params *iam.DeleteRolePolicyInput, optFns ...func(*iam.Options)) (*iam.DeleteRolePolicyOutput, error)
	GetRolePolicy(ctx context.Context, params *iam.GetRolePolicyInput, optFns ...func(*iam.Options)) (*iam.GetRolePolicyOutput, error)
}

func NewIAMClient(cfg aws.Config) *IAMProvider {
	return &IAMProvider{
		Client: iam.NewFromConfig(cfg),
	}
}

func (awsiam *IAMProvider) PutPolicy(accountName, roleName, policyName, policy string) error {
	i, err := awsiam.assumeRoleWithAccountName(accountName)
	if err != nil {
		return err
	}

	_, err = i.PutRolePolicy(context.TODO(), &iam.PutRolePolicyInput{
		PolicyDocument: aws.String(policy),
		PolicyName:     aws.String(policyName),
		RoleName:       aws.String(roleName),
	}, func(o *iam.Options) {
		o.Region = "eu-west-1"
	})

	if err != nil {
		return err
	}

	return nil
}

func (awsiam *IAMProvider) DeletePolicys(accountName, roleName string, InlinePolicysNames []string) error {
	i, err := awsiam.assumeRoleWithAccountName(accountName)
	if err != nil {
		return err
	}

	for _, policyName := range InlinePolicysNames {
		_, err := i.DeleteRolePolicy(context.TODO(), &iam.DeleteRolePolicyInput{RoleName: aws.String(roleName), PolicyName: aws.String(policyName)}, func(o *iam.Options) {
			o.Region = "eu-west-1"
		})
		if err != nil {
			return err
		}
		logrus.Infof("Delted Expired Inline Policy (%s) from Role %s", policyName, roleName)
	}

	return nil
}

func (awsiam *IAMProvider) GetCloudUserId(accountName string, roleName string) (string, error) {
	i, err := awsiam.assumeRoleWithAccountName(accountName)
	if err != nil {
		return "", err
	}

	roleOutput, err := i.GetRole(context.TODO(), &iam.GetRoleInput{RoleName: aws.String(roleName)}, func(o *iam.Options) {
		o.Region = "eu-west-1"
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(roleOutput.Role.RoleId), nil
}

func (awsiam *IAMProvider) FindPolicysForRole(accountName, roleName string) (map[string]string, error) {
	inlinePolicys := make(map[string]string)
	i, err := awsiam.assumeRoleWithAccountName(accountName)
	if err != nil {
		return nil, err
	}

	//GET all inline Policy's on the role
	listPolResp, err := i.ListRolePolicies(context.TODO(), &iam.ListRolePoliciesInput{RoleName: aws.String(roleName)}, func(o *iam.Options) {
		o.Region = "eu-west-1"
	})
	if err != nil {
		return nil, fmt.Errorf("func:FindPolicysForRole: error listing role polices: %s", err)
	}

	for _, policyName := range listPolResp.PolicyNames {
		policyResp, err := i.GetRolePolicy(context.TODO(), &iam.GetRolePolicyInput{RoleName: aws.String(roleName), PolicyName: aws.String(policyName)}, func(o *iam.Options) {
			o.Region = "eu-west-1"
		})
		if err != nil {
			return nil, fmt.Errorf("func:FindPolicysForRole: error getting inline policy from role: %s ", err)
		}

		policyDoc, _ := url.QueryUnescape(aws.ToString(policyResp.PolicyDocument))

		inlinePolicys[policyName] = policyDoc

	}

	return inlinePolicys, err
}

func (awsiam *IAMProvider) assumeRoleWithAccountName(accountName string) (IAMClientInterface, error) {
	i := awsiam.Client
	accountRoleArn, err := awsiam.Settings.GetRoleArn(accountName)
	if err != nil {
		return nil, err
	}

	if accountRoleArn != "" {
		cfg, err := assumeRole(accountRoleArn, *awsiam.STSProvider)
		if err != nil {
			return nil, fmt.Errorf("func:assumeRoleWithAccountName: Error assuming role %s.  AWS Error: %s", accountRoleArn, err.Error())
		}
		i = NewIAMClient(cfg).Client
	}

	return i, nil
}
