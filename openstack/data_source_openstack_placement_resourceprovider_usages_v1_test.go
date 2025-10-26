package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlacementV1ResourceProviderUsagesDataSource_basic(t *testing.T) {
	resourceName := "data.openstack_placement_resourceprovider_usages_v1.usages_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementV1ResourceProviderUsagesDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "resource_provider_id"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_provider_generation"),
					resource.TestCheckResourceAttrSet(resourceName, "usages.%"),
				),
			},
		},
	})
}

const testAccPlacementV1ResourceProviderUsagesDataSourceBasic = `
resource "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "test-provider-for-usages-ds"
}

resource "openstack_placement_resourceprovider_inventory_v1" "inventory_1" {
  resource_provider_id = openstack_placement_resourceprovider_v1.rp_1.id

  inventories = {
    VCPU = {
      total            = 16
      allocation_ratio = 16.0
      max_unit         = 16
      min_unit         = 1
      reserved         = 0
      step_size        = 1
    }
  }
}

data "openstack_placement_resourceprovider_usages_v1" "usages_1" {
  resource_provider_id = openstack_placement_resourceprovider_inventory_v1.inventory_1.resource_provider_id
}
`
