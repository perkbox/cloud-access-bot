package identitydata

type ArnData struct {
	ResourceType  string
	TmplFieldName string
}

type IamDefinitions map[string]IamServices

type IamServices struct {
	ServiceName             string               `json:"service_name"`
	Prefix                  string               `json:"prefix"`
	ServiceAuthorizationUrl string               `json:"service_authorization_url"`
	Privileges              map[string]Privilege `json:"privileges"`
	Resources               map[string]Resource  `json:"resources"`
}

type Privilege struct {
	Id          string
	Privilege   string `json:"privilege"`
	AccessLevel string `json:"access_level"`
}

type Resource struct {
	Resource string `json:"resources"`
	ArnTmpl  string `json:"arn"`
}
