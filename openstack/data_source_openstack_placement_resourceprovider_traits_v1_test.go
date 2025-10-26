package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlacementV1ResourceProviderTraitsDataSource_basic(t *testing.T) {
	resourceName := "data.openstack_placement_resourceprovider_traits_v1.traits_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementV1ResourceProviderTraitsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "resource_provider_id"),
					resource.TestCheckResourceAttr(resourceName, "traits.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "traits.*", "CUSTOM_TRAIT_1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "traits.*", "CUSTOM_TRAIT_2"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_provider_generation"),
				),
			},
		},
	})
}

const testAccPlacementV1ResourceProviderTraitsDataSourceBasic = `
resource "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "test-provider-for-traits-ds"
}

resource "openstack_placement_resourceprovider_traits_v1" "traits_1" {
  resource_provider_id = openstack_placement_resourceprovider_v1.rp_1.id

  traits = [
    "CUSTOM_TRAIT_1",
    "CUSTOM_TRAIT_2",
  ]
}

data "openstack_placement_resourceprovider_traits_v1" "traits_1" {
  resource_provider_id = openstack_placement_resourceprovider_traits_v1.traits_1.resource_provider_id
}
`
