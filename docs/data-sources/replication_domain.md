---
page_title: "pingdirectory_replication_domain Data Source - terraform-provider-pingdirectory"
subcategory: "Replication Domain"
description: |-
  Describes a Replication Domain.
---

# pingdirectory_replication_domain (Data Source)

Describes a Replication Domain.

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

data "pingdirectory_replication_domain" "myReplicationDomain" {
  name                          = "MyReplicationDomain"
  synchronization_provider_name = "MySynchronizationProvider"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of this config object.
- `synchronization_provider_name` (String) Name of the parent Synchronization Provider

### Read-Only

- `base_dn` (String) Specifies the base DN of the replicated data.
- `dependent_ops_replay_failure_wait_time` (String) Defines how long to wait before retrying certain operations, specifically operations that might have failed because they depend on an operation from a different server that has not yet replicated to this instance.
- `heartbeat_interval` (String) Specifies the heartbeat interval that the Directory Server will use when communicating with Replication Servers.
- `id` (String) The ID of this resource.
- `on_replay_failure_wait_for_dependent_ops_timeout` (String) Defines the maximum time to retry a failed operation. An operation will be retried only if it appears that the failure might be dependent on an earlier operation from a different server that hasn't replicated yet. The frequency of the retry is determined by the dependent-ops-replay-failure-wait-time property.
- `restricted` (Boolean) When set to true, changes are only replicated with server instances that belong to the same replication set.
- `server_id` (Number) Specifies a unique identifier for the Directory Server within the Replication Domain.
- `sync_hist_purge_delay` (String) The time in seconds after which historical information used in replication conflict resolution is purged. The information is removed from entries when they are modified after the purge delay has elapsed.
- `type` (String) The type of Replication Domain resource. Options are ['replication-domain']
- `window_size` (Number) Specifies the maximum number of replication updates the Directory Server can have outstanding from the Replication Server.
