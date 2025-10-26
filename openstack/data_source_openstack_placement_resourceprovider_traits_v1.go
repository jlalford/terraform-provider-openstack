package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePlacementResourceProviderTraitsV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePlacementResourceProviderTraitsV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"resource_provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"traits": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"resource_provider_generation": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourcePlacementResourceProviderTraitsV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	placementClient, err := config.PlacementV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack placement client: %s", err)
	}

	rpID := d.Get("resource_provider_id").(string)

	traits, err := resourceproviders.GetTraits(ctx, placementClient, rpID).Extract()
	if err != nil {
		return diag.Errorf("Error retrieving traits for resource provider %s: %s", rpID, err)
	}

	log.Printf("[DEBUG] Retrieved traits for resource provider %s: %#v", rpID, traits)

	d.SetId(rpID)
	d.Set("traits", traits.Traits)
	d.Set("resource_provider_generation", traits.ResourceProviderGeneration)
	d.Set("region", GetRegion(d, config))

	return nil
}
