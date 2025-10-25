---
subcategory: "Placement"
layout: "openstack"
page_title: "OpenStack: openstack_placement_resourceprovider_inventory_v1"
sidebar_current: "docs-openstack-datasource-placement-resourceprovider-inventory-v1"
description: |-
  Get information on OpenStack Placement Resource Provider Inventories.
---

# openstack\_placement\_resourceprovider\_inventory\_v1

Use this data source to get information about OpenStack Placement
Resource Provider Inventories.

~> **Note:** To manage resource provider inventories, use the
`openstack_placement_resourceprovider_inventory_v1` resource.

## Example Usage

### Get All Inventories

```hcl
data "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "compute-node-1"
}

data "openstack_placement_resourceprovider_inventory_v1" "inv_1" {
  resource_provider_id = data.openstack_placement_resourceprovider_v1.rp_1.id
}

output "vcpu_inventory" {
  value = data.openstack_placement_resourceprovider_inventory_v1.inv_1.inventories["VCPU"]
}

output "memory_inventory" {
  value = data.openstack_placement_resourceprovider_inventory_v1.inv_1.inventories["MEMORY_MB"]
}
```

### Process Inventory Information

```hcl
data "openstack_placement_resourceprovider_inventory_v1" "compute_inv" {
  resource_provider_id = "7c9f6d30-1d34-4d22-89e3-6c4d3a3c3d3e"
}

locals {
  vcpu_total = lookup(
    data.openstack_placement_resourceprovider_inventory_v1.compute_inv.inventories["VCPU"],
    "total",
    0
  )

  memory_total = lookup(
    data.openstack_placement_resourceprovider_inventory_v1.compute_inv.inventories["MEMORY_MB"],
    "total",
    0
  )
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to query the inventories.
    If omitted, the `region` argument of the provider is used.

* `resource_provider_id` - (Required) The UUID of the resource provider.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `resource_provider_id` - See Argument Reference above.
* `resource_provider_generation` - The generation of the resource provider.
* `inventories` - A map of resource class names to inventory details.
    Each inventory contains:
  * `allocation_ratio` - The allocation ratio for this resource class.
  * `max_unit` - The maximum amount of this resource that may be requested
      in a single allocation request.
  * `min_unit` - The minimum amount of this resource that must be requested
      in a single allocation request.
  * `reserved` - The amount of this resource that is reserved and not available
      for allocation.
  * `step_size` - The granularity of allocation requests for this resource.
  * `total` - The total amount of this resource.