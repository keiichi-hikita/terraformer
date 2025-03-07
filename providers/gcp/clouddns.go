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

package gcp

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"

	"google.golang.org/api/dns/v1"
)

var cloudDNSAllowEmptyValues = []string{}

var cloudDNSAdditionalFields = map[string]string{}

type CloudDNSGenerator struct {
	GCPService
}

func (g CloudDNSGenerator) createZonesResources(ctx context.Context, svc *dns.Service, project string) []terraform_utils.Resource {
	resources := []terraform_utils.Resource{}
	managedZonesListCall := svc.ManagedZones.List(project)
	err := managedZonesListCall.Pages(ctx, func(listDNS *dns.ManagedZonesListResponse) error {
		for _, zone := range listDNS.ManagedZones {
			resources = append(resources, terraform_utils.NewResource(
				zone.Name,
				zone.Name,
				"google_dns_managed_zone",
				"google",
				map[string]string{
					"name":    zone.Name,
					"project": project,
				},
				cloudDNSAllowEmptyValues,
				cloudDNSAdditionalFields,
			))
			records := g.createRecordsResources(ctx, svc, project, zone.Name)
			resources = append(resources, records...)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return []terraform_utils.Resource{}
	}
	return resources
}
func (CloudDNSGenerator) createRecordsResources(ctx context.Context, svc *dns.Service, project, zoneName string) []terraform_utils.Resource {
	resources := []terraform_utils.Resource{}
	managedRecordsListCall := svc.ResourceRecordSets.List(project, zoneName)
	err := managedRecordsListCall.Pages(ctx, func(listDNS *dns.ResourceRecordSetsListResponse) error {
		for _, record := range listDNS.Rrsets {
			resources = append(resources, terraform_utils.NewResource(
				fmt.Sprintf("%s/%s/%s", zoneName, record.Name, record.Type),
				strings.TrimSuffix(record.Name+"-"+record.Type, "."),
				"google_dns_record_set",
				"google",
				map[string]string{
					"name":         record.Name,
					"managed_zone": zoneName,
					"type":         record.Type,
					"project":      project,
				},
				cloudDNSAllowEmptyValues,
				cloudDNSAdditionalFields,
			))
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return []terraform_utils.Resource{}
	}
	return resources
}

// Generate TerraformResources from GCP API,
// create terraform resource for each zone + each record
func (g *CloudDNSGenerator) InitResources() error {
	project := g.GetArgs()["project"]
	ctx := context.Background()
	svc, err := dns.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	g.Resources = g.createZonesResources(ctx, svc, project)
	g.PopulateIgnoreKeys()
	return nil
}

func (g *CloudDNSGenerator) PostConvertHook() error {
	for i, resourceRecord := range g.Resources {
		if resourceRecord.InstanceInfo.Type == "google_dns_managed_zone" {
			continue
		}
		item := resourceRecord.Item
		zoneID := item["managed_zone"].(string)
		for _, resourceZone := range g.Resources {
			if resourceZone.InstanceInfo.Type != "google_dns_managed_zone" {
				continue
			}
			if zoneID == resourceZone.InstanceState.ID {
				g.Resources[i].Item["managed_zone"] = "${google_dns_managed_zone." + resourceZone.ResourceName + ".name}"
				name := g.Resources[i].Item["name"].(string)
				name = strings.Replace(name, resourceZone.Item["dns_name"].(string), "", -1)
				g.Resources[i].Item["name"] = name + "${google_dns_managed_zone." + resourceZone.ResourceName + ".dns_name}"
			}
		}
	}
	return nil
}
