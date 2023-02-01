---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_root_dn Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Root Dn.
---

# pingdirectory_root_dn (Resource)

Manages a Root Dn.

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
}

// This set is approximately the minimum set required for you to be able to run
// 'dsconfig get-root-dn-prop' successfully.  If you remove any of these permissions, 
// you risk loss of access to the RootDN permission object.
resource "pingdirectory_root_dn" "myrootdn" {
  default_root_privilege_name = ["bypass-acl", "config-read", "config-write", "modify-acl", "privilege-change", "use-admin-session"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `default_root_privilege_name` (Set of String) Specifies the names of the privileges that root users will be granted by default.

### Read-Only

- `id` (String) Placeholder name of this object required by Terraform.
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
# This resource is singleton, so the value of "id" doesn't matter - it is just a placeholder

terraform import pingdirectory_root_dn id
```