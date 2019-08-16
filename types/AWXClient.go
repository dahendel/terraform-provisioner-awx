package types

import (
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"net/http"
	"os"
	"gitlab.com/dhendel/awx-go"
	"strconv"
)

const (
	awxDefaultBaseURLString = "127.0.0.1"

	awxClientAttributeBaseURLString = "url"
	awxClientAttributeUsername = "username"
	awxClientAttributePassword = "password"
	awxClientAttributeSkipTLSVerify = "skip_tls_verify"

	awxClientEnvBaseURLString = "TF_AWX_URL"
	awxClientEnvUsername = "TF_AWX_USERNAME"
	awxClientEnvPassword = "TF_AWX_PASSWORD"
	awxClientEnvSkipTLSVerify = "TF_AWX_SKIP_TLS_VERIFY"
)

// AWXClientSettings represents settings to establish a awx api client
type AWXClientSettings struct {
	Username string
	Password string
	BaseURLString string
	SkipTLSVerify bool
}

type AWXClient struct {
	Client *awx.AWX
	o terraform.UIOutput
}

// NewAWXClientSchema returns new AWXClientSettings schema
func NewAWXClientSchema() *schema.Schema{
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				awxClientAttributeBaseURLString: &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
					DefaultFunc: func() (interface{}, error) {
						val := os.Getenv(awxClientEnvBaseURLString)
						if val != "" {
							return val, nil
						} else {
							return nil, fmt.Errorf("urlmust be set for provisioner")
						}
					},
				},
				awxClientAttributeSkipTLSVerify: &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					DefaultFunc: func() (interface{}, error) {
						val := os.Getenv(awxClientEnvSkipTLSVerify)
						if val != "" {
							skip, err := strconv.ParseBool(val)
							if err != nil {
								return false, err
							}
							return skip, nil
						} else {
							return false, nil
						}
					},
				},
				awxClientAttributeUsername: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					DefaultFunc: func() (interface{}, error) {
						val := os.Getenv(awxClientEnvUsername)
						if val != "" {
							return val, nil
						}
						return nil, fmt.Errorf("username must be set for provisioner")
					},
				},
				awxClientAttributePassword: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					DefaultFunc: func() (interface{}, error) {
						val := os.Getenv(awxClientEnvPassword)
						if val != "" {
							return val, nil
						}
						return nil, fmt.Errorf("password must be set for provisioner")
					},
				},
			},
		},
	}
}

func NewAWXClientSettingsFromInterface(i interface{}, ok bool) *AWXClientSettings {
	v := &AWXClientSettings{
		BaseURLString: awxDefaultBaseURLString,
	}
	if ok {
		vals := mapFromTypeSetList(i.(*schema.Set).List())
		v.BaseURLString = vals[awxClientAttributeBaseURLString].(string)
		v.Username = vals[awxClientAttributeUsername].(string)
		v.Password = vals[awxClientAttributePassword].(string)
		v.SkipTLSVerify = vals[awxClientAttributeSkipTLSVerify].(bool)
	}

	return v
}

func NewAWXClient(settings *AWXClientSettings, o terraform.UIOutput) *AWXClient {
	t := &http.Transport{TLSClientConfig: &tls.Config{
		InsecureSkipVerify: settings.SkipTLSVerify,
	}}

	a := awx.NewAWX(settings.BaseURLString, settings.Username, settings.Password, &http.Client{Transport: t})
	return &AWXClient{
		Client: a,
		o: o,
	}
}

