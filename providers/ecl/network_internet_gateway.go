// Copyright 2018 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ecl

import (
	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"github.com/nttcom/eclcloud"
	"github.com/nttcom/eclcloud/ecl"
	"github.com/nttcom/eclcloud/ecl/network/v2/internet_gateways"
	"github.com/nttcom/eclcloud/pagination"
)

type NetworkInternetGatewayGenerator struct {
	ECLService
}

// createResources iterate on all openstack_networking_secgroup_v2
func (g *NetworkInternetGatewayGenerator) createResources(list *pagination.Pager) []terraform_utils.Resource {
	resources := []terraform_utils.Resource{}

	list.EachPage(func(page pagination.Page) (bool, error) {
		igs, err := internet_gateways.ExtractInternetGateways(page)
		if err != nil {
			return false, err
		}

		for _, ig := range igs {
			name := ig.Name
			if ig.Name == "" {
				name = ig.ID
			}

			resource := terraform_utils.NewResource(
				ig.ID,
				name,
				"ecl_network_internet_gateway_v2",
				"ecl",
				map[string]string{},
				[]string{},
				map[string]string{},
			)
			resources = append(resources, resource)
		}

		return true, nil
	})

	return resources
}

// Generate TerraformResources from ECL API,
func (g *NetworkInternetGatewayGenerator) InitResources() error {
	opts, err := ecl.AuthOptionsFromEnv()
	if err != nil {
		return err
	}

	provider, err := ecl.AuthenticatedClient(opts)
	if err != nil {
		return err
	}

	client, err := ecl.NewNetworkV2(provider, eclcloud.EndpointOpts{
		Region: g.GetArgs()["region"],
	})
	if err != nil {
		return err
	}

	list := internet_gateways.List(client, internet_gateways.ListOpts{})

	g.Resources = g.createResources(&list)
	g.PopulateIgnoreKeys()

	return nil
}
