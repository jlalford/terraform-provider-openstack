package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKeyManagerSecretACLV1DataSource_basic(t *testing.T) {
	resourceName := "data.openstack_keymanager_secret_acl_v1.acl_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretACLV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "secret_id"),
					resource.TestCheckResourceAttr(resourceName, "read.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "read.0.project_access", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "read.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "read.0.updated"),
				),
			},
		},
	})
}

const testAccKeyManagerSecretACLV1DataSourceBasic = `
resource "openstack_keymanager_secret_v1" "secret_1" {
  name                 = "test-secret-for-acl-ds"
  payload              = "my-secret-payload"
  secret_type          = "passphrase"
  payload_content_type = "text/plain"
}

resource "openstack_keymanager_secret_acl_v1" "acl_1" {
  secret_id = openstack_keymanager_secret_v1.secret_1.id

  read {
    project_access = true
  }
}

data "openstack_keymanager_secret_acl_v1" "acl_1" {
  secret_id = openstack_keymanager_secret_acl_v1.acl_1.secret_id
}
`
