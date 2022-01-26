package utils

import (
	"reflect"
	"testing"
)

func Test_HashString(t *testing.T) {
	tests := []struct {
		Name         string
		StringToHash string
		ExpectedHash string
		Length       int
	}{
		{
			Name:         "Working Hash",
			StringToHash: "HashMe",
			ExpectedHash: "acdp",
			Length:       4,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			hash := HashString(tc.StringToHash, tc.Length)

			if !reflect.DeepEqual(tc.ExpectedHash, hash) {
				t.Errorf("error got %s expected %s", hash, tc.ExpectedHash)
			}

		})
	}
}
