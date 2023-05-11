---
page_title: "pingdirectory_default_fixed_time_log_rotation_policy Resource - terraform-provider-pingdirectory"
subcategory: "Log Rotation Policy"
description: |-
  Manages a Fixed Time Log Rotation Policy.
---

# pingdirectory_default_fixed_time_log_rotation_policy (Resource)

Manages a Fixed Time Log Rotation Policy.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `description` (String) A description for this Log Rotation Policy
- `time_of_day` (Set of String) Specifies the time of day at which log rotation should occur.

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


