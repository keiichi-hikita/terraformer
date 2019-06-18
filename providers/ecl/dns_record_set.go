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
	"github.com/nttcom/eclcloud/ecl/dns/v2/recordsets"
	"github.com/nttcom/eclcloud/ecl/dns/v2/zones"
	"github.com/nttcom/eclcloud/pagination"
)

type DNSRecordSetGenerator struct {
	ECLService
}

// createResources iterate on all ecl_dns_recordset_v2
func (g *DNSRecordSetGenerator) createResources(list *pagination.Pager, zoneID string) []terraform_utils.Resource {
	resources := []terraform_utils.Resource{}

	list.EachPage(func(page pagination.Page) (bool, error) {
		records, err := recordsets.ExtractRecordSets(page)
		if err != nil {
			return false, err
		}

		for _, r := range records {
			id := fmt.Sprintf("%s/%s", zoneID, r.ID)
			resource := terraform_utils.NewResource(
				id,
				id,
				"ecl_dns_recordset_v2",
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
func (g *DNSRecordSetGenerator) InitResources() error {
	opts, err := ecl.AuthOptionsFromEnv()
	if err != nil {
		return err
	}

	provider, err := ecl.AuthenticatedClient(opts)
	if err != nil {
		return err
	}

	client, err := ecl.NewDNSV2(provider, eclcloud.EndpointOpts{
		Region: g.GetArgs()["region"],
	})
	if err != nil {
		return err
	}

	list := zones.List(client, zones.ListOpts{})

	resources := []terraform_utils.Resource{}

	list.EachPage(func(page pagination.Page) (bool, error) {
		zns, err := zones.ExtractZones(page)
		if err != nil {
			return false, err
		}

		for _, z := range zns {
			zoneID := z.ID
			recordSetPage := recordsets.ListByZone(client, zoneID, recordsets.ListOpts{})

			recordSetResources := g.createResources(&recordSetPage, zoneID)
			resources = append(resources, recordSetResources...)
		}
		return true, nil
	})

	g.Resources = resources
	g.PopulateIgnoreKeys()

	return nil
}
