package awsproviderv2

import (
	"reflect"
	"testing"
)

func Test_ValidateResourcesFormat(t *testing.T) {
	tests := []struct {
		Name                    string
		Resources               []string
		ExpectedFailedResources []string
	}{
		{
			Name:                    "Catch Short ARN/Not complete",
			Resources:               []string{"arn:::"},
			ExpectedFailedResources: []string{"arn:::"},
		},
		{
			Name:                    "Catch WildCard",
			Resources:               []string{"*"},
			ExpectedFailedResources: []string{"*"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			vali := NewValidator()

			failedResourrces := vali.ValidateResourcesFormat(tc.Resources)

			if !reflect.DeepEqual(failedResourrces, tc.ExpectedFailedResources) {
				t.Errorf("Error Got %s, Expected %s", failedResourrces, tc.ExpectedFailedResources)
			}
		})
	}
}
