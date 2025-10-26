---
subcategory: "Placement"
layout: "openstack"
page_title: "OpenStack: openstack_placement_resourceprovider_v1"
sidebar_current: "docs-openstack-datasource-placement-resourceprovider-v1"
description: |-
  Get information on an OpenStack Placement Resource Provider.
---

# openstack\_placement\_resourceprovider\_v1

Use this data source to get information about an OpenStack Placement
Resource Provider.

## Example Usage

### By UUID

```hcl
data "openstack_placement_resourceprovider_v1" "rp_1" {
  uuid = "7c9f6d30-1d34-4d22-89e3-6c4d3a3c3d3e"
}
```

### By Name

```hcl
data "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "my-resource-provider"
}
```

### With Filters

```hcl
data "openstack_placement_resourceprovider_v1" "gpu_provider" {
  required = "CUSTOM_GPU"
}
```

### In Tree

```hcl
data "openstack_placement_resourceprovider_v1" "child_providers" {
  in_tree = "parent-provider-uuid"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to query the resource provider.
    If omitted, the `region` argument of the provider is used.

* `uuid` - (Optional) The UUID of the resource provider. If specified, this
    takes precedence over other filters.

* `name` - (Optional) The name of the resource provider to filter the list.

* `member_of` - (Optional) A string representing aggregate UUIDs to filter
    the list. Format: `in:uuid1,uuid2` or `!in:uuid3,uuid4`.

* `resources` - (Optional) A comma-separated list of strings indicating an
    amount of resource of a specified class that a provider must have the
    capacity and availability to serve. Format: `VCPU:4,DISK_GB:64,MEMORY_MB:2048`.

* `in_tree` - (Optional) A UUID of a resource provider. The returned resource
    providers will be in the same provider tree as the specified provider.

* `required` - (Optional) A comma-delimited list of string trait names.
    Resource providers must have all of these traits to be included in results.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `uuid` - The UUID of the resource provider.
* `name` - The name of the resource provider.
* `parent_provider_uuid` - The UUID of the immediate parent of the resource provider.
* `root_provider_uuid` - The UUID of the root provider in this provider tree.
* `generation` - The generation of the resource provider.
* `links` - A list of links associated with the resource provider.