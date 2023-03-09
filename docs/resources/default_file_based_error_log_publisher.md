---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingdirectory_default_file_based_error_log_publisher Resource - terraform-provider-pingdirectory"
subcategory: ""
description: |-
  Manages a File Based Error Log Publisher.
---

# pingdirectory_default_file_based_error_log_publisher (Resource)

Manages a File Based Error Log Publisher.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Optional

- `append` (Boolean) Specifies whether to append to existing log files.
- `asynchronous` (Boolean) Indicates whether the File Based Error Log Publisher will publish records asynchronously.
- `auto_flush` (Boolean) Specifies whether to flush the writer after every log record.
- `buffer_size` (String) Specifies the log file buffer size.
- `compression_mechanism` (String) Specifies the type of compression (if any) to use for log files that are written.
- `default_severity` (Set of String) Specifies the default severity levels for the logger.
- `description` (String) A description for this Log Publisher
- `enabled` (Boolean) Indicates whether the Log Publisher is enabled for use.
- `encrypt_log` (Boolean) Indicates whether log files should be encrypted so that their content is not available to unauthorized users.
- `encryption_settings_definition_id` (String) Specifies the ID of the encryption settings definition that should be used to encrypt the data. If this is not provided, the server's preferred encryption settings definition will be used. The "encryption-settings list" command can be used to obtain a list of the encryption settings definitions available in the server.
- `generify_message_strings_when_possible` (Boolean) Indicates whether to use the generified version of the log message string (which may use placeholders like %s for a string or %d for an integer), rather than the version of the message with those placeholders replaced with specific values that would normally be written to the log.
- `include_instance_name` (Boolean) Indicates whether log messages should include the instance name for the Directory Server.
- `include_product_name` (Boolean) Indicates whether log messages should include the product name for the Directory Server.
- `include_startup_id` (Boolean) Indicates whether log messages should include the startup ID for the Directory Server, which is a value assigned to the server instance at startup and may be used to identify when the server has been restarted.
- `include_thread_id` (Boolean) Indicates whether log messages should include the thread ID for the Directory Server in each log message. This ID can be used to correlate log messages from the same thread within a single log as well as generated by the same thread across different types of log files. More information about the thread with a specific ID can be obtained using the cn=JVM Stack Trace,cn=monitor entry.
- `log_file` (String) The file name to use for the log files generated by the File Based Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.
- `log_file_permissions` (String) The UNIX permissions of the log files created by this File Based Error Log Publisher.
- `logging_error_behavior` (String) Specifies the behavior that the server should exhibit if an error occurs during logging processing.
- `override_severity` (Set of String) Specifies the override severity levels for the logger based on the category of the messages.
- `queue_size` (Number) The maximum number of log records that can be stored in the asynchronous queue.
- `retention_policy` (Set of String) The retention policy to use for the File Based Error Log Publisher .
- `rotation_listener` (Set of String) A listener that should be notified whenever a log file is rotated out of service.
- `rotation_policy` (Set of String) The rotation policy to use for the File Based Error Log Publisher .
- `sign_log` (Boolean) Indicates whether the log should be cryptographically signed so that the log content cannot be altered in an undetectable manner.
- `time_interval` (String) Specifies the interval at which to check whether the log files need to be rotated.
- `timestamp_precision` (String) Specifies the smallest time unit to be included in timestamps.

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

