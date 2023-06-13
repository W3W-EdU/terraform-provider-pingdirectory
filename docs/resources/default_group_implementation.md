---
page_title: "pingdirectory_default_group_implementation Resource - terraform-provider-pingdirectory"
subcategory: "Group Implementation"
description: |-
  Manages a Group Implementation.
---

# pingdirectory_default_group_implementation (Resource)

Manages a Group Implementation.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.
- `type` (String) The type of Group Implementation resource. Options are ['static', 'virtual-static', 'dynamic']

### Optional

- `description` (String) A description for this Group Implementation
- `enabled` (Boolean) Indicates whether the Group Implementation is enabled.

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


