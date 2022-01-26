module github.com/perkbox/cloud-access-bot

go 1.17

replace github.com/slack-go/slack => github.com/xnok/slack v0.8.1-0.20210509200330-9b2b404dbde9

require (
	github.com/aws/aws-sdk-go v1.40.46
	github.com/aws/aws-sdk-go-v2 v1.10.0
	github.com/aws/aws-sdk-go-v2/config v1.9.0
	github.com/aws/aws-sdk-go-v2/credentials v1.5.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.3.0
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.6.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.6.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.10.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.17.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.8.0
	github.com/joho/godotenv v1.4.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/slack-go/slack v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.0.7 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.2.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.5.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.4.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.2.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.4.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.5.0 // indirect
	github.com/aws/smithy-go v1.8.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.0.0-20210423082822-04245dca01da // indirect
)
