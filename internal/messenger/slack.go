package messenger

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/perkbox/cloud-access-bot/internal"

	"github.com/slack-go/slack"
)

type Messenger struct {
	SlackClient slack.Client
}

func NewMessenger(Client *slack.Client) *Messenger {
	return &Messenger{
		SlackClient: *Client,
	}
}

func (m *Messenger) GenerateModal(modalType string, Accounts, LoginRoles []string, hasResourceFinder bool, privateMetadata string, selectedService string) (slack.ModalViewRequest, error) {
	switch modalType {
	case "firstView":
		firstViewTmplVals := Template{
			IsIamService: false,
			Accounts:     Accounts,
			LoginRoles:   LoginRoles,
		}
		return GetRequestAccessModal(firstViewTmplVals)

	case "accountSelectView":
		accountSelectViewTmplVals := Template{
			IsIamService: true,
			Accounts:     Accounts,
			LoginRoles:   LoginRoles,
		}
		return GetRequestAccessModal(accountSelectViewTmplVals)

	case "servicesView":
		servicesViewTmplVals := Template{
			IsIamService:     true,
			IsActionSelector: true,
			IsResourcesText:  true,
			SelectedService:  selectedService,
			PrivateMetadata:  privateMetadata,
			Accounts:         Accounts,
			LoginRoles:       LoginRoles,
		}
		if hasResourceFinder {
			servicesViewTmplVals.IsExternalResourcesSelector = true
			servicesViewTmplVals.IsResourcesText = false
		}

		return GetRequestAccessModal(servicesViewTmplVals)
	}

	return slack.ModalViewRequest{}, nil
}

func (m *Messenger) UpdateMessageFromMessageObj(requestId string, approvalMsgObj []internal.ApprovalMsgObj, msgContents []slack.Block) error {
	for _, msg := range approvalMsgObj {
		if _, _, _, err := m.SlackClient.UpdateMessage(
			msg.Channel, msg.Ts,
			slack.MsgOptionAttachments(slack.Attachment{Fields: []slack.AttachmentField{{}},
				Footer: requestId, Ts: json.Number(strconv.Itoa(int(time.Now().Unix())))}),
			slack.MsgOptionBlocks(msgContents...),
		); err != nil {
			return fmt.Errorf("func:UpdateMessageFromMessageObj: error updating message from audit object %s", err.Error())
		}
	}
	return nil
}

func (m *Messenger) PostBlockMessage(channelId string, msgContents []slack.Block, requestId string) (string, string, error) {
	return m.SlackClient.PostMessage(channelId, slack.MsgOptionAttachments(slack.Attachment{Fields: []slack.AttachmentField{{}},
		Footer: requestId, Ts: json.Number(strconv.Itoa(int(time.Now().Unix())))}),
		slack.MsgOptionBlocks(msgContents...))
}

func (m *Messenger) PostSimpleMessage(channelId string, msgText string, requestId string) error {
	_, _, err := m.SlackClient.PostMessage(channelId, slack.MsgOptionText(msgText, false), slack.MsgOptionAttachments(slack.Attachment{Fields: []slack.AttachmentField{{}},
		Footer: requestId, Ts: json.Number(strconv.Itoa(int(time.Now().Unix())))}))
	if err != nil {
		return fmt.Errorf("func:PostSimpleMessage: Error Posting Mesage to Requesting User:%s", err.Error())
	}
	return nil
}

func (m *Messenger) GetUserIdsFromGroup(groups []string) ([]string, error) {
	var approverIds []string
	grp, err := m.SlackClient.GetUserGroups()
	if err != nil {
		return nil, fmt.Errorf("func:GetUserIdsFromGroup: error getting users from group %s", err.Error())
	}

	for _, v := range grp {
		for _, approver := range groups {
			if strings.EqualFold(v.Name, approver) {
				approvers, err := m.SlackClient.GetUserGroupMembers(v.ID)
				if err != nil {
					return nil, err
				}
				approverIds = append(approverIds, approvers...)
			}
		}
	}

	return approverIds, nil
}
