package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPlacementV1ResourceProviderInventory_basic(t *testing.T) {
	var inventories resourceproviders.ResourceProviderInventories
	resourceName := "openstack_placement_resourceprovider_inventory_v1.inventory_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckPlacementV1ResourceProviderInventoryDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementV1ResourceProviderInventoryBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementV1ResourceProviderInventoryExists(t.Context(), resourceName, &inventories),
					resource.TestCheckResourceAttr(resourceName, "inventories.VCPU.total", "16"),
					resource.TestCheckResourceAttr(resourceName, "inventories.VCPU.allocation_ratio", "16.0"),
					resource.TestCheckResourceAttr(resourceName, "inventories.MEMORY_MB.total", "32768"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPlacementV1ResourceProviderInventory_update(t *testing.T) {
	var inventories resourceproviders.ResourceProviderInventories
	resourceName := "openstack_placement_resourceprovider_inventory_v1.inventory_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckPlacementV1ResourceProviderInventoryDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementV1ResourceProviderInventoryBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementV1ResourceProviderInventoryExists(t.Context(), resourceName, &inventories),
					resource.TestCheckResourceAttr(resourceName, "inventories.VCPU.total", "16"),
				),
			},
			{
				Config: testAccPlacementV1ResourceProviderInventoryUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementV1ResourceProviderInventoryExists(t.Context(), resourceName, &inventories),
					resource.TestCheckResourceAttr(resourceName, "inventories.VCPU.total", "32"),
					resource.TestCheckResourceAttr(resourceName, "inventories.DISK_GB.total", "1000"),
				),
			},
		},
	})
}

func testAccCheckPlacementV1ResourceProviderInventoryDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		placementClient, err := config.PlacementV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack placement client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_placement_resourceprovider_inventory_v1" {
				continue
			}

			rpID := rs.Primary.Attributes["resource_provider_id"]
			if rpID == "" {
				rpID = rs.Primary.ID
			}

			inventories, err := resourceproviders.GetInventories(ctx, placementClient, rpID).Extract()
			if err != nil {
				// If the resource provider doesn't exist, that's fine
				return nil
			}

			if len(inventories.Inventories) > 0 {
				return fmt.Errorf("Inventories still exist for resource provider %s", rpID)
			}
		}

		return nil
	}
}

func testAccCheckPlacementV1ResourceProviderInventoryExists(ctx context.Context, n string, inventories *resourceproviders.ResourceProviderInventories) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		placementClient, err := config.PlacementV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack placement client: %w", err)
		}

		rpID := rs.Primary.Attributes["resource_provider_id"]
		if rpID == "" {
			rpID = rs.Primary.ID
		}

		found, err := resourceproviders.GetInventories(ctx, placementClient, rpID).Extract()
		if err != nil {
			return err
		}

		if len(found.Inventories) == 0 {
			return errors.New("No inventories found")
		}

		*inventories = *found

		return nil
	}
}

const testAccPlacementV1ResourceProviderInventoryBasic = `
resource "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "test-provider-for-inventory"
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
    MEMORY_MB = {
      total            = 32768
      allocation_ratio = 1.5
      max_unit         = 32768
      min_unit         = 1
      reserved         = 512
      step_size        = 1
    }
  }
}
`

const testAccPlacementV1ResourceProviderInventoryUpdate = `
resource "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "test-provider-for-inventory"
}

resource "openstack_placement_resourceprovider_inventory_v1" "inventory_1" {
  resource_provider_id = openstack_placement_resourceprovider_v1.rp_1.id

  inventories = {
    VCPU = {
      total            = 32
      allocation_ratio = 16.0
      max_unit         = 32
      min_unit         = 1
      reserved         = 0
      step_size        = 1
    }
    MEMORY_MB = {
      total            = 32768
      allocation_ratio = 1.5
      max_unit         = 32768
      min_unit         = 1
      reserved         = 512
      step_size        = 1
    }
    DISK_GB = {
      total            = 1000
      allocation_ratio = 1.0
      max_unit         = 1000
      min_unit         = 1
      reserved         = 0
      step_size        = 1
    }
  }
}
`
