---
page_title: "pingdirectory_otp_delivery_mechanism Resource - terraform-provider-pingdirectory"
subcategory: "Otp Delivery Mechanism"
description: |-
  Manages a Otp Delivery Mechanism.
---

# pingdirectory_otp_delivery_mechanism (Resource)

Manages a Otp Delivery Mechanism.

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

resource "pingdirectory_otp_delivery_mechanism" "myOtpDeliveryMechanism" {
  id             = "MyOtpDeliveryMechanism"
  type           = "email"
  sender_address = "sender@example.com"
  enabled        = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `enabled` (Boolean) Indicates whether this OTP Delivery Mechanism is enabled for use in the server.
- `id` (String) Name of this object.
- `type` (String) The type of OTP Delivery Mechanism resource. Options are ['twilio', 'email', 'third-party']

### Optional

- `description` (String) A description for this OTP Delivery Mechanism
- `email_address_attribute_type` (String) The name or OID of the attribute that holds the email address to which the message should be sent.
- `email_address_json_field` (String) The name of the JSON field whose value is the email address to which the message should be sent. The email address must be contained in a top-level field whose value is a single string.
- `email_address_json_object_filter` (String) A JSON object filter that may be used to identify which email address value to use when sending the message.
- `extension_argument` (Set of String) The set of arguments used to customize the behavior for the Third Party OTP Delivery Mechanism. Each configuration property should be given in the form 'name=value'.
- `extension_class` (String) The fully-qualified name of the Java class providing the logic for the Third Party OTP Delivery Mechanism.
- `http_proxy_external_server` (String) A reference to an HTTP proxy server that should be used for requests sent to the Twilio service.
- `message_subject` (String) The subject to use for the e-mail message.
- `message_text_after_otp` (String) Any text that should appear in the message after the one-time password value.
- `message_text_before_otp` (String) Any text that should appear in the message before the one-time password value.
- `phone_number_attribute_type` (String) The name or OID of the attribute in the user's entry that holds the phone number to which the message should be sent.
- `phone_number_json_field` (String) The name of the JSON field whose value is the phone number to which the message should be sent. The phone number must be contained in a top-level field whose value is a single string.
- `phone_number_json_object_filter` (String) A JSON object filter that may be used to identify which phone number value to use when sending the message.
- `sender_address` (String) The e-mail address to use as the sender for the one-time password.
- `sender_phone_number` (Set of String) The outgoing phone number to use for the messages. Values must be phone numbers you have obtained for use with your Twilio account.
- `twilio_account_sid` (String) The unique identifier assigned to the Twilio account that will be used.
- `twilio_auth_token` (String, Sensitive) The auth token for the Twilio account that will be used.
- `twilio_auth_token_passphrase_provider` (String) The passphrase provider that may be used to obtain the auth token for the Twilio account that will be used.

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
# "otpDeliveryMechanismId" should be the id of the Otp Delivery Mechanism to be imported
terraform import pingdirectory_otp_delivery_mechanism.myOtpDeliveryMechanism otpDeliveryMechanismId
```
