terraform {
  required_providers {
    pingdirectory = {
      source = "pingidentity.com/terraform/pingdirectory"
    }
  }
}

provider "pingdirectory" {
  username = "cn=administrator"
  password = "2FederateM0re"
  ldap_host = "ldap://localhost:1389"
  https_host = "https://localhost:1443"
  default_user_password = "2FederateM0re"
}

resource "pingdirectory_user" "mahomes" {
  uid = "pm"
  sn = "Mahomes"
  given_name = "Patrick"
  mail = "pm@kcchiefs.com"
}

resource "pingdirectory_user" "knight" {
  uid = "hk"
  description = "the knight"
  sn = "Knight"
  given_name = "Hollow"
  mail = "hk@hallownest.com"
}

resource "pingdirectory_location" "drangleic" {
  name = "Drangleic"
  description = "Seek the king"
}

resource "pingdirectory_global_configuration" "global" {
  location = "Docker"
  encrypt_data = true
  sensitive_attribute = ["Delivered One-Time Password", "TOTP Shared Secret"]
  tracked_application = ["Requests by Root Users"]
  result_code_map = "Sun DS Compatible Behavior"
  #result_code_map = ""
}

resource "pingdirectory_blind_trust_manager_provider" "blindtest" {
  name = "Blind Test"
  enabled = true
  include_jvm_default_issuers = true
}

resource "pingdirectory_file_based_trust_manager_provider" "filetest" {
  name = "FileTest"
  enabled = true
  trust_store_file = "config/keystore"
  #trust_store_type = ""
}