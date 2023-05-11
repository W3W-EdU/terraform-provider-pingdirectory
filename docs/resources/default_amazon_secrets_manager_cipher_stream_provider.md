---
page_title: "pingdirectory_default_amazon_secrets_manager_cipher_stream_provider Resource - terraform-provider-pingdirectory"
subcategory: "Cipher Stream Provider"
description: |-
  Manages a Amazon Secrets Manager Cipher Stream Provider.
---

# pingdirectory_default_amazon_secrets_manager_cipher_stream_provider (Resource)

Manages a Amazon Secrets Manager Cipher Stream Provider.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `aws_external_server` (String) The external server with information to use when interacting with the AWS Secrets Manager.
- `description` (String) A description for this Cipher Stream Provider
- `enabled` (Boolean) Indicates whether this Cipher Stream Provider is enabled for use in the Directory Server.
- `encryption_metadata_file` (String) The path to a file that will hold metadata about the encryption performed by this Amazon Secrets Manager Cipher Stream Provider.
- `secret_field_name` (String) The name of the JSON field whose value is the passphrase that will be used to generate the encryption key for protecting the contents of the encryption settings database.
- `secret_id` (String) The Amazon Resource Name (ARN) or the user-friendly name of the secret to be retrieved.
- `secret_version_id` (String) The unique identifier for the version of the secret to be retrieved.
- `secret_version_stage` (String) The staging label for the version of the secret to be retrieved.

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


