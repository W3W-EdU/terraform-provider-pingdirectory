---
page_title: "pingdirectory_data_security_auditors Data Source - terraform-provider-pingdirectory"
subcategory: "Data Security Auditor"
description: |-
  Lists Data Security Auditor objects in the server configuration.
---

# pingdirectory_data_security_auditors (Data Source)

Lists Data Security Auditor objects in the server configuration.

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

data "pingdirectory_data_security_auditors" "list" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (String) SCIM filter used when searching the configuration.

### Read-Only

- `id` (String) The ID of this resource.
- `objects` (Set of Object) Data Security Auditor objects found in the configuration (see [below for nested schema](#nestedatt--objects))

<a id="nestedatt--objects"></a>
### Nested Schema for `objects`

Read-Only:

- `id` (String)
- `type` (String)
