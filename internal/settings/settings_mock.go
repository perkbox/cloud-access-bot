package settings

var MockSettings = Settings{
	Accounts: map[string]Account{
		"perkbox-mock": {
			IamRole:       "arn:aws:iam::123456789:role/cloud-access-bot",
			AccountNumber: "123456789",
		},
		"perkbox-mock2": {
			IamRole:       "arn:aws:iam::9876543221:role/cloud-access-bot",
			AccountNumber: "9876543221",
		},
	},
	IdentiyRegion:  "eu-west-1",
	Regions:        []string{"eu-region-1"},
	LoginRoles:     []string{"SSO-A", "SSO-B"},
	ApprovalGroups: []string{"devops"},
	DynamoDbTable:  "DynamoTableName",
	RequestCommand: "request",
	AuditCommand:   "audit",
}

func NewConfigMock() (Settings, error) {
	return MockSettings, nil
}
