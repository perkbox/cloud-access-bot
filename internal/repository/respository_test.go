package repository

import (
	"reflect"
	"testing"
	"time"

	"github.com/perkbox/cloud-access-bot/internal"
)

func Test_convertConversions(t *testing.T) {
	internalAudit := internal.AuditObject{
		UserId:      "ABC234USER",
		RequestId:   "REQID",
		Description: "A Description Goes Here",
		RequestTime: time.Date(
			2021, 10, 21, 21, 21, 21, 651387237, time.UTC),
		ApprovingUser:    "BOB",
		ApprovalMessages: nil,
		CloudUserId:      "ABC123:USERA@gmail.com",
		LoginRole:        "SSO-USER",
		AccountRole:      "",
		AccountId:        "123456789",
		Duration:         "120",
		Services:         []string{"s3", "sts"},
		Actions: map[string][]string{
			"s3":  {"GetBucket", "ReadBucket"},
			"sts": {"ReadSomething"},
		},
		Resources: map[string][]string{
			"s3":  {"AbucketHere", "SomethingsBucket"},
			"sts": {"arn::::::anSTSResource"},
		},
	}

	expectedRepoObj := repoAuditObject{
		PK:               internalAudit.UserId,
		SK:               internalAudit.RequestId,
		Description:      internalAudit.Description,
		RequestTime:      internalAudit.RequestTime,
		ApprovingUser:    internalAudit.ApprovingUser,
		ApprovalMessages: internalAudit.ApprovalMessages,
		CloudUserId:      internalAudit.CloudUserId,
		LoginRole:        internalAudit.LoginRole,
		AccountRole:      internalAudit.AccountRole,
		AccountId:        internalAudit.AccountId,
		Duration:         internalAudit.Duration,
		Services:         internalAudit.Services,
		Actions:          internalAudit.Actions,
		Resources:        internalAudit.Resources,
	}

	repoObj := convertToRepoObj(internalAudit)
	auditObj := repoObj.convertFromRepoObject()

	if !reflect.DeepEqual(repoObj, expectedRepoObj) {
		t.Errorf("TO REPO OBJ FAIL: Expected  %+v, \n Got: %+v", expectedRepoObj, repoObj)
	}

	if !reflect.DeepEqual(auditObj, internalAudit) {
		t.Errorf("TO COVERT REPO TO AUDIT OBJ FAIL: Expected  %+v, \n Got: %+v", auditObj, internalAudit)
	}
}
