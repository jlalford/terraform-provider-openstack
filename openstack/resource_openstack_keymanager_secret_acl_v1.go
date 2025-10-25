package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/acls"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeyManagerSecretACLV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeyManagerSecretACLV1Create,
		ReadContext:   resourceKeyManagerSecretACLV1Read,
		UpdateContext: resourceKeyManagerSecretACLV1Update,
		DeleteContext: resourceKeyManagerSecretACLV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"secret_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"read": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_access": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"users": {
							Type:     schema.TypeSet,
							Optional: true,
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

			// Computed
			"acl_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKeyManagerSecretACLV1Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	kmClient, err := config.KeyManagerV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack KeyManager client: %s", err)
	}

	secretID := d.Get("secret_id").(string)

	opts := expandKeyManagerSecretACLV1Read(d)

	log.Printf("[DEBUG] openstack_keymanager_secret_acl_v1 create options: %#v", opts)

	aclRef, err := acls.SetSecretACL(ctx, kmClient, secretID, opts).Extract()
	if err != nil {
		return diag.Errorf("Error creating secret ACL: %s", err)
	}

	d.SetId(secretID)
	d.Set("acl_ref", string(*aclRef))

	return resourceKeyManagerSecretACLV1Read(ctx, d, meta)
}

func resourceKeyManagerSecretACLV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	kmClient, err := config.KeyManagerV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack KeyManager client: %s", err)
	}

	secretID := d.Get("secret_id").(string)
	if secretID == "" && d.Id() != "" {
		// For import case, the ID is the secret ID
		secretID = d.Id()
		d.Set("secret_id", secretID)
	}

	acl, err := acls.GetSecretACL(ctx, kmClient, secretID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving secret ACL"))
	}

	log.Printf("[DEBUG] Retrieved secret ACL: %#v", acl)

	// Flatten ACL read details
	if readACL, ok := (*acl)["read"]; ok {
		read := []map[string]any{{
			"project_access": readACL.ProjectAccess,
			"users":          readACL.Users,
			"created":        readACL.Created.Format("2006-01-02T15:04:05"),
			"updated":        readACL.Updated.Format("2006-01-02T15:04:05"),
		}}
		d.Set("read", read)
	}

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceKeyManagerSecretACLV1Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	kmClient, err := config.KeyManagerV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack KeyManager client: %s", err)
	}

	secretID := d.Get("secret_id").(string)

	if d.HasChange("read") {
		opts := expandKeyManagerSecretACLV1Read(d)

		log.Printf("[DEBUG] openstack_keymanager_secret_acl_v1 update options: %#v", opts)

		aclRef, err := acls.UpdateSecretACL(ctx, kmClient, secretID, opts).Extract()
		if err != nil {
			return diag.Errorf("Error updating secret ACL: %s", err)
		}

		d.Set("acl_ref", string(*aclRef))
	}

	return resourceKeyManagerSecretACLV1Read(ctx, d, meta)
}

func resourceKeyManagerSecretACLV1Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	kmClient, err := config.KeyManagerV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack KeyManager client: %s", err)
	}

	secretID := d.Get("secret_id").(string)

	log.Printf("[DEBUG] Deleting secret ACL for secret %s", secretID)

	err = acls.DeleteSecretACL(ctx, kmClient, secretID).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting secret ACL"))
	}

	return nil
}

func expandKeyManagerSecretACLV1Read(d *schema.ResourceData) acls.SetOpts {
	opts := acls.SetOpts{}

	if v, ok := d.GetOk("read"); ok {
		readList := v.([]any)
		if len(readList) > 0 {
			readMap := readList[0].(map[string]any)

			projectAccess := readMap["project_access"].(bool)

			var users []string
			if usersSet, ok := readMap["users"].(*schema.Set); ok {
				users = make([]string, 0, usersSet.Len())
				for _, user := range usersSet.List() {
					users = append(users, user.(string))
				}
			}

			opts = append(opts, acls.SetOpt{
				Type:          "read",
				ProjectAccess: &projectAccess,
				Users:         &users,
			})
		}
	}

	return opts
}