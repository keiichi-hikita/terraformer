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
	"github.com/nttcom/eclcloud/ecl/computevolume/v2/volumes"
	"github.com/nttcom/eclcloud/pagination"
)

type ComputeVolumeGenerator struct {
	ECLService
}

// createResources iterate on all ecl_compute_volume
func (g *ComputeVolumeGenerator) createResources(list *pagination.Pager) []terraform_utils.Resource {
	resources := []terraform_utils.Resource{}

	list.EachPage(func(page pagination.Page) (bool, error) {
		volumes, err := volumes.ExtractVolumes(page)
		if err != nil {
			return false, err
		}

		for _, v := range volumes {
			// Use volume ID as a name if the volume doesn't have a name
			name := v.Name
			if v.Name == "" {
				name = v.ID
			}

			resource := terraform_utils.NewResource(
				v.ID,
				name,
				"ecl_compute_volume_v2",
				"ecl",
				map[string]string{},
				[]string{},
				map[string]string{},
			)

			resources = append(resources, resource)
		}

		return true, nil
	})

	// runtime.Breakpoint()
	return resources
}

// Generate TerraformResources from ECL API,
func (g *ComputeVolumeGenerator) InitResources() error {
	opts, err := ecl.AuthOptionsFromEnv()
	if err != nil {
		return err
	}

	provider, err := ecl.AuthenticatedClient(opts)
	if err != nil {
		return err
	}

	// client, err := newComputeVolumeClient(provider, eo)
	client, err := ecl.NewComputeVolumeV2(provider, eclcloud.EndpointOpts{
		Region: g.GetArgs()["region"],
	})
	if err != nil {
		return err
	}

	list := volumes.List(client, nil)

	g.Resources = g.createResources(&list)
	g.PopulateIgnoreKeys()

	return nil
}
