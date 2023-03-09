---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_default_referential_integrity_plugin Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Referential Integrity Plugin.
---

# pingdirectory_default_referential_integrity_plugin (Resource)

Manages a Referential Integrity Plugin.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `attribute_type` (Set of String) Specifies the attribute types for which referential integrity is to be maintained.
- `base_dn` (Set of String) Specifies the base DN that limits the scope within which referential integrity is maintained.
- `description` (String) A description for this Plugin
- `enabled` (Boolean) Indicates whether the plug-in is enabled for use.
- `invoke_for_internal_operations` (Boolean) Indicates whether the plug-in should be invoked for internal operations.
- `log_file` (String) Specifies the log file location where the update records are written when the plug-in is in background-mode processing.
- `plugin_type` (Set of String) Specifies the set of plug-in types for the plug-in, which specifies the times at which the plug-in is invoked.
- `update_interval` (String) Specifies the interval in seconds when referential integrity updates are made.

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

