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
package cmd

import (
	"log"

	ecl_terraforming "github.com/GoogleCloudPlatform/terraformer/providers/ecl"
	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"

	"github.com/spf13/cobra"
)

func newCmdECLImporter(options ImportOptions) *cobra.Command {
	// runtime.Breakpoint()
	cmd := &cobra.Command{
		Use:   "ecl",
		Short: "Import current State to terraform configuration from ECL",
		Long:  "Import current State to terraform configuration from ECL",
		RunE: func(cmd *cobra.Command, args []string) error {
			originalPathPattern := options.PathPattern
			for _, region := range options.Regions {
				provider := newECLProvider()
				options.PathPattern = originalPathPattern
				options.PathPattern += region + "/"
				log.Println(provider.GetName() + " importing region " + region)
				err := Import(provider, options, []string{region})
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.AddCommand(listCmd(newECLProvider()))
	cmd.PersistentFlags().BoolVarP(&options.Connect, "connect", "c", true, "")
	cmd.PersistentFlags().StringSliceVarP(&options.Resources, "resources", "r", []string{}, "computeServer,computeVolume")
	cmd.PersistentFlags().StringVarP(&options.PathPattern, "path-pattern", "p", DefaultPathPattern, "{output}/{provider}/custom/{service}/")
	cmd.PersistentFlags().StringVarP(&options.PathOutput, "path-output", "o", DefaultPathOutput, "")
	cmd.PersistentFlags().StringVarP(&options.State, "state", "s", DefaultState, "local or bucket")
	cmd.PersistentFlags().StringVarP(&options.Bucket, "bucket", "b", "", "gs://terraform-state")
	cmd.PersistentFlags().StringSliceVarP(&options.Filter, "filter", "f", []string{}, "ecl_compute_instance_v2=id1:id2:id4")
	cmd.PersistentFlags().StringSliceVarP(&options.Regions, "regions", "", []string{}, "RegionOne")

	return cmd
}

func newECLProvider() terraform_utils.ProviderGenerator {
	return &ecl_terraforming.ECLProvider{}
}
