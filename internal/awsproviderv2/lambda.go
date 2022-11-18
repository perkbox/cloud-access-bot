package awsproviderv2

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type LambdaProvider struct {
	Client      LambdaClientInterface
	STSProvider *STSProvider
	Regions     []string
}

type LambdaClientInterface interface {
	ListFunctions(ctx context.Context, params *lambda.ListFunctionsInput, optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error)
}

func NewLambdaClient(cfg aws.Config) *LambdaProvider {
	lambdaClient := lambda.NewFromConfig(cfg)

	return &LambdaProvider{
		Client: lambdaClient,
	}
}

func (lambdap *LambdaProvider) GetLambdaFunctions(accountRoleArn string) []string {
	p := lambdap
	if accountRoleArn != "" {
		cfg, err := assumeRole(accountRoleArn, *lambdap.STSProvider)
		if err != nil {
			logrus.Errorf("func:GetLambdaFunctions: Error assuming role %s.  AWS Error: %s", accountRoleArn, err.Error())
		}
		p = NewLambdaClient(cfg)
	}

	var functionNames []string
	var token *string

	for _, region := range lambdap.Regions {
		for {
			listFuncsResp, err := p.Client.ListFunctions(context.TODO(), &lambda.ListFunctionsInput{Marker: token}, func(o *lambda.Options) {
				o.Region = region
			})
			if err != nil {
				logrus.Errorf("func:GetLambdaFunctions: Error Fetching Functions.  AWS Error: %s", err.Error())
			}

			for _, v := range listFuncsResp.Functions {
				functionNames = append(functionNames, aws.ToString(v.FunctionName))
			}

			if listFuncsResp.NextMarker != nil {
				token = listFuncsResp.NextMarker
			} else {
				break
			}
		}
	}

	return functionNames
}
