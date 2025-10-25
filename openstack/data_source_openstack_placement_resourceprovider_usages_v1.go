package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePlacementResourceProviderUsagesV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePlacementResourceProviderUsagesV1Read,

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

			"resource_provider_generation": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"usages": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourcePlacementResourceProviderUsagesV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	placementClient, err := config.PlacementV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack placement client: %s", err)
	}

	rpID := d.Get("resource_provider_id").(string)

	usages, err := resourceproviders.GetUsages(ctx, placementClient, rpID).Extract()
	if err != nil {
		return diag.Errorf("Error retrieving usages for resource provider %s: %s", rpID, err)
	}

	log.Printf("[DEBUG] Retrieved usages for resource provider %s: %#v", rpID, usages)

	d.SetId(rpID)
	d.Set("resource_provider_generation", usages.ResourceProviderGeneration)
	d.Set("region", GetRegion(d, config))

	// Convert usages map to terraform format
	usagesMap := make(map[string]int)
	for resourceClass, usage := range usages.Usages {
		usagesMap[resourceClass] = usage
	}
	d.Set("usages", usagesMap)

	return nil
}
