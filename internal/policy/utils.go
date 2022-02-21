package policy

import (
	"bytes"
	"context"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

//inTimeSpan checks the current time is between the `start` and `end` times based on the current time(`check`)
func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func renderTmpl(tmplString string, tmplValues map[string]interface{}) (string, error) {
	tmpl, err := template.New("arn").Parse(tmplString)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, tmplValues)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

//Allows cross account role assumptions
func assumeRole(accountRoleArn string, stsprovider sts.Client) (aws.Config, error) {
	cnf, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(aws.NewCredentialsCache(
			stscreds.NewAssumeRoleProvider(
				&stsprovider,
				accountRoleArn,
			)),
		),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return cnf, nil
}
