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
		paginator := dynamodb.NewListTablesPaginator(p.Client, &dynamodb.ListTablesInput{}, func(o *dynamodb.ListTablesPaginatorOptions) {})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(context.TODO(), func(o *dynamodb.Options) { o.Region = region })
			if err != nil {
				logrus.Errorf("func:GetDynamoTableNames: AWS Error: %s in region %s", err.Error(), region)
				return []string{}
			}
			for _, tableName := range page.TableNames {
				tablesRsp = append(tablesRsp, tableName)
			}
		}
	}
	return tablesRsp
}