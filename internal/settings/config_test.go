package settings

import (
	"reflect"
	"sort"
	"testing"

	"gopkg.in/yaml.v3"
)

func Test_loadConfigFromBytes(t *testing.T) {
	yamlBytes, _ := yaml.Marshal(MockSettings)
	settings, _ := loadConfigFromBytes(yamlBytes)

	if !reflect.DeepEqual(settings, MockSettings) {
		t.Errorf("Got: %v, Expected: %v", settings, MockSettings)
	}
}

func TestSettings_GetAccountNames(t *testing.T) {
	mockSettings, _ := NewConfigMock()
	expected := []string{"perkbox-mock", "perkbox-mock2"}

	accountNames := mockSettings.GetAccountNames()

	if !array_sorted_equal(expected, accountNames) {
		t.Errorf("Got: %s, Expected: %s", accountNames, expected)
	}
}

func TestSettings_GetAccountNumFromName(t *testing.T) {
	mockSettings, _ := NewConfigMock()
	expected := "123456789"

	accountNum := mockSettings.GetAccountNumFromName("perkbox-mock")

	if !reflect.DeepEqual(expected, accountNum) {
		t.Errorf("Got: %s, Expected: %s", accountNum, expected)
	}
}

func TestSettings_GetRoleArn(t *testing.T) {
	mockSettings, _ := NewConfigMock()
	expected := "arn:aws:iam::123456789:role/cloud-access-bot"

	accountNames, _ := mockSettings.GetRoleArn("perkbox-mock")

	if !reflect.DeepEqual(expected, accountNames) {
		t.Errorf("Got: %s, Expected: %s", accountNames, expected)
	}
}

func TestSettings_GetAccountNameAccountNum(t *testing.T) {
	mockSettings, _ := NewConfigMock()
	expected := "perkbox-mock"

	accountName := mockSettings.GetAccountNameAccountNum("123456789")

	if !reflect.DeepEqual(expected, accountName) {
		t.Errorf("Got: %s, Expected: %s", accountName, expected)
	}
}

func array_sorted_equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	a_copy := make([]string, len(a))
	b_copy := make([]string, len(b))

	copy(a_copy, a)
	copy(b_copy, b)

	sort.Strings(a_copy)
	sort.Strings(b_copy)

	return reflect.DeepEqual(a_copy, b_copy)
}
