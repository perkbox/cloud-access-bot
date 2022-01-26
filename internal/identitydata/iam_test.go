package identitydata

import (
	"reflect"
	"testing"
)

func NewIamDefinitionsMock() *IamDefinitions {
	return &IamDefinitions{
		"testservice": {
			ServiceName:             "testservice",
			Prefix:                  "ts",
			ServiceAuthorizationUrl: "https://google.com",
			Privileges: map[string]Privilege{
				"readacccess": {
					Id:          "aaabbb111",
					Privilege:   "readacccess",
					AccessLevel: "Read",
				},
				"writeacccess": {
					Id:          "cccbbb1111",
					Privilege:   "writeacccess",
					AccessLevel: "Write",
				},
			},
			Resources: map[string]Resource{
				"mock": {
					ArnTmpl:  "arn:test::{{.mockField}}",
					Resource: "mock",
				},
			},
		},
	}
}

func Test_GetIamActions(t *testing.T) {
	expectedRespMap := make(map[string]string)
	expectedRespMap["readacccess"] = "aaabbb111"
	expectedRespMap["writeacccess"] = "cccbbb1111"

	iamDef := NewIamDefinitionsMock()
	actionMap := iamDef.GetActionsForService("testservice")

	if !reflect.DeepEqual(expectedRespMap, actionMap) {
		t.Errorf("error got %s expected %s", actionMap, expectedRespMap)
	}
}

func TestIamDefinitions_FindActionsById(t *testing.T) {
	iamDef := NewIamDefinitionsMock()
	expectedList := []string{"ts:readacccess"}

	respList := iamDef.FindActionsById([]string{"aaabbb111"})

	if !reflect.DeepEqual(expectedList, respList) {
		t.Errorf("error got %s expected %s", respList, expectedList)
	}
}
