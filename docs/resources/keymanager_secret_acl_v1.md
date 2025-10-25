---
subcategory: "Key Manager / Barbican"
layout: "openstack"
page_title: "OpenStack: openstack_keymanager_secret_acl_v1"
sidebar_current: "docs-openstack-resource-keymanager-secret-acl-v1"
description: |-
  Manages a V1 Barbican secret ACL resource within OpenStack.
---

# openstack\_keymanager\_secret\_acl\_v1

Manages a V1 Barbican secret ACL resource within OpenStack.

## Example Usage

### Basic Read ACL

```hcl
resource "openstack_keymanager_secret_v1" "secret_1" {
  name                 = "mysecret"
  payload              = "super-secret-password"
  payload_content_type = "text/plain"
  secret_type          = "passphrase"
}

resource "openstack_keymanager_secret_acl_v1" "acl_1" {
  secret_id = openstack_keymanager_secret_v1.secret_1.secret_ref

  read {
    project_access = false
    users = [
      "userid1",
      "userid2",
    ]
  }
}
```

### Project-level Access

```hcl
resource "openstack_keymanager_secret_v1" "secret_1" {
  name                 = "shared-secret"
  payload              = file("secret.txt")
  payload_content_type = "text/plain"
  secret_type          = "passphrase"
}

resource "openstack_keymanager_secret_acl_v1" "acl_1" {
  secret_id = openstack_keymanager_secret_v1.secret_1.secret_ref

  read {
    project_access = true
    users          = []
  }
}
```

### Update ACL Users

```hcl
variable "allowed_users" {
  type = list(string)
  default = ["user1", "user2", "user3"]
}

resource "openstack_keymanager_secret_v1" "secret_1" {
  name                 = "restricted-secret"
  payload              = "confidential-data"
  payload_content_type = "text/plain"
  secret_type          = "passphrase"
}

resource "openstack_keymanager_secret_acl_v1" "acl_1" {
  secret_id = openstack_keymanager_secret_v1.secret_1.secret_ref

  read {
    project_access = false
    users          = var.allowed_users
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the secret ACL.
    If omitted, the `region` argument of the provider is used. Changing this
    creates a new ACL.

* `secret_id` - (Required) The UUID or reference of the secret. Changing this
    creates a new ACL.

* `read` - (Required) The read ACL configuration. The `read` block is
    documented below.

The `read` block supports:

* `project_access` - (Optional) Whether to enable project-level access to
    the secret. Defaults to `true`. When `true`, all users in the project
    can read the secret.

* `users` - (Optional) A set of user IDs who have read access to the secret.
    This is typically a list of Keystone user UUIDs.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `secret_id` - See Argument Reference above.
* `read` - See Argument Reference above. Additionally exports:
  * `created` - The date the ACL was created.
  * `updated` - The date the ACL was last updated.
* `acl_ref` - The ACL reference / URL.

## Import

ACLs can be imported using the secret UUID, e.g.

```sh
terraform import openstack_keymanager_secret_acl_v1.acl_1 c3723b78-9f62-4b2f-8b2e-5c4c8a23fe38
```

## Notes

* Currently only read ACLs are supported by this resource. Write ACLs may be
  added in a future version.
* When deleting this resource, all ACLs for the secret are removed, effectively
  resetting the secret to its default project-level access.