---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_default_unique_attribute_plugin Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Unique Attribute Plugin.
---

# pingdirectory_default_unique_attribute_plugin (Resource)

Manages a Unique Attribute Plugin.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `base_dn` (Set of String) Specifies a base DN within which the attribute must be unique.
- `description` (String) A description for this Plugin
- `enabled` (Boolean) Indicates whether the plug-in is enabled for use.
- `filter` (String) Specifies the search filter to apply to determine if attribute uniqueness is enforced for the matching entries.
- `invoke_for_internal_operations` (Boolean) Indicates whether the plug-in should be invoked for internal operations.
- `multiple_attribute_behavior` (String) The behavior to exhibit if multiple attribute types are specified.
- `plugin_type` (Set of String) Specifies the set of plug-in types for the plug-in, which specifies the times at which the plug-in is invoked.
- `prevent_conflicts_with_soft_deleted_entries` (Boolean) Indicates whether this Unique Attribute Plugin should reject a change that would result in one or more conflicts, even if those conflicts only exist in soft-deleted entries.
- `type` (Set of String) Specifies the type of attributes to check for value uniqueness.

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

