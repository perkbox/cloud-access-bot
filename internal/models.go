package internal

import "time"

type AuditObject struct {
	UserId           string              `json:"userId"`
	RequestId        string              `json:"requestId"`
	Description      string              `json:"description"`
	RequestTime      time.Time           `json:"requestTime"`
	ApprovingUser    string              `json:"approvingUser"`
	ApprovalMessages []ApprovalMsgObj    `json:"approvalMessages"`
	CloudUserId      string              `json:"cloudUserId"`
	LoginRole        string              `json:"loginRole"`
	AccountRole      string              `json:"accountRole"`
	AccountId        string              `json:"accountId"`
	Duration         string              `json:"duration"`
	Services         []string            `json:"services"`
	Actions          map[string][]string `json:"actions"`
	Resources        map[string][]string `json:"resources"`
}

type ApprovalMsgObj struct {
	Ts      string `json:"ts"`
	Channel string `json:"channel"`
}
