---
page_title: "pingdirectory_uncached_attribute_criteria Data Source - terraform-provider-pingdirectory"
subcategory: "Uncached Attribute Criteria"
description: |-
  Describes a Uncached Attribute Criteria.
---

# pingdirectory_uncached_attribute_criteria (Data Source)

Describes a Uncached Attribute Criteria.

## Example Usage

```terraform
terraform {
  required_version = ">=1.1"
  required_providers {
    pingdirectory = {
      version = "~> 0.3.0"
      source  = "pingidentity/pingdirectory"
    }
  }
}

provider "pingdirectory" {
  username   = "cn=administrator"
  password   = "2FederateM0re"
  https_host = "https://localhost:1443"
  # Warning: The insecure_trust_all_tls attribute configures the provider to trust any certificate presented by the PingDirectory server.
  # It should not be used in production. If you need to specify trusted CA certificates, use the
  # ca_certificate_pem_files attribute to point to any number of trusted CA certificate files
  # in PEM format. If you do not specify certificates, the host's default root CA set will be used.
  # Example:
  # ca_certificate_pem_files = ["/example/path/to/cacert1.pem", "/example/path/to/cacert2.pem"]
  insecure_trust_all_tls = true
  product_version        = "9.3.0.0"
}

data "pingdirectory_uncached_attribute_criteria" "myUncachedAttributeCriteria" {
  name = "MyUncachedAttributeCriteria"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of this config object.

### Read-Only

- `attribute_type` (Set of String) Specifies the attribute types for attributes that may be written to the uncached-id2entry database.
- `description` (String) A description for this Uncached Attribute Criteria
- `enabled` (Boolean) Indicates whether this Uncached Attribute Criteria is enabled for use in the server.
- `extension_argument` (Set of String) The set of arguments used to customize the behavior for the Third Party Uncached Attribute Criteria. Each configuration property should be given in the form 'name=value'.
- `extension_class` (String) The fully-qualified name of the Java class providing the logic for the Third Party Uncached Attribute Criteria.
- `id` (String) The ID of this resource.
- `min_total_value_size` (String) Specifies the minimum total value size (i.e., the sum of the sizes of all values) that an attribute must have before it will be written into the uncached-id2entry database.
- `min_value_count` (Number) Specifies the minimum number of values that an attribute must have before it will be written into the uncached-id2entry database.
- `script_argument` (Set of String) The set of arguments used to customize the behavior for the Scripted Uncached Attribute Criteria. Each configuration property should be given in the form 'name=value'.
- `script_class` (String) The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Uncached Attribute Criteria.
- `type` (String) The type of Uncached Attribute Criteria resource. Options are ['default', 'groovy-scripted', 'simple', 'third-party']
