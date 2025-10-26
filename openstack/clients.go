package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
)

// PlacementV1Client returns a ServiceClient for making calls to the
// OpenStack Placement v1 API. An error will be returned if authentication
// or client creation was not possible.
func (config *Config) PlacementV1Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	client, err := config.InitializeClient(ctx, region)
	if err != nil {
		return nil, err
	}

	return openstack.NewPlacementV1(client, gophercloud.EndpointOpts{
		Region:       region,
		Availability: GetEndpointType(config),
	})
}