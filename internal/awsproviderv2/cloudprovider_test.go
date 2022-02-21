package awsproviderv2

import (
	"reflect"
	"sort"
	"testing"

	"github.com/perkbox/cloud-access-bot/internal/settings"
)

func NewMockCloudProvider() *ResourceFinder {
	settings, _ := settings.NewConfigMock()
	return &ResourceFinder{
		S3Provider:       &S3Provider{Client: S3Mock{}, Regions: settings.Regions},
		DynamodbProvider: &DynamodbProvider{Client: DynMock{}, Regions: settings.Regions},
		Settings:         settings,
	}
}

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
			vali := NewMockCloudProvider()

			failedResourrces := vali.ValidateResourcesFormat(tc.Resources)

			if !reflect.DeepEqual(failedResourrces, tc.ExpectedFailedResources) {
				t.Errorf("Error Got %s, Expected %s", failedResourrces, tc.ExpectedFailedResources)
			}
		})
	}
}

func Test_ResourceFinder(t *testing.T) {
	rFinder := NewMockCloudProvider()

	tests := []struct {
		Name      string
		Resource  string
		HasFinder bool
		Expected  []string
	}{
		{
			Name:      "Get S3 Resources",
			Resource:  "s3",
			HasFinder: true,
			Expected:  []string{"BucketA", "BucketB"},
		},
		{
			Name:      "Get Dynamo Resources",
			Resource:  "dynamodb",
			HasFinder: true,
			Expected:  []string{"TestDynTable1", "TestDynTable2"},
		},
		{
			Name:      "Get sts Resources(Should be empty)",
			Resource:  "sts",
			HasFinder: false,
			Expected:  []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {

			gotResp, gotHasFinder := rFinder.ResourceFinder(tc.Resource, "")
			if !array_sorted_equal(gotResp, tc.Expected) && gotHasFinder == tc.HasFinder {
				t.Errorf("GOT Resources %s, Expected Resources %s\n"+
					" Got Finder %v , Expected finder %v", gotResp, tc.Expected, gotHasFinder, tc.HasFinder)

			}
		})
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
