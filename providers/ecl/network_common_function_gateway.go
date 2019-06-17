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
	"github.com/nttcom/eclcloud/ecl/network/v2/common_function_gateways"
	"github.com/nttcom/eclcloud/pagination"
)

type NetworkCommonFunctionGatewayGenerator struct {
	ECLService
}

// createResources iterate on all openstack_networking_secgroup_v2
func (g *NetworkCommonFunctionGatewayGenerator) createResources(list *pagination.Pager) []terraform_utils.Resource {
	resources := []terraform_utils.Resource{}

	list.EachPage(func(page pagination.Page) (bool, error) {
		cfGws, err := common_function_gateways.ExtractCommonFunctionGateways(page)
		if err != nil {
			return false, err
		}

		for _, c := range cfGws {
			name := c.Name
			if c.Name == "" {
				name = c.ID
			}

			resource := terraform_utils.NewResource(
				c.ID,
				name,
				"ecl_network_common_function_gateway_v2",
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

// Generate TerraformResources from OpenStack API,
func (g *NetworkCommonFunctionGatewayGenerator) InitResources() error {
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

	list := common_function_gateways.List(client, common_function_gateways.ListOpts{})

	g.Resources = g.createResources(&list)
	g.PopulateIgnoreKeys()

	return nil
}
