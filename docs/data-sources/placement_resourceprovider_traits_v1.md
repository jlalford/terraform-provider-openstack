---
subcategory: "Placement"
layout: "openstack"
page_title: "OpenStack: openstack_placement_resourceprovider_traits_v1"
sidebar_current: "docs-openstack-datasource-placement-resourceprovider-traits-v1"
description: |-
  Get information on traits of an OpenStack Placement Resource Provider.
---

# openstack\_placement\_resourceprovider\_traits\_v1

Use this data source to get information about traits assigned to an OpenStack
Placement Resource Provider.

## Example Usage

```hcl
data "openstack_placement_resourceprovider_traits_v1" "traits_1" {
  resource_provider_id = "7c9f6d30-1d34-4d22-89e3-6c4d3a3c3d3e"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to query the resource provider traits.
    If omitted, the `region` argument of the provider is used.

* `resource_provider_id` - (Required) The UUID of the resource provider.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `resource_provider_id` - See Argument Reference above.
* `traits` - A set of trait names associated with the resource provider.
* `resource_provider_generation` - The generation of the resource provider.
