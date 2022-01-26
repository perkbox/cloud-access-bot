package repository

import (
	"time"

	"github.com/perkbox/cloud-access-bot/internal"
)

type repoAuditObject struct {
	PK               string                    `json:"pk"` //RequestingUserId   USER#USERID
	SK               string                    `json:"sk"` //UUID 			  REQUESTID#UUID
	Description      string                    `json:"description"`
	RequestTime      time.Time                 `json:"requestTime"`
	ApprovingUser    string                    `json:"approvingUser"`
	ApprovalMessages []internal.ApprovalMsgObj `json:"approvalMessages"`
	CloudUserId      string                    `json:"cloudUserId"`
	LoginRole        string                    `json:"loginRole"`
	AccountRole      string                    `json:"accountRole"`
	AccountId        string                    `json:"accountId"`
	Duration         string                    `json:"duration"`
	Services         []string                  `json:"services"`
	Actions          map[string][]string       `json:"actions"`
	Resources        map[string][]string       `json:"resources"`
}

func (o repoAuditObject) convertFromRepoObject() internal.AuditObject {
	return internal.AuditObject{
		UserId:           o.PK,
		RequestId:        o.SK,
		Description:      o.Description,
		RequestTime:      o.RequestTime,
		ApprovingUser:    o.ApprovingUser,
		ApprovalMessages: o.ApprovalMessages,
		CloudUserId:      o.CloudUserId,
		LoginRole:        o.LoginRole,
		AccountRole:      o.AccountRole,
		AccountId:        o.AccountId,
		Duration:         o.Duration,
		Services:         o.Services,
		Actions:          o.Actions,
		Resources:        o.Resources,
	}
}

func convertToRepoObj(object internal.AuditObject) repoAuditObject {
	return repoAuditObject{
		PK:               object.UserId,
		SK:               object.RequestId,
		Description:      object.Description,
		RequestTime:      object.RequestTime,
		ApprovingUser:    object.ApprovingUser,
		ApprovalMessages: object.ApprovalMessages,
		CloudUserId:      object.CloudUserId,
		LoginRole:        object.LoginRole,
		AccountRole:      object.AccountRole,
		AccountId:        object.AccountId,
		Duration:         object.Duration,
		Services:         object.Services,
		Actions:          object.Actions,
		Resources:        object.Resources,
	}
}
