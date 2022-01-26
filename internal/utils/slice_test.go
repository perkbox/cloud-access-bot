package utils

import (
	"reflect"
	"testing"
)

func Test_Contains(t *testing.T) {
	tests := []struct {
		Name     string
		Slice    []string
		Contains string
		Resp     bool
	}{
		{
			Name:     "Slice Does Contain String",
			Slice:    []string{"abc", "123", "amihere"},
			Contains: "amihere",
			Resp:     true,
		},
		{
			Name:     "Slice Doesnt Contain String",
			Slice:    []string{"abc", "123"},
			Contains: "amihere",
			Resp:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {

			funcresp := Contains(tc.Slice, tc.Contains)
			if funcresp != tc.Resp {
				t.Errorf("Got %+v, Expected %+v", funcresp, tc.Resp)
			}

		})
	}
}

func Test_RemoveDuplicateStr(t *testing.T) {
	tests := []struct {
		Name     string
		Slice    []string
		Expected []string
	}{
		{
			Name:     "Slice Doesnt Contain Duplicates",
			Slice:    []string{"abc", "123", "amihere"},
			Expected: []string{"abc", "123", "amihere"},
		},
		{
			Name:     "Slice Doesnt Contain String",
			Slice:    []string{"abc", "123", "123"},
			Expected: []string{"abc", "123"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {

			funcResp := RemoveDuplicateStr(tc.Slice)

			if !reflect.DeepEqual(tc.Expected, funcResp) {
				t.Errorf("Got %+v   Expected %+v", funcResp, tc.Expected)
			}

		})
	}
}

func Test_SplitFreeString(t *testing.T) {
	tests := []struct {
		Name        string
		StringInput string
		Expected    []string
	}{
		{
			Name:        "Slice Can be Split",
			StringInput: "Hello,New\nWorld",
			Expected:    []string{"Hello", "New", "World"},
		},
		{
			Name:        "Single Word Split",
			StringInput: "Hello",
			Expected:    []string{"Hello"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {

			funcResp := SplitFreeString(tc.StringInput)

			if !reflect.DeepEqual(tc.Expected, funcResp) {
				t.Errorf("Got %+v   Expected %+v", funcResp, tc.Expected)
			}

		})
	}
}
