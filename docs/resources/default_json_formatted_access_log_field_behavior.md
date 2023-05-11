---
page_title: "pingdirectory_default_json_formatted_access_log_field_behavior Resource - terraform-provider-pingdirectory"
subcategory: "Log Field Behavior"
description: |-
  Manages a Json Formatted Access Log Field Behavior.
---

# pingdirectory_default_json_formatted_access_log_field_behavior (Resource)

Manages a Json Formatted Access Log Field Behavior.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `default_behavior` (String) The default behavior that the server should exhibit for fields for which no explicit behavior is defined. If no default behavior is defined, the server will fall back to using the default behavior configured for the syntax used for each log field.
- `description` (String) A description for this Log Field Behavior
- `omit_field` (Set of String) The log fields that should be omitted entirely from log messages. Neither the field name nor value will be included.
- `omit_field_name` (Set of String) The names of any custom fields that should be omitted from log messages. This should generally only be used for fields that are not available through the omit-field property (for example, custom log fields defined in Server SDK extensions).
- `preserve_field` (Set of String) The log fields whose values should be logged with the intended value. The values for these fields will be preserved, although they may be sanitized for parsability or safety purposes (for example, to escape special characters in the value), and values that are too long may be truncated.
- `preserve_field_name` (Set of String) The names of any custom fields whose values should be preserved. This should generally only be used for fields that are not available through the preserve-field property (for example, custom log fields defined in Server SDK extensions).
- `redact_entire_value_field` (Set of String) The log fields whose values should be completely redacted in log messages. The field name will be included, but with a fixed value that does not reflect the actual value for the field.
- `redact_entire_value_field_name` (Set of String) The names of any custom fields whose values should be completely redacted. This should generally only be used for fields that are not available through the redact-entire-value-field property (for example, custom log fields defined in Server SDK extensions).
- `redact_value_components_field` (Set of String) The log fields whose values will include redacted components.
- `redact_value_components_field_name` (Set of String) The names of any custom fields for which to redact components within the value. This should generally only be used for fields that are not available through the redact-value-components-field property (for example, custom log fields defined in Server SDK extensions).
- `tokenize_entire_value_field` (Set of String) The log fields whose values should be completely tokenized in log messages. The field name will be included, but the value will be replaced with a token that does not reveal the actual value, but that is generated from the value.
- `tokenize_entire_value_field_name` (Set of String) The names of any custom fields whose values should be completely tokenized. This should generally only be used for fields that are not available through the tokenize-entire-value-field property (for example, custom log fields defined in Server SDK extensions).
- `tokenize_value_components_field` (Set of String) The log fields whose values will include tokenized components.
- `tokenize_value_components_field_name` (Set of String) The names of any custom fields for which to tokenize components within the value. This should generally only be used for fields that are not available through the tokenize-value-components-field property (for example, custom log fields defined in Server SDK extensions).

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


