package messenger

import (
	"embed"

	"github.com/perkbox/cloud-access-bot/internal"
)

// Slack Action and block id Constants which can be referenced else where
const (
	RequestDescriptionId             = "accessRequest"
	IamServicesSelectorActionID      = "awsServicesSelect"
	IamServiceActionSelectorActionID = "awsServiceActionSelector"
	IamResourcesSelectorActionID     = "awsResourcesSelector"
	AccountSelectorActionId          = "awsAccountSelector"
	LoginRoleSelector                = "awsLoginRoleSelector"
	TimeInputID                      = "awsTimeInput"
	ApprovedActionID                 = "approve"
	DenyActionID                     = "deny"
)

//go:embed assets/*
var slashCommandAssets embed.FS

// Template struct for options within Slack Model
type Template struct {
	SelectedService string //This is here to store the selected service into the BlockID as block_suggestion doesn't show the state

	LoginRoles      []string
	Accounts        []string
	PrivateMetadata string

	TimeInputID          string
	RequestDescriptionId string
	LoginRoleSelectorId  string
	AccountSelectorId    string
	ServiceActionId      string
	ActionsActionId      string
	ActionsBlockId       string
	ResourcesActionId    string
	ResourcesBlockId     string

	IsIamService                bool
	IsActionSelector            bool
	IsExternalResourcesSelector bool
	IsResourcesText             bool
}

// Options Slack Options block for rendering dynamic select lists
type Options struct {
	Options []Option `json:"options"`
}

type Option struct {
	Text  Text   `json:"text"`
	Value string `json:"value"`
}

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

//Approval Message Template vars
type approvalBlockVars struct {
	UserId          string              `json:"userId"`
	RequestId       string              `json:"requestId"`
	Description     string              `json:"description"`
	CloudUserId     string              `json:"cloudUserId"`
	LoginRole       string              `json:"loginRole"`
	AccountId       string              `json:"accountId"`
	Duration        string              `json:"duration"`
	Services        []string            `json:"services"`
	Actions         map[string][]string `json:"actions"`
	Resources       map[string][]string `json:"resources"`
	GotResponse     bool                `json:"gotResponse"`
	ResponseMSG     string              `json:"responseMsg"`
	ApproveActionId string              `json:"approveActionId"`
	DenyActionId    string              `json:"denyActionId"`
}

// auditObjtoApprovalBlockVars Type conversion function for rendering APPROVAL MESSAGE
func auditObjtoApprovalBlockVars(auditObj internal.AuditObject) approvalBlockVars {
	return approvalBlockVars{
		UserId:          auditObj.UserId,
		RequestId:       auditObj.RequestId,
		Description:     auditObj.Description,
		CloudUserId:     auditObj.CloudUserId,
		LoginRole:       auditObj.LoginRole,
		AccountId:       auditObj.AccountId,
		Duration:        auditObj.Duration,
		Services:        auditObj.Services,
		Actions:         auditObj.Actions,
		Resources:       auditObj.Resources,
		ApproveActionId: ApprovedActionID,
		DenyActionId:    DenyActionID,
	}
}
