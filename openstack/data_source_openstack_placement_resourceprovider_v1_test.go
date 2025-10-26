package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlacementV1ResourceProviderDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementV1ResourceProviderDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.openstack_placement_resourceprovider_v1.by_name", "id",
						"openstack_placement_resourceprovider_v1.rp_1", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_placement_resourceprovider_v1.by_name", "name", "test-provider-ds"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_placement_resourceprovider_v1.by_name", "generation"),
				),
			},
		},
	})
}

func TestAccPlacementV1ResourceProviderDataSource_byUUID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementV1ResourceProviderDataSourceByUUID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.openstack_placement_resourceprovider_v1.by_uuid", "id",
						"openstack_placement_resourceprovider_v1.rp_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_placement_resourceprovider_v1.by_uuid", "uuid",
						"openstack_placement_resourceprovider_v1.rp_1", "uuid"),
					resource.TestCheckResourceAttr(
						"data.openstack_placement_resourceprovider_v1.by_uuid", "name", "test-provider-uuid"),
				),
			},
		},
	})
}

const testAccPlacementV1ResourceProviderDataSourceBasic = `
resource "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "test-provider-ds"
}

data "openstack_placement_resourceprovider_v1" "by_name" {
  name = openstack_placement_resourceprovider_v1.rp_1.name
}
`

const testAccPlacementV1ResourceProviderDataSourceByUUID = `
resource "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "test-provider-uuid"
}

data "openstack_placement_resourceprovider_v1" "by_uuid" {
  uuid = openstack_placement_resourceprovider_v1.rp_1.id
}
`