---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_mirror_virtual_attribute Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Mirror Virtual Attribute.
---

# pingdirectory_mirror_virtual_attribute (Resource)

Manages a Mirror Virtual Attribute.

## Example Usage

```terraform
terraform {
  required_providers {
    pingdirectory = {
      source = "pingidentity/pingdirectory"
    }
  }
}

provider "pingdirectory" {
  username   = "cn=administrator"
  password   = "2FederateM0re"
  https_host = "https://localhost:1443"
  # Warning: The insecure_trust_all_tls attribute configures the provider to trust any certificate presented by the PingDirectory server.
  # It should not be used in production. If you need to specify trusted CA certificates, use the
  # ca_certificate_pem_files attribute to point to any number of trusted CA certificate files
  # in PEM format. If you do not specify certificates, the host's default root CA set will be used.
  # Example:
  # ca_certificate_pem_files = ["/example/path/to/cacert1.pem", "/example/path/to/cacert2.pem"]
  insecure_trust_all_tls = true
}

resource "pingdirectory_mirror_virtual_attribute" "myMirrorVirtualAttribute" {
  id               = "MyMirrorVirtualAttribute"
  source_attribute = "mail"
  enabled          = true
  attribute_type   = "name"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `attribute_type` (String) Specifies the attribute type for the attribute whose values are to be dynamically assigned by the virtual attribute.
- `enabled` (Boolean) Indicates whether the Virtual Attribute is enabled for use.
- `id` (String) Name of this object.
- `source_attribute` (String) Specifies the source attribute containing the values to use for this virtual attribute.

### Optional

- `allow_index_conflicts` (Boolean) Indicates whether the server should allow creating or altering this virtual attribute definition even if it conflicts with one or more indexes defined in the server.
- `base_dn` (Set of String) Specifies the base DNs for the branches containing entries that are eligible to use this virtual attribute.
- `bypass_access_control_for_searches` (Boolean) Indicates whether searches performed by this virtual attribute provider should be exempted from access control restrictions.
- `client_connection_policy` (Set of String) Specifies a set of client connection policies for which this Virtual Attribute should be generated. If this is undefined, then this Virtual Attribute will always be generated. If it is associated with one or more client connection policies, then this Virtual Attribute will be generated only for operations requested by clients assigned to one of those client connection policies.
- `conflict_behavior` (String) Specifies the behavior that the server is to exhibit for entries that already contain one or more real values for the associated attribute.
- `description` (String) A description for this Virtual Attribute
- `filter` (Set of String) Specifies the search filters to be applied against entries to determine if the virtual attribute is to be generated for those entries.
- `group_dn` (Set of String) Specifies the DNs of the groups whose members can be eligible to use this virtual attribute.
- `multiple_virtual_attribute_evaluation_order_index` (Number) Specifies the order in which virtual attribute definitions for the same attribute type will be evaluated when generating values for an entry.
- `multiple_virtual_attribute_merge_behavior` (String) Specifies the behavior that will be exhibited for cases in which multiple virtual attribute definitions apply to the same multivalued attribute type. This will be ignored for single-valued attribute types.
- `require_explicit_request_by_name` (Boolean) Indicates whether attributes of this type must be explicitly included by name in the list of requested attributes. Note that this will only apply to virtual attributes which are associated with an attribute type that is operational. It will be ignored for virtual attributes associated with a non-operational attribute type.
- `source_entry_dn_attribute` (String) Specifies the attribute containing the DN of another entry from which to obtain the source attribute providing the values for this virtual attribute.
- `source_entry_dn_map` (String) Specifies a DN map that will be used to identify the entry from which to obtain the source attribute providing the values for this virtual attribute.

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

## Import

Import is supported using the following syntax:

```shell
# "mirrorVirtualAttributeName" should be the name of the Mirror Virtual Attribute to be imported
terraform import pingdirectory_mirror_virtual_attribute mirrorVirtualAttributeName
```