package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/acls"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKeyManagerSecretACLV1() *schema.Resource {
	ret := &schema.Resource{
		ReadContext: dataSourceKeyManagerSecretACLV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"secret_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"read": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"users": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"updated": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"acl_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return ret
}

func dataSourceKeyManagerSecretACLV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	kmClient, err := config.KeyManagerV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack KeyManager client: %s", err)
	}

	secretID := d.Get("secret_id").(string)

	acl, err := acls.GetSecretACL(ctx, kmClient, secretID).Extract()
	if err != nil {
		return diag.Errorf("Error retrieving ACL for secret %s: %s", secretID, err)
	}

	log.Printf("[DEBUG] Retrieved ACL for secret %s: %#v", secretID, acl)

	d.SetId(secretID)

	// Set the read ACL details
	if readACL, ok := (*acl)["read"]; ok {
		readMap := map[string]any{
			"project_access": readACL.ProjectAccess,
			"users":          readACL.Users,
			"created":        readACL.Created.Format("2006-01-02T15:04:05"),
			"updated":        readACL.Updated.Format("2006-01-02T15:04:05"),
		}
		d.Set("read", []any{readMap})
	}

	d.Set("region", GetRegion(d, config))

	return nil
}
