---
page_title: "pingdirectory_default_uncached_attribute_criteria Resource - terraform-provider-pingdirectory"
subcategory: "Uncached Attribute Criteria"
description: |-
  Manages a Uncached Attribute Criteria.
---

# pingdirectory_default_uncached_attribute_criteria (Resource)

Manages a Uncached Attribute Criteria.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `attribute_type` (Set of String) Specifies the attribute types for attributes that may be written to the uncached-id2entry database.
- `description` (String) A description for this Uncached Attribute Criteria
- `enabled` (Boolean) Indicates whether this Uncached Attribute Criteria is enabled for use in the server.
- `extension_argument` (Set of String) The set of arguments used to customize the behavior for the Third Party Uncached Attribute Criteria. Each configuration property should be given in the form 'name=value'.
- `extension_class` (String) The fully-qualified name of the Java class providing the logic for the Third Party Uncached Attribute Criteria.
- `min_total_value_size` (String) Specifies the minimum total value size (i.e., the sum of the sizes of all values) that an attribute must have before it will be written into the uncached-id2entry database.
- `min_value_count` (Number) Specifies the minimum number of values that an attribute must have before it will be written into the uncached-id2entry database.
- `script_argument` (Set of String) The set of arguments used to customize the behavior for the Scripted Uncached Attribute Criteria. Each configuration property should be given in the form 'name=value'.
- `script_class` (String) The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Uncached Attribute Criteria.
- `type` (String) The type of Uncached Attribute Criteria resource. Options are ['default', 'groovy-scripted', 'simple', 'third-party']

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


