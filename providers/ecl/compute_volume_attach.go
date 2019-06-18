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
	"fmt"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"github.com/nttcom/eclcloud"
	"github.com/nttcom/eclcloud/ecl"
	"github.com/nttcom/eclcloud/ecl/compute/v2/extensions/volumeattach"
	"github.com/nttcom/eclcloud/ecl/compute/v2/servers"
	"github.com/nttcom/eclcloud/pagination"
)

type ComputeVolumeAttachGenerator struct {
	ECLService
}

// createResources iterate on all ecl_compute_instance_v2
func (g *ComputeVolumeAttachGenerator) createResources(list *pagination.Pager, serverID string) []terraform_utils.Resource {
	resources := []terraform_utils.Resource{}

	list.EachPage(func(page pagination.Page) (bool, error) {

		vas, err := volumeattach.ExtractVolumeAttachments(page)
		if err != nil {
			return false, err
		}

		if len(vas) > 0 {
			for _, va := range vas {
				volumeID := va.VolumeID
				id := fmt.Sprintf("%s/%s", volumeID, serverID)
				name := fmt.Sprintf("%s-attach", volumeID)

				resource := terraform_utils.NewResource(
					id,
					name,
					"ecl_compute_volume_attach_v2",
					"ecl",
					map[string]string{},
					[]string{},
					map[string]string{},
				)

				resources = append(resources, resource)
			}
		}
		return true, nil
	})

	return resources
}

// Generate TerraformResources from ECL API,
func (g *ComputeVolumeAttachGenerator) InitResources() error {
	opts, err := ecl.AuthOptionsFromEnv()
	if err != nil {
		return err
	}

	provider, err := ecl.AuthenticatedClient(opts)
	if err != nil {
		return err
	}

	client, err := ecl.NewComputeV2(provider, eclcloud.EndpointOpts{
		Region: g.GetArgs()["region"],
	})
	if err != nil {
		return err
	}

	list := servers.List(client, nil)

	resources := []terraform_utils.Resource{}
	list.EachPage(func(page pagination.Page) (bool, error) {
		svs, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}

		for _, s := range svs {
			serverID := s.ID
			volumeAttachPage := volumeattach.List(client, serverID)

			volumeAttachResources := g.createResources(&volumeAttachPage, serverID)
			resources = append(resources, volumeAttachResources...)
		}
		return true, nil
	})

	g.Resources = resources
	g.PopulateIgnoreKeys()

	return nil
}
