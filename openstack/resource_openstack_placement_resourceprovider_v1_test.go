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

func TestAccPlacementV1ResourceProvider_basic(t *testing.T) {
	var rp resourceproviders.ResourceProvider
	resourceName := "openstack_placement_resourceprovider_v1.rp_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckPlacementV1ResourceProviderDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementV1ResourceProviderBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementV1ResourceProviderExists(t.Context(), resourceName, &rp),
					resource.TestCheckResourceAttr(resourceName, "name", "test-resource-provider"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "generation"),
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

func TestAccPlacementV1ResourceProvider_withParent(t *testing.T) {
	var parent, child resourceproviders.ResourceProvider
	parentName := "openstack_placement_resourceprovider_v1.parent"
	childName := "openstack_placement_resourceprovider_v1.child"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckPlacementV1ResourceProviderDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementV1ResourceProviderWithParent,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementV1ResourceProviderExists(t.Context(), parentName, &parent),
					testAccCheckPlacementV1ResourceProviderExists(t.Context(), childName, &child),
					resource.TestCheckResourceAttr(parentName, "name", "parent-provider"),
					resource.TestCheckResourceAttr(childName, "name", "child-provider"),
					resource.TestCheckResourceAttrPtr(childName, "parent_provider_uuid", &parent.UUID),
					resource.TestCheckResourceAttrPtr(childName, "root_provider_uuid", &parent.UUID),
				),
			},
		},
	})
}

func TestAccPlacementV1ResourceProvider_update(t *testing.T) {
	var rp resourceproviders.ResourceProvider
	resourceName := "openstack_placement_resourceprovider_v1.rp_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckPlacementV1ResourceProviderDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementV1ResourceProviderBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementV1ResourceProviderExists(t.Context(), resourceName, &rp),
					resource.TestCheckResourceAttr(resourceName, "name", "test-resource-provider"),
				),
			},
			{
				Config: testAccPlacementV1ResourceProviderUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementV1ResourceProviderExists(t.Context(), resourceName, &rp),
					resource.TestCheckResourceAttr(resourceName, "name", "updated-resource-provider"),
				),
			},
		},
	})
}

func testAccCheckPlacementV1ResourceProviderDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		placementClient, err := config.PlacementV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack placement client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_placement_resourceprovider_v1" {
				continue
			}

			_, err := resourceproviders.Get(ctx, placementClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Resource Provider still exists")
			}
		}

		return nil
	}
}

func testAccCheckPlacementV1ResourceProviderExists(ctx context.Context, n string, rp *resourceproviders.ResourceProvider) resource.TestCheckFunc {
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

		found, err := resourceproviders.Get(ctx, placementClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.UUID != rs.Primary.ID {
			return errors.New("Resource Provider not found")
		}

		*rp = *found

		return nil
	}
}

const testAccPlacementV1ResourceProviderBasic = `
resource "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "test-resource-provider"
}
`

const testAccPlacementV1ResourceProviderUpdate = `
resource "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "updated-resource-provider"
}
`

const testAccPlacementV1ResourceProviderWithParent = `
resource "openstack_placement_resourceprovider_v1" "parent" {
  name = "parent-provider"
}

resource "openstack_placement_resourceprovider_v1" "child" {
  name                 = "child-provider"
  parent_provider_uuid = openstack_placement_resourceprovider_v1.parent.id
}
`