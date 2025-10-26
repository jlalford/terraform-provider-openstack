package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePlacementResourceProviderTraitsV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlacementResourceProviderTraitsV1Create,
		ReadContext:   resourcePlacementResourceProviderTraitsV1Read,
		UpdateContext: resourcePlacementResourceProviderTraitsV1Update,
		DeleteContext: resourcePlacementResourceProviderTraitsV1Delete,
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

			"traits": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Computed attributes
			"resource_provider_generation": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourcePlacementResourceProviderTraitsV1Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return resourcePlacementResourceProviderTraitsV1Update(ctx, d, meta)
}

func resourcePlacementResourceProviderTraitsV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	traits, err := resourceproviders.GetTraits(ctx, placementClient, rpID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving traits for resource provider"))
	}

	log.Printf("[DEBUG] Retrieved traits for resource provider %s: %#v", rpID, traits)

	d.SetId(rpID)
	d.Set("traits", traits.Traits)
	d.Set("resource_provider_generation", traits.ResourceProviderGeneration)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourcePlacementResourceProviderTraitsV1Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	placementClient, err := config.PlacementV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack placement client: %s", err)
	}

	rpID := d.Get("resource_provider_id").(string)

	// Convert set to string slice
	traitsSet := d.Get("traits").(*schema.Set)
	traits := make([]string, 0, traitsSet.Len())
	for _, trait := range traitsSet.List() {
		traits = append(traits, trait.(string))
	}

	// Get current generation from the resource provider
	// We need this for the update request
	rp, err := resourceproviders.Get(ctx, placementClient, rpID).Extract()
	if err != nil {
		return diag.Errorf("Error getting resource provider %s: %s", rpID, err)
	}

	updateOpts := resourceproviders.UpdateTraitsOpts{
		ResourceProviderGeneration: rp.Generation,
		Traits:                     traits,
	}

	log.Printf("[DEBUG] Updating traits for resource provider %s: %#v", rpID, updateOpts)

	_, err = resourceproviders.UpdateTraits(ctx, placementClient, rpID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating traits for resource provider %s: %s", rpID, err)
	}

	d.SetId(rpID)

	return resourcePlacementResourceProviderTraitsV1Read(ctx, d, meta)
}

func resourcePlacementResourceProviderTraitsV1Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	placementClient, err := config.PlacementV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack placement client: %s", err)
	}

	rpID := d.Get("resource_provider_id").(string)

	log.Printf("[DEBUG] Deleting traits for resource provider %s", rpID)

	err = resourceproviders.DeleteTraits(ctx, placementClient, rpID).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting traits for resource provider"))
	}

	return nil
}