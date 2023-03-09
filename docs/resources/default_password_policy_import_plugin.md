---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_default_password_policy_import_plugin Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Password Policy Import Plugin.
---

# pingdirectory_default_password_policy_import_plugin (Resource)

Manages a Password Policy Import Plugin.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `default_auth_password_storage_scheme` (Set of String) Specifies the names of password storage schemes that to be used for encoding passwords contained in attributes with the auth password syntax for entries that do not include the ds-pwp-password-policy-dn attribute specifying which password policy should be used to govern them.
- `default_user_password_storage_scheme` (Set of String) Specifies the names of the password storage schemes to be used for encoding passwords contained in attributes with the user password syntax for entries that do not include the ds-pwp-password-policy-dn attribute specifying which password policy is to be used to govern them.
- `description` (String) A description for this Plugin
- `enabled` (Boolean) Indicates whether the plug-in is enabled for use.
- `invoke_for_internal_operations` (Boolean) Indicates whether the plug-in should be invoked for internal operations.

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

