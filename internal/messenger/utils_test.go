package messenger

import (
	"reflect"
	"testing"

	"github.com/slack-go/slack"
)

func Test_SliceToOptions(t *testing.T) {
	InputSlice := []string{"getBucket", "readSomething"}

	ExpectedOpts := Options{Options: []Option{
		{
			Text: Text{
				Type: "Type",
				Text: "getBucket",
			},
			Value: "getBucket",
		},
		{
			Text: Text{
				Type: "Type",
				Text: "readSomething",
			},
			Value: "readSomething",
		},
	},
	}

	respOpts := SliceToOptions(InputSlice, "Type")
	if !reflect.DeepEqual(ExpectedOpts, respOpts) {
		t.Errorf("Expected  %+v\n, Got: %+v", ExpectedOpts, respOpts)
	}
}

func Test_GetValuesFromSelectedOptions(t *testing.T) {
	testOptBlocks := []slack.OptionBlockObject{
		{Value: "Vala"},
		{Value: "Valb"},
	}
	ExpectedSlic := []string{"Vala", "Valb"}

	respVals := GetValuesFromSelectedOptions(testOptBlocks)

	if !reflect.DeepEqual(respVals, ExpectedSlic) {
		t.Errorf("Expected  %+v\n, Got: %+v", ExpectedSlic, respVals)
	}
}

func Test_MapToOptions(t *testing.T) {
	vals := map[string]string{
		"GetBucket": "1",
	}

	ExpectedOpts := Options{Options: []Option{
		{
			Text: Text{
				Type: "Type",
				Text: "GetBucket",
			},
			Value: "1",
		},
	},
	}

	respOpts := MapToOptions(vals, "Type")

	if !reflect.DeepEqual(ExpectedOpts, respOpts) {
		t.Errorf("Expected  %+v\n, Got: %+v", ExpectedOpts, respOpts)
	}
}
