package policy

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/perkbox/cloud-access-bot/internal"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/perkbox/cloud-access-bot/internal/settings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
)

type IamPolicyMan struct {
	arnTemplates      map[string]string
	arnTmplFieldNames map[string]string
	*IAMProvider
}

func NewPolicyManager(cfg aws.Config, config settings.Settings, arnTmpl, arnTmplFieldName map[string]string) *IamPolicyMan {
	iamProv := NewIAMClient(cfg)
	iamProv.STSProvider = sts.NewFromConfig(cfg)
	iamProv.Settings = config

	return &IamPolicyMan{
		arnTmpl,
		arnTmplFieldName,
		iamProv,
	}
}

// IsPolicyExpired takes in a iam policy as a string and will marshall it into IamPolicy and check that the Statements aren't
// expired.
// Returns  true(Expired) , false(Not expired or no time params)
func (i IamPolicyMan) IsPolicyExpired(policy string) (bool, error) {
	var iampolicy IamPolicy
	policyBytes := []byte(policy)

	err := json.Unmarshal(policyBytes, &iampolicy)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling iam policy: %s", err)
	}

	for _, statement := range iampolicy.Statement {

		if statement.Condition == nil ||
			statement.Condition.DateLessThan == nil ||
			statement.Condition.DateGreaterThan == nil {
			return false, nil
		}

		if statement.Condition.DateLessThan.AwsCurrentTime == "" || statement.Condition.DateGreaterThan.AwsCurrentTime == "" {
			return false, nil
		}
		ltTime, _ := time.Parse("2006-01-02T15:04:05Z", statement.Condition.DateLessThan.AwsCurrentTime)
		gtTime, _ := time.Parse("2006-01-02T15:04:05Z", statement.Condition.DateGreaterThan.AwsCurrentTime)
		if !inTimeSpan(gtTime, ltTime, time.Now()) {
			return true, nil
		}
	}

	return false, nil
}

func (i *IamPolicyMan) GeneratePolicyFromAuditObj(curTime time.Time, object internal.AuditObject, tmpls, tmplFieldNames map[string]string) ([]byte, error) {
	var iamStatements []IamStatement

	ioverwrite := i
	if len(tmpls) != 0 && len(tmplFieldNames) != 0 {
		ioverwrite = &IamPolicyMan{
			arnTemplates:      tmpls,
			arnTmplFieldNames: tmplFieldNames,
			IAMProvider:       i.IAMProvider,
		}
	}

	arns, err := ioverwrite.generateArns(object.AccountId, object.Services, object.Resources)
	if err != nil {
		return nil, err
	}

	for _, service := range object.Services {
		if _, ok := object.Actions[service]; !ok {
			return nil, fmt.Errorf("error generating IAM statement. Unable to find actions for service: %s", service)
		}
		if _, ok := arns[service]; !ok {
			return nil, fmt.Errorf("error generating IAM statement. Unable to find actions for service: %s", service)
		}

		timeDuration, _ := strconv.Atoi(object.Duration)

		iamStatement := IamStatement{
			Effect:   "Allow",
			Action:   object.Actions[service],
			Resource: arns[service],
			Condition: &IamCondition{
				StringEqualsIgnoreCase: &StringEqualsIgnoreCase{AwsUserid: object.CloudUserId},
				DateGreaterThan:        &DateGreaterThan{AwsCurrentTime: curTime.UTC().Format("2006-01-02T15:04:05Z")},
				DateLessThan:           &DateLessThan{AwsCurrentTime: curTime.UTC().Add(time.Duration(timeDuration) * time.Minute).Format("2006-01-02T15:04:05Z")},
			},
		}

		iamStatements = append(iamStatements, iamStatement)
	}

	return json.Marshal(IamPolicy{
		Version:   "2012-10-17",
		Statement: iamStatements,
	})
}

func (i *IamPolicyMan) generateArns(AccountId string, Services []string, Resources map[string][]string) (map[string][]string, error) {
	var ARNS = make(map[string][]string)
	tmplVals := make(map[string]interface{})
	tmplVals["Account"] = AccountId
	tmplVals["Partition"] = "aws"
	tmplVals["Region"] = "*"

	for _, service := range Services {
		ARNS[service] = []string{}
		if serviceResources, ok := Resources[service]; ok {
			for _, resource := range serviceResources {
				//For manually entered resources wildcard them and add them to the resources on the policy
				if arn.IsARN(resource) {
					serviceArns := ARNS[service]
					serviceArns = append(serviceArns, resource, wildcardARN(resource))
					ARNS[service] = serviceArns
					continue
				}

				//For resources which SHOULD have a template run through the below logic to template them and add them to the
				// map[string][]string
				if _, ok := i.arnTemplates[service]; !ok {
					return nil, fmt.Errorf("error generating arns, no template found for service %s", service)
				}
				if _, ok := i.arnTmplFieldNames[service]; !ok {
					return nil, fmt.Errorf("error generating arns, no template feild name found for service %s", service)
				}

				tmplVals[i.arnTmplFieldNames[service]] = resource
				iamArn, _ := renderTmpl(i.arnTemplates[service], tmplVals)

				serviceArns := ARNS[service]
				serviceArns = append(serviceArns, iamArn, wildcardARN(iamArn))
				ARNS[service] = serviceArns
			}
		}
	}

	return ARNS, nil
}

func wildcardARN(arn string) string {
	if arn != "" {
		if !strings.Contains(arn, "/*") {
			return fmt.Sprintf("%s/*", arn)
		}
	}
	return ""
}
