package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePlacementResourceProviderV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlacementResourceProviderV1Create,
		ReadContext:   resourcePlacementResourceProviderV1Read,
		UpdateContext: resourcePlacementResourceProviderV1Update,
		DeleteContext: resourcePlacementResourceProviderV1Delete,
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

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"parent_provider_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed attributes
			"generation": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"root_provider_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourcePlacementResourceProviderV1Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	placementClient, err := config.PlacementV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack placement client: %s", err)
	}

	createOpts := resourceproviders.CreateOpts{
		Name: d.Get("name").(string),
	}

	if v, ok := d.GetOk("uuid"); ok {
		createOpts.UUID = v.(string)
	}

	if v, ok := d.GetOk("parent_provider_uuid"); ok {
		createOpts.ParentProviderUUID = v.(string)
	}

	log.Printf("[DEBUG] openstack_placement_resourceprovider_v1 create options: %#v", createOpts)

	rp, err := resourceproviders.Create(ctx, placementClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Unable to create openstack_placement_resourceprovider_v1: %s", err)
	}

	d.SetId(rp.UUID)

	return resourcePlacementResourceProviderV1Read(ctx, d, meta)
}

func resourcePlacementResourceProviderV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	placementClient, err := config.PlacementV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack placement client: %s", err)
	}

	rp, err := resourceproviders.Get(ctx, placementClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_placement_resourceprovider_v1"))
	}

	log.Printf("[DEBUG] Retrieved openstack_placement_resourceprovider_v1 %s: %#v", d.Id(), rp)

	d.Set("name", rp.Name)
	d.Set("uuid", rp.UUID)
	d.Set("parent_provider_uuid", rp.ParentProviderUUID)
	d.Set("root_provider_uuid", rp.RootProviderUUID)
	d.Set("generation", rp.Generation)
	d.Set("region", GetRegion(d, config))

	// Set links
	links := make([]map[string]any, len(rp.Links))
	for i, link := range rp.Links {
		links[i] = map[string]any{
			"href": link.Href,
			"rel":  link.Rel,
		}
	}
	d.Set("links", links)

	return nil
}

func resourcePlacementResourceProviderV1Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	placementClient, err := config.PlacementV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack placement client: %s", err)
	}

	var updateOpts resourceproviders.UpdateOpts

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("parent_provider_uuid") {
		parentUUID := d.Get("parent_provider_uuid").(string)
		updateOpts.ParentProviderUUID = &parentUUID
	}

	log.Printf("[DEBUG] openstack_placement_resourceprovider_v1 %s update options: %#v", d.Id(), updateOpts)

	_, err = resourceproviders.Update(ctx, placementClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating openstack_placement_resourceprovider_v1 %s: %s", d.Id(), err)
	}

	return resourcePlacementResourceProviderV1Read(ctx, d, meta)
}

func resourcePlacementResourceProviderV1Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	placementClient, err := config.PlacementV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack placement client: %s", err)
	}

	log.Printf("[DEBUG] Deleting openstack_placement_resourceprovider_v1 %s", d.Id())

	err = resourceproviders.Delete(ctx, placementClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_placement_resourceprovider_v1"))
	}

	return nil
}