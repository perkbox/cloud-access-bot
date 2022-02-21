package internal_test

import (
	"testing"
	"time"

	"github.com/perkbox/cloud-access-bot/internal"

	"github.com/perkbox/cloud-access-bot/internal/policy"

	"github.com/perkbox/cloud-access-bot/internal/identitydata"
)

func Test_GeneratePolicyFromAuditObj(t *testing.T) {
	Ser := internal.Service{IdentityData: identitydata.NewIamDefinitions(), CloudIdentityManager: &policy.IamPolicyMan{}}

	auditObj := internal.AuditObject{
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
	}

	expectedPolicy := `{"Version":"2012-10-17","Statement":[{"Sid":"","Effect":"Allow","Action":["GetBucket","ReadBucket"],"Resource":["arn:aws:s3:::AbucketHere","arn:aws:s3:::AbucketHere/*","arn:aws:s3:::SomethingsBucket","arn:aws:s3:::SomethingsBucket/*"],"Condition":{"StringEqualsIgnoreCase":{"aws:userid":"abc123:email@gmail.com"},"DateGreaterThan":{"aws:CurrentTime":"2021-10-21T21:21:21Z"},"DateLessThan":{"aws:CurrentTime":"2021-10-21T22:21:21Z"}}},{"Sid":"","Effect":"Allow","Action":["ReadSomething"],"Resource":["arn::::::anSTSResource","arn::::::anSTSResource/*"],"Condition":{"StringEqualsIgnoreCase":{"aws:userid":"abc123:email@gmail.com"},"DateGreaterThan":{"aws:CurrentTime":"2021-10-21T21:21:21Z"},"DateLessThan":{"aws:CurrentTime":"2021-10-21T22:21:21Z"}}}]}`

	pol, err := Ser.GeneratePolicyFromAuditObj(auditObj)
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(pol) != expectedPolicy {
		t.Errorf("Error Policy doesnt match generareted, Got %s\n, Expected %s\n", pol, expectedPolicy)
	}
}
