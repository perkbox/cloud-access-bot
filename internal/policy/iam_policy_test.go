package policy

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/perkbox/cloud-access-bot/internal"

	"github.com/aws/aws-sdk-go-v2/aws"
	localconfig "github.com/perkbox/cloud-access-bot/internal/settings"
)

func TestIamPolicy_IsPolicyExpired(t *testing.T) {
	tests := []struct {
		Name      string
		IsExpired bool
		Policy    IamPolicy
	}{
		{
			Name:      "Expired Policy",
			IsExpired: true,
			Policy: IamPolicy{
				Version: "2012-10-17",
				Statement: []IamStatement{
					{
						Sid:    "TestPolicy",
						Effect: "Allow",
						Condition: &IamCondition{
							DateGreaterThan: &DateGreaterThan{AwsCurrentTime: "2021-09-30T09:57:42Z"},
							DateLessThan:    &DateLessThan{AwsCurrentTime: "2021-09-30T11:10:42Z"},
						},
					},
				},
			},
		},
		{
			Name:      "Active Policy",
			IsExpired: false,
			Policy: IamPolicy{
				Version: "2012-10-17",
				Statement: []IamStatement{
					{
						Sid:    "TestPolicy",
						Effect: "Allow",
						Condition: &IamCondition{
							DateGreaterThan: &DateGreaterThan{AwsCurrentTime: "2021-09-30T11:07:42Z"},
							DateLessThan:    &DateLessThan{AwsCurrentTime: time.Now().UTC().Add(time.Duration(60) * time.Minute).Format("2006-01-02T15:04:05Z")},
						},
					},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {

			policyMan := NewPolicyManager(aws.Config{}, localconfig.Settings{}, nil, nil)
			polciy, _ := json.Marshal(tc.Policy)

			resp, err := policyMan.IsPolicyExpired(string(polciy))
			if err != nil {
				t.Errorf(err.Error())
			}

			if resp != tc.IsExpired {
				t.Errorf("Error expected IsExpired to be %v but got %v", tc.IsExpired, resp)
			}
		})
	}
}

func TestIamPolicyMan_GeneratePolicy(t *testing.T) {
	policyman := IamPolicyMan{}

	tests := []struct {
		Name              string
		arnTemplates      map[string]string
		arnTmplFieldNames map[string]string
		auditObj          internal.AuditObject
		expectedPolicy    IamPolicy
	}{
		{
			Name:              "SQS FIFO ARN",
			arnTemplates:      map[string]string{},
			arnTmplFieldNames: map[string]string{},
			auditObj: internal.AuditObject{
				RequestTime: time.Date(
					2021, 10, 21, 21, 21, 21, 651387237, time.UTC),
				CloudUserId: "abc123:email@gmail.com",
				AccountId:   "123456789",
				Duration:    "60",
				Services:    []string{"sqs"},
				Actions: map[string][]string{
					"sqs": {"sqsDOSOmething"},
				},
				Resources: map[string][]string{
					"sqs": {"arn:aws:sqs:eu-west-1:123456789:event-dlq.fifo"},
				},
			},
			expectedPolicy: IamPolicy{
				Version: "2012-10-17",
				Statement: []IamStatement{
					{
						Effect:   "Allow",
						Action:   []string{"sqsDOSOmething"},
						Resource: []string{"arn:aws:sqs:eu-west-1:123456789:event-dlq.fifo", "arn:aws:sqs:eu-west-1:123456789:event-dlq.fifo/*"},
						Condition: &IamCondition{
							StringEqualsIgnoreCase: &StringEqualsIgnoreCase{AwsUserid: "abc123:email@gmail.com"},
							DateGreaterThan:        &DateGreaterThan{AwsCurrentTime: "2021-10-21T21:21:21Z"},
							DateLessThan:           &DateLessThan{AwsCurrentTime: "2021-10-21T22:21:21Z"},
						},
					},
				},
			},
		},

		{
			Name:              "Mixed Templated and input arn policy Gen",
			arnTemplates:      map[string]string{"s3": "arn:{{.Partition}}:s3:::{{.BucketName}}"},
			arnTmplFieldNames: map[string]string{"s3": "BucketName"},
			auditObj: internal.AuditObject{
				RequestTime: time.Date(
					2021, 10, 21, 21, 21, 21, 651387237, time.UTC),
				CloudUserId: "abc123:email@gmail.com",
				AccountId:   "123456789",
				Duration:    "60",
				Services:    []string{"s3", "sts"},
				Actions: map[string][]string{
					"s3":  {"GetBucket", "ReadBucket"},
					"sts": {"ReadSomething"},
				},
				Resources: map[string][]string{
					"s3":  {"AbucketHere", "SomethingsBucket"},
					"sts": {"arn::::::anSTSResource"},
				},
			},
			expectedPolicy: IamPolicy{
				Version: "2012-10-17",
				Statement: []IamStatement{
					{
						Effect:   "Allow",
						Action:   []string{"GetBucket", "ReadBucket"},
						Resource: []string{"arn:aws:s3:::AbucketHere", "arn:aws:s3:::AbucketHere/*", "arn:aws:s3:::SomethingsBucket", "arn:aws:s3:::SomethingsBucket/*"},
						Condition: &IamCondition{
							StringEqualsIgnoreCase: &StringEqualsIgnoreCase{AwsUserid: "abc123:email@gmail.com"},
							DateGreaterThan:        &DateGreaterThan{AwsCurrentTime: "2021-10-21T21:21:21Z"},
							DateLessThan:           &DateLessThan{AwsCurrentTime: "2021-10-21T22:21:21Z"},
						},
					},
					{
						Effect:   "Allow",
						Action:   []string{"ReadSomething"},
						Resource: []string{"arn::::::anSTSResource", "arn::::::anSTSResource/*"},
						Condition: &IamCondition{
							StringEqualsIgnoreCase: &StringEqualsIgnoreCase{AwsUserid: "abc123:email@gmail.com"},
							DateGreaterThan:        &DateGreaterThan{AwsCurrentTime: "2021-10-21T21:21:21Z"},
							DateLessThan:           &DateLessThan{AwsCurrentTime: "2021-10-21T22:21:21Z"},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			// TEST HERE --------------
			policy, err := policyman.GeneratePolicyFromAuditObj(tc.auditObj.RequestTime, tc.auditObj, tc.arnTemplates, tc.arnTmplFieldNames)
			if err != nil {
				t.Errorf(err.Error())
			}

			expectedPolicyJson, _ := json.Marshal(tc.expectedPolicy)

			if !reflect.DeepEqual(policy, expectedPolicyJson) {
				t.Errorf("Error Policy doesnt match generareted, Got %s\n, Expected %s", policy, expectedPolicyJson)
			}
		})
	}
}

func Test_generateArn(t *testing.T) {
	policyman := IamPolicyMan{
		arnTemplates: map[string]string{
			"s3": "arn:{{.Partition}}:s3:::{{.BucketName}}",
		},
		arnTmplFieldNames: map[string]string{
			"s3": "BucketName",
		},
	}
	services := []string{"s3", "sts"}
	resources := map[string][]string{
		"s3":  {"AbucketHere", "SomethingsBucket"},
		"sts": {"arn::::::anSTSResource"},
	}

	expected := map[string][]string{
		"s3":  {"arn:aws:s3:::AbucketHere", "arn:aws:s3:::AbucketHere/*", "arn:aws:s3:::SomethingsBucket", "arn:aws:s3:::SomethingsBucket/*"},
		"sts": {"arn::::::anSTSResource", "arn::::::anSTSResource/*"},
	}

	arns, err := policyman.generateArns("123456789", services, resources)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(arns, expected) {
		t.Errorf("Not Equal Expected %v\n Got: %v\n", expected, arns)
	}
}
