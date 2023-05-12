---
page_title: "pingdirectory_deliver_password_reset_token_extended_operation_handler Resource - terraform-provider-pingdirectory"
subcategory: "Extended Operation Handler"
description: |-
  Manages a Deliver Password Reset Token Extended Operation Handler.
---

# pingdirectory_deliver_password_reset_token_extended_operation_handler (Resource)

Manages a Deliver Password Reset Token Extended Operation Handler.

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

resource "pingdirectory_deliver_password_reset_token_extended_operation_handler" "myDeliverPasswordResetTokenExtendedOperationHandler" {
  id                               = "MyDeliverPasswordResetTokenExtendedOperationHandler"
  password_generator               = "Passphrase"
  default_token_delivery_mechanism = ["my_example_delivery_mechanism"]
  enabled                          = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `default_token_delivery_mechanism` (Set of String) The set of delivery mechanisms that may be used to deliver password reset tokens to users for requests that do not specify one or more preferred delivery mechanisms.
- `enabled` (Boolean) Indicates whether the Extended Operation Handler is enabled (that is, whether the types of extended operations are allowed in the server).
- `id` (String) Name of this object.
- `password_generator` (String) The password generator that will be used to create the password reset token values to be delivered to the end user.

### Optional

- `description` (String) A description for this Extended Operation Handler
- `password_reset_token_validity_duration` (String) The maximum length of time that a password reset token should be considered valid.

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
# "deliverPasswordResetTokenExtendedOperationHandlerId" should be the id of the Deliver Password Reset Token Extended Operation Handler to be imported
terraform import pingdirectory_deliver_password_reset_token_extended_operation_handler.myDeliverPasswordResetTokenExtendedOperationHandler deliverPasswordResetTokenExtendedOperationHandlerId
```
