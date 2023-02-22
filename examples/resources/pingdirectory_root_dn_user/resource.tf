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

resource "pingdirectory_root_dn_user" "myRootDnUser" {
  id                              = "MyRootDnUser"
  inherit_default_root_privileges = true
  search_result_entry_limit       = 0
  time_limit_seconds              = 0
  look_through_entry_limit        = 0
  idle_time_limit_seconds         = 0
  password_policy                 = "Root Password Policy"
  require_secure_authentication   = false
  require_secure_connections      = false
}