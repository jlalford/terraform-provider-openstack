---
subcategory: "Key Manager / Barbican"
layout: "openstack"
page_title: "OpenStack: openstack_keymanager_secret_acl_v1"
sidebar_current: "docs-openstack-datasource-keymanager-secret-acl-v1"
description: |-
  Get information on a V1 Barbican secret ACL resource within OpenStack.
---

# openstack\_keymanager\_secret\_acl\_v1

Use this data source to get information about the ACL of an existing Barbican secret.

## Example Usage

```hcl
data "openstack_keymanager_secret_acl_v1" "secret_acl" {
  secret_id = "d4a6b0f7-7e1a-4b1c-9c8e-5f3d2a1b0c9d"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V1 KeyManager client.
  A KeyManager client is needed to fetch the secret ACL. If omitted, the `region`
  argument of the provider is used.

* `secret_id` - (Required) The UUID of the secret for which to retrieve the ACL.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `secret_id` - See Argument Reference above.
* `read` - The read ACL settings. This is a list containing a single map with
  the following fields:
  * `project_access` - Whether the secret is accessible project-wide.
  * `users` - The list of user IDs that are allowed to access the secret.
  * `created` - The date the ACL settings were created.
  * `updated` - The date the ACL settings were last updated.
