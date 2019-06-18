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
	"os"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"

	"github.com/pkg/errors"
)

type ECLProvider struct {
	terraform_utils.Provider
	region string
}

const eclProviderVersion = "~>1.0.0"

func (p ECLProvider) GetProviderData(arg ...string) map[string]interface{} {
	return map[string]interface{}{
		"provider": map[string]interface{}{
			"ecl": map[string]interface{}{
				"version": eclProviderVersion,
				"region":  p.region,
			},
		},
	}
}

// check projectName in env params
func (p *ECLProvider) Init(args []string) error {
	p.region = args[0]
	// terraform work with env param OS_REGION_NAME
	err := os.Setenv("OS_REGION_NAME", p.region)
	if err != nil {
		return err
	}
	return nil
}

func (p *ECLProvider) GetName() string {
	return "ecl"
}

func (p *ECLProvider) InitService(serviceName string) error {
	var isSupported bool
	if _, isSupported = p.GetSupportedService()[serviceName]; !isSupported {
		return errors.New("ecl: " + serviceName + " not supported service")
	}
	p.Service = p.GetSupportedService()[serviceName]
	p.Service.SetName(serviceName)
	p.Service.SetProviderName(p.GetName())
	p.Service.SetArgs(map[string]string{
		"region": p.region,
	})
	return nil
}

// GetSupportedService return map of support service for ECL
func (p *ECLProvider) GetSupportedService() map[string]terraform_utils.ServiceGenerator {
	return map[string]terraform_utils.ServiceGenerator{
		"computeKeypair":               &ComputeKeypairGenerator{},
		"computeServer":                &ComputeServerGenerator{},
		"computeVolumeAttach":          &ComputeVolumeAttachGenerator{},
		"computeVolume":                &ComputeVolumeGenerator{},
		"dnsZone":                      &DNSZoneGenerator{},
		"dnsRecordSet":                 &DNSRecordSetGenerator{},
		"networkCommonFunctionGateway": &NetworkCommonFunctionGatewayGenerator{},
		"networkGatewayInterface":      &NetworkGatewayInterfaceGenerator{},
		"networkInternetGateway":       &NetworkInternetGatewayGenerator{},
		"networkNetwork":               &NetworkNetworkGenerator{},
		"networkPort":                  &NetworkPortGenerator{},
		"networkPublicIP":              &NetworkPublicIPGenerator{},
		"networkStaticRoute":           &NetworkStaticRouteGenerator{},
		"networkSubnet":                &NetworkSubnetGenerator{},
		"storageVirtualStorage":        &StorageVirtualStorageGenerator{},
		"storageVolume":                &StorageVolumeGenerator{},
	}
}

func (ECLProvider) GetResourceConnections() map[string]map[string][]string {
	return map[string]map[string][]string{
		"computeServer": {
			"computeKeypair": []string{"key_pair", "self_link"},
		},
		"computeVolumeAttach": {
			"computeServer": []string{"server_id", "self_link"},
			"computeVolume": []string{"volume_id", "self_link"},
		},
		"dnsRecordSet": {
			"dnsZone": []string{"zone_id", "self_link"},
		},
		"networkPublicIP": {
			"networkInternetGateway": []string{"internet_gw_id", "self_link"},
		},
		"networkStaticRoute": {
			"networkInternetGateway": []string{"internet_gw_id", "self_link"},
		},
		"networkSubnet": {
			"networkNetwork": []string{"network_id", "self_link"},
		},
		"storageVolume": {
			"storageVirtualStorage": []string{"virtual_storage_id", "self_link"},
		},
	}
}
