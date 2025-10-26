---
subcategory: "Placement"
layout: "openstack"
page_title: "OpenStack: openstack_placement_resourceprovider_usages_v1"
sidebar_current: "docs-openstack-datasource-placement-resourceprovider-usages-v1"
description: |-
  Get information on resource usages of an OpenStack Placement Resource Provider.
---

# openstack\_placement\_resourceprovider\_usages\_v1

Use this data source to get information about resource usages of an OpenStack
Placement Resource Provider.

## Example Usage

```hcl
data "openstack_placement_resourceprovider_usages_v1" "usages_1" {
  resource_provider_id = "7c9f6d30-1d34-4d22-89e3-6c4d3a3c3d3e"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to query the resource provider usages.
    If omitted, the `region` argument of the provider is used.

* `resource_provider_id` - (Required) The UUID of the resource provider.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `resource_provider_id` - See Argument Reference above.
* `usages` - A map of resource class names to usage amounts. Each key is a
    resource class name (e.g., "VCPU", "MEMORY_MB", "DISK_GB") and the value is
    the integer amount of that resource class currently in use.
* `resource_provider_generation` - The generation of the resource provider.
