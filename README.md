# AWX Provisioner

This is still a work in progres! Execute Tower/AWX job templates from terraform with a provisioner

# Current Functionality

- Execute Job template with ID

# TODO

- [ ] Fix Loop for printing job logs

- [ ] Support passing inventories to job template

- [ ] Support passing custom variables (survey data)

- [ ] Add new features

# Usage

You can use [AWX Provider Provider](https://gitlab.com/dhendel/terraform-provider-awx) to setup the job template

```
resource "awx_job_template" "elasticsearch-install" {
  name          = "${upper(terraform.workspace)}-ELK"
  description   = "ELK Installation Template"
  project_id    = awx_project.elk-project.id
  job_type      = "run"
  inventory_id  = awx_inventory.elkstack.id
  playbook      = "site.yml"
  job_tags      = "install"
  credential_id = "3"
  //  extra_credential_ids = [4, 5]
}

resource "null_resource" "tower_job" {
  depends_on = [module.demo1]
  //  triggers = {
  //    ips = "${module.demo1.virtual_machine_default_ips}"
  //  }
  
  provisioner "awx" {
    awx_settings {
      url      = var.tower_endpoint
      username = var.tower_user
      password = var.tower_password
    }

    job_template {
      template_id = awx_job_template.install-app.id
    }
  }
}
```