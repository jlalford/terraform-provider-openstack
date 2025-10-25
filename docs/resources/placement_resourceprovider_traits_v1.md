---
subcategory: "Placement"
layout: "openstack"
page_title: "OpenStack: openstack_placement_resourceprovider_traits_v1"
sidebar_current: "docs-openstack-resource-placement-resourceprovider-traits-v1"
description: |-
  Manages V1 Placement resource provider traits within OpenStack.
---

# openstack\_placement\_resourceprovider\_traits\_v1

Manages V1 Placement resource provider traits within OpenStack.

## Example Usage

### Basic Traits

```hcl
resource "openstack_placement_resourceprovider_v1" "rp_1" {
  name = "my-resource-provider"
}

resource "openstack_placement_resourceprovider_traits_v1" "traits_1" {
  resource_provider_id = openstack_placement_resourceprovider_v1.rp_1.id

  traits = [
    "CUSTOM_MY_TRAIT",
    "HW_CPU_X86_VMX",
    "HW_GPU_API_DIRECTX_V10"
  ]
}
```

### Update Traits

```hcl
# Create a resource provider
resource "openstack_placement_resourceprovider_v1" "compute_node" {
  name = "compute-node-1"
}

# Manage its traits
resource "openstack_placement_resourceprovider_traits_v1" "compute_traits" {
  resource_provider_id = openstack_placement_resourceprovider_v1.compute_node.id

  traits = [
    "COMPUTE_TRUSTED_CERTS",
    "HW_CPU_X86_SGX",
    "HW_NIC_OFFLOAD_TSO",
    "CUSTOM_GOLD"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to manage the resource provider
    traits. If omitted, the `region` argument of the provider is used.
    Changing this creates new traits.

* `resource_provider_id` - (Required) The UUID of the resource provider.
    Changing this creates new traits.

* `traits` - (Required) A set of trait names to associate with the resource
    provider. Trait names must be strings of characters from the set of
    letters A-Z and a-z, the numbers 0-9, and the underscore character.
    Standard traits are defined by OpenStack, custom traits should be
    prefixed with `CUSTOM_`.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `resource_provider_id` - See Argument Reference above.
* `traits` - See Argument Reference above.
* `resource_provider_generation` - The generation of the resource provider.

## Import

Resource provider traits can be imported using the resource provider UUID, e.g.

```sh
terraform import openstack_placement_resourceprovider_traits_v1.traits_1 7c9f6d30-1d34-4d22-89e3-6c4d3a3c3d3e
```