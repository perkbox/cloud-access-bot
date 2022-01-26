package awsproviderv2

import "github.com/aws/aws-sdk-go-v2/aws/arn"

type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateResourcesFormat(resources []string) []string {
	var failedResources []string
	for _, r := range resources {
		if !arn.IsARN(r) && r != "" {
			failedResources = append(failedResources, r)
		}
	}

	return failedResources
}
