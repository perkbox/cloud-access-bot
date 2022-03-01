package commands

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

func NewRequestCommandHandler(cfg settings.Settings, service *internal.Service, eventhandler *socketmode.SocketmodeHandler) SlashCommandController {
	c := SlashCommandController{
		EventHandler: eventhandler,
		Service:      *service,
		Settings:     cfg,
	}

	// Handle and log connection to slack confirmation
	c.EventHandler.Handle(
		socketmode.RequestTypeHello,
		func(event *socketmode.Event, client *socketmode.Client) {
			logrus.Infof("Successfully Connected to Slack")
		},
	)

	// Register callback for the Request command, start of Request process
	c.EventHandler.HandleSlashCommand(
		fmt.Sprintf("/%s", c.Settings.GetRequestCommand()),
		c.handleRequestStart,
	)

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
		messenger.ApprovedActionID,
		c.handleReqApproval,
	)
	// Deny
	c.EventHandler.HandleInteractionBlockAction(
		messenger.DenyActionID,
		c.handleReqApproval,
	)

	return c
}

func (c *SlashCommandController) SuggestServices(evt *socketmode.Event, clt *socketmode.Client) {
	var payloadOpts messenger.Options
	suggestCallback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Errorf("ERROR converting event to Slash Command: %v", ok)
		return
	}

	switch suggestCallback.ActionID {
	case messenger.IamServicesSelectorActionID:
		payloadOpts = messenger.SliceToOptions(c.Service.GetServicesWithFilter(strings.ToLower(suggestCallback.Value)), "plain_text")

	case messenger.IamServiceActionSelectorActionID:
		selService := strings.Split(suggestCallback.BlockID, ":")[1]
		payloadOpts = messenger.MapToOptions(c.Service.GetActionsWithFilter(selService, strings.ToLower(suggestCallback.Value)), "plain_text")

	case messenger.IamResourcesSelectorActionID:
		selService := strings.Split(suggestCallback.BlockID, ":")[1]
		resources, _ := c.Service.GetCloudResourcesForService(suggestCallback.Value, selService, suggestCallback.View.PrivateMetadata)
		payloadOpts = messenger.MapToOptions(resources, "plain_text")
	default:
		logrus.Warnf("Unknown action")
		return
	}

	clt.Send(socketmode.Response{
		EnvelopeID: evt.Request.EnvelopeID,
		Payload:    payloadOpts,
	})
}

func (c *SlashCommandController) updateViewAccountSelect(evt *socketmode.Event, clt *socketmode.Client) {
	actionCallback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Errorf("Error converting event to Update View")
		return
	}

	client := clt.GetApiClient()
	clt.Ack(*evt.Request)

	viewBody, err := c.Service.Messenger.GenerateModal("accountSelectView", c.Settings.GetAccountNames(), c.Settings.GetLoginRoles(), false, "", "")
	if err != nil {
		logrus.Errorf("Error Getting Modal accountSelectView Err: %s", err)
		return
	}

	_, err = client.UpdateView(viewBody, actionCallback.View.ExternalID, "", actionCallback.View.ID)
	if err != nil {
		logrus.Errorf("Error Updating View.. Err: %s", err.Error())
		return
	}
}

func (c *SlashCommandController) updateViewServices(evt *socketmode.Event, clt *socketmode.Client) {
	actionCallback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Errorf("ERROR converting event to Update View")
		return
	}
	client := clt.GetApiClient()

	selService := actionCallback.View.State.Values[messenger.IamServicesSelectorActionID][messenger.IamServicesSelectorActionID].SelectedOption.Value
	selAccount := actionCallback.View.State.Values[messenger.AccountSelectorActionId][messenger.AccountSelectorActionId].SelectedOption.Value
	clt.Ack(*evt.Request)

	_, hasResourceFinder := c.Service.GetCloudResourcesForService("", selService, selAccount)
	viewBody, err := c.Service.Messenger.GenerateModal("servicesView", c.Settings.GetAccountNames(), c.Settings.GetLoginRoles(), hasResourceFinder, selAccount, selService)
	if err != nil {
		logrus.Errorf("Error Getting Modal servicesView Err: %s", err)
		return
	}

	_, err = client.UpdateView(viewBody, actionCallback.View.ExternalID, "", actionCallback.View.ID)
	if err != nil {
		logrus.WithField("User", actionCallback.User.Name).Errorf("Error Updating View Err: %s", err.Error())
		return
	}
}

func (c *SlashCommandController) handleRequestStart(evt *socketmode.Event, clt *socketmode.Client) {
	// we need to cast our socket mode.Event into a Slash Command
	command, ok := evt.Data.(slack.SlashCommand)
	if !ok {
		logrus.Errorf("ERROR converting event to Slash Command")
		return
	}

	clt.Ack(*evt.Request)
	client := clt.GetApiClient()

	viewBody, err := c.Service.Messenger.GenerateModal("firstView", c.Settings.GetAccountNames(), c.Settings.GetLoginRoles(), false, "", "")
	if err != nil {
		logrus.Errorf("Error Getting Modal firstView Err: %s", err)
		return
	}
	_, err = client.OpenView(command.TriggerID, viewBody)
	if err != nil {
		logrus.Errorf("Error opening slack model Err: %s", err.Error())
		return
	}

	if err != nil {
		logrus.Errorf("ERROR while sending message for /request: Err: %s", err)
		return
	}
}

func (c *SlashCommandController) requestModelSubmitted(evt *socketmode.Event, clt *socketmode.Client) {
	var (
		approvalMsgs    []internal.ApprovalMsgObj
		policyResources []string
	)

	viewCallabck, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Errorf("ERROR converting event to Slash Command")
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
		logrus.Errorf("Error getting User info for %s Err: %s", viewCallabck.User.ID, err.Error())
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
		logrus.Errorf("Error generating template for approval block Err:%s\n", err.Error())
	}

	approvers, err := c.Service.Messenger.GetUserIdsFromGroup(c.Settings.ApprovalGroups)
	if err != nil {
		logrus.Errorf("Error getting users from group:  Err:%s\n", err.Error())
	}

	for _, approver := range approvers {
		respChan, timestamp, err := c.Service.Messenger.PostBlockMessage(approver, blocks, auditObj.RequestId)
		if err != nil {
			logrus.Fatalf("Error posting approval message: Err:%s", err.Error())
			return
		}
		msgObj := internal.ApprovalMsgObj{Ts: timestamp, Channel: respChan}
		approvalMsgs = append(approvalMsgs, msgObj)
	}

	auditObj.ApprovalMessages = approvalMsgs

	err = c.Service.SetAuditObj(auditObj)
	if err != nil {
		logrus.Errorf("errors Setting obj in cache Err:%s", err.Error())
		return
	}

	err = c.Service.Messenger.PostSimpleMessage(viewCallabck.User.ID, "Request raised and sent to approvers", auditObj.RequestId)
	if err != nil {
		logrus.Errorf("Error sending approval sent confirmation Err:%s ", err.Error())
	}
}

func (c *SlashCommandController) handleReqApproval(evt *socketmode.Event, clt *socketmode.Client) {
	var (
		approverMsgText  string
		requesterMsgText string
	)

	approvalCallback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Errorf("ERROR converting event to Slash Command")
		return
	}

	requestId := strings.Split(approvalCallback.ActionCallback.BlockActions[0].Value, ":")[0]
	userId := strings.Split(approvalCallback.ActionCallback.BlockActions[0].Value, ":")[1]

	cachedObject, err := c.Service.GetAuditObj(userId, requestId)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	clt.Ack(*evt.Request)

	switch approvalCallback.ActionCallback.BlockActions[0].ActionID {
	case messenger.ApprovedActionID:
		approverMsgText = fmt.Sprintf(":white_check_mark: Request Approved by <@%s>", approvalCallback.User.ID)
		err := c.Service.Repo.UpdateApprovingUser(cachedObject.UserId, cachedObject.RequestId, approvalCallback.User.ID)
		if err != nil {
			logrus.Errorf("Error Updating Requesting User %s", err.Error())
		}
		requesterMsgText = "Request Approved, Policy Applied"
	case messenger.DenyActionID:
		approverMsgText = fmt.Sprintf(":no_entry_sign: Request Denied by  <@%s>", approvalCallback.User.ID)
		requesterMsgText = "Request Denied, Please raise a new request"
	}

	responseMSG, _ := messenger.GetRequestApprovalBlocks(cachedObject, true, approverMsgText)

	if err := c.Service.Messenger.UpdateMessageFromMessageObj(cachedObject.RequestId, cachedObject.ApprovalMessages, responseMSG); err != nil {
		logrus.Errorf("error updating message from audit object  Err:%s", err.Error())
	}

	cloudAccountName := c.Settings.GetAccountNameAccountNum(cachedObject.AccountId)
	c.Service.FindExpiredPermissions(cloudAccountName, cachedObject.LoginRole, true)

	if approvalCallback.ActionCallback.BlockActions[0].ActionID == "approve" {

		policyDoc, err := c.Service.GeneratePolicyFromAuditObj(cachedObject)
		if err != nil {
			logrus.Errorf("Error building policy. Err: Err:%s", err.Error())
			return
		}

		err = c.Service.CloudIdentityManager.PutPolicy(cloudAccountName, cachedObject.LoginRole, cachedObject.RequestId, string(policyDoc))
		if err != nil {
			logrus.Errorf("Error building policy. Err: %s", err.Error())
			return
		}
	}

	err = c.Service.Messenger.PostSimpleMessage(cachedObject.UserId, requesterMsgText, cachedObject.RequestId)
	if err != nil {
		logrus.Errorf("Error Posting Mesage to Requesting User  Err:%s", err.Error())
	}
}
