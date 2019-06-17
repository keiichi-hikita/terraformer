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
	"github.com/nttcom/eclcloud/ecl/network/v2/networks"
	// "github.com/nttcom/eclcloud/ecl/network/v2/extensions/security/rules"
	"github.com/nttcom/eclcloud/pagination"
)

type NetworkNetworkGenerator struct {
	ECLService
}

// createResources iterate on all openstack_networking_secgroup_v2
func (g *NetworkNetworkGenerator) createResources(list *pagination.Pager) []terraform_utils.Resource {
	resources := []terraform_utils.Resource{}

	list.EachPage(func(page pagination.Page) (bool, error) {
		networks, err := networks.ExtractNetworks(page)
		if err != nil {
			return false, err
		}

		for _, n := range networks {
			name := n.Name
			if n.Name == "" {
				name = n.ID
			}

			resource := terraform_utils.NewResource(
				n.ID,
				name,
				"ecl_network_network_v2",
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
func (g *NetworkNetworkGenerator) InitResources() error {
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

	list := networks.List(client, networks.ListOpts{})

	g.Resources = g.createResources(&list)
	g.PopulateIgnoreKeys()

	return nil
}

// func (g *NetworkNetworkGenerator) PostConvertHook() error {
// 	for i, r := range g.Resources {
// 		if r.InstanceInfo.Type != "ecl_network_network_v2" {
// 			continue
// 		}
// 		for _, sg := range g.Resources {
// 			if sg.InstanceInfo.Type != "openstack_networking_secgroup_v2" {
// 				continue
// 			}
// 			if r.InstanceState.Attributes["security_group_id"] == sg.InstanceState.Attributes["id"] {
// 				g.Resources[i].Item["security_group_id"] = "${openstack_networking_secgroup_v2." + sg.ResourceName + ".id}"
// 			}
// 		}
// 	}

// 	return nil
// }
