package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/slack-go/slack"

	"github.com/sirupsen/logrus"

	"github.com/perkbox/cloud-access-bot/internal/utils"
)

type Cloud interface {
	ResourceFinder(service string, accountName string) ([]string, bool)
	ValidateResourcesFormat(resources []string) []string
}

type Messenger interface {
	PostSimpleMessage(channelId string, msgText string, requestId string) error
	PostBlockMessage(channelId string, msgContents []slack.Block, requestId string) (string, string, error)
	GetUserIdsFromGroup(groups []string) ([]string, error)
	UpdateMessageFromMessageObj(requestId string, approvalMsgObj []ApprovalMsgObj, msgContents []slack.Block) error
	GenerateModal(modalType string, Accounts, LoginRoles []string, hasResourceFinder bool, privateMetadata string, selectedService string) slack.ModalViewRequest
}

type Repo interface {
	QueryAuditObjs(UserID string) ([]AuditObject, error)
	GetAuditObj(UserID, RequestId string) (AuditObject, error)
	SetAuditObj(requestObj AuditObject) error
	UpdateApprovingUser(UserID, RequestId, approvingUser string) error
}

type IdentityData interface {
	GetResourceTmplDetails(service string) (string, string)
	GetActionsForService(serviceName string) map[string]string
	FindActionsById(ids []string) []string
	GetIamServices() []string
}

type CloudIdentityManager interface {
	IsPolicyExpired(policy string) (bool, error)
	GeneratePolicyFromAuditObj(curTime time.Time, object AuditObject, tmpls, tmplFieldNmaes map[string]string) ([]byte, error)

	GetCloudUserId(accountName string, roleName string) (string, error)
	PutPolicy(accountName, roleName, policyName, policy string) error
	FindPolicysForRole(accountName, roleName string) (map[string]string, error)
	DeletePolicys(accountName, roleName string, policysNames []string) error
}

type Service struct {
	Cloud                Cloud
	Messenger            Messenger
	Repo                 Repo
	CloudIdentityManager CloudIdentityManager
	IdentityData         IdentityData
}

func NewService(cloud Cloud, repo Repo, cim CloudIdentityManager, identitydata IdentityData, messenger Messenger) *Service {
	return &Service{
		cloud,
		messenger,
		repo,
		cim,
		identitydata,
	}
}

//GetServicesWithFilter Gets Services with a filter and returns them as a list. Will return an empty []string if nothing is found.
func (s *Service) GetServicesWithFilter(filter string) []string {
	services := []string{}

	servicesList := s.IdentityData.GetIamServices()

	for _, ser := range servicesList {
		if strings.Contains(ser, filter) {
			services = append(services, ser)
		}
	}

	return services
}

//GetActionsWithFilter Gets Actions for a selected service with a filter. Returns a map[string]string the key is the service name while the value
//is a unique id for each action. Will return an empty map[string]string if there is nothing found.
func (s *Service) GetActionsWithFilter(service string, filter string) map[string]string {
	actions := make(map[string]string)

	actionsMap := s.IdentityData.GetActionsForService(service)

	for k, v := range actionsMap {
		if strings.Contains(strings.ToLower(k), strings.ToLower(filter)) {
			actions[k] = v
		}
	}

	return actions
}

// GetAuditObj Gets the Audit & Message data in the repository based on the inputted UserId and RequestId
func (s *Service) GetAuditObj(UserId, RequestID string) (AuditObject, error) {
	return s.Repo.GetAuditObj(UserId, RequestID)
}

// SetAuditObj Sets the Audit & Message data in the repository
func (s *Service) SetAuditObj(object AuditObject) error {
	return s.Repo.SetAuditObj(object)
}

func (s *Service) GetCloudUserId(accountName string, roleName string) (string, error) {
	return s.CloudIdentityManager.GetCloudUserId(accountName, roleName)
}

// GetCloudResourcesForService Overwrite account used by client in the individuals clients to keep the functions
// ordered and as simple as possible in the service interface
func (s *Service) GetCloudResourcesForService(filter, service, accountname string) (map[string]string, bool) {
	resources, hasFinder := s.Cloud.ResourceFinder(service, accountname)

	resourcesNoDups := utils.RemoveDuplicateStr(resources)

	hashMap := make(map[string]string)
	for _, table := range resourcesNoDups {
		if strings.Contains(table, filter) {
			hashMap[table] = utils.HashString(table, 6)
		}
	}

	return hashMap, hasFinder
}

func (s *Service) FindSelectedCloudResoucesNames(service, accountname string, selected []string) []string {
	var resoruceNames []string
	resources, _ := s.GetCloudResourcesForService("", service, accountname)

	for _, selRes := range selected {
		for resource, hash := range resources {
			if hash == selRes {
				resoruceNames = append(resoruceNames, resource)
			}
		}
	}

	return resoruceNames
}

func (s *Service) FindExpiredPermissions(accountName, role string, delete bool) {
	var expiredPols []string

	forRole, err := s.CloudIdentityManager.FindPolicysForRole(accountName, role)
	if err != nil {
		fmt.Println(err.Error())
	}
	for name, pol := range forRole {
		isExpired, err := s.CloudIdentityManager.IsPolicyExpired(pol)
		if err != nil {
			fmt.Println(err.Error())
		}
		if isExpired {
			expiredPols = append(expiredPols, name)
		}
	}

	if delete {
		err = s.CloudIdentityManager.DeletePolicys(accountName, role, expiredPols)
		if err != nil {
			logrus.Errorf("Error Deleting Policy Err: %s", err.Error())
		}
	}

}

func (s *Service) GeneratePolicyFromAuditObj(object AuditObject) ([]byte, error) {
	arnTemplates := make(map[string]string)
	arnTmplFieldNames := make(map[string]string)

	for _, service := range object.Services {
		tmpl, tmplfield := s.IdentityData.GetResourceTmplDetails(service)
		if tmpl != "" && tmplfield != "" {
			arnTemplates[service] = tmpl
			arnTmplFieldNames[service] = tmplfield
		}
	}

	return s.CloudIdentityManager.GeneratePolicyFromAuditObj(object.RequestTime, object, arnTemplates, arnTmplFieldNames)
}
