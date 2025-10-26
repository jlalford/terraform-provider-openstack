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

func TestAccPlacementV1ResourceProviderTraits_basic(t *testing.T) {
	var traits resourceproviders.ResourceProviderTraits
	resourceName := "openstack_placement_resourceprovider_traits_v1.traits_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckPlacementV1ResourceProviderTraitsDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementV1ResourceProviderTraitsBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementV1ResourceProviderTraitsExists(t.Context(), resourceName, &traits),
					resource.TestCheckResourceAttr(resourceName, "traits.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "traits.*", "CUSTOM_TRAIT_1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "traits.*", "CUSTOM_TRAIT_2"),
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

func TestAccPlacementV1ResourceProviderTraits_update(t *testing.T) {
	var traits resourceproviders.ResourceProviderTraits
	resourceName := "openstack_placement_resourceprovider_traits_v1.traits_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckPlacementV1ResourceProviderTraitsDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementV1ResourceProviderTraitsBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementV1ResourceProviderTraitsExists(t.Context(), resourceName, &traits),
					resource.TestCheckResourceAttr(resourceName, "traits.#", "2"),
				),
			},
			{
				Config: testAccPlacementV1ResourceProviderTraitsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementV1ResourceProviderTraitsExists(t.Context(), resourceName, &traits),
					resource.TestCheckResourceAttr(resourceName, "traits.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "traits.*", "CUSTOM_TRAIT_1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "traits.*", "CUSTOM_TRAIT_3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "traits.*", "HW_CPU_X86_VMX"),
				),
			},
		},
	})
}

func testAccCheckPlacementV1ResourceProviderTraitsDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		placementClient, err := config.PlacementV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack placement client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_placement_resourceprovider_traits_v1" {
				continue
			}

			rpID := rs.Primary.Attributes["resource_provider_id"]
			if rpID == "" {
				rpID = rs.Primary.ID
			}

			traits, err := resourceproviders.GetTraits(ctx, placementClient, rpID).Extract()
			if err != nil {
				// If the resource provider doesn't exist, that's fine
				return nil
			}

			if len(traits.Traits) > 0 {
				return fmt.Errorf("Traits still exist for resource provider %s", rpID)
			}
		}

		return nil
	}
}

func testAccCheckPlacementV1ResourceProviderTraitsExists(ctx context.Context, n string, traits *resourceproviders.ResourceProviderTraits) resource.TestCheckFunc {
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

		found, err := resourceproviders.GetTraits(ctx, placementClient, rpID).Extract()
		if err != nil {
			return err
		}

		if len(found.Traits) == 0 {
			return errors.New("No traits found")
		}

		*traits = *found

		return nil
	}
}

const testAccPlacementV1ResourceProviderTraitsBasic = `
resource "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "test-provider-for-traits"
}

resource "openstack_placement_resourceprovider_traits_v1" "traits_1" {
  resource_provider_id = openstack_placement_resourceprovider_v1.rp_1.id

  traits = [
    "CUSTOM_TRAIT_1",
    "CUSTOM_TRAIT_2",
  ]
}
`

const testAccPlacementV1ResourceProviderTraitsUpdate = `
resource "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "test-provider-for-traits"
}

resource "openstack_placement_resourceprovider_traits_v1" "traits_1" {
  resource_provider_id = openstack_placement_resourceprovider_v1.rp_1.id

  traits = [
    "CUSTOM_TRAIT_1",
    "CUSTOM_TRAIT_3",
    "HW_CPU_X86_VMX",
  ]
}
`