---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_http_connection_handler Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a Http Connection Handler.
---

# pingdirectory_http_connection_handler (Resource)

Manages a Http Connection Handler.

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

resource "pingdirectory_http_connection_handler" "http" {
  id                     = "example"
  description            = "Description of http connection handler"
  listen_port            = 2443
  enabled                = true
  http_servlet_extension = ["Available or Degraded State", "Available State"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `enabled` (Boolean) Indicates whether the Connection Handler is enabled.
- `id` (String) Name of this object.
- `listen_port` (Number) Specifies the port number on which the HTTP Connection Handler will listen for connections from clients.

### Optional

- `accept_backlog` (Number) Specifies the number of concurrent outstanding connection attempts that the connection handler should allow. The default value should be acceptable in most cases, but it may need to be increased in environments that may attempt to establish large numbers of connections simultaneously.
- `allow_tcp_reuse_address` (Boolean) Indicates whether the server should attempt to reuse socket descriptors. This may be useful in environments with a high rate of connection establishment and termination.
- `correlation_id_request_header` (Set of String) Specifies the set of HTTP request headers that may contain a value to be used as the correlation ID. Example values are "Correlation-Id", "X-Amzn-Trace-Id", and "X-Request-Id".
- `correlation_id_response_header` (String) Specifies the name of the HTTP response header that will contain a correlation ID value. Example values are "Correlation-Id", "X-Amzn-Trace-Id", and "X-Request-Id".
- `description` (String) A description for this Connection Handler
- `enable_multipart_mime_parameters` (Boolean) Determines whether request form parameters submitted in multipart/ form-data (RFC 2388) format should be processed as request parameters.
- `http_operation_log_publisher` (Set of String) Specifies the set of HTTP operation loggers that should be used to log information about requests and responses for operations processed through this HTTP Connection Handler.
- `http_request_header_size` (Number) Specifies the maximum buffer size of an http request including the request uri and all of the request headers.
- `http_servlet_extension` (Set of String) Specifies information about servlets that will be provided via this connection handler.
- `idle_time_limit` (String) Specifies the maximum idle time for a connection. The max idle time is applied when waiting for a new request to be received on a connection, when reading the headers and content of a request, or when writing the headers and content of a response.
- `keep_stats` (Boolean) Indicates whether to enable statistics collection for this connection handler.
- `key_manager_provider` (String) Specifies the key manager provider that will be used to obtain the certificate to present to HTTPS clients.
- `listen_address` (String) Specifies the address on which to listen for connections from HTTP clients. If no value is defined, the server will listen on all addresses on all interfaces.
- `low_resources_connection_threshold` (Number) Specifies the number of connections, which if exceeded, places this handler in a low resource state where a different idle time limit is applied on the connections.
- `low_resources_idle_time_limit` (String) Specifies the maximum idle time for a connection when this handler is in a low resource state as defined by low-resource-connections. The max idle time is applied when waiting for a new request to be received on a connection, when reading the headers and content of a request, or when writing the headers and content of a response.
- `num_request_handlers` (Number) Specifies the number of threads that will be used for accepting connections and reading requests from clients.
- `response_header` (Set of String) Specifies HTTP header fields and values added to response headers for all requests.
- `ssl_cert_nickname` (String) Specifies the nickname (also called the alias) of the certificate that the HTTP Connection Handler should use when performing SSL communication.
- `ssl_cipher_suite` (Set of String) Specifies the names of the SSL cipher suites that are allowed for use in SSL communication. The set of supported cipher suites can be viewed via the ssl context monitor entry.
- `ssl_client_auth_policy` (String) Specifies the policy that the HTTP Connection Handler should use regarding client SSL certificates. In order for a client certificate to be accepted it must be known to the trust-manager-provider associated with this HTTP Connection Handler. Client certificates received by the HTTP Connection Handler are by default used for TLS mutual authentication only, as there is no support for user authentication.
- `ssl_protocol` (Set of String) Specifies the names of the SSL protocols that are allowed for use in SSL communication. The set of supported ssl protocols can be viewed via the ssl context monitor entry.
- `trust_manager_provider` (String) Specifies the trust manager provider that will be used to validate any certificates presented by HTTPS clients.
- `use_correlation_id_header` (Boolean) If enabled, a correlation ID header will be added to outgoing HTTP responses.
- `use_forwarded_headers` (Boolean) Indicates whether to use "Forwarded" and "X-Forwarded-*" request headers to override corresponding HTTP request information available during request processing.
- `use_ssl` (Boolean) Indicates whether the HTTP Connection Handler should use SSL.
- `web_application_extension` (Set of String) Specifies information about web applications that will be provided via this connection handler.

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
# "connectionHandlerName" should be the name of the http connection handler to be imported

terraform import pingdirectory_http_connection_handler connectionHandlerName
```