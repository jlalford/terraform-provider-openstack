---
subcategory: "Placement"
layout: "openstack"
page_title: "OpenStack: openstack_placement_resourceprovider_inventory_v1"
sidebar_current: "docs-openstack-resource-placement-resourceprovider-inventory-v1"
description: |-
  Manages V1 Placement resource provider inventories within OpenStack.
---

# openstack\_placement\_resourceprovider\_inventory\_v1

Manages V1 Placement resource provider inventories within OpenStack.

## Example Usage

### Basic Inventory

```hcl
resource "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "my-compute-node"
}

resource "openstack_placement_resourceprovider_inventory_v1" "inventory_1" {
  resource_provider_id = openstack_placement_resourceprovider_v1.rp_1.id

  inventories = {
    VCPU = {
      total            = 16
      allocation_ratio = 16.0
      max_unit         = 16
      min_unit         = 1
      reserved         = 0
      step_size        = 1
    }
    MEMORY_MB = {
      total            = 32768
      allocation_ratio = 1.5
      max_unit         = 32768
      min_unit         = 1
      reserved         = 512
      step_size        = 1
    }
    DISK_GB = {
      total            = 1000
      allocation_ratio = 1.0
      max_unit         = 1000
      min_unit         = 1
      reserved         = 0
      step_size        = 1
    }
  }
}
```

### GPU Resource Provider

```hcl
resource "openstack_placement_resourceprovider_v1" "gpu_node" {
  name = "gpu-compute-node"
}

resource "openstack_placement_resourceprovider_inventory_v1" "gpu_inventory" {
  resource_provider_id = openstack_placement_resourceprovider_v1.gpu_node.id

  inventories = {
    VCPU = {
      total            = 32
      allocation_ratio = 4.0
      max_unit         = 32
      min_unit         = 1
      reserved         = 0
      step_size        = 1
    }
    VGPU = {
      total            = 4
      allocation_ratio = 1.0
      max_unit         = 4
      min_unit         = 1
      reserved         = 0
      step_size        = 1
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to manage the resource provider
    inventory. If omitted, the `region` argument of the provider is used.
    Changing this creates a new inventory.

* `resource_provider_id` - (Required) The UUID of the resource provider.
    Changing this creates a new inventory.

* `inventories` - (Required) A map of inventory records, keyed by resource class
    name. Each inventory record is a map with the following fields:

    * `total` - (Required) The actual amount of the resource that the provider
        can accommodate.

    * `allocation_ratio` - (Required) It is used in determining whether consumption
        of the resource of the provider can exceed physical constraints. For example,
        for a vCPU resource with: allocation_ratio = 16.0, total = 8. Overall
        capacity is calculated as: capacity = total * allocation_ratio, i.e. 128 vCPUs.

    * `max_unit` - (Required) A maximum amount any single allocation against an
        inventory can have.

    * `min_unit` - (Optional) A minimum amount any single allocation against an
        inventory can have. Defaults to 1.

    * `reserved` - (Optional) The amount of the resource a provider has reserved
        for its own use. Defaults to 0.

    * `step_size` - (Optional) A representation of the divisible amount of the
        resource that may be requested. For example, step_size = 5 means that only
        values divisible by 5 (5, 10, 15, etc.) can be requested. Defaults to 1.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `resource_provider_id` - See Argument Reference above.
* `inventories` - See Argument Reference above.
* `resource_provider_generation` - The generation of the resource provider.

## Import

Resource provider inventories can be imported using the resource provider UUID, e.g.

```sh
terraform import openstack_placement_resourceprovider_inventory_v1.inventory_1 7c9f6d30-1d34-4d22-89e3-6c4d3a3c3d3e
```
