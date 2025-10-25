package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/acls"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccKeyManagerV1SecretACL_basic(t *testing.T) {
	var acl acls.ACL
	resourceName := "openstack_keymanager_secret_acl_v1.acl_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckKeyManager(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckKeyManagerV1SecretACLDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerV1SecretACLBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyManagerV1SecretACLExists(t.Context(), resourceName, &acl),
					resource.TestCheckResourceAttr(resourceName, "read.0.project_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "read.0.users.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "acl_ref"),
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

func TestAccKeyManagerV1SecretACL_withUsers(t *testing.T) {
	var acl acls.ACL
	resourceName := "openstack_keymanager_secret_acl_v1.acl_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckKeyManager(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckKeyManagerV1SecretACLDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerV1SecretACLWithUsers,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyManagerV1SecretACLExists(t.Context(), resourceName, &acl),
					resource.TestCheckResourceAttr(resourceName, "read.0.project_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "read.0.users.#", "1"),
				),
			},
		},
	})
}

func TestAccKeyManagerV1SecretACL_update(t *testing.T) {
	var acl acls.ACL
	resourceName := "openstack_keymanager_secret_acl_v1.acl_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckKeyManager(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckKeyManagerV1SecretACLDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerV1SecretACLBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyManagerV1SecretACLExists(t.Context(), resourceName, &acl),
					resource.TestCheckResourceAttr(resourceName, "read.0.project_access", "false"),
				),
			},
			{
				Config: testAccKeyManagerV1SecretACLUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyManagerV1SecretACLExists(t.Context(), resourceName, &acl),
					resource.TestCheckResourceAttr(resourceName, "read.0.project_access", "true"),
				),
			},
		},
	})
}

func testAccCheckKeyManagerV1SecretACLDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		kmClient, err := config.KeyManagerV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack keymanager client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_keymanager_secret_acl_v1" {
				continue
			}

			secretID := rs.Primary.Attributes["secret_id"]
			if secretID == "" {
				secretID = rs.Primary.ID
			}

			acl, err := acls.GetSecretACL(ctx, kmClient, secretID).Extract()
			if err != nil {
				// If the secret doesn't exist, that's fine
				return nil
			}

			// Check if read ACL exists
			if readACL, ok := (*acl)["read"]; ok {
				// After deletion, ACL should revert to project-level access
				if !readACL.ProjectAccess || len(readACL.Users) > 0 {
					return errors.New("ACL still has custom configuration")
				}
			}
		}

		return nil
	}
}

func testAccCheckKeyManagerV1SecretACLExists(ctx context.Context, n string, acl *acls.ACL) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		kmClient, err := config.KeyManagerV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack keymanager client: %w", err)
		}

		secretID := rs.Primary.Attributes["secret_id"]
		if secretID == "" {
			secretID = rs.Primary.ID
		}

		found, err := acls.GetSecretACL(ctx, kmClient, secretID).Extract()
		if err != nil {
			return err
		}

		*acl = *found

		return nil
	}
}

const testAccKeyManagerV1SecretACLBasic = `
resource "openstack_keymanager_secret_v1" "secret_1" {
  name                 = "test-secret-for-acl"
  payload              = "super-secret"
  payload_content_type = "text/plain"
  secret_type          = "passphrase"
}

resource "openstack_keymanager_secret_acl_v1" "acl_1" {
  secret_id = openstack_keymanager_secret_v1.secret_1.secret_ref

  read {
    project_access = false
  }
}
`

const testAccKeyManagerV1SecretACLUpdate = `
resource "openstack_keymanager_secret_v1" "secret_1" {
  name                 = "test-secret-for-acl"
  payload              = "super-secret"
  payload_content_type = "text/plain"
  secret_type          = "passphrase"
}

resource "openstack_keymanager_secret_acl_v1" "acl_1" {
  secret_id = openstack_keymanager_secret_v1.secret_1.secret_ref

  read {
    project_access = true
  }
}
`

const testAccKeyManagerV1SecretACLWithUsers = `
data "openstack_identity_auth_scope_v3" "scope" {}

resource "openstack_keymanager_secret_v1" "secret_1" {
  name                 = "test-secret-for-acl"
  payload              = "super-secret"
  payload_content_type = "text/plain"
  secret_type          = "passphrase"
}

resource "openstack_keymanager_secret_acl_v1" "acl_1" {
  secret_id = openstack_keymanager_secret_v1.secret_1.secret_ref

  read {
    project_access = false
    users = [
      data.openstack_identity_auth_scope_v3.scope.user_id,
    ]
  }
}
`