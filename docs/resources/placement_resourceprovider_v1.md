---
subcategory: "Placement"
layout: "openstack"
page_title: "OpenStack: openstack_placement_resourceprovider_v1"
sidebar_current: "docs-openstack-resource-placement-resourceprovider-v1"
description: |-
  Manages a V1 Placement resource provider resource within OpenStack.
---

# openstack\_placement\_resourceprovider\_v1

Manages a V1 Placement resource provider resource within OpenStack.

## Example Usage

### Basic Resource Provider

```hcl
resource "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "my-resource-provider"
}
```

### Resource Provider with Parent

```hcl
resource "openstack_placement_resourceprovider_v1" "parent" {
  name = "parent-provider"
}

resource "openstack_placement_resourceprovider_v1" "child" {
  name                 = "child-provider"
  parent_provider_uuid = openstack_placement_resourceprovider_v1.parent.id
}
```

### Resource Provider with explicit UUID

```hcl
resource "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "my-resource-provider"
  uuid = "7c9f6d30-1d34-4d22-89e3-6c4d3a3c3d3e"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the resource provider.
    If omitted, the `region` argument of the provider is used. Changing this
    creates a new resource provider.

* `name` - (Required) The name of the resource provider.

* `uuid` - (Optional) The UUID of the resource provider. If omitted, one will
    be generated. Changing this creates a new resource provider.

* `parent_provider_uuid` - (Optional) The UUID of the immediate parent of the
    resource provider. Requires microversion 1.14 or above.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `uuid` - See Argument Reference above.
* `parent_provider_uuid` - See Argument Reference above.
* `generation` - The generation of the resource provider.
* `root_provider_uuid` - The UUID of the root provider in this provider tree.
* `links` - A list of links associated with the resource provider.