package main

import (
	"context"
	"fmt"
	"gitlab.com/dhendel/awx-go"
	"navihealth.us/tf/terraform-provisioner-awx/types"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"strconv"
	"time"
)


type provisioner struct {
	awxSettings *types.AWXClientSettings
	awxJobTemplateSettings *types.JobTemplateSettings
}

// Provisioner describes this provisioner configuration.
func Provisioner() terraform.ResourceProvisioner {
	return &schema.Provisioner{
		Schema: map[string]*schema.Schema{
			"awx_settings":  types.NewAWXClientSchema(),
			"job_template":  types.NewJobTemplateSchema(),
		},
		ApplyFunc:    applyFn,
	}
}


func applyFn(ctx context.Context) error {

	o := ctx.Value(schema.ProvOutputKey).(terraform.UIOutput)
	d := ctx.Value(schema.ProvConfigDataKey).(*schema.ResourceData)

	p, err := decodeConfig(d)
	if err != nil {
		return err
	}

	awx := types.NewAWXClient(p.awxSettings, o)


	if err != nil {
		return err
	}

	jl, err := awx.Client.JobTemplateService.Launch(p.awxJobTemplateSettings.TemplateID, map[string]interface{}{}, map[string]string{})

	if err != nil {
		return err
	}

	jb := awx.Client.JobService
	var jobStatus string
	jobStatus = checkJobStatus(jb, jl.ID)

	var counter int
	counter = 1

	for jobStatus != "failed" ||  jobStatus != "successful" || jobStatus != "error" || jobStatus != "canceled" {
		if jobStatus == "pending" {
			o.Output(fmt.Sprintf("Job %d is pending\n", jl.ID))
		}

		for jobStatus == "running" {
			events, _, err := jb.GetJobEvents(jl.ID, map[string]string{
				"counter":      strconv.Itoa(counter),
			})

			if err != nil {
				panic(err)
			}
			if len(events) > 0 {
				event := events[0]

				if event.Stdout != "" {
					o.Output(fmt.Sprintf("%s\n", event.Stdout))
				}

			}
			jobStatus = checkJobStatus(jb, jl.ID)
			if jobStatus != "running" {
				return nil
			}
			counter = counter + 1
		}
		jobStatus = checkJobStatus(jb, jl.ID)
		time.Sleep(5 * time.Second)
	}


	return nil
}

func decodeConfig(d *schema.ResourceData) (*provisioner, error) {
	p := &provisioner{
		awxSettings: types.NewAWXClientSettingsFromInterface(d.GetOk("awx_settings")),
		awxJobTemplateSettings: types.NewJobTemplateSchemaFromInterface(d.GetOk("job_template")),
	}

	return  p, nil
}

func checkJobStatus(j *awx.JobService, id int) string {
	job, err := j.GetJob(id, map[string]string{})
	if err != nil {
		panic(err)
	}

	return job.Status
}