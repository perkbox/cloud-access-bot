package settings

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Accounts       map[string]Account
	IdentiyRegion  string   `yaml:"identiyRegion"`
	Regions        []string `yaml:"regions"`
	LoginRoles     []string `yaml:"loginRoles"`
	ApprovalGroups []string `yaml:"approvalGroups"`
	DynamoDbTable  string   `yaml:"dynamoDbTable"`
	RequestCommand string   `yaml:"request_command"`
	AuditCommand   string   `yaml:"audit_command"`
}

type Account struct {
	IamRole       string `yaml:"iam_role"`
	AccountNumber string `yaml:"account_number"`
	AccountAlias  string `yaml:"account_alias"`
}

func loadConfigFromBytes(data []byte) (Settings, error) {
	var Conf Settings

	if err := yaml.Unmarshal(data, &Conf); err != nil {
		return Settings{}, err
	}

	return Conf, nil
}

func (s Settings) GetAccountNameAccountNum(accountNum string) string {
	for accName, v := range s.Accounts {
		if accountNum == v.AccountNumber {
			return accName
		}
	}
	return ""
}

func (s Settings) GetLoginRoles() []string {
	return s.LoginRoles
}

func (s Settings) GetRoleArn(accountName string) (string, error) {

	if val, ok := s.Accounts[accountName]; ok {
		if val.IamRole != "" {
			return val.IamRole, nil
		}
		return "", errors.New("error IamRole isnt Set")
	}

	return "", fmt.Errorf("error finding account in config, %s", accountName)
}
func (s Settings) GetAccountNames() []string {
	var accountNames []string
	for acc := range s.Accounts {
		accountNames = append(accountNames, acc)
	}
	return accountNames
}

func (s Settings) GetAccountNumFromName(accountName string) string {
	for accName, v := range s.Accounts {
		if accountName == accName {
			return v.AccountNumber
		}
	}
	return ""
}

func (s Settings) GetAuditCommand() string {
	if s.AuditCommand == "" {
		return "audit"
	}
	return s.AuditCommand
}

func (s Settings) GetRequestCommand() string {
	if s.RequestCommand == "" {
		return "request"
	}
	return s.RequestCommand
}

func (s Settings) GetDynamodbTable() string {
	if s.DynamoDbTable == "" {
		return "request_access_bot"
	}
	return s.DynamoDbTable
}
