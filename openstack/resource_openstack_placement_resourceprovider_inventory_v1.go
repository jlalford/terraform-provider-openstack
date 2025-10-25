package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePlacementResourceProviderInventoryV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlacementResourceProviderInventoryV1Create,
		ReadContext:   resourcePlacementResourceProviderInventoryV1Read,
		UpdateContext: resourcePlacementResourceProviderInventoryV1Update,
		DeleteContext: resourcePlacementResourceProviderInventoryV1Delete,
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

			"resource_provider_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"inventories": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allocation_ratio": {
							Type:     schema.TypeFloat,
							Required: true,
						},
						"max_unit": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"min_unit": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"reserved": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"step_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"total": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			// Computed attributes
			"resource_provider_generation": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourcePlacementResourceProviderInventoryV1Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return resourcePlacementResourceProviderInventoryV1Update(ctx, d, meta)
}

func resourcePlacementResourceProviderInventoryV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	placementClient, err := config.PlacementV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack placement client: %s", err)
	}

	rpID := d.Get("resource_provider_id").(string)
	if rpID == "" && d.Id() != "" {
		// For import case, the ID is the resource provider ID
		rpID = d.Id()
		d.Set("resource_provider_id", rpID)
	}

	inventories, err := resourceproviders.GetInventories(ctx, placementClient, rpID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving inventories for resource provider"))
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

func resourcePlacementResourceProviderInventoryV1Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	placementClient, err := config.PlacementV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack placement client: %s", err)
	}

	rpID := d.Get("resource_provider_id").(string)

	// Get current generation from the resource provider
	rp, err := resourceproviders.Get(ctx, placementClient, rpID).Extract()
	if err != nil {
		return diag.Errorf("Error getting resource provider %s: %s", rpID, err)
	}

	// Convert terraform inventories to gophercloud format
	inventoriesRaw := d.Get("inventories").(map[string]any)
	inventories := make(map[string]resourceproviders.Inventory)

	for resourceClass, invRaw := range inventoriesRaw {
		invMap := invRaw.(map[string]any)
		inventory := resourceproviders.Inventory{
			AllocationRatio: float32(invMap["allocation_ratio"].(float64)),
			MaxUnit:         invMap["max_unit"].(int),
			MinUnit:         invMap["min_unit"].(int),
			Reserved:        invMap["reserved"].(int),
			StepSize:        invMap["step_size"].(int),
			Total:           invMap["total"].(int),
		}
		inventories[resourceClass] = inventory
	}

	updateOpts := resourceproviders.UpdateInventoriesOpts{
		ResourceProviderGeneration: rp.Generation,
		Inventories:                inventories,
	}

	log.Printf("[DEBUG] Updating inventories for resource provider %s: %#v", rpID, updateOpts)

	_, err = resourceproviders.UpdateInventories(ctx, placementClient, rpID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating inventories for resource provider %s: %s", rpID, err)
	}

	d.SetId(rpID)

	return resourcePlacementResourceProviderInventoryV1Read(ctx, d, meta)
}

func resourcePlacementResourceProviderInventoryV1Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	placementClient, err := config.PlacementV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack placement client: %s", err)
	}

	rpID := d.Get("resource_provider_id").(string)

	log.Printf("[DEBUG] Deleting inventories for resource provider %s", rpID)

	err = resourceproviders.DeleteInventories(ctx, placementClient, rpID).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting inventories for resource provider"))
	}

	return nil
}
