package main

import (
	"context"
	"fmt"
	"gitlab.com/dhendel/awx-go"
	awxp "gitlab.com/dhendel/terraform-provisioner-awx/awx"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"strconv"
	"time"
)


type provisioner struct {
	awxSettings *awxp.ClientSettings
	awxJobTemplateSettings *awxp.JobTemplateSettings
}

// Provisioner describes this provisioner configuration.
func Provisioner() terraform.ResourceProvisioner {

	return &schema.Provisioner{
		Schema: map[string]*schema.Schema{
			"awx_settings": awxp.NewAWXClientSchema(),
			"job_template": awxp.NewJobTemplateSchema(),
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

	awxClient := awxp.NewAWXClient(p.awxSettings, o)


	if err != nil {
		return err
	}

	launchPayload := &awx.JobLaunchOpts{
		Inventory: p.awxJobTemplateSettings.InventoryID,
	}


	jl, err := awxClient.Client.JobTemplateService.Launch(p.awxJobTemplateSettings.TemplateID, launchPayload, map[string]string{})

	if err != nil {
		return err
	}

	jb := awxClient.Client.JobService
	var jobStatus string
	jobStatus = checkJobStatus(jb, jl.ID)

	for jobStatus == "pending" {
		fmt.Printf("Job %d is pending\n", jl.ID)
		jobStatus = checkJobStatus(jb, jl.ID)
		time.Sleep(5 * time.Second)
	}

	var counter int
	counter = 1

	for jobStatus == "running" {
		printEvents(jb, counter, jl.ID, o)
		jobStatus = checkJobStatus(jb, jl.ID)
		// Increment to get the next event
		counter = counter + 1
	}

	if jobStatus != "successful" {
		return fmt.Errorf("Job %d %s consult ansible tower for more details", jl.ID, jobStatus)
	}

	o.Output(fmt.Sprintf("Job %d %s", jl.ID, jobStatus))

	return nil
}

func decodeConfig(d *schema.ResourceData) (*provisioner, error) {
	p := &provisioner{
		awxSettings:            awxp.NewAWXClientSettingsFromInterface(d.GetOk("awx_settings")),
		awxJobTemplateSettings: awxp.NewJobTemplateSchemaFromInterface(d.GetOk("job_template")),
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

func printEvents(jb *awx.JobService, counter, jobID int, output terraform.UIOutput) {
	events, _, err := jb.GetJobEvents(jobID, map[string]string{
		"counter":      strconv.Itoa(counter),
	})

	if err != nil {
		panic(err)
	}
	if len(events) > 0 {
		event := events[0]

		if event.Stdout != "" {
			output.Output(fmt.Sprintf("%s\n", event.Stdout))
		}

	}
}