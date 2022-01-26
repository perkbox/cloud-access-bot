package policy

import (
	"encoding/json"
	"fmt"
)

type IamPolicy struct {
	Version   string         `json:"Version"`
	Statement []IamStatement `json:"Statement"`
}

type ListOrString []string

type IamStatement struct {
	Sid       string        `json:"Sid"`
	Effect    string        `json:"Effect"`
	Action    ListOrString  `json:"Action"`
	Resource  ListOrString  `json:"Resource"`
	Condition *IamCondition `json:"Condition"`
}

type IamCondition struct {
	StringLike             *StringLike             `json:"StringLike,omitempty"`
	StringEqualsIgnoreCase *StringEqualsIgnoreCase `json:"StringEqualsIgnoreCase,omitempty"`
	DateGreaterThan        *DateGreaterThan        `json:"DateGreaterThan,omitempty"`
	DateLessThan           *DateLessThan           `json:"DateLessThan,omitempty"`
}

type StringEqualsIgnoreCase struct {
	AwsUserid string `json:"aws:userid,omitempty"`
}

type StringLike struct {
	AwsUserid string `json:"aws:userid,omitempty"`
}

type DateGreaterThan struct {
	AwsCurrentTime string `json:"aws:CurrentTime,omitempty"`
}

type DateLessThan struct {
	AwsCurrentTime string `json:"aws:CurrentTime,omitempty"`
}

func (a *ListOrString) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err == nil {
		if val, ok := s.(string); ok {
			*a = []string{val}
		} else if _, ok1 := s.([]interface{}); ok1 {
			z := make([]string, len(s.([]interface{})))
			for i, v := range s.([]interface{}) {
				z[i] = fmt.Sprint(v)
			}
			*a = z
		}

		return nil
	}
	return json.Unmarshal(b, a)
}
