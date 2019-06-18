# Terraformer -- Enterprise Cloud(ECL) provider dedicated page

CLI tool to generate `tf` and `tfstate` files from existing infrastructure
(reverse Terraform).

*   Disclaimer: This is not an official page of Terraformer product.
*   Status: beta - need improve documentations, bugs etc..
*   Created by: Keiichi Hikita

# What is ECL provider (means this branch)

- Additional functionality for Terraformer to export existing ECL tenant as `tf` and `tfstate` files.
- Of course you can use other provider's functionality, because this branch includes original functionality(like GCP, AWS...) as well.

# Original README

- [Original Terraformer README](https://github.com/GoogleCloudPlatform/terraformer#use-with-datadog)

# Table of Contents

- [Capabilities](#capabilities)
- [Installation](#installation)
- [Use With ECL](#use-with-ecl)

## Capabilities

1.  Generate `tf` + `tfstate` files from existing infrastructure for all
    supported objects by resource.
2.  Remote state can be uploaded to a GCS bucket.
3.  Connect between resources with `terraform_remote_state` (local and bucket).
4.  Compatible with terraform 0.12 syntax.
5.  Save `tf` files with custom folder tree pattern.
6.  Import by resource name and type.

Terraformer use terraform providers and built for easy to add new supported resources.
For upgrade resources with new fields you need upgrade only terraform providers.
```
Import current State to terraform configuration from `Enterprise Cloud`

Usage:
   import ecl [flags]
   import ecl [command]

Available Commands:
  list        List supported resources for google provider

Flags:
  -b, --bucket string         gs://terraform-state
  -c, --connect                (default true)
  -f, --filter strings        google_compute_firewall=id1:id2:id4
  -h, --help                  help for google
  -o, --path-output string     (default "generated")
  -p, --path-pattern string   {output}/{provider}/custom/{service}/ (default "{output}/{provider}/{service}/")
      --projects strings
  -r, --resources strings     firewalls,networks
  -s, --state string          local or bucket (default "local")
  -z, --zone string
```

## Installation

### From binary

1. You can download binary of `terraformer` which has `Enterpise Cloud functionality` from [here](https://github.com/keiichi-hikita/terraformer/releases/tag/v0.7.4-ecl) as `terraformer.zip` .

2. Copy your Terraform provider's plugin(s) to folder
    `~/.terraform.d/plugins/darwin_amd64/`, as appropriate.
    
    If you already use terraform provider for ECL, you can copy the provider binary from your working directory.

### From source

1.  Run `git clone <this repository>`
2.  Run `GO111MODULE=on go mod vendor`
3.  Run `go build -v`
4. Copy your Terraform provider's plugin(s) to folder
    `~/.terraform.d/plugins/darwin_amd64/`, as appropriate.
    
    If you already use terraform provider for ECL, you can copy the provider binary from your working directory.


## Use with ECL

Example:

```
terraformer import ecl --resources=computeKeypair,computeServer,networkNetwork,networkSubnet --connect=true --region=jp1
```

**Note:**

**You need to set environment variables which are used by ECL provider's authorization(like OS_AUTH_URL ...) prior to execute above.**

List of supported ECL resources:

*   `computeKeypair`
    * `ecl_compute_keypair_v2`
*   `computeServer`
    * `ecl_compute_instance_v2`
*   `computeVolumeAttach`
    * `ecl_compute_volume_attach_v2`
*   `computeVolume`
    * `ecl_compute_volume_v2`
*   `dnsZone`
    * `ecl_dns_zone_v2`
*   `dnsRecordSet`
    * `ecl_dns_recordset_v2`
*   `networkCommonFunctionGateway`
    * `ecl_network_common_function_gateway_v2`
*   `networkGatewayInterface`
    * `ecl_network_gateway_interface_v2`
*   `networkInternetGateway`
    * `ecl_network_internet_gateway_v2`
*   `networkNetwork`
    * `ecl_network_network_v2`
*   `networkPort`
    * `ecl_network_port_v2`
*   `networkPublicIP`
    * `ecl_network_public_ip_v2`
*   `networkStaticRoute`
    * `ecl_network_static_route_v2`
*   `networkSubnet`
    * `ecl_network_subnet_v2`
*   `storageVirtualStorage`
    * `ecl_storage_virtualstorage_v1`
*   `storageVolume`
    * `ecl_storage_volume_v1`
