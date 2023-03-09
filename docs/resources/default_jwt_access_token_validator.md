---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_default_jwt_access_token_validator Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Jwt Access Token Validator.
---

# pingdirectory_default_jwt_access_token_validator (Resource)

Manages a Jwt Access Token Validator.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `allowed_content_encryption_algorithm` (Set of String) Specifies an allow list of JWT content encryption algorithms that will be accepted by the JWT Access Token Validator.
- `allowed_key_encryption_algorithm` (Set of String) Specifies an allow list of JWT key encryption algorithms that will be accepted by the JWT Access Token Validator. This setting is only used if encryption-key-pair is set.
- `allowed_signing_algorithm` (Set of String) Specifies an allow list of JWT signing algorithms that will be accepted by the JWT Access Token Validator.
- `authorization_server` (String) Specifies the external server that will be used to aid in validating access tokens. In most cases this will be the Authorization Server that minted the token.
- `client_id_claim_name` (String) The name of the token claim that contains the OAuth2 client Id.
- `clock_skew_grace_period` (String) Specifies the amount of clock skew that is tolerated by the JWT Access Token Validator when evaluating whether a token is within its valid time interval. The duration specified by this parameter will be subtracted from the token's not-before (nbf) time and added to the token's expiration (exp) time, if present, to allow for any time difference between the local server's clock and the token issuer's clock.
- `description` (String) A description for this Access Token Validator
- `enabled` (Boolean) Indicates whether this Access Token Validator is enabled for use in Directory Server.
- `encryption_key_pair` (String) The public-private key pair that is used to encrypt the JWT payload. If specified, the JWT Access Token Validator will use the private key to decrypt the JWT payload, and the public key must be exported to the Authorization Server that is issuing access tokens.
- `evaluation_order_index` (Number) When multiple JWT Access Token Validators are defined for a single Directory Server, this property determines the evaluation order for determining the correct validator class for an access token received by the Directory Server. Values of this property must be unique among all JWT Access Token Validators defined within Directory Server but not necessarily contiguous. JWT Access Token Validators with a smaller value will be evaluated first to determine if they are able to validate the access token.
- `identity_mapper` (String) Specifies the name of the Identity Mapper that should be used for associating user entries with Bearer token subject names. The claim name from which to obtain the subject (i.e. the currently logged-in user) may be configured using the subject-claim-name property.
- `jwks_endpoint_path` (String) The relative path to JWKS endpoint from which to retrieve one or more public signing keys that may be used to validate the signature of an incoming JWT access token. This path is relative to the base_url property defined for the validator's external authorization server. If jwks-endpoint-path is specified, the JWT Access Token Validator will not consult locally stored certificates for validating token signatures.
- `scope_claim_name` (String) The name of the token claim that contains the scopes granted by the token.
- `signing_certificate` (Set of String) Specifies the locally stored certificates that may be used to validate the signature of an incoming JWT access token. If this property is specified, the JWT Access Token Validator will not use a JWKS endpoint to retrieve public keys.
- `subject_claim_name` (String) The name of the token claim that contains the subject, i.e. the logged-in user in an access token. This property goes hand-in-hand with the identity-mapper property and tells the Identity Mapper which field to use to look up the user entry on the server.

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

