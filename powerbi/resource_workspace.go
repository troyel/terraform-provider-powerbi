package powerbi

import (
	"github.com/codecutout/terraform-provider-powerbi/powerbi/internal/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ResourceWorkspace represents a Power BI workspace
func ResourceWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: createWorkspace,
		Read:   readWorkspace,
		Update: updateWorkspace,
		Delete: deleteWorkspace,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the workspace.",
			},
		},
	}
}

func createWorkspace(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	resp, err := client.CreateGroup(api.CreateGroupRequest{
		Name: d.Get("name").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return readWorkspace(d, meta)
}

func readWorkspace(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	workspace, err := client.GetGroup(d.Id())
	if err != nil {
		return err
	}

	if workspace == nil {
		d.SetId("")
	} else {
		d.SetId(workspace.ID)
		d.Set("name", workspace.Name)
	}

	return nil
}

func updateWorkspace(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	err := client.UpdateGroupAsAdmin(api.UpdateGroupAsAdminRequest{
		GroupID: d.Id(),
		Name:    d.Get("name").(string),
	})
	if err != nil {
		return err
	}

	return readWorkspace(d, meta)
}

func deleteWorkspace(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	return client.DeleteGroup(d.Id())
}
