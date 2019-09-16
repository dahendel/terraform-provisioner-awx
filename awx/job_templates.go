package awx

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.com/dhendel/awx-go"
	"os"
)

const (
	awxJobTemplateAttributeJobTemplateID = "template_id"
	awxJobTemplateAttributeExtraArgs     = "extra_args"
	awxJobTemplateAttributeInventoryID   = "inventory_id"
	awxJobTemplateEnvID                  = "TF_AWX_JOB_TEMPLATE_ID"
)

type JobTemplateSettings struct {
	TemplateID  int
	InventoryID int
	ExtraArgs   map[string]interface{}
}

type JobTemplate struct {
	ID                    int           `json:"id"`
	Type                  string        `json:"type"`
	URL                   string        `json:"url"`
	Related               *awx.Related   `json:"related"`
	Inventory             int           `json:"inventory"`
	SurveyEnabled         bool          `json:"survey_enabled"`
}

type Related struct {
	Launch                       string `json:"launch"`
	SurveySpec                   string `json:"survey_spec"`
}

type Inventory struct {
	ID                           int    `json:"id"`
	Name                         string `json:"name"`
	Description                  string `json:"description"`
	HasActiveFailures            bool   `json:"has_active_failures"`
	TotalHosts                   int    `json:"total_hosts"`
	HostsWithActiveFailures      int    `json:"hosts_with_active_failures"`
	TotalGroups                  int    `json:"total_groups"`
	GroupsWithActiveFailures     int    `json:"groups_with_active_failures"`
	HasInventorySources          bool   `json:"has_inventory_sources"`
	TotalInventorySources        int    `json:"total_inventory_sources"`
	InventorySourcesWithFailures int    `json:"inventory_sources_with_failures"`
	OrganizationID               int    `json:"organization_id"`
	Kind                         string `json:"kind"`
}


func NewJobTemplateSchema()  *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				awxJobTemplateAttributeJobTemplateID: &schema.Schema{
					Type: schema.TypeInt,
					Required: true,
					DefaultFunc: func() (i interface{}, e error) {
						if val := os.Getenv(awxJobTemplateEnvID); val != "" {
							return val, nil
						}
						return nil, fmt.Errorf("%s or env variable %s must be set for awx provisioner", awxJobTemplateAttributeJobTemplateID, awxJobTemplateEnvID)
					},
				},
				awxJobTemplateAttributeExtraArgs: &schema.Schema{
					Type: schema.TypeMap,
					Optional: true,
				},
				awxJobTemplateAttributeInventoryID: &schema.Schema{
					Type: schema.TypeInt,
					Optional: true,
				},
			},
		},
	}
}

func NewJobTemplateSchemaFromInterface(i interface{}, ok bool) *JobTemplateSettings {
	v := &JobTemplateSettings{}
	if ok {
		vals := mapFromTypeSetList(i.(*schema.Set).List())
		v.TemplateID = vals[awxJobTemplateAttributeJobTemplateID].(int)
		v.InventoryID = vals[awxJobTemplateAttributeInventoryID].(int)
		v.ExtraArgs = vals[awxJobTemplateAttributeExtraArgs].(map[string]interface{})
	}
	return v
}