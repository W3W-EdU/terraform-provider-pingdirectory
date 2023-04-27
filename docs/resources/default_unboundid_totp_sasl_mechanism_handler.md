---
page_title: "pingdirectory_default_unboundid_totp_sasl_mechanism_handler Resource - terraform-provider-pingdirectory"
subcategory: "Sasl Mechanism Handler"
description: |-
  Manages a Unboundid Totp Sasl Mechanism Handler.
---

# pingdirectory_default_unboundid_totp_sasl_mechanism_handler (Resource)

Manages a Unboundid Totp Sasl Mechanism Handler.

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
  product_version        = "9.2.0.0"
}

resource "pingdirectory_default_unboundid_totp_sasl_mechanism_handler" "myUnboundidTotpSaslMechanismHandler" {
  id              = "MyUnboundidTotpSaslMechanismHandler"
  identity_mapper = "Exact Match"
  enabled         = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `adjacent_intervals_to_check` (Number) The number of adjacent time intervals (both before and after the current time) that should be checked when performing authentication.
- `description` (String) A description for this SASL Mechanism Handler
- `enabled` (Boolean) Indicates whether the SASL mechanism handler is enabled for use.
- `identity_mapper` (String) The identity mapper that should be used to identify the user(s) targeted in the authentication and/or authorization identities contained in the bind request. This will only be used for "u:"-style identities.
- `prevent_totp_reuse` (Boolean) Indicates whether to prevent clients from re-using TOTP passwords.
- `require_static_password` (Boolean) Indicates whether to require a static password (as might be held in the userPassword attribute, or whatever password attribute is defined in the password policy governing the user) in addition to the one-time password.
- `shared_secret_attribute_type` (String) The name or OID of the attribute that will be used to hold the shared secret key used during TOTP processing.
- `time_interval_duration` (String) The duration of the time interval used for TOTP processing.

### Read-Only

- `last_updated` (String) Timestamp of the last Terraform update of this resource.
- `notifications` (Set of String) Notifications returned by the PingDirectory Configuration API.
- `required_actions` (Set of Object) Required actions returned by the PingDirectory Configuration API. (see [below for nested schema](#nestedatt--required_actions))

<a id="nestedatt--required_actions"></a>
### Nested Schema for `required_actions`

Read-Only:

- `property` (String)
- `synopsis` (String)
- `type` (String)

## Import

Import is supported using the following syntax:

```shell
# "unboundidTotpSaslMechanismHandlerId" should be the id of the Unboundid Totp Sasl Mechanism Handler to be imported
terraform import pingdirectory_default_unboundid_totp_sasl_mechanism_handler.myUnboundidTotpSaslMechanismHandler unboundidTotpSaslMechanismHandlerId
```
