---
page_title: "pingdirectory_failure_lockout_action Data Source - terraform-provider-pingdirectory"
subcategory: "Failure Lockout Action"
description: |-
  Describes a Failure Lockout Action.
---

# pingdirectory_failure_lockout_action (Data Source)

Describes a Failure Lockout Action.

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

data "pingdirectory_failure_lockout_action" "myFailureLockoutAction" {
  name = "MyFailureLockoutAction"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of this config object.

### Read-Only

- `allow_blocking_delay` (Boolean) Indicates whether to delay the response for authentication attempts even if that delay may block the thread being used to process the attempt.
- `delay` (String) The length of time to delay the bind response for accounts with too many failed authentication attempts.
- `description` (String) A description for this Failure Lockout Action
- `generate_account_status_notification` (Boolean) When the `type` attribute is set to:
  - `delay-bind-response`: Indicates whether to generate an account status notification for cases in which a bind response is delayed because of failure lockout.
  - `no-operation`: Indicates whether to generate an account status notification for cases in which this failure lockout action is invoked for a bind attempt with too many outstanding authentication failures.
- `id` (String) The ID of this resource.
- `type` (String) The type of Failure Lockout Action resource. Options are ['delay-bind-response', 'no-operation', 'lock-account']
