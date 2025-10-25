package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePlacementResourceProviderV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePlacementResourceProviderV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"member_of": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"resources": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"in_tree": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"required": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed attributes
			"parent_provider_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"root_provider_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"generation": {
				Type:     schema.TypeInt,
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

func dataSourcePlacementResourceProviderV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	placementClient, err := config.PlacementV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack placement client: %s", err)
	}

	// If UUID is specified, get directly
	if v, ok := d.GetOk("uuid"); ok {
		uuid := v.(string)
		rp, err := resourceproviders.Get(ctx, placementClient, uuid).Extract()
		if err != nil {
			return diag.Errorf("Error retrieving resource provider %s: %s", uuid, err)
		}
		return dataSourcePlacementResourceProviderV1Attributes(d, config, rp)
	}

	// Otherwise, use list with filters
	listOpts := resourceproviders.ListOpts{
		Name:     d.Get("name").(string),
		MemberOf: d.Get("member_of").(string),
		Resources: d.Get("resources").(string),
		InTree:   d.Get("in_tree").(string),
		Required: d.Get("required").(string),
	}

	log.Printf("[DEBUG] openstack_placement_resourceprovider_v1 list options: %#v", listOpts)

	allPages, err := resourceproviders.List(placementClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to list resource providers: %s", err)
	}

	allRPs, err := resourceproviders.ExtractResourceProviders(allPages)
	if err != nil {
		return diag.Errorf("Unable to extract resource providers: %s", err)
	}

	if len(allRPs) < 1 {
		return diag.Errorf("No resource provider found matching criteria")
	}

	if len(allRPs) > 1 {
		return diag.Errorf("More than one resource provider found (%d)", len(allRPs))
	}

	return dataSourcePlacementResourceProviderV1Attributes(d, config, &allRPs[0])
}

func dataSourcePlacementResourceProviderV1Attributes(d *schema.ResourceData, config *Config, rp *resourceproviders.ResourceProvider) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieved openstack_placement_resourceprovider_v1: %#v", rp)

	d.SetId(rp.UUID)
	d.Set("uuid", rp.UUID)
	d.Set("name", rp.Name)
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