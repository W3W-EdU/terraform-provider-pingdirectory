---
page_title: "pingdirectory_default_logging_change_subscription_handler Resource - terraform-provider-pingdirectory"
subcategory: "Change Subscription Handler"
description: |-
  Manages a Logging Change Subscription Handler.
---

# pingdirectory_default_logging_change_subscription_handler (Resource)

Manages a Logging Change Subscription Handler.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `change_subscription` (Set of String) The set of change subscriptions for which this change subscription handler should be notified. If no values are provided then it will be notified for all change subscriptions defined in the server.
- `description` (String) A description for this Change Subscription Handler
- `enabled` (Boolean) Indicates whether this change subscription handler is enabled within the server.
- `log_file` (String) Specifies the log file in which the change notification messages will be written.

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


