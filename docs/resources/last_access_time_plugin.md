---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_last_access_time_plugin Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Last Access Time Plugin.
---

# pingdirectory_last_access_time_plugin (Resource)

Manages a Last Access Time Plugin.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `description` (String) A description for this Plugin
- `enabled` (Boolean) Indicates whether the plug-in is enabled for use.
- `invoke_for_failed_binds` (Boolean) Indicates whether to update the last access time for an entry targeted by a bind operation if the bind is unsuccessful.
- `invoke_for_internal_operations` (Boolean) Indicates whether the plug-in should be invoked for internal operations.
- `max_search_result_entries_to_update` (Number) Specifies the maximum number of entries that should be updated in a search operation. Only search result entries actually returned to the client may have their last access time updated, but because a single search operation may return a very large number of entries, the plugin will only update entries if no more than a specified number of entries are updated.
- `max_update_frequency` (String) Specifies the maximum frequency with which last access time values should be written for an entry. This may help limit the rate of internal write operations processed in the server.
- `operation_type` (Set of String) Specifies the types of operations that should result in access time updates.
- `request_criteria` (String) Specifies a set of request criteria that may be used to indicate whether to apply access time updates for the associated operation.

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

