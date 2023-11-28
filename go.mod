module github.com/perkbox/cloud-access-bot

go 1.21

require (
	github.com/aws/aws-sdk-go v1.48.4
	github.com/aws/aws-sdk-go-v2 v1.23.1
	github.com/aws/aws-sdk-go-v2/config v1.25.5
	github.com/aws/aws-sdk-go-v2/credentials v1.16.4
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.12.3
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.14.3
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.25.3
	github.com/aws/aws-sdk-go-v2/service/iam v1.27.3
	github.com/aws/aws-sdk-go-v2/service/lambda v1.48.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.45.0
	github.com/aws/aws-sdk-go-v2/service/sns v1.25.4
	github.com/aws/aws-sdk-go-v2/service/sqs v1.28.3
	github.com/aws/aws-sdk-go-v2/service/sts v1.25.5
	github.com/aws/smithy-go v1.17.0
	github.com/banzaicloud/logrus-runtime-formatter v0.0.0-20190729070250-5ae5475bae5e
	github.com/joho/godotenv v1.5.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.9.3
	github.com/slack-go/slack v0.12.3
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.5.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.14.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.7.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.2.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.17.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.2.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.16.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.17.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.20.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)