package controllers

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/perkbox/cloud-access-bot/internal/messenger"

	uuid "github.com/satori/go.uuid"

	"github.com/perkbox/cloud-access-bot/internal"

	"github.com/perkbox/cloud-access-bot/internal/utils"

	"github.com/perkbox/cloud-access-bot/internal/settings"

	"github.com/sirupsen/logrus"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// SlashCommandController We create a structure to let us use dependency injection
type SlashCommandController struct {
	EventHandler *socketmode.SocketmodeHandler
	Service      internal.Service
	Settings     settings.Settings
}

func NewSlashCommandController(cfg settings.Settings, service *internal.Service, eventhandler *socketmode.SocketmodeHandler) SlashCommandController {
	c := SlashCommandController{
		EventHandler: eventhandler,
		Service:      *service,
		Settings:     cfg,
	}

	// Handle and log connection to slack confirmation
	c.EventHandler.Handle(
		socketmode.RequestTypeHello,
		func(event *socketmode.Event, client *socketmode.Client) {
			logrus.Infof("Sucessfully Connected to Slack")
		},
	)

	// Register callback for the Request command, start of Request process
	c.EventHandler.HandleSlashCommand(
		fmt.Sprintf("/%s", c.Settings.GetRequestCommand()),
		c.handleRequestStart,
	)
	logrus.Infof("Slash command %s Registered", c.Settings.GetRequestCommand())

	// Event called when user submits the model once completed all fields
	c.EventHandler.HandleInteraction(
		slack.InteractionTypeViewSubmission,
		c.requestModelSubmitted,
	)

	// Model Update Event, Triggered User Selects an account from the model.
	c.EventHandler.HandleInteractionBlockAction(
		messenger.AccountSelectorActionId,
		c.updateViewAccountSelect,
	)

	// Model Update Event, triggered when a cloud Service is selected.
	c.EventHandler.HandleInteractionBlockAction(
		messenger.IamServicesSelectorActionID,
		c.updateViewServices,
	)

	// Event triggered by auto complete searching, shows users results for live searchable fields
	c.EventHandler.HandleInteraction(
		slack.InteractionTypeBlockSuggestion,
		c.SuggestServices,
	)

	// Approve or Deny Permission Request
	// Approve
	c.EventHandler.HandleInteractionBlockAction(
		"approve",
		c.handleReqApproval,
	)
	// Deny
	c.EventHandler.HandleInteractionBlockAction(
		"deny",
		c.handleReqApproval,
	)

	return c
}

func (c SlashCommandController) SuggestServices(evt *socketmode.Event, clt *socketmode.Client) {
	var paylaodOpts messenger.Options
	suggestCallabck, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Errorf("ERROR converting event to Slash Command: %v", ok)
		return
	}

	switch suggestCallabck.ActionID {
	case messenger.IamServicesSelectorActionID:
		paylaodOpts = messenger.SliceToOptions(c.Service.GetServicesWithFilter(strings.ToLower(suggestCallabck.Value)), "plain_text")

	case messenger.IamServiceActionSelectorActionID:
		selService := strings.Split(suggestCallabck.BlockID, ":")[1]
		paylaodOpts = messenger.MapToOptions(c.Service.GetActionsWithFilter(selService, strings.ToLower(suggestCallabck.Value)), "plain_text")

	case messenger.IamResourcesSelectorActionID:
		selService := strings.Split(suggestCallabck.BlockID, ":")[1]
		resources, _ := c.Service.GetCloudResourcesForService(suggestCallabck.Value, selService, suggestCallabck.View.PrivateMetadata)
		paylaodOpts = messenger.MapToOptions(resources, "plain_text")
	default:
		logrus.Warnf("Unknown action")
		return
	}

	payload := socketmode.Response{
		EnvelopeID: evt.Request.EnvelopeID,
		Payload:    paylaodOpts,
	}

	clt.Send(payload)
}

func (c SlashCommandController) updateViewAccountSelect(evt *socketmode.Event, clt *socketmode.Client) {
	actionCallabck, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Errorf("ERROR converting event to Update View: %v", ok)
		return
	}

	client := clt.GetApiClient()
	clt.Ack(*evt.Request)

	viewBody := c.Service.Messenger.GenerateModal("accountSelectView", c.Settings.GetAccountNames(), c.Settings.GetLoginRoles(), false, "", "")

	_, err := client.UpdateView(viewBody, actionCallabck.View.ExternalID, "", actionCallabck.View.ID)
	if err != nil {
		logrus.Errorf("Error Updating View %s", err.Error())
		return
	}
}

func (c SlashCommandController) updateViewServices(evt *socketmode.Event, clt *socketmode.Client) {
	actionCallabck, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Errorf("ERROR converting event to Update View: %v", ok)
		return
	}
	client := clt.GetApiClient()

	selService := actionCallabck.View.State.Values[messenger.IamServicesSelectorActionID][messenger.IamServicesSelectorActionID].SelectedOption.Value
	selAccount := actionCallabck.View.State.Values[messenger.AccountSelectorActionId][messenger.AccountSelectorActionId].SelectedOption.Value
	clt.Ack(*evt.Request)

	_, hasResourceFinder := c.Service.GetCloudResourcesForService("", selService, selAccount)
	viewBody := c.Service.Messenger.GenerateModal("servicesView", c.Settings.GetAccountNames(), c.Settings.GetLoginRoles(), hasResourceFinder, selAccount, selService)

	_, err := client.UpdateView(viewBody, actionCallabck.View.ExternalID, "", actionCallabck.View.ID)
	if err != nil {
		logrus.Errorf("Error Updating View %s", err.Error())
		return
	}
}

func (c SlashCommandController) handleRequestStart(evt *socketmode.Event, clt *socketmode.Client) {
	// we need to cast our socket mode.Event into a Slash Command
	command, ok := evt.Data.(slack.SlashCommand)
	if !ok {
		logrus.Errorf("ERROR converting event to Slash Command: %v", ok)
		return
	}

	clt.Ack(*evt.Request)
	client := clt.GetApiClient()

	viewBody := c.Service.Messenger.GenerateModal("firstView", c.Settings.GetAccountNames(), c.Settings.GetLoginRoles(), false, "", "")

	_, err := client.OpenView(command.TriggerID, viewBody)
	if err != nil {
		logrus.Errorf("Error opening slack model Err: %s", err.Error())
		return
	}

	if err != nil {
		logrus.Errorf("ERROR while sending message for /request: %v", err)
		return
	}

}

func (c SlashCommandController) requestModelSubmitted(evt *socketmode.Event, clt *socketmode.Client) {
	var (
		approvalMsgs    []internal.ApprovalMsgObj
		policyResources []string
	)

	viewCallabck, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Errorf("ERROR converting event to Slash Command: %v", ok)
		return
	}

	requestDescription := viewCallabck.View.State.Values[messenger.RequestDescriptionId][messenger.RequestDescriptionId].Value
	requestDuration := viewCallabck.View.State.Values[messenger.TimeInputID][messenger.TimeInputID].Value
	loginRole := viewCallabck.View.State.Values[messenger.LoginRoleSelector][messenger.LoginRoleSelector].SelectedOption.Value
	selAccount := viewCallabck.View.State.Values[messenger.AccountSelectorActionId][messenger.AccountSelectorActionId].SelectedOption.Value
	selService := viewCallabck.View.State.Values[messenger.IamServicesSelectorActionID][messenger.IamServicesSelectorActionID].SelectedOption.Value
	actionIds := messenger.GetValuesFromSelectedOptions(viewCallabck.View.State.Values[fmt.Sprintf("%s:%s", messenger.IamServiceActionSelectorActionID, selService)][messenger.IamServiceActionSelectorActionID].SelectedOptions)
	selActions := c.Service.IdentityData.FindActionsById(actionIds)

	selResources := messenger.GetValuesFromSelectedOptions(viewCallabck.View.State.Values[fmt.Sprintf("%s:%s", messenger.IamResourcesSelectorActionID, selService)][messenger.IamResourcesSelectorActionID].SelectedOptions)

	if len(selResources) != 0 {
		policyResources = c.Service.FindSelectedCloudResoucesNames(selService, viewCallabck.View.PrivateMetadata, selResources)
	}
	if len(selResources) == 0 {
		policyResources = utils.SplitFreeString(viewCallabck.View.State.Values[fmt.Sprintf("%s:%s", messenger.IamResourcesSelectorActionID, selService)][messenger.IamResourcesSelectorActionID].Value)
		invalidResources := c.Service.Cloud.ValidateResourcesFormat(policyResources)
		if len(invalidResources) > 0 {
			resp := slack.NewErrorsViewSubmissionResponse(map[string]string{
				fmt.Sprintf("%s:%s", messenger.IamResourcesSelectorActionID, selService): fmt.Sprintf("Invalid Resource Input: %+v", invalidResources),
			})
			clt.Ack(*evt.Request, resp)
			return
		}
	}

	clt.Ack(*evt.Request)

	client := clt.GetApiClient()
	id, err := client.GetUserInfo(viewCallabck.User.ID)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	roleId, err := c.Service.GetCloudUserId(viewCallabck.View.PrivateMetadata, loginRole)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	//----------- Replace CarriageReturn new line with literal new line char
	re := regexp.MustCompile(`\r?\n`)
	requestDescription = re.ReplaceAllString(requestDescription, `\n`)
	//-------------------------------------------------------

	auditObj := internal.AuditObject{
		RequestId:   uuid.NewV4().String(),
		RequestTime: time.Now(),
		Description: requestDescription,
		Duration:    requestDuration,
		UserId:      viewCallabck.User.ID,
		CloudUserId: fmt.Sprintf("%s:%s", roleId, id.Profile.Email),
		LoginRole:   loginRole,
		AccountId:   c.Settings.GetAccountNumFromName(selAccount),
		Services:    []string{selService},
		Actions: map[string][]string{
			selService: selActions,
		},
		Resources: map[string][]string{
			selService: policyResources,
		},
	}

	blocks, err := messenger.GetRequestApprovalBlocks(auditObj, false, "")
	if err != nil {
		logrus.WithField("func", "requestModelSubmitted").Errorf(" Error generating template for approval block %s\n", err.Error())
	}

	approvers, err := c.Service.Messenger.GetUserIdsFromGroup(c.Settings.ApprovalGroups)
	if err != nil {
		logrus.WithField("func", "requestModelSubmitted").Errorf(" Error getting users from group: %s\n", err.Error())
	}

	for _, approver := range approvers {
		respChan, timestamp, err := c.Service.Messenger.PostBlockMessage(approver, blocks, auditObj.RequestId)
		if err != nil {
			logrus.WithField("func", "requestModelSubmitted").Fatalf("Error posting aproval message: %s", err.Error())
			return
		}
		msgObj := internal.ApprovalMsgObj{Ts: timestamp, Channel: respChan}
		approvalMsgs = append(approvalMsgs, msgObj)
	}

	auditObj.ApprovalMessages = approvalMsgs

	err = c.Service.SetAuditObj(auditObj)
	if err != nil {
		logrus.Errorf("errors Setting obj in cache %s", err.Error())
	}

	err = c.Service.Messenger.PostSimpleMessage(viewCallabck.User.ID, "Request raised and sent to approvers", auditObj.RequestId)
	if err != nil {
		logrus.WithField("func", "handleinteaction").Errorf("Error sending messages: %s ", err.Error())
	}
}

func (c SlashCommandController) handleReqApproval(evt *socketmode.Event, clt *socketmode.Client) {
	var (
		approverMsgText  string
		requesterMsgText string
	)

	approvalCallabck, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.WithField("func", "handleReqApproval").Errorf("ERROR converting event to Slash Command: %v", ok)
		return
	}

	requestId := strings.Split(approvalCallabck.ActionCallback.BlockActions[0].Value, ":")[0]
	userId := strings.Split(approvalCallabck.ActionCallback.BlockActions[0].Value, ":")[1]

	cachedObject, err := c.Service.GetAuditObj(userId, requestId)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	clt.Ack(*evt.Request)

	switch approvalCallabck.ActionCallback.BlockActions[0].ActionID {
	case "approve":
		approverMsgText = fmt.Sprintf(":white_check_mark: Request Approved by <@%s>", approvalCallabck.User.ID)
		err := c.Service.Repo.UpdateApprovingUser(cachedObject.UserId, cachedObject.RequestId, approvalCallabck.User.ID)
		if err != nil {
			logrus.Errorf("Error Updating Requesting User %s", err.Error())
		}
		requesterMsgText = "Request Approved, Policy Applied"
	case "deny":
		approverMsgText = fmt.Sprintf(":no_entry_sign: Request Denied by  <@%s>", approvalCallabck.User.ID)
		requesterMsgText = "Request Denied, Please raise a new request"
	}

	responseMSG, _ := messenger.GetRequestApprovalBlocks(cachedObject, true, approverMsgText)

	if err := c.Service.Messenger.UpdateMessageFromMessageObj(cachedObject.RequestId, cachedObject.ApprovalMessages, responseMSG); err != nil {
		logrus.Errorf("error updating message from audit object %s", err.Error())
	}

	cloudAccountName := c.Settings.GetAccountNameAccountNum(cachedObject.AccountId)
	c.Service.FindExpiredPermissions(cloudAccountName, cachedObject.LoginRole, true)

	if approvalCallabck.ActionCallback.BlockActions[0].ActionID == "approve" {

		policyDoc, err := c.Service.GeneratePolicyFromAuditObj(cachedObject)
		if err != nil {
			logrus.Errorf("Error building policy. Err: %s", err.Error())
		}

		err = c.Service.CloudIdentityManager.PutPolicy(cloudAccountName, cachedObject.LoginRole, cachedObject.RequestId, string(policyDoc))
		if err != nil {
			logrus.Errorf("Error building policy. Err: %s", err.Error())
		}
	}

	err = c.Service.Messenger.PostSimpleMessage(cachedObject.UserId, requesterMsgText, cachedObject.RequestId)
	if err != nil {
		logrus.WithField("func", "handleReqApproval").Errorf("Error Posting Mesage to Requesting User:%s", err.Error())
	}
}
