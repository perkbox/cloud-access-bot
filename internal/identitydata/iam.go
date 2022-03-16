package identitydata

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/perkbox/cloud-access-bot/internal/utils"
	"github.com/sirupsen/logrus"
)

//go:embed assets/*
var pToolAssets embed.FS

func NewIamDefinitions() *IamDefinitions {
	var iamDefinitions IamDefinitions

	plan, err := pToolAssets.ReadFile("assets/iam-definition.json")
	if err != nil {
		logrus.Fatalf("func:NewIamDefinitions, error reading in file from assets Err: %s ", err)
	}

	err = json.Unmarshal(plan, &iamDefinitions)
	if err != nil {
		logrus.Fatalf("func:NewIamDefinitions, error unmarhaslling iam-definition.json Err: %s", err)
	}

	for serK, serV := range iamDefinitions {
		for k := range serV.Privileges {
			p := iamDefinitions[serK].Privileges[k]
			p.Id = utils.HashString(fmt.Sprintf("%s:%s", serV.ServiceName, k), 4)
			iamDefinitions[serK].Privileges[k] = p
		}
	}

	return &iamDefinitions
}

func (i IamDefinitions) GetResourceTmplDetails(service string) (string, string) {
	var data ArnData

	switch service {
	case "dynamodb":
		data = ArnData{ResourceType: "table", TmplFieldName: "TableName"}
	case "s3":
		data = ArnData{ResourceType: "bucket", TmplFieldName: "BucketName"}
	case "lambda":
		data = ArnData{ResourceType: "function", TmplFieldName: "FunctionName"}
	case "sqs":
		data = ArnData{ResourceType: "queue", TmplFieldName: "QueueName"}
	case "sns":
		data = ArnData{ResourceType: "topic", TmplFieldName: "TopicName"}
	}

	if tmplstr, ok := i[service].Resources[data.ResourceType]; ok {
		return tmplstr.ArnTmpl, data.TmplFieldName
	}

	return "", ""
}

func (i IamDefinitions) GetIamServices() []string {
	var IamServiceNames []string

	for service := range i {
		IamServiceNames = append(IamServiceNames, service)
	}

	return IamServiceNames
}

//ACTIONS---------------------------------

func (i IamDefinitions) GetActionsForService(serviceName string) map[string]string {
	IamActions := make(map[string]string)

	if v, ok := i[serviceName]; ok {
		for _, v := range v.Privileges {

			IamActions[v.Privilege] = v.Id

		}
	}

	return IamActions
}

func (i IamDefinitions) FindActionsById(ids []string) []string {
	var actions []string
	for _, serV := range i {
		for _, privV := range serV.Privileges {
			if utils.Contains(ids, privV.Id) {
				actions = append(actions, fmt.Sprintf("%s:%s", serV.Prefix, privV.Privilege))
			}
		}
	}

	return actions
}
