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
	"github.com/nttcom/eclcloud/ecl/network/v2/ports"
	"github.com/nttcom/eclcloud/pagination"
)

type NetworkPortGenerator struct {
	ECLService
}

// createResources iterate on all ecl_network_port_v2
func (g *NetworkPortGenerator) createResources(list *pagination.Pager) []terraform_utils.Resource {
	resources := []terraform_utils.Resource{}

	list.EachPage(func(page pagination.Page) (bool, error) {
		ports, err := ports.ExtractPorts(page)
		if err != nil {
			return false, err
		}

		for _, p := range ports {
			name := p.Name
			if p.Name == "" {
				name = p.ID
			}

			resource := terraform_utils.NewResource(
				p.ID,
				name,
				"ecl_network_port_v2",
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
func (g *NetworkPortGenerator) InitResources() error {
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

	list := ports.List(client, ports.ListOpts{})

	g.Resources = g.createResources(&list)
	g.PopulateIgnoreKeys()

	return nil
}
