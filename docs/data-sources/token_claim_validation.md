---
page_title: "pingdirectory_token_claim_validation Data Source - terraform-provider-pingdirectory"
subcategory: "Token Claim Validation"
description: |-
  Describes a Token Claim Validation.
---

# pingdirectory_token_claim_validation (Data Source)

Describes a Token Claim Validation.

## Example Usage

```terraform
terraform {
  required_version = ">=1.1"
  required_providers {
    pingdirectory = {
      version = "~> 0.3.0"
      source  = "pingidentity/pingdirectory"
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
  product_version        = "9.3.0.0"
}

data "pingdirectory_token_claim_validation" "myTokenClaimValidation" {
  name                    = "MyTokenClaimValidation"
  id_token_validator_name = "MyIdTokenValidator"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id_token_validator_name` (String) Name of the parent ID Token Validator
- `name` (String) Name of this config object.

### Read-Only

- `all_required_value` (Set of String) The set of all values that the claim must have to be considered valid.
- `any_required_value` (Set of String) The set of values that the claim may have to be considered valid.
- `claim_name` (String) The name of the claim to be validated.
- `description` (String) A description for this Token Claim Validation
- `id` (String) The ID of this resource.
- `required_value` (String) Specifies the boolean claim's required value.
- `type` (String) The type of Token Claim Validation resource. Options are ['string-array', 'boolean', 'string']
