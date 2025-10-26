package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePlacementResourceProviderInventoryV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePlacementResourceProviderInventoryV1Read,

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

			"inventories": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allocation_ratio": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"max_unit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_unit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"reserved": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"step_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePlacementResourceProviderInventoryV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	placementClient, err := config.PlacementV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack placement client: %s", err)
	}

	rpID := d.Get("resource_provider_id").(string)

	inventories, err := resourceproviders.GetInventories(ctx, placementClient, rpID).Extract()
	if err != nil {
		return diag.Errorf("Error retrieving inventories for resource provider %s: %s", rpID, err)
	}

	log.Printf("[DEBUG] Retrieved inventories for resource provider %s: %#v", rpID, inventories)

	d.SetId(rpID)
	d.Set("resource_provider_generation", inventories.ResourceProviderGeneration)
	d.Set("region", GetRegion(d, config))

	// Convert inventories map to terraform format
	invMap := make(map[string]any)
	for resourceClass, inventory := range inventories.Inventories {
		invMap[resourceClass] = map[string]any{
			"allocation_ratio": inventory.AllocationRatio,
			"max_unit":         inventory.MaxUnit,
			"min_unit":         inventory.MinUnit,
			"reserved":         inventory.Reserved,
			"step_size":        inventory.StepSize,
			"total":            inventory.Total,
		}
	}
	d.Set("inventories", invMap)

	return nil
}