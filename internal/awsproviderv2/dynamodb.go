package awsproviderv2

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamodbProvider struct {
	Client      dynamoClientInterface
	STSProvider *STSProvider
	Regions     []string
}

func NewDynamoDBCClient(cfg aws.Config) *DynamodbProvider {
	dynamocClient := dynamodb.NewFromConfig(cfg)
	return &DynamodbProvider{
		Client: dynamocClient,
	}
}

type dynamoClientInterface interface {
	dynamodb.ListTablesAPIClient
	dynamodb.DescribeTableAPIClient
}

func (dyn *DynamodbProvider) GetDynamoTableNames(accountRoleArn string) []string {
	p := dyn
	if accountRoleArn != "" {
		cfg, err := assumeRole(accountRoleArn, *dyn.STSProvider)
		if err != nil {
			logrus.Errorf("func:GetDynamoTableNames: Error assuming role %s.  AWS Error: %s", accountRoleArn, err.Error())
		}
		p = NewDynamoDBCClient(cfg)
	}

	var tablesRsp []string
	for _, region := range dyn.Regions {
		tables, err := p.Client.ListTables(context.TODO(), &dynamodb.ListTablesInput{}, func(o *dynamodb.Options) {
			o.Region = region
		})
		if err != nil {
			logrus.Errorf("func:GetDynamoTableNames: Error Fetching Dynamo Tables from region %s with role %s", region, accountRoleArn)
			continue
		}

		tablesRsp = append(tablesRsp, tables.TableNames...)
	}

	return tablesRsp
}
