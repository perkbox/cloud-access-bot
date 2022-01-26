package messenger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/perkbox/cloud-access-bot/internal"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func GetRequestAccessModal(tmplVals Template) slack.ModalViewRequest {
	tmplVals.TimeInputID = TimeInputID
	tmplVals.AccountSelectorId = AccountSelectorActionId
	tmplVals.ServiceActionId = IamServicesSelectorActionID
	tmplVals.ActionsActionId = IamServiceActionSelectorActionID
	tmplVals.LoginRoleSelectorId = LoginRoleSelector
	tmplVals.RequestDescriptionId = RequestDescriptionId
	tmplVals.ActionsBlockId = fmt.Sprintf("%s:%s", IamServiceActionSelectorActionID, tmplVals.SelectedService)
	tmplVals.ResourcesActionId = IamResourcesSelectorActionID
	tmplVals.ResourcesBlockId = fmt.Sprintf("%s:%s", IamResourcesSelectorActionID, tmplVals.SelectedService)

	tmpl, _ := renderTemplate(slashCommandAssets, "assets/acessmodal.json", tmplVals)
	str, err := ioutil.ReadAll(&tmpl)
	if err != nil {
		logrus.Errorf(err.Error())
	}

	view := slack.ModalViewRequest{}
	err = json.Unmarshal(str, &view)

	if err != nil {
		logrus.Errorf("---ERROR MARSHALLING func:GetBlocks %s", err.Error())
	}

	return view
}

func GetRequestApprovalBlocks(auditObj internal.AuditObject, gotResponse bool, responseMSG string) ([]slack.Block, error) {
	approvalTmplVals := auditObjtoApprovalBlockVars(auditObj)
	approvalTmplVals.GotResponse = gotResponse
	approvalTmplVals.ResponseMSG = responseMSG

	tmpl, err := renderTemplate(slashCommandAssets, "assets/requestapproval.json", approvalTmplVals)
	if err != nil {
		return nil, fmt.Errorf("error rending template  for approval messages  err: %s", err.Error())
	}
	// we convert the view into a message struct
	view := slack.Msg{}

	str, err := ioutil.ReadAll(&tmpl)
	if err != nil {
		return nil, fmt.Errorf("error reading processed approval msg template  err: %s", err.Error())
	}

	err = json.Unmarshal(str, &view)
	if err != nil {
		return nil, fmt.Errorf("error umarhsalling approval msg into slack.Block  err: %s", err.Error())

	}

	return view.Blocks.BlockSet, nil
}
