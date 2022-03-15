package messenger

import (
	"testing"

	"github.com/slack-go/slack"
)

type MockClient struct {
	slack.Client
}

func (api *MockClient) GetUserGroups(options ...slack.GetUserGroupsOption) ([]slack.UserGroup, error) {

	return nil, nil
}

func TestName(t *testing.T) {
	test := slack.Client{}

	NewMessenger(&test)

}
