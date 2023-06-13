package plugin

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &pluginResource{}
	_ resource.ResourceWithConfigure   = &pluginResource{}
	_ resource.ResourceWithImportState = &pluginResource{}
	_ resource.Resource                = &defaultPluginResource{}
	_ resource.ResourceWithConfigure   = &defaultPluginResource{}
	_ resource.ResourceWithImportState = &defaultPluginResource{}
)

// Create a Plugin resource
func NewPluginResource() resource.Resource {
	return &pluginResource{}
}

func NewDefaultPluginResource() resource.Resource {
	return &defaultPluginResource{}
}

// pluginResource is the resource implementation.
type pluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPluginResource is the resource implementation.
type defaultPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *pluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plugin"
}

func (r *defaultPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *pluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type pluginResourceModel struct {
	Id                                                   types.String `tfsdk:"id"`
	LastUpdated                                          types.String `tfsdk:"last_updated"`
	Notifications                                        types.Set    `tfsdk:"notifications"`
	RequiredActions                                      types.Set    `tfsdk:"required_actions"`
	ResourceType                                         types.String `tfsdk:"resource_type"`
	PassThroughAuthenticationHandler                     types.String `tfsdk:"pass_through_authentication_handler"`
	Type                                                 types.Set    `tfsdk:"type"`
	MultipleAttributeBehavior                            types.String `tfsdk:"multiple_attribute_behavior"`
	ScriptClass                                          types.String `tfsdk:"script_class"`
	PreventConflictsWithSoftDeletedEntries               types.Bool   `tfsdk:"prevent_conflicts_with_soft_deleted_entries"`
	UpdateInterval                                       types.String `tfsdk:"update_interval"`
	ScriptArgument                                       types.Set    `tfsdk:"script_argument"`
	Delay                                                types.String `tfsdk:"delay"`
	SourceAttribute                                      types.String `tfsdk:"source_attribute"`
	TargetAttribute                                      types.String `tfsdk:"target_attribute"`
	ValuePattern                                         types.Set    `tfsdk:"value_pattern"`
	MultipleValuePatternBehavior                         types.String `tfsdk:"multiple_value_pattern_behavior"`
	MultiValuedAttributeBehavior                         types.String `tfsdk:"multi_valued_attribute_behavior"`
	TargetAttributeExistsDuringInitialPopulationBehavior types.String `tfsdk:"target_attribute_exists_during_initial_population_behavior"`
	UpdateSourceAttributeBehavior                        types.String `tfsdk:"update_source_attribute_behavior"`
	SourceAttributeRemovalBehavior                       types.String `tfsdk:"source_attribute_removal_behavior"`
	UpdateTargetAttributeBehavior                        types.String `tfsdk:"update_target_attribute_behavior"`
	IncludeBaseDN                                        types.Set    `tfsdk:"include_base_dn"`
	ExcludeBaseDN                                        types.Set    `tfsdk:"exclude_base_dn"`
	IncludeFilter                                        types.Set    `tfsdk:"include_filter"`
	ExcludeFilter                                        types.Set    `tfsdk:"exclude_filter"`
	UpdatedEntryNewlyMatchesCriteriaBehavior             types.String `tfsdk:"updated_entry_newly_matches_criteria_behavior"`
	UpdatedEntryNoLongerMatchesCriteriaBehavior          types.String `tfsdk:"updated_entry_no_longer_matches_criteria_behavior"`
	ContextName                                          types.String `tfsdk:"context_name"`
	AgentxAddress                                        types.String `tfsdk:"agentx_address"`
	AgentxPort                                           types.Int64  `tfsdk:"agentx_port"`
	NumWorkerThreads                                     types.Int64  `tfsdk:"num_worker_threads"`
	SessionTimeout                                       types.String `tfsdk:"session_timeout"`
	ConnectRetryMaxWait                                  types.String `tfsdk:"connect_retry_max_wait"`
	PingInterval                                         types.String `tfsdk:"ping_interval"`
	ExtensionClass                                       types.String `tfsdk:"extension_class"`
	ReferralBaseURL                                      types.Set    `tfsdk:"referral_base_url"`
	SourceDN                                             types.String `tfsdk:"source_dn"`
	TargetDN                                             types.String `tfsdk:"target_dn"`
	EnableAttributeMapping                               types.Bool   `tfsdk:"enable_attribute_mapping"`
	MapAttribute                                         types.Set    `tfsdk:"map_attribute"`
	EnableControlMapping                                 types.Bool   `tfsdk:"enable_control_mapping"`
	AlwaysMapResponses                                   types.Bool   `tfsdk:"always_map_responses"`
	Server                                               types.Set    `tfsdk:"server"`
	ExtensionArgument                                    types.Set    `tfsdk:"extension_argument"`
	DatetimeAttribute                                    types.String `tfsdk:"datetime_attribute"`
	DatetimeJSONField                                    types.String `tfsdk:"datetime_json_field"`
	ServerAccessMode                                     types.String `tfsdk:"server_access_mode"`
	NumMostExpensivePhasesShown                          types.Int64  `tfsdk:"num_most_expensive_phases_shown"`
	DatetimeFormat                                       types.String `tfsdk:"datetime_format"`
	CustomDatetimeFormat                                 types.String `tfsdk:"custom_datetime_format"`
	DnMap                                                types.Set    `tfsdk:"dn_map"`
	BindDNPattern                                        types.String `tfsdk:"bind_dn_pattern"`
	SearchBaseDN                                         types.String `tfsdk:"search_base_dn"`
	SearchFilterPattern                                  types.String `tfsdk:"search_filter_pattern"`
	InitialConnections                                   types.Int64  `tfsdk:"initial_connections"`
	MaxConnections                                       types.Int64  `tfsdk:"max_connections"`
	CustomTimezone                                       types.String `tfsdk:"custom_timezone"`
	ExpirationOffset                                     types.String `tfsdk:"expiration_offset"`
	PurgeBehavior                                        types.String `tfsdk:"purge_behavior"`
	LogInterval                                          types.String `tfsdk:"log_interval"`
	SuppressIfIdle                                       types.Bool   `tfsdk:"suppress_if_idle"`
	HeaderPrefixPerColumn                                types.Bool   `tfsdk:"header_prefix_per_column"`
	EmptyInsteadOfZero                                   types.Bool   `tfsdk:"empty_instead_of_zero"`
	LinesBetweenHeader                                   types.Int64  `tfsdk:"lines_between_header"`
	IncludedLDAPStat                                     types.Set    `tfsdk:"included_ldap_stat"`
	IncludedResourceStat                                 types.Set    `tfsdk:"included_resource_stat"`
	HistogramFormat                                      types.String `tfsdk:"histogram_format"`
	HistogramOpType                                      types.Set    `tfsdk:"histogram_op_type"`
	Scope                                                types.String `tfsdk:"scope"`
	IncludeAttribute                                     types.Set    `tfsdk:"include_attribute"`
	GaugeInfo                                            types.String `tfsdk:"gauge_info"`
	LogFileFormat                                        types.String `tfsdk:"log_file_format"`
	LogFile                                              types.String `tfsdk:"log_file"`
	LogFilePermissions                                   types.String `tfsdk:"log_file_permissions"`
	Append                                               types.Bool   `tfsdk:"append"`
	RotationPolicy                                       types.Set    `tfsdk:"rotation_policy"`
	RotationListener                                     types.Set    `tfsdk:"rotation_listener"`
	RetentionPolicy                                      types.Set    `tfsdk:"retention_policy"`
	LoggingErrorBehavior                                 types.String `tfsdk:"logging_error_behavior"`
	OutputFile                                           types.String `tfsdk:"output_file"`
	PreviousFileExtension                                types.String `tfsdk:"previous_file_extension"`
	ApiURL                                               types.String `tfsdk:"api_url"`
	AuthURL                                              types.String `tfsdk:"auth_url"`
	OAuthClientID                                        types.String `tfsdk:"oauth_client_id"`
	OAuthClientSecret                                    types.String `tfsdk:"oauth_client_secret"`
	OAuthClientSecretPassphraseProvider                  types.String `tfsdk:"oauth_client_secret_passphrase_provider"`
	EnvironmentID                                        types.String `tfsdk:"environment_id"`
	HttpProxyExternalServer                              types.String `tfsdk:"http_proxy_external_server"`
	IncludedLocalEntryBaseDN                             types.Set    `tfsdk:"included_local_entry_base_dn"`
	ConnectionCriteria                                   types.String `tfsdk:"connection_criteria"`
	PollingInterval                                      types.String `tfsdk:"polling_interval"`
	TryLocalBind                                         types.Bool   `tfsdk:"try_local_bind"`
	OverrideLocalPassword                                types.Bool   `tfsdk:"override_local_password"`
	UpdateLocalPassword                                  types.Bool   `tfsdk:"update_local_password"`
	UpdateLocalPasswordDN                                types.String `tfsdk:"update_local_password_dn"`
	AllowLaxPassThroughAuthenticationPasswords           types.Bool   `tfsdk:"allow_lax_pass_through_authentication_passwords"`
	IgnoredPasswordPolicyStateErrorCondition             types.Set    `tfsdk:"ignored_password_policy_state_error_condition"`
	UserMappingLocalAttribute                            types.Set    `tfsdk:"user_mapping_local_attribute"`
	UserMappingRemoteJSONField                           types.Set    `tfsdk:"user_mapping_remote_json_field"`
	AdditionalUserMappingSCIMFilter                      types.String `tfsdk:"additional_user_mapping_scim_filter"`
	InvokeGCDayOfWeek                                    types.Set    `tfsdk:"invoke_gc_day_of_week"`
	InvokeGCTimeUtc                                      types.Set    `tfsdk:"invoke_gc_time_utc"`
	DelayAfterAlert                                      types.String `tfsdk:"delay_after_alert"`
	DelayPostGC                                          types.String `tfsdk:"delay_post_gc"`
	PeerServerPriorityIndex                              types.Int64  `tfsdk:"peer_server_priority_index"`
	PluginType                                           types.Set    `tfsdk:"plugin_type"`
	MaxUpdatesPerSecond                                  types.Int64  `tfsdk:"max_updates_per_second"`
	NumDeleteThreads                                     types.Int64  `tfsdk:"num_delete_threads"`
	AttributeType                                        types.Set    `tfsdk:"attribute_type"`
	Filter                                               types.Set    `tfsdk:"filter"`
	NumThreads                                           types.Int64  `tfsdk:"num_threads"`
	BaseDN                                               types.Set    `tfsdk:"base_dn"`
	LowerBound                                           types.Int64  `tfsdk:"lower_bound"`
	UpperBound                                           types.Int64  `tfsdk:"upper_bound"`
	FilterPrefix                                         types.String `tfsdk:"filter_prefix"`
	FilterSuffix                                         types.String `tfsdk:"filter_suffix"`
	CollectionInterval                                   types.String `tfsdk:"collection_interval"`
	PerApplicationLDAPStats                              types.String `tfsdk:"per_application_ldap_stats"`
	LdapChangelogInfo                                    types.String `tfsdk:"ldap_changelog_info"`
	StatusSummaryInfo                                    types.String `tfsdk:"status_summary_info"`
	LocalDBBackendInfo                                   types.String `tfsdk:"local_db_backend_info"`
	ReplicationInfo                                      types.String `tfsdk:"replication_info"`
	EntryCacheInfo                                       types.String `tfsdk:"entry_cache_info"`
	HostInfo                                             types.Set    `tfsdk:"host_info"`
	IncludedLDAPApplication                              types.Set    `tfsdk:"included_ldap_application"`
	RequestCriteria                                      types.String `tfsdk:"request_criteria"`
	InvokeForInternalOperations                          types.Bool   `tfsdk:"invoke_for_internal_operations"`
	Description                                          types.String `tfsdk:"description"`
	Enabled                                              types.Bool   `tfsdk:"enabled"`
}

type defaultPluginResourceModel struct {
	Id                                                   types.String `tfsdk:"id"`
	LastUpdated                                          types.String `tfsdk:"last_updated"`
	Notifications                                        types.Set    `tfsdk:"notifications"`
	RequiredActions                                      types.Set    `tfsdk:"required_actions"`
	ResourceType                                         types.String `tfsdk:"resource_type"`
	PassThroughAuthenticationHandler                     types.String `tfsdk:"pass_through_authentication_handler"`
	Type                                                 types.Set    `tfsdk:"type"`
	MultipleAttributeBehavior                            types.String `tfsdk:"multiple_attribute_behavior"`
	ScriptClass                                          types.String `tfsdk:"script_class"`
	PreventConflictsWithSoftDeletedEntries               types.Bool   `tfsdk:"prevent_conflicts_with_soft_deleted_entries"`
	ProfileSampleInterval                                types.String `tfsdk:"profile_sample_interval"`
	ExcludeAttribute                                     types.Set    `tfsdk:"exclude_attribute"`
	UpdateInterval                                       types.String `tfsdk:"update_interval"`
	ScriptArgument                                       types.Set    `tfsdk:"script_argument"`
	Delay                                                types.String `tfsdk:"delay"`
	SourceAttribute                                      types.String `tfsdk:"source_attribute"`
	TargetAttribute                                      types.String `tfsdk:"target_attribute"`
	ProfileDirectory                                     types.String `tfsdk:"profile_directory"`
	ValuePattern                                         types.Set    `tfsdk:"value_pattern"`
	MultipleValuePatternBehavior                         types.String `tfsdk:"multiple_value_pattern_behavior"`
	MultiValuedAttributeBehavior                         types.String `tfsdk:"multi_valued_attribute_behavior"`
	TargetAttributeExistsDuringInitialPopulationBehavior types.String `tfsdk:"target_attribute_exists_during_initial_population_behavior"`
	UpdateSourceAttributeBehavior                        types.String `tfsdk:"update_source_attribute_behavior"`
	SourceAttributeRemovalBehavior                       types.String `tfsdk:"source_attribute_removal_behavior"`
	UpdateTargetAttributeBehavior                        types.String `tfsdk:"update_target_attribute_behavior"`
	IncludeBaseDN                                        types.Set    `tfsdk:"include_base_dn"`
	ExcludeBaseDN                                        types.Set    `tfsdk:"exclude_base_dn"`
	IncludeFilter                                        types.Set    `tfsdk:"include_filter"`
	ExcludeFilter                                        types.Set    `tfsdk:"exclude_filter"`
	UpdatedEntryNewlyMatchesCriteriaBehavior             types.String `tfsdk:"updated_entry_newly_matches_criteria_behavior"`
	UpdatedEntryNoLongerMatchesCriteriaBehavior          types.String `tfsdk:"updated_entry_no_longer_matches_criteria_behavior"`
	EnableProfilingOnStartup                             types.Bool   `tfsdk:"enable_profiling_on_startup"`
	ProfileAction                                        types.String `tfsdk:"profile_action"`
	ContextName                                          types.String `tfsdk:"context_name"`
	DefaultUserPasswordStorageScheme                     types.Set    `tfsdk:"default_user_password_storage_scheme"`
	DefaultAuthPasswordStorageScheme                     types.Set    `tfsdk:"default_auth_password_storage_scheme"`
	AgentxAddress                                        types.String `tfsdk:"agentx_address"`
	AgentxPort                                           types.Int64  `tfsdk:"agentx_port"`
	NumWorkerThreads                                     types.Int64  `tfsdk:"num_worker_threads"`
	SessionTimeout                                       types.String `tfsdk:"session_timeout"`
	ConnectRetryMaxWait                                  types.String `tfsdk:"connect_retry_max_wait"`
	PingInterval                                         types.String `tfsdk:"ping_interval"`
	ExtensionClass                                       types.String `tfsdk:"extension_class"`
	ReferralBaseURL                                      types.Set    `tfsdk:"referral_base_url"`
	SourceDN                                             types.String `tfsdk:"source_dn"`
	TargetDN                                             types.String `tfsdk:"target_dn"`
	EnableAttributeMapping                               types.Bool   `tfsdk:"enable_attribute_mapping"`
	MapAttribute                                         types.Set    `tfsdk:"map_attribute"`
	RetainFilesSparselyByAge                             types.Bool   `tfsdk:"retain_files_sparsely_by_age"`
	Sanitize                                             types.Bool   `tfsdk:"sanitize"`
	EnableControlMapping                                 types.Bool   `tfsdk:"enable_control_mapping"`
	AlwaysMapResponses                                   types.Bool   `tfsdk:"always_map_responses"`
	Server                                               types.Set    `tfsdk:"server"`
	ExtensionArgument                                    types.Set    `tfsdk:"extension_argument"`
	EncryptionSettingsDefinitionID                       types.String `tfsdk:"encryption_settings_definition_id"`
	DatetimeAttribute                                    types.String `tfsdk:"datetime_attribute"`
	DatetimeJSONField                                    types.String `tfsdk:"datetime_json_field"`
	ServerAccessMode                                     types.String `tfsdk:"server_access_mode"`
	NumMostExpensivePhasesShown                          types.Int64  `tfsdk:"num_most_expensive_phases_shown"`
	DatetimeFormat                                       types.String `tfsdk:"datetime_format"`
	CustomDatetimeFormat                                 types.String `tfsdk:"custom_datetime_format"`
	DnMap                                                types.Set    `tfsdk:"dn_map"`
	BindDNPattern                                        types.String `tfsdk:"bind_dn_pattern"`
	SearchBaseDN                                         types.String `tfsdk:"search_base_dn"`
	SearchFilterPattern                                  types.String `tfsdk:"search_filter_pattern"`
	InitialConnections                                   types.Int64  `tfsdk:"initial_connections"`
	MaxConnections                                       types.Int64  `tfsdk:"max_connections"`
	CustomTimezone                                       types.String `tfsdk:"custom_timezone"`
	ExpirationOffset                                     types.String `tfsdk:"expiration_offset"`
	PurgeBehavior                                        types.String `tfsdk:"purge_behavior"`
	LogInterval                                          types.String `tfsdk:"log_interval"`
	ChangelogPasswordEncryptionKey                       types.String `tfsdk:"changelog_password_encryption_key"`
	SuppressIfIdle                                       types.Bool   `tfsdk:"suppress_if_idle"`
	HeaderPrefixPerColumn                                types.Bool   `tfsdk:"header_prefix_per_column"`
	EmptyInsteadOfZero                                   types.Bool   `tfsdk:"empty_instead_of_zero"`
	LinesBetweenHeader                                   types.Int64  `tfsdk:"lines_between_header"`
	IncludedLDAPStat                                     types.Set    `tfsdk:"included_ldap_stat"`
	IncludedResourceStat                                 types.Set    `tfsdk:"included_resource_stat"`
	HistogramFormat                                      types.String `tfsdk:"histogram_format"`
	HistogramOpType                                      types.Set    `tfsdk:"histogram_op_type"`
	Scope                                                types.String `tfsdk:"scope"`
	HistogramCategoryBoundary                            types.Set    `tfsdk:"histogram_category_boundary"`
	IncludeAttribute                                     types.Set    `tfsdk:"include_attribute"`
	GaugeInfo                                            types.String `tfsdk:"gauge_info"`
	LogFileFormat                                        types.String `tfsdk:"log_file_format"`
	LogFile                                              types.String `tfsdk:"log_file"`
	LogFilePermissions                                   types.String `tfsdk:"log_file_permissions"`
	Append                                               types.Bool   `tfsdk:"append"`
	RotationPolicy                                       types.Set    `tfsdk:"rotation_policy"`
	RotationListener                                     types.Set    `tfsdk:"rotation_listener"`
	RetentionPolicy                                      types.Set    `tfsdk:"retention_policy"`
	LoggingErrorBehavior                                 types.String `tfsdk:"logging_error_behavior"`
	OutputFile                                           types.String `tfsdk:"output_file"`
	PreviousFileExtension                                types.String `tfsdk:"previous_file_extension"`
	IncludeQueueTime                                     types.Bool   `tfsdk:"include_queue_time"`
	SeparateMonitorEntryPerTrackedApplication            types.Bool   `tfsdk:"separate_monitor_entry_per_tracked_application"`
	ChangelogPasswordEncryptionKeyPassphraseProvider     types.String `tfsdk:"changelog_password_encryption_key_passphrase_provider"`
	ApiURL                                               types.String `tfsdk:"api_url"`
	AuthURL                                              types.String `tfsdk:"auth_url"`
	OAuthClientID                                        types.String `tfsdk:"oauth_client_id"`
	OAuthClientSecret                                    types.String `tfsdk:"oauth_client_secret"`
	OAuthClientSecretPassphraseProvider                  types.String `tfsdk:"oauth_client_secret_passphrase_provider"`
	EnvironmentID                                        types.String `tfsdk:"environment_id"`
	HttpProxyExternalServer                              types.String `tfsdk:"http_proxy_external_server"`
	IncludedLocalEntryBaseDN                             types.Set    `tfsdk:"included_local_entry_base_dn"`
	ConnectionCriteria                                   types.String `tfsdk:"connection_criteria"`
	PollingInterval                                      types.String `tfsdk:"polling_interval"`
	TryLocalBind                                         types.Bool   `tfsdk:"try_local_bind"`
	OverrideLocalPassword                                types.Bool   `tfsdk:"override_local_password"`
	UpdateLocalPassword                                  types.Bool   `tfsdk:"update_local_password"`
	UpdateLocalPasswordDN                                types.String `tfsdk:"update_local_password_dn"`
	AllowLaxPassThroughAuthenticationPasswords           types.Bool   `tfsdk:"allow_lax_pass_through_authentication_passwords"`
	IgnoredPasswordPolicyStateErrorCondition             types.Set    `tfsdk:"ignored_password_policy_state_error_condition"`
	UserMappingLocalAttribute                            types.Set    `tfsdk:"user_mapping_local_attribute"`
	UserMappingRemoteJSONField                           types.Set    `tfsdk:"user_mapping_remote_json_field"`
	AdditionalUserMappingSCIMFilter                      types.String `tfsdk:"additional_user_mapping_scim_filter"`
	InvokeGCDayOfWeek                                    types.Set    `tfsdk:"invoke_gc_day_of_week"`
	InvokeGCTimeUtc                                      types.Set    `tfsdk:"invoke_gc_time_utc"`
	DelayAfterAlert                                      types.String `tfsdk:"delay_after_alert"`
	DelayPostGC                                          types.String `tfsdk:"delay_post_gc"`
	PeerServerPriorityIndex                              types.Int64  `tfsdk:"peer_server_priority_index"`
	PluginType                                           types.Set    `tfsdk:"plugin_type"`
	MaxUpdatesPerSecond                                  types.Int64  `tfsdk:"max_updates_per_second"`
	NumDeleteThreads                                     types.Int64  `tfsdk:"num_delete_threads"`
	AttributeType                                        types.Set    `tfsdk:"attribute_type"`
	Filter                                               types.Set    `tfsdk:"filter"`
	NumThreads                                           types.Int64  `tfsdk:"num_threads"`
	BaseDN                                               types.Set    `tfsdk:"base_dn"`
	LowerBound                                           types.Int64  `tfsdk:"lower_bound"`
	UpperBound                                           types.Int64  `tfsdk:"upper_bound"`
	FilterPrefix                                         types.String `tfsdk:"filter_prefix"`
	FilterSuffix                                         types.String `tfsdk:"filter_suffix"`
	SampleInterval                                       types.String `tfsdk:"sample_interval"`
	CollectionInterval                                   types.String `tfsdk:"collection_interval"`
	LdapInfo                                             types.String `tfsdk:"ldap_info"`
	ServerInfo                                           types.String `tfsdk:"server_info"`
	PerApplicationLDAPStats                              types.String `tfsdk:"per_application_ldap_stats"`
	LdapChangelogInfo                                    types.String `tfsdk:"ldap_changelog_info"`
	StatusSummaryInfo                                    types.String `tfsdk:"status_summary_info"`
	GenerateCollectorFiles                               types.Bool   `tfsdk:"generate_collector_files"`
	LocalDBBackendInfo                                   types.String `tfsdk:"local_db_backend_info"`
	ReplicationInfo                                      types.String `tfsdk:"replication_info"`
	EntryCacheInfo                                       types.String `tfsdk:"entry_cache_info"`
	HostInfo                                             types.Set    `tfsdk:"host_info"`
	IncludedLDAPApplication                              types.Set    `tfsdk:"included_ldap_application"`
	MaxUpdateFrequency                                   types.String `tfsdk:"max_update_frequency"`
	OperationType                                        types.Set    `tfsdk:"operation_type"`
	InvokeForFailedBinds                                 types.Bool   `tfsdk:"invoke_for_failed_binds"`
	MaxSearchResultEntriesToUpdate                       types.Int64  `tfsdk:"max_search_result_entries_to_update"`
	RequestCriteria                                      types.String `tfsdk:"request_criteria"`
	InvokeForInternalOperations                          types.Bool   `tfsdk:"invoke_for_internal_operations"`
	Description                                          types.String `tfsdk:"description"`
	Enabled                                              types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *pluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	pluginSchema(ctx, req, resp, false)
}

func (r *defaultPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	pluginSchema(ctx, req, resp, true)
}

func pluginSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Plugin.",
		Attributes: map[string]schema.Attribute{
			"resource_type": schema.StringAttribute{
				Description: "The type of Plugin resource. Options are ['last-access-time', 'stats-collector', 'internal-search-rate', 'modifiable-password-policy-state', 'seven-bit-clean', 'clean-up-expired-pingfederate-persistent-access-grants', 'periodic-gc', 'ping-one-pass-through-authentication', 'changelog-password-encryption', 'processing-time-histogram', 'search-shutdown', 'periodic-stats-logger', 'purge-expired-data', 'change-subscription-notification', 'sub-operation-timing', 'third-party', 'encrypt-attribute-values', 'pass-through-authentication', 'dn-mapper', 'monitor-history', 'referral-on-update', 'simple-to-external-bind', 'custom', 'snmp-subagent', 'password-policy-import', 'profiler', 'clean-up-inactive-pingfederate-persistent-sessions', 'composed-attribute', 'ldap-result-code-tracker', 'attribute-mapper', 'delay', 'clean-up-expired-pingfederate-persistent-sessions', 'groovy-scripted', 'last-mod', 'pluggable-pass-through-authentication', 'referential-integrity', 'unique-attribute']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"internal-search-rate", "modifiable-password-policy-state", "seven-bit-clean", "clean-up-expired-pingfederate-persistent-access-grants", "periodic-gc", "ping-one-pass-through-authentication", "search-shutdown", "periodic-stats-logger", "purge-expired-data", "sub-operation-timing", "third-party", "pass-through-authentication", "dn-mapper", "referral-on-update", "simple-to-external-bind", "snmp-subagent", "clean-up-inactive-pingfederate-persistent-sessions", "composed-attribute", "attribute-mapper", "delay", "clean-up-expired-pingfederate-persistent-sessions", "groovy-scripted", "pluggable-pass-through-authentication", "referential-integrity", "unique-attribute"}...),
				},
			},
			"pass_through_authentication_handler": schema.StringAttribute{
				Description: "The component used to manage authentication with the external authentication service.",
				Optional:    true,
			},
			"type": schema.SetAttribute{
				Description: "Specifies the type of attributes to check for value uniqueness.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"multiple_attribute_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit if multiple attribute types are specified.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Plugin.",
				Optional:    true,
			},
			"prevent_conflicts_with_soft_deleted_entries": schema.BoolAttribute{
				Description: "Indicates whether this Unique Attribute Plugin should reject a change that would result in one or more conflicts, even if those conflicts only exist in soft-deleted entries.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"update_interval": schema.StringAttribute{
				Description: "Specifies the interval in seconds when referential integrity updates are made.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Plugin. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"delay": schema.StringAttribute{
				Description: "The delay to inject for operations matching the associated criteria.",
				Optional:    true,
			},
			"source_attribute": schema.StringAttribute{
				Description: "Specifies the source attribute type that may appear in client requests which should be remapped to the target attribute. Note that the source attribute type must be defined in the server schema and must not be equal to the target attribute type.",
				Optional:    true,
			},
			"target_attribute": schema.StringAttribute{
				Description: "Specifies the target attribute type to which the source attribute type should be mapped. Note that the target attribute type must be defined in the server schema and must not be equal to the source attribute type.",
				Optional:    true,
			},
			"value_pattern": schema.SetAttribute{
				Description: "Specifies a pattern for constructing the values to use for the target attribute type.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"multiple_value_pattern_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit if the plugin is configured with multiple value patterns.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"multi_valued_attribute_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for source attributes that have multiple values.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"target_attribute_exists_during_initial_population_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit if the target attribute exists when initially populating the entry with composed values (whether during an LDIF import, an add operation, or an invocation of the populate composed attribute values task).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_source_attribute_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify and modify DN operations that update one or more of the source attributes used in any of the value patterns.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_attribute_removal_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify and modify DN operations that update an entry to remove source attributes in such a way that this plugin would no longer generate any composed values for that entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_target_attribute_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify and modify DN operations that attempt to update the set of values for the target attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"include_base_dn": schema.SetAttribute{
				Description: "The set of base DNs below which composed values may be generated.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"exclude_base_dn": schema.SetAttribute{
				Description: "The set of base DNs below which composed values will not be generated.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"include_filter": schema.SetAttribute{
				Description: "The set of search filters that identify entries for which composed values may be generated.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"exclude_filter": schema.SetAttribute{
				Description: "The set of search filters that identify entries for which composed values will not be generated.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_entry_newly_matches_criteria_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify or modify DN operations that update an entry that previously did not satisfy either the base DN or filter criteria, but now do satisfy that criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_entry_no_longer_matches_criteria_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify or modify DN operations that update an entry that previously satisfied the base DN and filter criteria, but now no longer satisfies that criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"context_name": schema.StringAttribute{
				Description: "The SNMP context name for this sub-agent. The context name must not be longer than 30 ASCII characters. Each server in a topology must have a unique SNMP context name.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"agentx_address": schema.StringAttribute{
				Description: "The hostname or IP address of the SNMP master agent.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"agentx_port": schema.Int64Attribute{
				Description: "The port number on which the SNMP master agent will be contacted.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"num_worker_threads": schema.Int64Attribute{
				Description: "The number of worker threads to use to handle SNMP requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"session_timeout": schema.StringAttribute{
				Description: "Specifies the maximum amount of time to wait for a session to the master agent to be established.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connect_retry_max_wait": schema.StringAttribute{
				Description: "The maximum amount of time to wait between attempts to establish a connection to the master agent.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ping_interval": schema.StringAttribute{
				Description: "The amount of time between consecutive pings sent by the sub-agent on its connection to the master agent. A value of zero disables the sending of pings by the sub-agent.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Plugin.",
				Optional:    true,
			},
			"referral_base_url": schema.SetAttribute{
				Description: "Specifies the base URL to use for the referrals generated by this plugin. It should include only the scheme, address, and port to use to communicate with the target server (e.g., \"ldap://server.example.com:389/\").",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"source_dn": schema.StringAttribute{
				Description: "Specifies the source DN that may appear in client requests which should be remapped to the target DN. Note that the source DN must not be equal to the target DN.",
				Optional:    true,
			},
			"target_dn": schema.StringAttribute{
				Description: "Specifies the DN to which the source DN should be mapped. Note that the target DN must not be equal to the source DN.",
				Optional:    true,
			},
			"enable_attribute_mapping": schema.BoolAttribute{
				Description: "Indicates whether DN mapping should be applied to the values of attributes with appropriate syntaxes.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"map_attribute": schema.SetAttribute{
				Description: "Specifies a set of specific attributes for which DN mapping should be applied. This will only be applicable if the enable-attribute-mapping property has a value of \"true\". Any attributes listed must be defined in the server schema with either the distinguished name syntax or the name and optional UID syntax.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_control_mapping": schema.BoolAttribute{
				Description: "Indicates whether DN mapping should be applied to DNs that may be present in specific controls. DN mapping will only be applied for control types which are specifically supported by the DN mapper plugin.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"always_map_responses": schema.BoolAttribute{
				Description: "Indicates whether DNs in response messages containing the target DN should always be remapped back to the source DN. If this is \"false\", then mapping will be performed for a response message only if one or more elements of the associated request are mapped. Otherwise, the mapping will be performed for all responses regardless of whether the mapping was applied to the request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"server": schema.SetAttribute{
				Description: "Specifies the LDAP external server(s) to which authentication attempts should be forwarded.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Plugin. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"datetime_attribute": schema.StringAttribute{
				Description: "The LDAP attribute that determines when data should be deleted. This could store the expiration time, or it could store the creation time and the expiration-offset property specifies the duration before data is deleted.",
				Optional:    true,
			},
			"datetime_json_field": schema.StringAttribute{
				Description: "The top-level JSON field within the configured datetime-attribute that determines when data should be deleted. This could store the expiration time, or it could store the creation time and the expiration-offset property specifies the duration before data is deleted.",
				Optional:    true,
			},
			"server_access_mode": schema.StringAttribute{
				Description: "Specifies the manner in which external servers should be used for pass-through authentication attempts if multiple servers are defined.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"num_most_expensive_phases_shown": schema.Int64Attribute{
				Description: "This controls how many of the most expensive phases are included per operation type in the monitor entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"datetime_format": schema.StringAttribute{
				Description: "Specifies the format of the datetime stored within the entry that determines when data should be purged.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_datetime_format": schema.StringAttribute{
				Description: "When the datetime-format property is configured with a value of \"custom\", this specifies the format (using a string compatible with the java.text.SimpleDateFormat class) that will be used to search for expired data.",
				Optional:    true,
			},
			"dn_map": schema.SetAttribute{
				Description: "Specifies one or more DN mappings that may be used to transform bind DNs before attempting to bind to the external servers.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"bind_dn_pattern": schema.StringAttribute{
				Description: "A pattern to use to construct the bind DN for the simple bind request to send to the remote server. This may consist of a combination of static text and attribute values and other directives enclosed in curly braces.  For example, the value \"cn={cn},ou=People,dc=example,dc=com\" indicates that the remote bind DN should be constructed from the text \"cn=\" followed by the value of the local entry's cn attribute followed by the text \"ou=People,dc=example,dc=com\". If an attribute contains the value to use as the bind DN for pass-through authentication, then the pattern may simply be the name of that attribute in curly braces (e.g., if the seeAlso attribute contains the bind DN for the target user, then a bind DN pattern of \"{seeAlso}\" would be appropriate).  Note that a bind DN pattern can be used to construct a bind DN that is not actually a valid LDAP distinguished name. For example, if authentication is being passed through to a Microsoft Active Directory server, then a bind DN pattern could be used to construct a user principal name (UPN) as an alternative to a distinguished name.",
				Optional:    true,
			},
			"search_base_dn": schema.StringAttribute{
				Description: "The base DN to use when searching for the user entry using a filter constructed from the pattern defined in the search-filter-pattern property. If no base DN is specified, the null DN will be used as the search base DN.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"search_filter_pattern": schema.StringAttribute{
				Description: "A pattern to use to construct a filter to use when searching an external server for the entry of the user as whom to bind. For example, \"(mail={uid:ldapFilterEscape}@example.com)\" would construct a search filter to search for a user whose entry in the local server contains a uid attribute whose value appears before \"@example.com\" in the mail attribute in the external server. Note that the \"ldapFilterEscape\" modifier should almost always be used with attributes specified in the pattern.",
				Optional:    true,
			},
			"initial_connections": schema.Int64Attribute{
				Description: "Specifies the initial number of connections to establish to each external server against which authentication may be attempted.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_connections": schema.Int64Attribute{
				Description: "Specifies the maximum number of connections to maintain to each external server against which authentication may be attempted. This value must be greater than or equal to the value for the initial-connections property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"custom_timezone": schema.StringAttribute{
				Description: "Specifies the time zone to use when generating a date string using the configured custom-datetime-format value. The provided value must be accepted by java.util.TimeZone.getTimeZone.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expiration_offset": schema.StringAttribute{
				Description: "The duration to wait after the value specified in datetime-attribute (and optionally datetime-json-field) before purging the data.",
				Optional:    true,
			},
			"purge_behavior": schema.StringAttribute{
				Description: "Specifies whether to delete expired entries or attribute values. By default entries are deleted.",
				Optional:    true,
			},
			"log_interval": schema.StringAttribute{
				Description: "The duration between statistics collection and logging. A new line is logged to the output for each interval. Setting this value too small can have an impact on performance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"suppress_if_idle": schema.BoolAttribute{
				Description: "If the server is idle during the specified interval, then do not log any output if this property is set to true. The server is idle if during the interval, no new connections were established, no operations were processed, and no operations are pending.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"header_prefix_per_column": schema.BoolAttribute{
				Description: "This property controls whether the header prefix, which applies to a group of columns, appears at the start of each column header or only the first column in a group.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"empty_instead_of_zero": schema.BoolAttribute{
				Description: "This property controls whether a value in the output is shown as empty if the value is zero.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"lines_between_header": schema.Int64Attribute{
				Description: "The number of lines to log between logging the header line that summarizes the columns in the table.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"included_ldap_stat": schema.SetAttribute{
				Description: "Specifies the types of statistics related to LDAP connections and operation processing that should be included in the output.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"included_resource_stat": schema.SetAttribute{
				Description: "Specifies whether statistics related to resource utilization such as JVM memory.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"histogram_format": schema.StringAttribute{
				Description: "The format of the data in the processing time histogram.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"histogram_op_type": schema.SetAttribute{
				Description: "Specifies the operation type(s) to use when outputting the response time histogram data. The order of the operations here determines the order of the columns in the output. Use the per-application-ldap-stats setting to further control this.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"scope": schema.StringAttribute{
				Description: "The scope to use for the search.",
				Optional:    true,
			},
			"include_attribute": schema.SetAttribute{
				Description: "The name of an attribute that should be included in the results. This may include any token which is allowed as a requested attribute in search requests, including the name of an attribute, an asterisk (to indicate all user attributes), a plus sign (to indicate all operational attributes), an object class name preceded with an at symbol (to indicate all attributes associated with that object class), an attribute name preceded by a caret (to indicate that attribute should be excluded), or an object class name preceded by a caret and an at symbol (to indicate that all attributes associated with that object class should be excluded).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"gauge_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include for Gauges.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"log_file_format": schema.StringAttribute{
				Description: "Specifies the format to use when logging server statistics.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"log_file": schema.StringAttribute{
				Description: "The file name to use for the log files generated by the Periodic Stats Logger Plugin. The path to the file can be specified either as relative to the server root or as an absolute path.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"log_file_permissions": schema.StringAttribute{
				Description: "The UNIX permissions of the log files created by this Periodic Stats Logger Plugin.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"append": schema.BoolAttribute{
				Description: "Specifies whether to append to existing log files.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"rotation_policy": schema.SetAttribute{
				Description: "The rotation policy to use for the Periodic Stats Logger Plugin .",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"rotation_listener": schema.SetAttribute{
				Description: "A listener that should be notified whenever a log file is rotated out of service.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"retention_policy": schema.SetAttribute{
				Description: "The retention policy to use for the Periodic Stats Logger Plugin .",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"logging_error_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if an error occurs during logging processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"output_file": schema.StringAttribute{
				Description: "The path of an LDIF file that should be created with the results of the search.",
				Optional:    true,
			},
			"previous_file_extension": schema.StringAttribute{
				Description: "An extension that should be appended to the name of an existing output file rather than deleting it. If a file already exists with the full previous file name, then it will be deleted before the current file is renamed to become the previous file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"api_url": schema.StringAttribute{
				Description: "Specifies the API endpoint for the PingOne web service.",
				Optional:    true,
			},
			"auth_url": schema.StringAttribute{
				Description: "Specifies the API endpoint for the PingOne authentication service.",
				Optional:    true,
			},
			"oauth_client_id": schema.StringAttribute{
				Description: "Specifies the OAuth Client ID used to authenticate connections to the PingOne API.",
				Optional:    true,
			},
			"oauth_client_secret": schema.StringAttribute{
				Description: "Specifies the OAuth Client Secret used to authenticate connections to the PingOne API.",
				Optional:    true,
				Sensitive:   true,
			},
			"oauth_client_secret_passphrase_provider": schema.StringAttribute{
				Description: "Specifies a passphrase provider that can be used to obtain the OAuth Client Secret used to authenticate connections to the PingOne API.",
				Optional:    true,
			},
			"environment_id": schema.StringAttribute{
				Description: "Specifies the PingOne Environment that will be associated with this PingOne Pass Through Authentication Plugin.",
				Optional:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "A reference to an HTTP proxy server that should be used for requests sent to the PingOne service.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"included_local_entry_base_dn": schema.SetAttribute{
				Description: "The base DNs for the local users whose authentication attempts may be passed through to the PingOne service.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_criteria": schema.StringAttribute{
				Description: "A reference to connection criteria that will be used to indicate which bind requests should be passed through to the PingOne service.",
				Optional:    true,
			},
			"polling_interval": schema.StringAttribute{
				Description: "This specifies how often the plugin should check for expired data. It also controls the offset of peer servers (see the peer-server-priority-index for more information).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"try_local_bind": schema.BoolAttribute{
				Description: "Indicates whether to attempt the bind in the local server first, or to only send it to the PingOne service.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"override_local_password": schema.BoolAttribute{
				Description: "Indicates whether to attempt the authentication in the PingOne service if the local user entry includes a password. This property will only be used if try-local-bind is true.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"update_local_password": schema.BoolAttribute{
				Description: "Indicates whether to overwrite the user's local password if the local bind fails but the authentication attempt succeeds when attempted in the PingOne service.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"update_local_password_dn": schema.StringAttribute{
				Description: "This is the DN of the user that will be used to overwrite the user's local password if update-local-password is set. The DN put here should be added to 'ignore-changes-by-dn' in the appropriate Sync Source.",
				Optional:    true,
			},
			"allow_lax_pass_through_authentication_passwords": schema.BoolAttribute{
				Description: "Indicates whether to overwrite the user's local password even if the password used to authenticate to the PingOne service would have failed validation if the user attempted to set it directly.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ignored_password_policy_state_error_condition": schema.SetAttribute{
				Description: "A set of password policy state error conditions that should not be enforced when authentication succeeds when attempted in the PingOne service. This option can only be used if try-local-bind is true.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"user_mapping_local_attribute": schema.SetAttribute{
				Description: "The names of the attributes in the local user entry whose values must match the values of the corresponding fields in the PingOne service.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"user_mapping_remote_json_field": schema.SetAttribute{
				Description: "The names of the fields in the PingOne service whose values must match the values of the corresponding attributes in the local user entry, as specified in the user-mapping-local-attribute property.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"additional_user_mapping_scim_filter": schema.StringAttribute{
				Description: "An optional SCIM filter that will be ANDed with the filter created to identify the account in the PingOne service that corresponds to the local entry. Only the \"eq\", \"sw\", \"and\", and \"or\" filter types may be used.",
				Optional:    true,
			},
			"invoke_gc_day_of_week": schema.SetAttribute{
				Description: "Specifies the days of the week which the Periodic GC Plugin should run. If no values are provided, then the plugin will run every day at the specified time.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"invoke_gc_time_utc": schema.SetAttribute{
				Description: "Specifies the times of the day at which garbage collection may be explicitly invoked. The times should be specified in \"HH:MM\" format, with \"HH\" as a two-digit numeric value between 00 and 23 representing the hour of the day, and MM as a two-digit numeric value between 00 and 59 representing the minute of the hour. All times will be interpreted in the UTC time zone.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"delay_after_alert": schema.StringAttribute{
				Description: "Specifies the length of time that the Directory Server should wait after sending the \"force-gc-starting\" administrative alert before actually invoking the garbage collection processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"delay_post_gc": schema.StringAttribute{
				Description: "Specifies the length of time that the Directory Server should wait after successfully completing the garbage collection processing, before removing the \"force-gc-starting\" administrative alert, which marks the server as unavailable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"peer_server_priority_index": schema.Int64Attribute{
				Description: "In a replicated environment, this determines the order in which peer servers should attempt to purge data.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"plugin_type": schema.SetAttribute{
				Description: "Specifies the set of plug-in types for the plug-in, which specifies the times at which the plug-in is invoked.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"max_updates_per_second": schema.Int64Attribute{
				Description: "This setting smooths out the performance impact on the server by throttling the purging to the specified maximum number of updates per second. To avoid a large backlog, this value should be set comfortably above the average rate that expired data is generated. When purge-behavior is set to subtree-delete-entries, then deletion of the entire subtree is considered a single update for the purposes of throttling.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"num_delete_threads": schema.Int64Attribute{
				Description: "The number of threads used to delete expired entries.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"attribute_type": schema.SetAttribute{
				Description: "Specifies the name or OID of an attribute type for which values should be checked to ensure that they are 7-bit clean.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"filter": schema.SetAttribute{
				Description: "A filter that may be used to identify entries that should support the ds-pwp-modifiable-state-json operational attribute.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"num_threads": schema.Int64Attribute{
				Description: "Specifies the number of concurrent threads that should be used to process the search operations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"base_dn": schema.SetAttribute{
				Description: "A base DN that may be used to identify entries that should support the ds-pwp-modifiable-state-json operational attribute.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"lower_bound": schema.Int64Attribute{
				Description: "Specifies the lower bound for the numeric value which will be inserted into the search filter.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"upper_bound": schema.Int64Attribute{
				Description: "Specifies the upper bound for the numeric value which will be inserted into the search filter.",
				Optional:    true,
			},
			"filter_prefix": schema.StringAttribute{
				Description: "Specifies a prefix which will be used in front of the randomly-selected numeric value in all search filters used. If no upper bound is defined, then this should contain the entire filter string.",
				Optional:    true,
			},
			"filter_suffix": schema.StringAttribute{
				Description: "Specifies a suffix which will be used after of the randomly-selected numeric value in all search filters used. If no upper bound is defined, then this should be omitted.",
				Optional:    true,
			},
			"collection_interval": schema.StringAttribute{
				Description: "Some of the calculated statistics, such as the average and maximum queue sizes, can use multiple samples within a log interval. This value controls how often samples are gathered, and setting this value too small can have an adverse impact on performance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"per_application_ldap_stats": schema.StringAttribute{
				Description: "Controls whether per application LDAP statistics are included in the output for selected LDAP operation statistics.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_changelog_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include for the LDAP changelog.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status_summary_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include about the status summary monitor entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"local_db_backend_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include about the Local DB Backends.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"replication_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include about replication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"entry_cache_info": schema.StringAttribute{
				Description: "Specifies the level of detail to include for each entry cache.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host_info": schema.SetAttribute{
				Description: "Specifies the level of detail to include about the host system resource utilization including CPU, memory, disk and network activity.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"included_ldap_application": schema.SetAttribute{
				Description: "If statistics should not be included for all applications, this property names the subset of applications that should be included.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a set of request criteria that may be used to indicate whether to apply access time updates for the associated operation.",
				Optional:    true,
			},
			"invoke_for_internal_operations": schema.BoolAttribute{
				Description: "Indicates whether the plug-in should be invoked for internal operations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Plugin",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the plug-in is enabled for use.",
				Required:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["resource_type"].(schema.StringAttribute)
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"last-access-time", "stats-collector", "internal-search-rate", "modifiable-password-policy-state", "seven-bit-clean", "clean-up-expired-pingfederate-persistent-access-grants", "periodic-gc", "ping-one-pass-through-authentication", "changelog-password-encryption", "processing-time-histogram", "search-shutdown", "periodic-stats-logger", "purge-expired-data", "change-subscription-notification", "sub-operation-timing", "third-party", "encrypt-attribute-values", "pass-through-authentication", "dn-mapper", "monitor-history", "referral-on-update", "simple-to-external-bind", "custom", "snmp-subagent", "password-policy-import", "profiler", "clean-up-inactive-pingfederate-persistent-sessions", "composed-attribute", "ldap-result-code-tracker", "attribute-mapper", "delay", "clean-up-expired-pingfederate-persistent-sessions", "groovy-scripted", "last-mod", "pluggable-pass-through-authentication", "referential-integrity", "unique-attribute"}...),
		}
		schemaDef.Attributes["resource_type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		schemaDef.Attributes["profile_sample_interval"] = schema.StringAttribute{
			Description: "Specifies the sample interval in milliseconds to be used when capturing profiling information in the server.",
			Optional:    true,
		}
		schemaDef.Attributes["exclude_attribute"] = schema.SetAttribute{
			Description: "Specifies the name or OID of an attribute type which may be updated in a modify or modify DN operation without causing the modifiersName and modifyTimestamp values to be updated for that entry.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["profile_directory"] = schema.StringAttribute{
			Description: "Specifies the path to the directory where profile information is to be written. This path may be either an absolute path or a path that is relative to the root of the Directory Server instance.",
			Optional:    true,
		}
		schemaDef.Attributes["enable_profiling_on_startup"] = schema.BoolAttribute{
			Description: "Indicates whether the profiler plug-in is to start collecting data automatically when the Directory Server is started.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["profile_action"] = schema.StringAttribute{
			Description: "Specifies the action that should be taken by the profiler.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["default_user_password_storage_scheme"] = schema.SetAttribute{
			Description: "Specifies the names of the password storage schemes to be used for encoding passwords contained in attributes with the user password syntax for entries that do not include the ds-pwp-password-policy-dn attribute specifying which password policy is to be used to govern them.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["default_auth_password_storage_scheme"] = schema.SetAttribute{
			Description: "Specifies the names of password storage schemes that to be used for encoding passwords contained in attributes with the auth password syntax for entries that do not include the ds-pwp-password-policy-dn attribute specifying which password policy should be used to govern them.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["retain_files_sparsely_by_age"] = schema.BoolAttribute{
			Description: "Retain some older files to give greater perspective on how monitoring information has changed over time.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["sanitize"] = schema.BoolAttribute{
			Description: "Server monitoring data can include a small amount of personally identifiable information in the form of LDAP DNs and search filters. Setting this property to true will redact this information from the monitor files. This should only be used when necessary, as it reduces the information available in the archive and can increase the time to find the source of support issues.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["encryption_settings_definition_id"] = schema.StringAttribute{
			Description: "Specifies the ID of the encryption settings definition that should be used to encrypt the data. If this is not provided, the server's preferred encryption settings definition will be used. The \"encryption-settings list\" command can be used to obtain a list of the encryption settings definitions available in the server.",
			Optional:    true,
		}
		schemaDef.Attributes["changelog_password_encryption_key"] = schema.StringAttribute{
			Description: "A passphrase that may be used to generate the key for encrypting passwords stored in the changelog. The same passphrase also needs to be set (either through the \"changelog-password-decryption-key\" property or the \"changelog-password-decryption-key-passphrase-provider\" property) in the Global Sync Configuration in the Data Sync Server.",
			Optional:    true,
			Sensitive:   true,
		}
		schemaDef.Attributes["histogram_category_boundary"] = schema.SetAttribute{
			Description: "Specifies the boundary values that will be used to separate the processing times into categories. Values should be specified as durations, and all values must be greater than zero.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["include_queue_time"] = schema.BoolAttribute{
			Description: "Indicates whether operation processing times should include the time spent waiting on the work queue. This will only be available if the work queue is configured to monitor the queue time.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["separate_monitor_entry_per_tracked_application"] = schema.BoolAttribute{
			Description: "When enabled, separate monitor entries will be included for each application defined in the Global Configuration's tracked-application property.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["changelog_password_encryption_key_passphrase_provider"] = schema.StringAttribute{
			Description: "A passphrase provider that may be used to obtain the passphrase that will be used to generate the key for encrypting passwords stored in the changelog. The same passphrase also needs to be set (either through the \"changelog-password-decryption-key\" property or the \"changelog-password-decryption-key-passphrase-provider\" property) in the Global Sync Configuration in the Data Sync Server.",
			Optional:    true,
		}
		schemaDef.Attributes["sample_interval"] = schema.StringAttribute{
			Description: "The duration between statistics collections. Setting this value too small can have an impact on performance. This value should be a multiple of collection-interval.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["ldap_info"] = schema.StringAttribute{
			Description: "Specifies the level of detail to include about the LDAP connection handlers.",
			Optional:    true,
		}
		schemaDef.Attributes["server_info"] = schema.StringAttribute{
			Description: "Specifies whether statistics related to resource utilization such as JVM memory and CPU/Network/Disk utilization.",
			Optional:    true,
		}
		schemaDef.Attributes["generate_collector_files"] = schema.BoolAttribute{
			Description: "Indicates whether this plugin should store metric samples on disk for use by the Data Metrics Server. If the Stats Collector Plugin is only being used to collect metrics for one or more StatsD Monitoring Endpoints, then this can be set to false to prevent unnecessary I/O.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["max_update_frequency"] = schema.StringAttribute{
			Description: "Specifies the maximum frequency with which last access time values should be written for an entry. This may help limit the rate of internal write operations processed in the server.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["operation_type"] = schema.SetAttribute{
			Description: "Specifies the types of operations that should result in access time updates.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["invoke_for_failed_binds"] = schema.BoolAttribute{
			Description: "Indicates whether to update the last access time for an entry targeted by a bind operation if the bind is unsuccessful.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["max_search_result_entries_to_update"] = schema.Int64Attribute{
			Description: "Specifies the maximum number of entries that should be updated in a search operation. Only search result entries actually returned to the client may have their last access time updated, but because a single search operation may return a very large number of entries, the plugin will only update entries if no more than a specified number of entries are updated.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		}
		config.SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *pluginResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPluginResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanPlugin(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model defaultPluginResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.AdditionalUserMappingSCIMFilter) && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'additional_user_mapping_scim_filter' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'additional_user_mapping_scim_filter', the 'resource_type' attribute must be one of ['ping-one-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.DatetimeAttribute) && model.ResourceType.ValueString() != "purge-expired-data" {
		resp.Diagnostics.AddError("Attribute 'datetime_attribute' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'datetime_attribute', the 'resource_type' attribute must be one of ['purge-expired-data']")
	}
	if internaltypes.IsDefined(model.ApiURL) && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'api_url' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'api_url', the 'resource_type' attribute must be one of ['ping-one-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.GaugeInfo) && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'gauge_info' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'gauge_info', the 'resource_type' attribute must be one of ['periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.ExcludeBaseDN) && model.ResourceType.ValueString() != "composed-attribute" {
		resp.Diagnostics.AddError("Attribute 'exclude_base_dn' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'exclude_base_dn', the 'resource_type' attribute must be one of ['composed-attribute']")
	}
	if internaltypes.IsDefined(model.HistogramOpType) && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'histogram_op_type' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'histogram_op_type', the 'resource_type' attribute must be one of ['periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.SourceAttributeRemovalBehavior) && model.ResourceType.ValueString() != "composed-attribute" {
		resp.Diagnostics.AddError("Attribute 'source_attribute_removal_behavior' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'source_attribute_removal_behavior', the 'resource_type' attribute must be one of ['composed-attribute']")
	}
	if internaltypes.IsDefined(model.Type) && model.ResourceType.ValueString() != "unique-attribute" {
		resp.Diagnostics.AddError("Attribute 'type' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'type', the 'resource_type' attribute must be one of ['unique-attribute']")
	}
	if internaltypes.IsDefined(model.PollingInterval) && model.ResourceType.ValueString() != "clean-up-expired-pingfederate-persistent-access-grants" && model.ResourceType.ValueString() != "purge-expired-data" && model.ResourceType.ValueString() != "clean-up-expired-pingfederate-persistent-sessions" && model.ResourceType.ValueString() != "clean-up-inactive-pingfederate-persistent-sessions" {
		resp.Diagnostics.AddError("Attribute 'polling_interval' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'polling_interval', the 'resource_type' attribute must be one of ['clean-up-expired-pingfederate-persistent-access-grants', 'purge-expired-data', 'clean-up-expired-pingfederate-persistent-sessions', 'clean-up-inactive-pingfederate-persistent-sessions']")
	}
	if internaltypes.IsDefined(model.EnableControlMapping) && model.ResourceType.ValueString() != "attribute-mapper" && model.ResourceType.ValueString() != "dn-mapper" {
		resp.Diagnostics.AddError("Attribute 'enable_control_mapping' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'enable_control_mapping', the 'resource_type' attribute must be one of ['attribute-mapper', 'dn-mapper']")
	}
	if internaltypes.IsDefined(model.DefaultAuthPasswordStorageScheme) && model.ResourceType.ValueString() != "password-policy-import" {
		resp.Diagnostics.AddError("Attribute 'default_auth_password_storage_scheme' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'default_auth_password_storage_scheme', the 'resource_type' attribute must be one of ['password-policy-import']")
	}
	if internaltypes.IsDefined(model.SampleInterval) && model.ResourceType.ValueString() != "stats-collector" {
		resp.Diagnostics.AddError("Attribute 'sample_interval' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'sample_interval', the 'resource_type' attribute must be one of ['stats-collector']")
	}
	if internaltypes.IsDefined(model.ConnectionCriteria) && model.ResourceType.ValueString() != "pass-through-authentication" && model.ResourceType.ValueString() != "delay" && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" && model.ResourceType.ValueString() != "simple-to-external-bind" && model.ResourceType.ValueString() != "pluggable-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'connection_criteria' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'connection_criteria', the 'resource_type' attribute must be one of ['pass-through-authentication', 'delay', 'ping-one-pass-through-authentication', 'simple-to-external-bind', 'pluggable-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.ChangelogPasswordEncryptionKeyPassphraseProvider) && model.ResourceType.ValueString() != "changelog-password-encryption" {
		resp.Diagnostics.AddError("Attribute 'changelog_password_encryption_key_passphrase_provider' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'changelog_password_encryption_key_passphrase_provider', the 'resource_type' attribute must be one of ['changelog-password-encryption']")
	}
	if internaltypes.IsDefined(model.RequestCriteria) && model.ResourceType.ValueString() != "pass-through-authentication" && model.ResourceType.ValueString() != "last-access-time" && model.ResourceType.ValueString() != "delay" && model.ResourceType.ValueString() != "groovy-scripted" && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" && model.ResourceType.ValueString() != "simple-to-external-bind" && model.ResourceType.ValueString() != "pluggable-pass-through-authentication" && model.ResourceType.ValueString() != "sub-operation-timing" && model.ResourceType.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'request_criteria' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'request_criteria', the 'resource_type' attribute must be one of ['pass-through-authentication', 'last-access-time', 'delay', 'groovy-scripted', 'ping-one-pass-through-authentication', 'simple-to-external-bind', 'pluggable-pass-through-authentication', 'sub-operation-timing', 'third-party']")
	}
	if internaltypes.IsDefined(model.IncludedLocalEntryBaseDN) && model.ResourceType.ValueString() != "pass-through-authentication" && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" && model.ResourceType.ValueString() != "pluggable-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'included_local_entry_base_dn' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'included_local_entry_base_dn', the 'resource_type' attribute must be one of ['pass-through-authentication', 'ping-one-pass-through-authentication', 'pluggable-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.ValuePattern) && model.ResourceType.ValueString() != "composed-attribute" {
		resp.Diagnostics.AddError("Attribute 'value_pattern' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'value_pattern', the 'resource_type' attribute must be one of ['composed-attribute']")
	}
	if internaltypes.IsDefined(model.SessionTimeout) && model.ResourceType.ValueString() != "snmp-subagent" {
		resp.Diagnostics.AddError("Attribute 'session_timeout' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'session_timeout', the 'resource_type' attribute must be one of ['snmp-subagent']")
	}
	if internaltypes.IsDefined(model.IncludeQueueTime) && model.ResourceType.ValueString() != "processing-time-histogram" {
		resp.Diagnostics.AddError("Attribute 'include_queue_time' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'include_queue_time', the 'resource_type' attribute must be one of ['processing-time-histogram']")
	}
	if internaltypes.IsDefined(model.MultipleValuePatternBehavior) && model.ResourceType.ValueString() != "composed-attribute" {
		resp.Diagnostics.AddError("Attribute 'multiple_value_pattern_behavior' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'multiple_value_pattern_behavior', the 'resource_type' attribute must be one of ['composed-attribute']")
	}
	if internaltypes.IsDefined(model.TargetDN) && model.ResourceType.ValueString() != "dn-mapper" {
		resp.Diagnostics.AddError("Attribute 'target_dn' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'target_dn', the 'resource_type' attribute must be one of ['dn-mapper']")
	}
	if internaltypes.IsDefined(model.HistogramFormat) && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'histogram_format' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'histogram_format', the 'resource_type' attribute must be one of ['periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.LogFileFormat) && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'log_file_format' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'log_file_format', the 'resource_type' attribute must be one of ['periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.ContextName) && model.ResourceType.ValueString() != "snmp-subagent" {
		resp.Diagnostics.AddError("Attribute 'context_name' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'context_name', the 'resource_type' attribute must be one of ['snmp-subagent']")
	}
	if internaltypes.IsDefined(model.EnableProfilingOnStartup) && model.ResourceType.ValueString() != "profiler" {
		resp.Diagnostics.AddError("Attribute 'enable_profiling_on_startup' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'enable_profiling_on_startup', the 'resource_type' attribute must be one of ['profiler']")
	}
	if internaltypes.IsDefined(model.PluginType) && model.ResourceType.ValueString() != "encrypt-attribute-values" && model.ResourceType.ValueString() != "pass-through-authentication" && model.ResourceType.ValueString() != "internal-search-rate" && model.ResourceType.ValueString() != "seven-bit-clean" && model.ResourceType.ValueString() != "periodic-gc" && model.ResourceType.ValueString() != "dn-mapper" && model.ResourceType.ValueString() != "referral-on-update" && model.ResourceType.ValueString() != "custom" && model.ResourceType.ValueString() != "changelog-password-encryption" && model.ResourceType.ValueString() != "composed-attribute" && model.ResourceType.ValueString() != "processing-time-histogram" && model.ResourceType.ValueString() != "ldap-result-code-tracker" && model.ResourceType.ValueString() != "attribute-mapper" && model.ResourceType.ValueString() != "delay" && model.ResourceType.ValueString() != "groovy-scripted" && model.ResourceType.ValueString() != "change-subscription-notification" && model.ResourceType.ValueString() != "last-mod" && model.ResourceType.ValueString() != "referential-integrity" && model.ResourceType.ValueString() != "unique-attribute" && model.ResourceType.ValueString() != "sub-operation-timing" && model.ResourceType.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'plugin_type' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'plugin_type', the 'resource_type' attribute must be one of ['encrypt-attribute-values', 'pass-through-authentication', 'internal-search-rate', 'seven-bit-clean', 'periodic-gc', 'dn-mapper', 'referral-on-update', 'custom', 'changelog-password-encryption', 'composed-attribute', 'processing-time-histogram', 'ldap-result-code-tracker', 'attribute-mapper', 'delay', 'groovy-scripted', 'change-subscription-notification', 'last-mod', 'referential-integrity', 'unique-attribute', 'sub-operation-timing', 'third-party']")
	}
	if internaltypes.IsDefined(model.SearchFilterPattern) && model.ResourceType.ValueString() != "pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'search_filter_pattern' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'search_filter_pattern', the 'resource_type' attribute must be one of ['pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.IncludeFilter) && model.ResourceType.ValueString() != "composed-attribute" {
		resp.Diagnostics.AddError("Attribute 'include_filter' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'include_filter', the 'resource_type' attribute must be one of ['composed-attribute']")
	}
	if internaltypes.IsDefined(model.UpdatedEntryNewlyMatchesCriteriaBehavior) && model.ResourceType.ValueString() != "composed-attribute" {
		resp.Diagnostics.AddError("Attribute 'updated_entry_newly_matches_criteria_behavior' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'updated_entry_newly_matches_criteria_behavior', the 'resource_type' attribute must be one of ['composed-attribute']")
	}
	if internaltypes.IsDefined(model.BaseDN) && model.ResourceType.ValueString() != "search-shutdown" && model.ResourceType.ValueString() != "internal-search-rate" && model.ResourceType.ValueString() != "modifiable-password-policy-state" && model.ResourceType.ValueString() != "seven-bit-clean" && model.ResourceType.ValueString() != "clean-up-expired-pingfederate-persistent-access-grants" && model.ResourceType.ValueString() != "purge-expired-data" && model.ResourceType.ValueString() != "referral-on-update" && model.ResourceType.ValueString() != "clean-up-expired-pingfederate-persistent-sessions" && model.ResourceType.ValueString() != "referential-integrity" && model.ResourceType.ValueString() != "unique-attribute" && model.ResourceType.ValueString() != "clean-up-inactive-pingfederate-persistent-sessions" {
		resp.Diagnostics.AddError("Attribute 'base_dn' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'base_dn', the 'resource_type' attribute must be one of ['search-shutdown', 'internal-search-rate', 'modifiable-password-policy-state', 'seven-bit-clean', 'clean-up-expired-pingfederate-persistent-access-grants', 'purge-expired-data', 'referral-on-update', 'clean-up-expired-pingfederate-persistent-sessions', 'referential-integrity', 'unique-attribute', 'clean-up-inactive-pingfederate-persistent-sessions']")
	}
	if internaltypes.IsDefined(model.PurgeBehavior) && model.ResourceType.ValueString() != "purge-expired-data" {
		resp.Diagnostics.AddError("Attribute 'purge_behavior' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'purge_behavior', the 'resource_type' attribute must be one of ['purge-expired-data']")
	}
	if internaltypes.IsDefined(model.EnableAttributeMapping) && model.ResourceType.ValueString() != "dn-mapper" {
		resp.Diagnostics.AddError("Attribute 'enable_attribute_mapping' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'enable_attribute_mapping', the 'resource_type' attribute must be one of ['dn-mapper']")
	}
	if internaltypes.IsDefined(model.NumDeleteThreads) && model.ResourceType.ValueString() != "clean-up-expired-pingfederate-persistent-access-grants" && model.ResourceType.ValueString() != "purge-expired-data" && model.ResourceType.ValueString() != "clean-up-expired-pingfederate-persistent-sessions" && model.ResourceType.ValueString() != "clean-up-inactive-pingfederate-persistent-sessions" {
		resp.Diagnostics.AddError("Attribute 'num_delete_threads' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'num_delete_threads', the 'resource_type' attribute must be one of ['clean-up-expired-pingfederate-persistent-access-grants', 'purge-expired-data', 'clean-up-expired-pingfederate-persistent-sessions', 'clean-up-inactive-pingfederate-persistent-sessions']")
	}
	if internaltypes.IsDefined(model.ReferralBaseURL) && model.ResourceType.ValueString() != "referral-on-update" {
		resp.Diagnostics.AddError("Attribute 'referral_base_url' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'referral_base_url', the 'resource_type' attribute must be one of ['referral-on-update']")
	}
	if internaltypes.IsDefined(model.Delay) && model.ResourceType.ValueString() != "delay" {
		resp.Diagnostics.AddError("Attribute 'delay' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'delay', the 'resource_type' attribute must be one of ['delay']")
	}
	if internaltypes.IsDefined(model.PerApplicationLDAPStats) && model.ResourceType.ValueString() != "stats-collector" && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'per_application_ldap_stats' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'per_application_ldap_stats', the 'resource_type' attribute must be one of ['stats-collector', 'periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.OperationType) && model.ResourceType.ValueString() != "last-access-time" {
		resp.Diagnostics.AddError("Attribute 'operation_type' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'operation_type', the 'resource_type' attribute must be one of ['last-access-time']")
	}
	if internaltypes.IsDefined(model.FilterSuffix) && model.ResourceType.ValueString() != "internal-search-rate" {
		resp.Diagnostics.AddError("Attribute 'filter_suffix' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'filter_suffix', the 'resource_type' attribute must be one of ['internal-search-rate']")
	}
	if internaltypes.IsDefined(model.UpdateInterval) && model.ResourceType.ValueString() != "referential-integrity" {
		resp.Diagnostics.AddError("Attribute 'update_interval' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'update_interval', the 'resource_type' attribute must be one of ['referential-integrity']")
	}
	if internaltypes.IsDefined(model.ExtensionClass) && model.ResourceType.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_class' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'extension_class', the 'resource_type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.MaxSearchResultEntriesToUpdate) && model.ResourceType.ValueString() != "last-access-time" {
		resp.Diagnostics.AddError("Attribute 'max_search_result_entries_to_update' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'max_search_result_entries_to_update', the 'resource_type' attribute must be one of ['last-access-time']")
	}
	if internaltypes.IsDefined(model.PassThroughAuthenticationHandler) && model.ResourceType.ValueString() != "pluggable-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'pass_through_authentication_handler' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'pass_through_authentication_handler', the 'resource_type' attribute must be one of ['pluggable-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.IncludeBaseDN) && model.ResourceType.ValueString() != "composed-attribute" {
		resp.Diagnostics.AddError("Attribute 'include_base_dn' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'include_base_dn', the 'resource_type' attribute must be one of ['composed-attribute']")
	}
	if internaltypes.IsDefined(model.AllowLaxPassThroughAuthenticationPasswords) && model.ResourceType.ValueString() != "pass-through-authentication" && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" && model.ResourceType.ValueString() != "pluggable-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'allow_lax_pass_through_authentication_passwords' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'allow_lax_pass_through_authentication_passwords', the 'resource_type' attribute must be one of ['pass-through-authentication', 'ping-one-pass-through-authentication', 'pluggable-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.IncludeAttribute) && model.ResourceType.ValueString() != "search-shutdown" && model.ResourceType.ValueString() != "last-mod" {
		resp.Diagnostics.AddError("Attribute 'include_attribute' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'include_attribute', the 'resource_type' attribute must be one of ['search-shutdown', 'last-mod']")
	}
	if internaltypes.IsDefined(model.Server) && model.ResourceType.ValueString() != "pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'server' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'server', the 'resource_type' attribute must be one of ['pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.UpdatedEntryNoLongerMatchesCriteriaBehavior) && model.ResourceType.ValueString() != "composed-attribute" {
		resp.Diagnostics.AddError("Attribute 'updated_entry_no_longer_matches_criteria_behavior' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'updated_entry_no_longer_matches_criteria_behavior', the 'resource_type' attribute must be one of ['composed-attribute']")
	}
	if internaltypes.IsDefined(model.LinesBetweenHeader) && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'lines_between_header' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'lines_between_header', the 'resource_type' attribute must be one of ['periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.OAuthClientID) && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'oauth_client_id' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'oauth_client_id', the 'resource_type' attribute must be one of ['ping-one-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.DatetimeJSONField) && model.ResourceType.ValueString() != "purge-expired-data" {
		resp.Diagnostics.AddError("Attribute 'datetime_json_field' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'datetime_json_field', the 'resource_type' attribute must be one of ['purge-expired-data']")
	}
	if internaltypes.IsDefined(model.AgentxPort) && model.ResourceType.ValueString() != "snmp-subagent" {
		resp.Diagnostics.AddError("Attribute 'agentx_port' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'agentx_port', the 'resource_type' attribute must be one of ['snmp-subagent']")
	}
	if internaltypes.IsDefined(model.PreventConflictsWithSoftDeletedEntries) && model.ResourceType.ValueString() != "unique-attribute" {
		resp.Diagnostics.AddError("Attribute 'prevent_conflicts_with_soft_deleted_entries' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'prevent_conflicts_with_soft_deleted_entries', the 'resource_type' attribute must be one of ['unique-attribute']")
	}
	if internaltypes.IsDefined(model.DatetimeFormat) && model.ResourceType.ValueString() != "purge-expired-data" {
		resp.Diagnostics.AddError("Attribute 'datetime_format' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'datetime_format', the 'resource_type' attribute must be one of ['purge-expired-data']")
	}
	if internaltypes.IsDefined(model.UpdateSourceAttributeBehavior) && model.ResourceType.ValueString() != "composed-attribute" {
		resp.Diagnostics.AddError("Attribute 'update_source_attribute_behavior' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'update_source_attribute_behavior', the 'resource_type' attribute must be one of ['composed-attribute']")
	}
	if internaltypes.IsDefined(model.UpdateLocalPasswordDN) && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" && model.ResourceType.ValueString() != "pluggable-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'update_local_password_dn' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'update_local_password_dn', the 'resource_type' attribute must be one of ['ping-one-pass-through-authentication', 'pluggable-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.LowerBound) && model.ResourceType.ValueString() != "internal-search-rate" {
		resp.Diagnostics.AddError("Attribute 'lower_bound' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'lower_bound', the 'resource_type' attribute must be one of ['internal-search-rate']")
	}
	if internaltypes.IsDefined(model.DnMap) && model.ResourceType.ValueString() != "pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'dn_map' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'dn_map', the 'resource_type' attribute must be one of ['pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.PeerServerPriorityIndex) && model.ResourceType.ValueString() != "clean-up-expired-pingfederate-persistent-access-grants" && model.ResourceType.ValueString() != "purge-expired-data" && model.ResourceType.ValueString() != "clean-up-expired-pingfederate-persistent-sessions" && model.ResourceType.ValueString() != "clean-up-inactive-pingfederate-persistent-sessions" {
		resp.Diagnostics.AddError("Attribute 'peer_server_priority_index' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'peer_server_priority_index', the 'resource_type' attribute must be one of ['clean-up-expired-pingfederate-persistent-access-grants', 'purge-expired-data', 'clean-up-expired-pingfederate-persistent-sessions', 'clean-up-inactive-pingfederate-persistent-sessions']")
	}
	if internaltypes.IsDefined(model.EnvironmentID) && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'environment_id' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'environment_id', the 'resource_type' attribute must be one of ['ping-one-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.ExtensionArgument) && model.ResourceType.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_argument' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'extension_argument', the 'resource_type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.ScriptArgument) && model.ResourceType.ValueString() != "groovy-scripted" {
		resp.Diagnostics.AddError("Attribute 'script_argument' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'script_argument', the 'resource_type' attribute must be one of ['groovy-scripted']")
	}
	if internaltypes.IsDefined(model.LocalDBBackendInfo) && model.ResourceType.ValueString() != "stats-collector" && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'local_db_backend_info' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'local_db_backend_info', the 'resource_type' attribute must be one of ['stats-collector', 'periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.LogFile) && model.ResourceType.ValueString() != "periodic-stats-logger" && model.ResourceType.ValueString() != "monitor-history" && model.ResourceType.ValueString() != "referential-integrity" {
		resp.Diagnostics.AddError("Attribute 'log_file' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'log_file', the 'resource_type' attribute must be one of ['periodic-stats-logger', 'monitor-history', 'referential-integrity']")
	}
	if internaltypes.IsDefined(model.RetainFilesSparselyByAge) && model.ResourceType.ValueString() != "monitor-history" {
		resp.Diagnostics.AddError("Attribute 'retain_files_sparsely_by_age' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'retain_files_sparsely_by_age', the 'resource_type' attribute must be one of ['monitor-history']")
	}
	if internaltypes.IsDefined(model.TryLocalBind) && model.ResourceType.ValueString() != "pass-through-authentication" && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" && model.ResourceType.ValueString() != "pluggable-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'try_local_bind' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'try_local_bind', the 'resource_type' attribute must be one of ['pass-through-authentication', 'ping-one-pass-through-authentication', 'pluggable-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.NumWorkerThreads) && model.ResourceType.ValueString() != "snmp-subagent" {
		resp.Diagnostics.AddError("Attribute 'num_worker_threads' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'num_worker_threads', the 'resource_type' attribute must be one of ['snmp-subagent']")
	}
	if internaltypes.IsDefined(model.IncludedLDAPStat) && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'included_ldap_stat' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'included_ldap_stat', the 'resource_type' attribute must be one of ['periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.HeaderPrefixPerColumn) && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'header_prefix_per_column' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'header_prefix_per_column', the 'resource_type' attribute must be one of ['periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.ExcludeAttribute) && model.ResourceType.ValueString() != "last-mod" {
		resp.Diagnostics.AddError("Attribute 'exclude_attribute' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'exclude_attribute', the 'resource_type' attribute must be one of ['last-mod']")
	}
	if internaltypes.IsDefined(model.MultipleAttributeBehavior) && model.ResourceType.ValueString() != "unique-attribute" {
		resp.Diagnostics.AddError("Attribute 'multiple_attribute_behavior' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'multiple_attribute_behavior', the 'resource_type' attribute must be one of ['unique-attribute']")
	}
	if internaltypes.IsDefined(model.TargetAttributeExistsDuringInitialPopulationBehavior) && model.ResourceType.ValueString() != "composed-attribute" {
		resp.Diagnostics.AddError("Attribute 'target_attribute_exists_during_initial_population_behavior' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'target_attribute_exists_during_initial_population_behavior', the 'resource_type' attribute must be one of ['composed-attribute']")
	}
	if internaltypes.IsDefined(model.Filter) && model.ResourceType.ValueString() != "search-shutdown" && model.ResourceType.ValueString() != "modifiable-password-policy-state" && model.ResourceType.ValueString() != "purge-expired-data" && model.ResourceType.ValueString() != "unique-attribute" {
		resp.Diagnostics.AddError("Attribute 'filter' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'filter', the 'resource_type' attribute must be one of ['search-shutdown', 'modifiable-password-policy-state', 'purge-expired-data', 'unique-attribute']")
	}
	if internaltypes.IsDefined(model.EncryptionSettingsDefinitionID) && model.ResourceType.ValueString() != "encrypt-attribute-values" {
		resp.Diagnostics.AddError("Attribute 'encryption_settings_definition_id' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'encryption_settings_definition_id', the 'resource_type' attribute must be one of ['encrypt-attribute-values']")
	}
	if internaltypes.IsDefined(model.ProfileDirectory) && model.ResourceType.ValueString() != "profiler" {
		resp.Diagnostics.AddError("Attribute 'profile_directory' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'profile_directory', the 'resource_type' attribute must be one of ['profiler']")
	}
	if internaltypes.IsDefined(model.OAuthClientSecret) && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'oauth_client_secret' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'oauth_client_secret', the 'resource_type' attribute must be one of ['ping-one-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.ExcludeFilter) && model.ResourceType.ValueString() != "composed-attribute" {
		resp.Diagnostics.AddError("Attribute 'exclude_filter' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'exclude_filter', the 'resource_type' attribute must be one of ['composed-attribute']")
	}
	if internaltypes.IsDefined(model.HistogramCategoryBoundary) && model.ResourceType.ValueString() != "processing-time-histogram" {
		resp.Diagnostics.AddError("Attribute 'histogram_category_boundary' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'histogram_category_boundary', the 'resource_type' attribute must be one of ['processing-time-histogram']")
	}
	if internaltypes.IsDefined(model.PreviousFileExtension) && model.ResourceType.ValueString() != "search-shutdown" {
		resp.Diagnostics.AddError("Attribute 'previous_file_extension' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'previous_file_extension', the 'resource_type' attribute must be one of ['search-shutdown']")
	}
	if internaltypes.IsDefined(model.TargetAttribute) && model.ResourceType.ValueString() != "attribute-mapper" {
		resp.Diagnostics.AddError("Attribute 'target_attribute' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'target_attribute', the 'resource_type' attribute must be one of ['attribute-mapper']")
	}
	if internaltypes.IsDefined(model.Sanitize) && model.ResourceType.ValueString() != "monitor-history" {
		resp.Diagnostics.AddError("Attribute 'sanitize' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'sanitize', the 'resource_type' attribute must be one of ['monitor-history']")
	}
	if internaltypes.IsDefined(model.LogInterval) && model.ResourceType.ValueString() != "periodic-stats-logger" && model.ResourceType.ValueString() != "monitor-history" {
		resp.Diagnostics.AddError("Attribute 'log_interval' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'log_interval', the 'resource_type' attribute must be one of ['periodic-stats-logger', 'monitor-history']")
	}
	if internaltypes.IsDefined(model.InvokeGCDayOfWeek) && model.ResourceType.ValueString() != "periodic-gc" {
		resp.Diagnostics.AddError("Attribute 'invoke_gc_day_of_week' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'invoke_gc_day_of_week', the 'resource_type' attribute must be one of ['periodic-gc']")
	}
	if internaltypes.IsDefined(model.ReplicationInfo) && model.ResourceType.ValueString() != "stats-collector" && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'replication_info' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'replication_info', the 'resource_type' attribute must be one of ['stats-collector', 'periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.AuthURL) && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'auth_url' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'auth_url', the 'resource_type' attribute must be one of ['ping-one-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.ServerAccessMode) && model.ResourceType.ValueString() != "pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'server_access_mode' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'server_access_mode', the 'resource_type' attribute must be one of ['pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.UpdateTargetAttributeBehavior) && model.ResourceType.ValueString() != "composed-attribute" {
		resp.Diagnostics.AddError("Attribute 'update_target_attribute_behavior' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'update_target_attribute_behavior', the 'resource_type' attribute must be one of ['composed-attribute']")
	}
	if internaltypes.IsDefined(model.MaxConnections) && model.ResourceType.ValueString() != "pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'max_connections' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'max_connections', the 'resource_type' attribute must be one of ['pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.SearchBaseDN) && model.ResourceType.ValueString() != "pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'search_base_dn' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'search_base_dn', the 'resource_type' attribute must be one of ['pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.IncludedResourceStat) && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'included_resource_stat' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'included_resource_stat', the 'resource_type' attribute must be one of ['periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.ExpirationOffset) && model.ResourceType.ValueString() != "purge-expired-data" && model.ResourceType.ValueString() != "clean-up-inactive-pingfederate-persistent-sessions" {
		resp.Diagnostics.AddError("Attribute 'expiration_offset' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'expiration_offset', the 'resource_type' attribute must be one of ['purge-expired-data', 'clean-up-inactive-pingfederate-persistent-sessions']")
	}
	if internaltypes.IsDefined(model.CustomTimezone) && model.ResourceType.ValueString() != "purge-expired-data" {
		resp.Diagnostics.AddError("Attribute 'custom_timezone' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'custom_timezone', the 'resource_type' attribute must be one of ['purge-expired-data']")
	}
	if internaltypes.IsDefined(model.SuppressIfIdle) && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'suppress_if_idle' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'suppress_if_idle', the 'resource_type' attribute must be one of ['periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.MultiValuedAttributeBehavior) && model.ResourceType.ValueString() != "composed-attribute" {
		resp.Diagnostics.AddError("Attribute 'multi_valued_attribute_behavior' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'multi_valued_attribute_behavior', the 'resource_type' attribute must be one of ['composed-attribute']")
	}
	if internaltypes.IsDefined(model.SourceAttribute) && model.ResourceType.ValueString() != "attribute-mapper" {
		resp.Diagnostics.AddError("Attribute 'source_attribute' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'source_attribute', the 'resource_type' attribute must be one of ['attribute-mapper']")
	}
	if internaltypes.IsDefined(model.EntryCacheInfo) && model.ResourceType.ValueString() != "stats-collector" && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'entry_cache_info' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'entry_cache_info', the 'resource_type' attribute must be one of ['stats-collector', 'periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.AgentxAddress) && model.ResourceType.ValueString() != "snmp-subagent" {
		resp.Diagnostics.AddError("Attribute 'agentx_address' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'agentx_address', the 'resource_type' attribute must be one of ['snmp-subagent']")
	}
	if internaltypes.IsDefined(model.AttributeType) && model.ResourceType.ValueString() != "encrypt-attribute-values" && model.ResourceType.ValueString() != "seven-bit-clean" && model.ResourceType.ValueString() != "referential-integrity" && model.ResourceType.ValueString() != "composed-attribute" {
		resp.Diagnostics.AddError("Attribute 'attribute_type' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'attribute_type', the 'resource_type' attribute must be one of ['encrypt-attribute-values', 'seven-bit-clean', 'referential-integrity', 'composed-attribute']")
	}
	if internaltypes.IsDefined(model.FilterPrefix) && model.ResourceType.ValueString() != "internal-search-rate" {
		resp.Diagnostics.AddError("Attribute 'filter_prefix' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'filter_prefix', the 'resource_type' attribute must be one of ['internal-search-rate']")
	}
	if internaltypes.IsDefined(model.IncludedLDAPApplication) && model.ResourceType.ValueString() != "stats-collector" && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'included_ldap_application' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'included_ldap_application', the 'resource_type' attribute must be one of ['stats-collector', 'periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.DelayAfterAlert) && model.ResourceType.ValueString() != "periodic-gc" {
		resp.Diagnostics.AddError("Attribute 'delay_after_alert' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'delay_after_alert', the 'resource_type' attribute must be one of ['periodic-gc']")
	}
	if internaltypes.IsDefined(model.UpdateLocalPassword) && model.ResourceType.ValueString() != "pass-through-authentication" && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" && model.ResourceType.ValueString() != "pluggable-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'update_local_password' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'update_local_password', the 'resource_type' attribute must be one of ['pass-through-authentication', 'ping-one-pass-through-authentication', 'pluggable-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.CustomDatetimeFormat) && model.ResourceType.ValueString() != "purge-expired-data" {
		resp.Diagnostics.AddError("Attribute 'custom_datetime_format' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'custom_datetime_format', the 'resource_type' attribute must be one of ['purge-expired-data']")
	}
	if internaltypes.IsDefined(model.NumThreads) && model.ResourceType.ValueString() != "internal-search-rate" {
		resp.Diagnostics.AddError("Attribute 'num_threads' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'num_threads', the 'resource_type' attribute must be one of ['internal-search-rate']")
	}
	if internaltypes.IsDefined(model.HttpProxyExternalServer) && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'http_proxy_external_server' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'http_proxy_external_server', the 'resource_type' attribute must be one of ['ping-one-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.IgnoredPasswordPolicyStateErrorCondition) && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" && model.ResourceType.ValueString() != "pluggable-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'ignored_password_policy_state_error_condition' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'ignored_password_policy_state_error_condition', the 'resource_type' attribute must be one of ['ping-one-pass-through-authentication', 'pluggable-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.RotationListener) && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'rotation_listener' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'rotation_listener', the 'resource_type' attribute must be one of ['periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.MapAttribute) && model.ResourceType.ValueString() != "dn-mapper" {
		resp.Diagnostics.AddError("Attribute 'map_attribute' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'map_attribute', the 'resource_type' attribute must be one of ['dn-mapper']")
	}
	if internaltypes.IsDefined(model.StatusSummaryInfo) && model.ResourceType.ValueString() != "stats-collector" && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'status_summary_info' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'status_summary_info', the 'resource_type' attribute must be one of ['stats-collector', 'periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.SeparateMonitorEntryPerTrackedApplication) && model.ResourceType.ValueString() != "processing-time-histogram" {
		resp.Diagnostics.AddError("Attribute 'separate_monitor_entry_per_tracked_application' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'separate_monitor_entry_per_tracked_application', the 'resource_type' attribute must be one of ['processing-time-histogram']")
	}
	if internaltypes.IsDefined(model.ScriptClass) && model.ResourceType.ValueString() != "groovy-scripted" {
		resp.Diagnostics.AddError("Attribute 'script_class' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'script_class', the 'resource_type' attribute must be one of ['groovy-scripted']")
	}
	if internaltypes.IsDefined(model.MaxUpdateFrequency) && model.ResourceType.ValueString() != "last-access-time" {
		resp.Diagnostics.AddError("Attribute 'max_update_frequency' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'max_update_frequency', the 'resource_type' attribute must be one of ['last-access-time']")
	}
	if internaltypes.IsDefined(model.CollectionInterval) && model.ResourceType.ValueString() != "stats-collector" && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'collection_interval' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'collection_interval', the 'resource_type' attribute must be one of ['stats-collector', 'periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.EmptyInsteadOfZero) && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'empty_instead_of_zero' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'empty_instead_of_zero', the 'resource_type' attribute must be one of ['periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.BindDNPattern) && model.ResourceType.ValueString() != "pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'bind_dn_pattern' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'bind_dn_pattern', the 'resource_type' attribute must be one of ['pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.SourceDN) && model.ResourceType.ValueString() != "dn-mapper" {
		resp.Diagnostics.AddError("Attribute 'source_dn' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'source_dn', the 'resource_type' attribute must be one of ['dn-mapper']")
	}
	if internaltypes.IsDefined(model.InvokeForFailedBinds) && model.ResourceType.ValueString() != "last-access-time" {
		resp.Diagnostics.AddError("Attribute 'invoke_for_failed_binds' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'invoke_for_failed_binds', the 'resource_type' attribute must be one of ['last-access-time']")
	}
	if internaltypes.IsDefined(model.RotationPolicy) && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'rotation_policy' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'rotation_policy', the 'resource_type' attribute must be one of ['periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.LoggingErrorBehavior) && model.ResourceType.ValueString() != "periodic-stats-logger" && model.ResourceType.ValueString() != "monitor-history" {
		resp.Diagnostics.AddError("Attribute 'logging_error_behavior' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'logging_error_behavior', the 'resource_type' attribute must be one of ['periodic-stats-logger', 'monitor-history']")
	}
	if internaltypes.IsDefined(model.PingInterval) && model.ResourceType.ValueString() != "snmp-subagent" {
		resp.Diagnostics.AddError("Attribute 'ping_interval' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'ping_interval', the 'resource_type' attribute must be one of ['snmp-subagent']")
	}
	if internaltypes.IsDefined(model.Description) && model.ResourceType.ValueString() != "last-access-time" && model.ResourceType.ValueString() != "stats-collector" && model.ResourceType.ValueString() != "internal-search-rate" && model.ResourceType.ValueString() != "modifiable-password-policy-state" && model.ResourceType.ValueString() != "seven-bit-clean" && model.ResourceType.ValueString() != "periodic-gc" && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" && model.ResourceType.ValueString() != "changelog-password-encryption" && model.ResourceType.ValueString() != "processing-time-histogram" && model.ResourceType.ValueString() != "search-shutdown" && model.ResourceType.ValueString() != "periodic-stats-logger" && model.ResourceType.ValueString() != "purge-expired-data" && model.ResourceType.ValueString() != "change-subscription-notification" && model.ResourceType.ValueString() != "sub-operation-timing" && model.ResourceType.ValueString() != "third-party" && model.ResourceType.ValueString() != "encrypt-attribute-values" && model.ResourceType.ValueString() != "pass-through-authentication" && model.ResourceType.ValueString() != "dn-mapper" && model.ResourceType.ValueString() != "monitor-history" && model.ResourceType.ValueString() != "referral-on-update" && model.ResourceType.ValueString() != "simple-to-external-bind" && model.ResourceType.ValueString() != "custom" && model.ResourceType.ValueString() != "snmp-subagent" && model.ResourceType.ValueString() != "password-policy-import" && model.ResourceType.ValueString() != "profiler" && model.ResourceType.ValueString() != "composed-attribute" && model.ResourceType.ValueString() != "ldap-result-code-tracker" && model.ResourceType.ValueString() != "attribute-mapper" && model.ResourceType.ValueString() != "delay" && model.ResourceType.ValueString() != "groovy-scripted" && model.ResourceType.ValueString() != "last-mod" && model.ResourceType.ValueString() != "pluggable-pass-through-authentication" && model.ResourceType.ValueString() != "referential-integrity" && model.ResourceType.ValueString() != "unique-attribute" {
		resp.Diagnostics.AddError("Attribute 'description' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'description', the 'resource_type' attribute must be one of ['last-access-time', 'stats-collector', 'internal-search-rate', 'modifiable-password-policy-state', 'seven-bit-clean', 'periodic-gc', 'ping-one-pass-through-authentication', 'changelog-password-encryption', 'processing-time-histogram', 'search-shutdown', 'periodic-stats-logger', 'purge-expired-data', 'change-subscription-notification', 'sub-operation-timing', 'third-party', 'encrypt-attribute-values', 'pass-through-authentication', 'dn-mapper', 'monitor-history', 'referral-on-update', 'simple-to-external-bind', 'custom', 'snmp-subagent', 'password-policy-import', 'profiler', 'composed-attribute', 'ldap-result-code-tracker', 'attribute-mapper', 'delay', 'groovy-scripted', 'last-mod', 'pluggable-pass-through-authentication', 'referential-integrity', 'unique-attribute']")
	}
	if internaltypes.IsDefined(model.HostInfo) && model.ResourceType.ValueString() != "stats-collector" && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'host_info' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'host_info', the 'resource_type' attribute must be one of ['stats-collector', 'periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.DefaultUserPasswordStorageScheme) && model.ResourceType.ValueString() != "password-policy-import" {
		resp.Diagnostics.AddError("Attribute 'default_user_password_storage_scheme' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'default_user_password_storage_scheme', the 'resource_type' attribute must be one of ['password-policy-import']")
	}
	if internaltypes.IsDefined(model.UpperBound) && model.ResourceType.ValueString() != "internal-search-rate" {
		resp.Diagnostics.AddError("Attribute 'upper_bound' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'upper_bound', the 'resource_type' attribute must be one of ['internal-search-rate']")
	}
	if internaltypes.IsDefined(model.MaxUpdatesPerSecond) && model.ResourceType.ValueString() != "clean-up-expired-pingfederate-persistent-access-grants" && model.ResourceType.ValueString() != "purge-expired-data" && model.ResourceType.ValueString() != "clean-up-expired-pingfederate-persistent-sessions" && model.ResourceType.ValueString() != "clean-up-inactive-pingfederate-persistent-sessions" {
		resp.Diagnostics.AddError("Attribute 'max_updates_per_second' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'max_updates_per_second', the 'resource_type' attribute must be one of ['clean-up-expired-pingfederate-persistent-access-grants', 'purge-expired-data', 'clean-up-expired-pingfederate-persistent-sessions', 'clean-up-inactive-pingfederate-persistent-sessions']")
	}
	if internaltypes.IsDefined(model.LdapInfo) && model.ResourceType.ValueString() != "stats-collector" {
		resp.Diagnostics.AddError("Attribute 'ldap_info' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'ldap_info', the 'resource_type' attribute must be one of ['stats-collector']")
	}
	if internaltypes.IsDefined(model.ChangelogPasswordEncryptionKey) && model.ResourceType.ValueString() != "changelog-password-encryption" {
		resp.Diagnostics.AddError("Attribute 'changelog_password_encryption_key' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'changelog_password_encryption_key', the 'resource_type' attribute must be one of ['changelog-password-encryption']")
	}
	if internaltypes.IsDefined(model.LdapChangelogInfo) && model.ResourceType.ValueString() != "stats-collector" && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'ldap_changelog_info' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'ldap_changelog_info', the 'resource_type' attribute must be one of ['stats-collector', 'periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.OAuthClientSecretPassphraseProvider) && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'oauth_client_secret_passphrase_provider' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'oauth_client_secret_passphrase_provider', the 'resource_type' attribute must be one of ['ping-one-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.UserMappingRemoteJSONField) && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'user_mapping_remote_json_field' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'user_mapping_remote_json_field', the 'resource_type' attribute must be one of ['ping-one-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.UserMappingLocalAttribute) && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'user_mapping_local_attribute' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'user_mapping_local_attribute', the 'resource_type' attribute must be one of ['ping-one-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.ConnectRetryMaxWait) && model.ResourceType.ValueString() != "snmp-subagent" {
		resp.Diagnostics.AddError("Attribute 'connect_retry_max_wait' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'connect_retry_max_wait', the 'resource_type' attribute must be one of ['snmp-subagent']")
	}
	if internaltypes.IsDefined(model.Scope) && model.ResourceType.ValueString() != "search-shutdown" {
		resp.Diagnostics.AddError("Attribute 'scope' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'scope', the 'resource_type' attribute must be one of ['search-shutdown']")
	}
	if internaltypes.IsDefined(model.OverrideLocalPassword) && model.ResourceType.ValueString() != "pass-through-authentication" && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" && model.ResourceType.ValueString() != "pluggable-pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'override_local_password' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'override_local_password', the 'resource_type' attribute must be one of ['pass-through-authentication', 'ping-one-pass-through-authentication', 'pluggable-pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.ProfileSampleInterval) && model.ResourceType.ValueString() != "profiler" {
		resp.Diagnostics.AddError("Attribute 'profile_sample_interval' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'profile_sample_interval', the 'resource_type' attribute must be one of ['profiler']")
	}
	if internaltypes.IsDefined(model.GenerateCollectorFiles) && model.ResourceType.ValueString() != "stats-collector" {
		resp.Diagnostics.AddError("Attribute 'generate_collector_files' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'generate_collector_files', the 'resource_type' attribute must be one of ['stats-collector']")
	}
	if internaltypes.IsDefined(model.ProfileAction) && model.ResourceType.ValueString() != "profiler" {
		resp.Diagnostics.AddError("Attribute 'profile_action' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'profile_action', the 'resource_type' attribute must be one of ['profiler']")
	}
	if internaltypes.IsDefined(model.OutputFile) && model.ResourceType.ValueString() != "search-shutdown" {
		resp.Diagnostics.AddError("Attribute 'output_file' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'output_file', the 'resource_type' attribute must be one of ['search-shutdown']")
	}
	if internaltypes.IsDefined(model.RetentionPolicy) && model.ResourceType.ValueString() != "periodic-stats-logger" && model.ResourceType.ValueString() != "monitor-history" {
		resp.Diagnostics.AddError("Attribute 'retention_policy' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'retention_policy', the 'resource_type' attribute must be one of ['periodic-stats-logger', 'monitor-history']")
	}
	if internaltypes.IsDefined(model.InvokeForInternalOperations) && model.ResourceType.ValueString() != "last-access-time" && model.ResourceType.ValueString() != "internal-search-rate" && model.ResourceType.ValueString() != "seven-bit-clean" && model.ResourceType.ValueString() != "periodic-gc" && model.ResourceType.ValueString() != "ping-one-pass-through-authentication" && model.ResourceType.ValueString() != "changelog-password-encryption" && model.ResourceType.ValueString() != "processing-time-histogram" && model.ResourceType.ValueString() != "change-subscription-notification" && model.ResourceType.ValueString() != "sub-operation-timing" && model.ResourceType.ValueString() != "third-party" && model.ResourceType.ValueString() != "encrypt-attribute-values" && model.ResourceType.ValueString() != "pass-through-authentication" && model.ResourceType.ValueString() != "dn-mapper" && model.ResourceType.ValueString() != "referral-on-update" && model.ResourceType.ValueString() != "custom" && model.ResourceType.ValueString() != "snmp-subagent" && model.ResourceType.ValueString() != "password-policy-import" && model.ResourceType.ValueString() != "composed-attribute" && model.ResourceType.ValueString() != "ldap-result-code-tracker" && model.ResourceType.ValueString() != "attribute-mapper" && model.ResourceType.ValueString() != "delay" && model.ResourceType.ValueString() != "groovy-scripted" && model.ResourceType.ValueString() != "last-mod" && model.ResourceType.ValueString() != "pluggable-pass-through-authentication" && model.ResourceType.ValueString() != "referential-integrity" && model.ResourceType.ValueString() != "unique-attribute" {
		resp.Diagnostics.AddError("Attribute 'invoke_for_internal_operations' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'invoke_for_internal_operations', the 'resource_type' attribute must be one of ['last-access-time', 'internal-search-rate', 'seven-bit-clean', 'periodic-gc', 'ping-one-pass-through-authentication', 'changelog-password-encryption', 'processing-time-histogram', 'change-subscription-notification', 'sub-operation-timing', 'third-party', 'encrypt-attribute-values', 'pass-through-authentication', 'dn-mapper', 'referral-on-update', 'custom', 'snmp-subagent', 'password-policy-import', 'composed-attribute', 'ldap-result-code-tracker', 'attribute-mapper', 'delay', 'groovy-scripted', 'last-mod', 'pluggable-pass-through-authentication', 'referential-integrity', 'unique-attribute']")
	}
	if internaltypes.IsDefined(model.InitialConnections) && model.ResourceType.ValueString() != "pass-through-authentication" {
		resp.Diagnostics.AddError("Attribute 'initial_connections' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'initial_connections', the 'resource_type' attribute must be one of ['pass-through-authentication']")
	}
	if internaltypes.IsDefined(model.DelayPostGC) && model.ResourceType.ValueString() != "periodic-gc" {
		resp.Diagnostics.AddError("Attribute 'delay_post_gc' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'delay_post_gc', the 'resource_type' attribute must be one of ['periodic-gc']")
	}
	if internaltypes.IsDefined(model.NumMostExpensivePhasesShown) && model.ResourceType.ValueString() != "sub-operation-timing" {
		resp.Diagnostics.AddError("Attribute 'num_most_expensive_phases_shown' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'num_most_expensive_phases_shown', the 'resource_type' attribute must be one of ['sub-operation-timing']")
	}
	if internaltypes.IsDefined(model.AlwaysMapResponses) && model.ResourceType.ValueString() != "attribute-mapper" && model.ResourceType.ValueString() != "dn-mapper" {
		resp.Diagnostics.AddError("Attribute 'always_map_responses' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'always_map_responses', the 'resource_type' attribute must be one of ['attribute-mapper', 'dn-mapper']")
	}
	if internaltypes.IsDefined(model.ServerInfo) && model.ResourceType.ValueString() != "stats-collector" {
		resp.Diagnostics.AddError("Attribute 'server_info' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'server_info', the 'resource_type' attribute must be one of ['stats-collector']")
	}
	if internaltypes.IsDefined(model.LogFilePermissions) && model.ResourceType.ValueString() != "periodic-stats-logger" && model.ResourceType.ValueString() != "monitor-history" {
		resp.Diagnostics.AddError("Attribute 'log_file_permissions' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'log_file_permissions', the 'resource_type' attribute must be one of ['periodic-stats-logger', 'monitor-history']")
	}
	if internaltypes.IsDefined(model.Append) && model.ResourceType.ValueString() != "periodic-stats-logger" {
		resp.Diagnostics.AddError("Attribute 'append' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'append', the 'resource_type' attribute must be one of ['periodic-stats-logger']")
	}
	if internaltypes.IsDefined(model.InvokeGCTimeUtc) && model.ResourceType.ValueString() != "periodic-gc" {
		resp.Diagnostics.AddError("Attribute 'invoke_gc_time_utc' not supported by pingdirectory_plugin resources with 'resource_type' '"+model.ResourceType.ValueString()+"'",
			"When using attribute 'invoke_gc_time_utc', the 'resource_type' attribute must be one of ['periodic-gc']")
	}
}

// Add optional fields to create request for internal-search-rate plugin
func addOptionalInternalSearchRatePluginFields(ctx context.Context, addRequest *client.AddInternalSearchRatePluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.PluginType) {
		var slice []string
		plan.PluginType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginPluginTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginPluginTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PluginType = enumSlice
	}
	if internaltypes.IsDefined(plan.NumThreads) {
		addRequest.NumThreads = plan.NumThreads.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.LowerBound) {
		addRequest.LowerBound = plan.LowerBound.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.UpperBound) {
		addRequest.UpperBound = plan.UpperBound.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.FilterSuffix) {
		addRequest.FilterSuffix = plan.FilterSuffix.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for modifiable-password-policy-state plugin
func addOptionalModifiablePasswordPolicyStatePluginFields(ctx context.Context, addRequest *client.AddModifiablePasswordPolicyStatePluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for seven-bit-clean plugin
func addOptionalSevenBitCleanPluginFields(ctx context.Context, addRequest *client.AddSevenBitCleanPluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.PluginType) {
		var slice []string
		plan.PluginType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginPluginTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginPluginTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PluginType = enumSlice
	}
	if internaltypes.IsDefined(plan.AttributeType) {
		var slice []string
		plan.AttributeType.ElementsAs(ctx, &slice, false)
		addRequest.AttributeType = slice
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for clean-up-expired-pingfederate-persistent-access-grants plugin
func addOptionalCleanUpExpiredPingfederatePersistentAccessGrantsPluginFields(ctx context.Context, addRequest *client.AddCleanUpExpiredPingfederatePersistentAccessGrantsPluginRequest, plan pluginResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PollingInterval) {
		addRequest.PollingInterval = plan.PollingInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.PeerServerPriorityIndex) {
		addRequest.PeerServerPriorityIndex = plan.PeerServerPriorityIndex.ValueInt64Pointer()
	}
	// Treat this set as a single string
	if internaltypes.IsDefined(plan.BaseDN) && len(plan.BaseDN.Elements()) > 0 {
		addRequest.BaseDN = plan.BaseDN.Elements()[0].(types.String).ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MaxUpdatesPerSecond) {
		addRequest.MaxUpdatesPerSecond = plan.MaxUpdatesPerSecond.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.NumDeleteThreads) {
		addRequest.NumDeleteThreads = plan.NumDeleteThreads.ValueInt64Pointer()
	}
	return nil
}

// Add optional fields to create request for periodic-gc plugin
func addOptionalPeriodicGcPluginFields(ctx context.Context, addRequest *client.AddPeriodicGcPluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.PluginType) {
		var slice []string
		plan.PluginType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginPluginTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginPluginTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PluginType = enumSlice
	}
	if internaltypes.IsDefined(plan.InvokeGCDayOfWeek) {
		var slice []string
		plan.InvokeGCDayOfWeek.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginInvokeGCDayOfWeekProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginInvokeGCDayOfWeekPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.InvokeGCDayOfWeek = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DelayAfterAlert) {
		addRequest.DelayAfterAlert = plan.DelayAfterAlert.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DelayPostGC) {
		addRequest.DelayPostGC = plan.DelayPostGC.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for ping-one-pass-through-authentication plugin
func addOptionalPingOnePassThroughAuthenticationPluginFields(ctx context.Context, addRequest *client.AddPingOnePassThroughAuthenticationPluginRequest, plan pluginResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OAuthClientSecret) {
		addRequest.OAuthClientSecret = plan.OAuthClientSecret.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OAuthClientSecretPassphraseProvider) {
		addRequest.OAuthClientSecretPassphraseProvider = plan.OAuthClientSecretPassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpProxyExternalServer) {
		addRequest.HttpProxyExternalServer = plan.HttpProxyExternalServer.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludedLocalEntryBaseDN) {
		var slice []string
		plan.IncludedLocalEntryBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedLocalEntryBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.TryLocalBind) {
		addRequest.TryLocalBind = plan.TryLocalBind.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.OverrideLocalPassword) {
		addRequest.OverrideLocalPassword = plan.OverrideLocalPassword.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.UpdateLocalPassword) {
		addRequest.UpdateLocalPassword = plan.UpdateLocalPassword.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdateLocalPasswordDN) {
		addRequest.UpdateLocalPasswordDN = plan.UpdateLocalPasswordDN.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AllowLaxPassThroughAuthenticationPasswords) {
		addRequest.AllowLaxPassThroughAuthenticationPasswords = plan.AllowLaxPassThroughAuthenticationPasswords.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IgnoredPasswordPolicyStateErrorCondition) {
		var slice []string
		plan.IgnoredPasswordPolicyStateErrorCondition.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginIgnoredPasswordPolicyStateErrorConditionProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginIgnoredPasswordPolicyStateErrorConditionPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.IgnoredPasswordPolicyStateErrorCondition = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AdditionalUserMappingSCIMFilter) {
		addRequest.AdditionalUserMappingSCIMFilter = plan.AdditionalUserMappingSCIMFilter.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for search-shutdown plugin
func addOptionalSearchShutdownPluginFields(ctx context.Context, addRequest *client.AddSearchShutdownPluginRequest, plan pluginResourceModel) error {
	// Treat this set as a single string
	if internaltypes.IsDefined(plan.BaseDN) && len(plan.BaseDN.Elements()) > 0 {
		addRequest.BaseDN = plan.BaseDN.Elements()[0].(types.String).ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PreviousFileExtension) {
		addRequest.PreviousFileExtension = plan.PreviousFileExtension.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for periodic-stats-logger plugin
func addOptionalPeriodicStatsLoggerPluginFields(ctx context.Context, addRequest *client.AddPeriodicStatsLoggerPluginRequest, plan pluginResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogInterval) {
		addRequest.LogInterval = plan.LogInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CollectionInterval) {
		addRequest.CollectionInterval = plan.CollectionInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.SuppressIfIdle) {
		addRequest.SuppressIfIdle = plan.SuppressIfIdle.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.HeaderPrefixPerColumn) {
		addRequest.HeaderPrefixPerColumn = plan.HeaderPrefixPerColumn.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EmptyInsteadOfZero) {
		addRequest.EmptyInsteadOfZero = plan.EmptyInsteadOfZero.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LinesBetweenHeader) {
		addRequest.LinesBetweenHeader = plan.LinesBetweenHeader.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.IncludedLDAPStat) {
		var slice []string
		plan.IncludedLDAPStat.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginIncludedLDAPStatProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginIncludedLDAPStatPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.IncludedLDAPStat = enumSlice
	}
	if internaltypes.IsDefined(plan.IncludedResourceStat) {
		var slice []string
		plan.IncludedResourceStat.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginIncludedResourceStatProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginIncludedResourceStatPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.IncludedResourceStat = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HistogramFormat) {
		histogramFormat, err := client.NewEnumpluginHistogramFormatPropFromValue(plan.HistogramFormat.ValueString())
		if err != nil {
			return err
		}
		addRequest.HistogramFormat = histogramFormat
	}
	if internaltypes.IsDefined(plan.HistogramOpType) {
		var slice []string
		plan.HistogramOpType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginHistogramOpTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginHistogramOpTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.HistogramOpType = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PerApplicationLDAPStats) {
		perApplicationLDAPStats, err := client.NewEnumpluginPeriodicStatsLoggerPerApplicationLDAPStatsPropFromValue(plan.PerApplicationLDAPStats.ValueString())
		if err != nil {
			return err
		}
		addRequest.PerApplicationLDAPStats = perApplicationLDAPStats
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.StatusSummaryInfo) {
		statusSummaryInfo, err := client.NewEnumpluginStatusSummaryInfoPropFromValue(plan.StatusSummaryInfo.ValueString())
		if err != nil {
			return err
		}
		addRequest.StatusSummaryInfo = statusSummaryInfo
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LdapChangelogInfo) {
		ldapChangelogInfo, err := client.NewEnumpluginLdapChangelogInfoPropFromValue(plan.LdapChangelogInfo.ValueString())
		if err != nil {
			return err
		}
		addRequest.LdapChangelogInfo = ldapChangelogInfo
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.GaugeInfo) {
		gaugeInfo, err := client.NewEnumpluginGaugeInfoPropFromValue(plan.GaugeInfo.ValueString())
		if err != nil {
			return err
		}
		addRequest.GaugeInfo = gaugeInfo
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFileFormat) {
		logFileFormat, err := client.NewEnumpluginLogFileFormatPropFromValue(plan.LogFileFormat.ValueString())
		if err != nil {
			return err
		}
		addRequest.LogFileFormat = logFileFormat
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumpluginLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LocalDBBackendInfo) {
		localDBBackendInfo, err := client.NewEnumpluginLocalDBBackendInfoPropFromValue(plan.LocalDBBackendInfo.ValueString())
		if err != nil {
			return err
		}
		addRequest.LocalDBBackendInfo = localDBBackendInfo
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReplicationInfo) {
		replicationInfo, err := client.NewEnumpluginReplicationInfoPropFromValue(plan.ReplicationInfo.ValueString())
		if err != nil {
			return err
		}
		addRequest.ReplicationInfo = replicationInfo
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EntryCacheInfo) {
		entryCacheInfo, err := client.NewEnumpluginEntryCacheInfoPropFromValue(plan.EntryCacheInfo.ValueString())
		if err != nil {
			return err
		}
		addRequest.EntryCacheInfo = entryCacheInfo
	}
	if internaltypes.IsDefined(plan.HostInfo) {
		var slice []string
		plan.HostInfo.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginHostInfoProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginHostInfoPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.HostInfo = enumSlice
	}
	if internaltypes.IsDefined(plan.IncludedLDAPApplication) {
		var slice []string
		plan.IncludedLDAPApplication.ElementsAs(ctx, &slice, false)
		addRequest.IncludedLDAPApplication = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for purge-expired-data plugin
func addOptionalPurgeExpiredDataPluginFields(ctx context.Context, addRequest *client.AddPurgeExpiredDataPluginRequest, plan pluginResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DatetimeJSONField) {
		addRequest.DatetimeJSONField = plan.DatetimeJSONField.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DatetimeFormat) {
		datetimeFormat, err := client.NewEnumpluginDatetimeFormatPropFromValue(plan.DatetimeFormat.ValueString())
		if err != nil {
			return err
		}
		addRequest.DatetimeFormat = datetimeFormat
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CustomDatetimeFormat) {
		addRequest.CustomDatetimeFormat = plan.CustomDatetimeFormat.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CustomTimezone) {
		addRequest.CustomTimezone = plan.CustomTimezone.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PurgeBehavior) {
		purgeBehavior, err := client.NewEnumpluginPurgeBehaviorPropFromValue(plan.PurgeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.PurgeBehavior = purgeBehavior
	}
	// Treat this set as a single string
	if internaltypes.IsDefined(plan.BaseDN) && len(plan.BaseDN.Elements()) > 0 {
		addRequest.BaseDN = plan.BaseDN.Elements()[0].(types.String).ValueStringPointer()
	}
	// Treat this set as a single string
	if internaltypes.IsDefined(plan.Filter) && len(plan.Filter.Elements()) > 0 {
		addRequest.Filter = plan.Filter.Elements()[0].(types.String).ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PollingInterval) {
		addRequest.PollingInterval = plan.PollingInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MaxUpdatesPerSecond) {
		addRequest.MaxUpdatesPerSecond = plan.MaxUpdatesPerSecond.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.PeerServerPriorityIndex) {
		addRequest.PeerServerPriorityIndex = plan.PeerServerPriorityIndex.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.NumDeleteThreads) {
		addRequest.NumDeleteThreads = plan.NumDeleteThreads.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for sub-operation-timing plugin
func addOptionalSubOperationTimingPluginFields(ctx context.Context, addRequest *client.AddSubOperationTimingPluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.PluginType) {
		var slice []string
		plan.PluginType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginPluginTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginPluginTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PluginType = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.NumMostExpensivePhasesShown) {
		addRequest.NumMostExpensivePhasesShown = plan.NumMostExpensivePhasesShown.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for third-party plugin
func addOptionalThirdPartyPluginFields(ctx context.Context, addRequest *client.AddThirdPartyPluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for pass-through-authentication plugin
func addOptionalPassThroughAuthenticationPluginFields(ctx context.Context, addRequest *client.AddPassThroughAuthenticationPluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.PluginType) {
		var slice []string
		plan.PluginType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginPluginTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginPluginTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PluginType = enumSlice
	}
	if internaltypes.IsDefined(plan.TryLocalBind) {
		addRequest.TryLocalBind = plan.TryLocalBind.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.OverrideLocalPassword) {
		addRequest.OverrideLocalPassword = plan.OverrideLocalPassword.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.UpdateLocalPassword) {
		addRequest.UpdateLocalPassword = plan.UpdateLocalPassword.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AllowLaxPassThroughAuthenticationPasswords) {
		addRequest.AllowLaxPassThroughAuthenticationPasswords = plan.AllowLaxPassThroughAuthenticationPasswords.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerAccessMode) {
		serverAccessMode, err := client.NewEnumpluginServerAccessModePropFromValue(plan.ServerAccessMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.ServerAccessMode = serverAccessMode
	}
	if internaltypes.IsDefined(plan.IncludedLocalEntryBaseDN) {
		var slice []string
		plan.IncludedLocalEntryBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedLocalEntryBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.DnMap) {
		var slice []string
		plan.DnMap.ElementsAs(ctx, &slice, false)
		addRequest.DnMap = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindDNPattern) {
		addRequest.BindDNPattern = plan.BindDNPattern.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchBaseDN) {
		addRequest.SearchBaseDN = plan.SearchBaseDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchFilterPattern) {
		addRequest.SearchFilterPattern = plan.SearchFilterPattern.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InitialConnections) {
		addRequest.InitialConnections = plan.InitialConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxConnections) {
		addRequest.MaxConnections = plan.MaxConnections.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for dn-mapper plugin
func addOptionalDnMapperPluginFields(ctx context.Context, addRequest *client.AddDnMapperPluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.PluginType) {
		var slice []string
		plan.PluginType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginPluginTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginPluginTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PluginType = enumSlice
	}
	if internaltypes.IsDefined(plan.EnableAttributeMapping) {
		addRequest.EnableAttributeMapping = plan.EnableAttributeMapping.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MapAttribute) {
		var slice []string
		plan.MapAttribute.ElementsAs(ctx, &slice, false)
		addRequest.MapAttribute = slice
	}
	if internaltypes.IsDefined(plan.EnableControlMapping) {
		addRequest.EnableControlMapping = plan.EnableControlMapping.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlwaysMapResponses) {
		addRequest.AlwaysMapResponses = plan.AlwaysMapResponses.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for referral-on-update plugin
func addOptionalReferralOnUpdatePluginFields(ctx context.Context, addRequest *client.AddReferralOnUpdatePluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.PluginType) {
		var slice []string
		plan.PluginType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginPluginTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginPluginTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PluginType = enumSlice
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for simple-to-external-bind plugin
func addOptionalSimpleToExternalBindPluginFields(ctx context.Context, addRequest *client.AddSimpleToExternalBindPluginRequest, plan pluginResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for snmp-subagent plugin
func addOptionalSnmpSubagentPluginFields(ctx context.Context, addRequest *client.AddSnmpSubagentPluginRequest, plan pluginResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ContextName) {
		addRequest.ContextName = plan.ContextName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AgentxAddress) {
		addRequest.AgentxAddress = plan.AgentxAddress.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AgentxPort) {
		addRequest.AgentxPort = plan.AgentxPort.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.NumWorkerThreads) {
		addRequest.NumWorkerThreads = plan.NumWorkerThreads.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SessionTimeout) {
		addRequest.SessionTimeout = plan.SessionTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectRetryMaxWait) {
		addRequest.ConnectRetryMaxWait = plan.ConnectRetryMaxWait.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PingInterval) {
		addRequest.PingInterval = plan.PingInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for clean-up-inactive-pingfederate-persistent-sessions plugin
func addOptionalCleanUpInactivePingfederatePersistentSessionsPluginFields(ctx context.Context, addRequest *client.AddCleanUpInactivePingfederatePersistentSessionsPluginRequest, plan pluginResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PollingInterval) {
		addRequest.PollingInterval = plan.PollingInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.PeerServerPriorityIndex) {
		addRequest.PeerServerPriorityIndex = plan.PeerServerPriorityIndex.ValueInt64Pointer()
	}
	// Treat this set as a single string
	if internaltypes.IsDefined(plan.BaseDN) && len(plan.BaseDN.Elements()) > 0 {
		addRequest.BaseDN = plan.BaseDN.Elements()[0].(types.String).ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MaxUpdatesPerSecond) {
		addRequest.MaxUpdatesPerSecond = plan.MaxUpdatesPerSecond.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.NumDeleteThreads) {
		addRequest.NumDeleteThreads = plan.NumDeleteThreads.ValueInt64Pointer()
	}
	return nil
}

// Add optional fields to create request for composed-attribute plugin
func addOptionalComposedAttributePluginFields(ctx context.Context, addRequest *client.AddComposedAttributePluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.PluginType) {
		var slice []string
		plan.PluginType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginPluginTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginPluginTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PluginType = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleValuePatternBehavior) {
		multipleValuePatternBehavior, err := client.NewEnumpluginMultipleValuePatternBehaviorPropFromValue(plan.MultipleValuePatternBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleValuePatternBehavior = multipleValuePatternBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultiValuedAttributeBehavior) {
		multiValuedAttributeBehavior, err := client.NewEnumpluginMultiValuedAttributeBehaviorPropFromValue(plan.MultiValuedAttributeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultiValuedAttributeBehavior = multiValuedAttributeBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TargetAttributeExistsDuringInitialPopulationBehavior) {
		targetAttributeExistsDuringInitialPopulationBehavior, err := client.NewEnumpluginTargetAttributeExistsDuringInitialPopulationBehaviorPropFromValue(plan.TargetAttributeExistsDuringInitialPopulationBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.TargetAttributeExistsDuringInitialPopulationBehavior = targetAttributeExistsDuringInitialPopulationBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdateSourceAttributeBehavior) {
		updateSourceAttributeBehavior, err := client.NewEnumpluginUpdateSourceAttributeBehaviorPropFromValue(plan.UpdateSourceAttributeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.UpdateSourceAttributeBehavior = updateSourceAttributeBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SourceAttributeRemovalBehavior) {
		sourceAttributeRemovalBehavior, err := client.NewEnumpluginSourceAttributeRemovalBehaviorPropFromValue(plan.SourceAttributeRemovalBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.SourceAttributeRemovalBehavior = sourceAttributeRemovalBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdateTargetAttributeBehavior) {
		updateTargetAttributeBehavior, err := client.NewEnumpluginUpdateTargetAttributeBehaviorPropFromValue(plan.UpdateTargetAttributeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.UpdateTargetAttributeBehavior = updateTargetAttributeBehavior
	}
	if internaltypes.IsDefined(plan.IncludeBaseDN) {
		var slice []string
		plan.IncludeBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludeBaseDN = slice
	}
	if internaltypes.IsDefined(plan.ExcludeBaseDN) {
		var slice []string
		plan.ExcludeBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeBaseDN = slice
	}
	if internaltypes.IsDefined(plan.IncludeFilter) {
		var slice []string
		plan.IncludeFilter.ElementsAs(ctx, &slice, false)
		addRequest.IncludeFilter = slice
	}
	if internaltypes.IsDefined(plan.ExcludeFilter) {
		var slice []string
		plan.ExcludeFilter.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeFilter = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdatedEntryNewlyMatchesCriteriaBehavior) {
		updatedEntryNewlyMatchesCriteriaBehavior, err := client.NewEnumpluginUpdatedEntryNewlyMatchesCriteriaBehaviorPropFromValue(plan.UpdatedEntryNewlyMatchesCriteriaBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.UpdatedEntryNewlyMatchesCriteriaBehavior = updatedEntryNewlyMatchesCriteriaBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdatedEntryNoLongerMatchesCriteriaBehavior) {
		updatedEntryNoLongerMatchesCriteriaBehavior, err := client.NewEnumpluginUpdatedEntryNoLongerMatchesCriteriaBehaviorPropFromValue(plan.UpdatedEntryNoLongerMatchesCriteriaBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.UpdatedEntryNoLongerMatchesCriteriaBehavior = updatedEntryNoLongerMatchesCriteriaBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for attribute-mapper plugin
func addOptionalAttributeMapperPluginFields(ctx context.Context, addRequest *client.AddAttributeMapperPluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.PluginType) {
		var slice []string
		plan.PluginType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginPluginTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginPluginTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PluginType = enumSlice
	}
	if internaltypes.IsDefined(plan.EnableControlMapping) {
		addRequest.EnableControlMapping = plan.EnableControlMapping.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlwaysMapResponses) {
		addRequest.AlwaysMapResponses = plan.AlwaysMapResponses.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for delay plugin
func addOptionalDelayPluginFields(ctx context.Context, addRequest *client.AddDelayPluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.PluginType) {
		var slice []string
		plan.PluginType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginPluginTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginPluginTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PluginType = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for clean-up-expired-pingfederate-persistent-sessions plugin
func addOptionalCleanUpExpiredPingfederatePersistentSessionsPluginFields(ctx context.Context, addRequest *client.AddCleanUpExpiredPingfederatePersistentSessionsPluginRequest, plan pluginResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PollingInterval) {
		addRequest.PollingInterval = plan.PollingInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.PeerServerPriorityIndex) {
		addRequest.PeerServerPriorityIndex = plan.PeerServerPriorityIndex.ValueInt64Pointer()
	}
	// Treat this set as a single string
	if internaltypes.IsDefined(plan.BaseDN) && len(plan.BaseDN.Elements()) > 0 {
		addRequest.BaseDN = plan.BaseDN.Elements()[0].(types.String).ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MaxUpdatesPerSecond) {
		addRequest.MaxUpdatesPerSecond = plan.MaxUpdatesPerSecond.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.NumDeleteThreads) {
		addRequest.NumDeleteThreads = plan.NumDeleteThreads.ValueInt64Pointer()
	}
	return nil
}

// Add optional fields to create request for groovy-scripted plugin
func addOptionalGroovyScriptedPluginFields(ctx context.Context, addRequest *client.AddGroovyScriptedPluginRequest, plan pluginResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for pluggable-pass-through-authentication plugin
func addOptionalPluggablePassThroughAuthenticationPluginFields(ctx context.Context, addRequest *client.AddPluggablePassThroughAuthenticationPluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.IncludedLocalEntryBaseDN) {
		var slice []string
		plan.IncludedLocalEntryBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedLocalEntryBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.TryLocalBind) {
		addRequest.TryLocalBind = plan.TryLocalBind.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.OverrideLocalPassword) {
		addRequest.OverrideLocalPassword = plan.OverrideLocalPassword.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.UpdateLocalPassword) {
		addRequest.UpdateLocalPassword = plan.UpdateLocalPassword.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdateLocalPasswordDN) {
		addRequest.UpdateLocalPasswordDN = plan.UpdateLocalPasswordDN.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AllowLaxPassThroughAuthenticationPasswords) {
		addRequest.AllowLaxPassThroughAuthenticationPasswords = plan.AllowLaxPassThroughAuthenticationPasswords.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IgnoredPasswordPolicyStateErrorCondition) {
		var slice []string
		plan.IgnoredPasswordPolicyStateErrorCondition.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginIgnoredPasswordPolicyStateErrorConditionProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginIgnoredPasswordPolicyStateErrorConditionPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.IgnoredPasswordPolicyStateErrorCondition = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for referential-integrity plugin
func addOptionalReferentialIntegrityPluginFields(ctx context.Context, addRequest *client.AddReferentialIntegrityPluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.PluginType) {
		var slice []string
		plan.PluginType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginPluginTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginPluginTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PluginType = enumSlice
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFile) {
		addRequest.LogFile = plan.LogFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdateInterval) {
		addRequest.UpdateInterval = plan.UpdateInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for unique-attribute plugin
func addOptionalUniqueAttributePluginFields(ctx context.Context, addRequest *client.AddUniqueAttributePluginRequest, plan pluginResourceModel) error {
	if internaltypes.IsDefined(plan.PluginType) {
		var slice []string
		plan.PluginType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginPluginTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginPluginTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PluginType = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleAttributeBehavior) {
		multipleAttributeBehavior, err := client.NewEnumpluginMultipleAttributeBehaviorPropFromValue(plan.MultipleAttributeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleAttributeBehavior = multipleAttributeBehavior
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.PreventConflictsWithSoftDeletedEntries) {
		addRequest.PreventConflictsWithSoftDeletedEntries = plan.PreventConflictsWithSoftDeletedEntries.ValueBoolPointer()
	}
	// Treat this set as a single string
	if internaltypes.IsDefined(plan.Filter) && len(plan.Filter.Elements()) > 0 {
		addRequest.Filter = plan.Filter.Elements()[0].(types.String).ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		addRequest.InvokeForInternalOperations = plan.InvokeForInternalOperations.ValueBoolPointer()
	}
	return nil
}

// Populate any sets that have a nil ElementType, to avoid a nil pointer when setting the state
func populatePluginNilSets(ctx context.Context, model *pluginResourceModel) {
	if model.ValuePattern.ElementType(ctx) == nil {
		model.ValuePattern = types.SetNull(types.StringType)
	}
	if model.UserMappingRemoteJSONField.ElementType(ctx) == nil {
		model.UserMappingRemoteJSONField = types.SetNull(types.StringType)
	}
	if model.InvokeGCDayOfWeek.ElementType(ctx) == nil {
		model.InvokeGCDayOfWeek = types.SetNull(types.StringType)
	}
	if model.Server.ElementType(ctx) == nil {
		model.Server = types.SetNull(types.StringType)
	}
	if model.IgnoredPasswordPolicyStateErrorCondition.ElementType(ctx) == nil {
		model.IgnoredPasswordPolicyStateErrorCondition = types.SetNull(types.StringType)
	}
	if model.IncludedLDAPApplication.ElementType(ctx) == nil {
		model.IncludedLDAPApplication = types.SetNull(types.StringType)
	}
	if model.IncludedLocalEntryBaseDN.ElementType(ctx) == nil {
		model.IncludedLocalEntryBaseDN = types.SetNull(types.StringType)
	}
	if model.RotationPolicy.ElementType(ctx) == nil {
		model.RotationPolicy = types.SetNull(types.StringType)
	}
	if model.ReferralBaseURL.ElementType(ctx) == nil {
		model.ReferralBaseURL = types.SetNull(types.StringType)
	}
	if model.ExcludeFilter.ElementType(ctx) == nil {
		model.ExcludeFilter = types.SetNull(types.StringType)
	}
	if model.IncludeAttribute.ElementType(ctx) == nil {
		model.IncludeAttribute = types.SetNull(types.StringType)
	}
	if model.UserMappingLocalAttribute.ElementType(ctx) == nil {
		model.UserMappingLocalAttribute = types.SetNull(types.StringType)
	}
	if model.AttributeType.ElementType(ctx) == nil {
		model.AttributeType = types.SetNull(types.StringType)
	}
	if model.IncludedLDAPStat.ElementType(ctx) == nil {
		model.IncludedLDAPStat = types.SetNull(types.StringType)
	}
	if model.RotationListener.ElementType(ctx) == nil {
		model.RotationListener = types.SetNull(types.StringType)
	}
	if model.BaseDN.ElementType(ctx) == nil {
		model.BaseDN = types.SetNull(types.StringType)
	}
	if model.IncludedResourceStat.ElementType(ctx) == nil {
		model.IncludedResourceStat = types.SetNull(types.StringType)
	}
	if model.IncludeBaseDN.ElementType(ctx) == nil {
		model.IncludeBaseDN = types.SetNull(types.StringType)
	}
	if model.DnMap.ElementType(ctx) == nil {
		model.DnMap = types.SetNull(types.StringType)
	}
	if model.InvokeGCTimeUtc.ElementType(ctx) == nil {
		model.InvokeGCTimeUtc = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
	if model.PluginType.ElementType(ctx) == nil {
		model.PluginType = types.SetNull(types.StringType)
	}
	if model.HistogramOpType.ElementType(ctx) == nil {
		model.HistogramOpType = types.SetNull(types.StringType)
	}
	if model.RetentionPolicy.ElementType(ctx) == nil {
		model.RetentionPolicy = types.SetNull(types.StringType)
	}
	if model.MapAttribute.ElementType(ctx) == nil {
		model.MapAttribute = types.SetNull(types.StringType)
	}
	if model.ScriptArgument.ElementType(ctx) == nil {
		model.ScriptArgument = types.SetNull(types.StringType)
	}
	if model.Type.ElementType(ctx) == nil {
		model.Type = types.SetNull(types.StringType)
	}
	if model.ExcludeBaseDN.ElementType(ctx) == nil {
		model.ExcludeBaseDN = types.SetNull(types.StringType)
	}
	if model.Filter.ElementType(ctx) == nil {
		model.Filter = types.SetNull(types.StringType)
	}
	if model.HostInfo.ElementType(ctx) == nil {
		model.HostInfo = types.SetNull(types.StringType)
	}
	if model.IncludeFilter.ElementType(ctx) == nil {
		model.IncludeFilter = types.SetNull(types.StringType)
	}
}

// Populate any sets that have a nil ElementType, to avoid a nil pointer when setting the state
func populatePluginNilSetsDefault(ctx context.Context, model *defaultPluginResourceModel) {
	if model.ValuePattern.ElementType(ctx) == nil {
		model.ValuePattern = types.SetNull(types.StringType)
	}
	if model.UserMappingRemoteJSONField.ElementType(ctx) == nil {
		model.UserMappingRemoteJSONField = types.SetNull(types.StringType)
	}
	if model.InvokeGCDayOfWeek.ElementType(ctx) == nil {
		model.InvokeGCDayOfWeek = types.SetNull(types.StringType)
	}
	if model.Server.ElementType(ctx) == nil {
		model.Server = types.SetNull(types.StringType)
	}
	if model.IgnoredPasswordPolicyStateErrorCondition.ElementType(ctx) == nil {
		model.IgnoredPasswordPolicyStateErrorCondition = types.SetNull(types.StringType)
	}
	if model.IncludedLDAPApplication.ElementType(ctx) == nil {
		model.IncludedLDAPApplication = types.SetNull(types.StringType)
	}
	if model.IncludedLocalEntryBaseDN.ElementType(ctx) == nil {
		model.IncludedLocalEntryBaseDN = types.SetNull(types.StringType)
	}
	if model.RotationPolicy.ElementType(ctx) == nil {
		model.RotationPolicy = types.SetNull(types.StringType)
	}
	if model.ReferralBaseURL.ElementType(ctx) == nil {
		model.ReferralBaseURL = types.SetNull(types.StringType)
	}
	if model.ExcludeFilter.ElementType(ctx) == nil {
		model.ExcludeFilter = types.SetNull(types.StringType)
	}
	if model.IncludeAttribute.ElementType(ctx) == nil {
		model.IncludeAttribute = types.SetNull(types.StringType)
	}
	if model.UserMappingLocalAttribute.ElementType(ctx) == nil {
		model.UserMappingLocalAttribute = types.SetNull(types.StringType)
	}
	if model.DefaultAuthPasswordStorageScheme.ElementType(ctx) == nil {
		model.DefaultAuthPasswordStorageScheme = types.SetNull(types.StringType)
	}
	if model.HistogramCategoryBoundary.ElementType(ctx) == nil {
		model.HistogramCategoryBoundary = types.SetNull(types.StringType)
	}
	if model.AttributeType.ElementType(ctx) == nil {
		model.AttributeType = types.SetNull(types.StringType)
	}
	if model.IncludedLDAPStat.ElementType(ctx) == nil {
		model.IncludedLDAPStat = types.SetNull(types.StringType)
	}
	if model.RotationListener.ElementType(ctx) == nil {
		model.RotationListener = types.SetNull(types.StringType)
	}
	if model.BaseDN.ElementType(ctx) == nil {
		model.BaseDN = types.SetNull(types.StringType)
	}
	if model.IncludedResourceStat.ElementType(ctx) == nil {
		model.IncludedResourceStat = types.SetNull(types.StringType)
	}
	if model.IncludeBaseDN.ElementType(ctx) == nil {
		model.IncludeBaseDN = types.SetNull(types.StringType)
	}
	if model.DnMap.ElementType(ctx) == nil {
		model.DnMap = types.SetNull(types.StringType)
	}
	if model.InvokeGCTimeUtc.ElementType(ctx) == nil {
		model.InvokeGCTimeUtc = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
	if model.PluginType.ElementType(ctx) == nil {
		model.PluginType = types.SetNull(types.StringType)
	}
	if model.HistogramOpType.ElementType(ctx) == nil {
		model.HistogramOpType = types.SetNull(types.StringType)
	}
	if model.RetentionPolicy.ElementType(ctx) == nil {
		model.RetentionPolicy = types.SetNull(types.StringType)
	}
	if model.MapAttribute.ElementType(ctx) == nil {
		model.MapAttribute = types.SetNull(types.StringType)
	}
	if model.ScriptArgument.ElementType(ctx) == nil {
		model.ScriptArgument = types.SetNull(types.StringType)
	}
	if model.ExcludeAttribute.ElementType(ctx) == nil {
		model.ExcludeAttribute = types.SetNull(types.StringType)
	}
	if model.Type.ElementType(ctx) == nil {
		model.Type = types.SetNull(types.StringType)
	}
	if model.ExcludeBaseDN.ElementType(ctx) == nil {
		model.ExcludeBaseDN = types.SetNull(types.StringType)
	}
	if model.Filter.ElementType(ctx) == nil {
		model.Filter = types.SetNull(types.StringType)
	}
	if model.DefaultUserPasswordStorageScheme.ElementType(ctx) == nil {
		model.DefaultUserPasswordStorageScheme = types.SetNull(types.StringType)
	}
	if model.OperationType.ElementType(ctx) == nil {
		model.OperationType = types.SetNull(types.StringType)
	}
	if model.HostInfo.ElementType(ctx) == nil {
		model.HostInfo = types.SetNull(types.StringType)
	}
	if model.IncludeFilter.ElementType(ctx) == nil {
		model.IncludeFilter = types.SetNull(types.StringType)
	}
}

// Read a LastAccessTimePluginResponse object into the model struct
func readLastAccessTimePluginResponseDefault(ctx context.Context, r *client.LastAccessTimePluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("last-access-time")
	state.Id = types.StringValue(r.Id)
	state.MaxUpdateFrequency = internaltypes.StringTypeOrNil(r.MaxUpdateFrequency, internaltypes.IsEmptyString(expectedValues.MaxUpdateFrequency))
	config.CheckMismatchedPDFormattedAttributes("max_update_frequency",
		expectedValues.MaxUpdateFrequency, state.MaxUpdateFrequency, diagnostics)
	state.OperationType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginOperationTypeProp(r.OperationType))
	state.InvokeForFailedBinds = internaltypes.BoolTypeOrNil(r.InvokeForFailedBinds)
	state.MaxSearchResultEntriesToUpdate = internaltypes.Int64TypeOrNil(r.MaxSearchResultEntriesToUpdate)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a StatsCollectorPluginResponse object into the model struct
func readStatsCollectorPluginResponseDefault(ctx context.Context, r *client.StatsCollectorPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("stats-collector")
	state.Id = types.StringValue(r.Id)
	state.SampleInterval = types.StringValue(r.SampleInterval)
	config.CheckMismatchedPDFormattedAttributes("sample_interval",
		expectedValues.SampleInterval, state.SampleInterval, diagnostics)
	state.CollectionInterval = types.StringValue(r.CollectionInterval)
	config.CheckMismatchedPDFormattedAttributes("collection_interval",
		expectedValues.CollectionInterval, state.CollectionInterval, diagnostics)
	state.LdapInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLdapInfoProp(r.LdapInfo), internaltypes.IsEmptyString(expectedValues.LdapInfo))
	state.ServerInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginServerInfoProp(r.ServerInfo), internaltypes.IsEmptyString(expectedValues.ServerInfo))
	state.PerApplicationLDAPStats = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginStatsCollectorPerApplicationLDAPStatsProp(r.PerApplicationLDAPStats), internaltypes.IsEmptyString(expectedValues.PerApplicationLDAPStats))
	state.LdapChangelogInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLdapChangelogInfoProp(r.LdapChangelogInfo), internaltypes.IsEmptyString(expectedValues.LdapChangelogInfo))
	state.StatusSummaryInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginStatusSummaryInfoProp(r.StatusSummaryInfo), internaltypes.IsEmptyString(expectedValues.StatusSummaryInfo))
	state.GenerateCollectorFiles = internaltypes.BoolTypeOrNil(r.GenerateCollectorFiles)
	state.LocalDBBackendInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLocalDBBackendInfoProp(r.LocalDBBackendInfo), internaltypes.IsEmptyString(expectedValues.LocalDBBackendInfo))
	state.ReplicationInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginReplicationInfoProp(r.ReplicationInfo), internaltypes.IsEmptyString(expectedValues.ReplicationInfo))
	state.EntryCacheInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginEntryCacheInfoProp(r.EntryCacheInfo), internaltypes.IsEmptyString(expectedValues.EntryCacheInfo))
	state.HostInfo = internaltypes.GetStringSet(
		client.StringSliceEnumpluginHostInfoProp(r.HostInfo))
	state.IncludedLDAPApplication = internaltypes.GetStringSet(r.IncludedLDAPApplication)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a InternalSearchRatePluginResponse object into the model struct
func readInternalSearchRatePluginResponse(ctx context.Context, r *client.InternalSearchRatePluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("internal-search-rate")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.NumThreads = types.Int64Value(r.NumThreads)
	baseDNValues := []string{r.BaseDN}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.LowerBound = internaltypes.Int64TypeOrNil(r.LowerBound)
	state.UpperBound = internaltypes.Int64TypeOrNil(r.UpperBound)
	state.FilterPrefix = types.StringValue(r.FilterPrefix)
	state.FilterSuffix = internaltypes.StringTypeOrNil(r.FilterSuffix, internaltypes.IsEmptyString(expectedValues.FilterSuffix))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a InternalSearchRatePluginResponse object into the model struct
func readInternalSearchRatePluginResponseDefault(ctx context.Context, r *client.InternalSearchRatePluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("internal-search-rate")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.NumThreads = types.Int64Value(r.NumThreads)
	baseDNValues := []string{r.BaseDN}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.LowerBound = internaltypes.Int64TypeOrNil(r.LowerBound)
	state.UpperBound = internaltypes.Int64TypeOrNil(r.UpperBound)
	state.FilterPrefix = types.StringValue(r.FilterPrefix)
	state.FilterSuffix = internaltypes.StringTypeOrNil(r.FilterSuffix, internaltypes.IsEmptyString(expectedValues.FilterSuffix))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a ModifiablePasswordPolicyStatePluginResponse object into the model struct
func readModifiablePasswordPolicyStatePluginResponse(ctx context.Context, r *client.ModifiablePasswordPolicyStatePluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("modifiable-password-policy-state")
	state.Id = types.StringValue(r.Id)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a ModifiablePasswordPolicyStatePluginResponse object into the model struct
func readModifiablePasswordPolicyStatePluginResponseDefault(ctx context.Context, r *client.ModifiablePasswordPolicyStatePluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("modifiable-password-policy-state")
	state.Id = types.StringValue(r.Id)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a SevenBitCleanPluginResponse object into the model struct
func readSevenBitCleanPluginResponse(ctx context.Context, r *client.SevenBitCleanPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("seven-bit-clean")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.AttributeType = internaltypes.GetStringSet(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a SevenBitCleanPluginResponse object into the model struct
func readSevenBitCleanPluginResponseDefault(ctx context.Context, r *client.SevenBitCleanPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("seven-bit-clean")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.AttributeType = internaltypes.GetStringSet(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse object into the model struct
func readCleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse(ctx context.Context, r *client.CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("clean-up-expired-pingfederate-persistent-access-grants")
	state.Id = types.StringValue(r.Id)
	state.PollingInterval = types.StringValue(r.PollingInterval)
	config.CheckMismatchedPDFormattedAttributes("polling_interval",
		expectedValues.PollingInterval, state.PollingInterval, diagnostics)
	state.PeerServerPriorityIndex = internaltypes.Int64TypeOrNil(r.PeerServerPriorityIndex)
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.MaxUpdatesPerSecond = types.Int64Value(r.MaxUpdatesPerSecond)
	state.NumDeleteThreads = types.Int64Value(r.NumDeleteThreads)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse object into the model struct
func readCleanUpExpiredPingfederatePersistentAccessGrantsPluginResponseDefault(ctx context.Context, r *client.CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("clean-up-expired-pingfederate-persistent-access-grants")
	state.Id = types.StringValue(r.Id)
	state.PollingInterval = types.StringValue(r.PollingInterval)
	config.CheckMismatchedPDFormattedAttributes("polling_interval",
		expectedValues.PollingInterval, state.PollingInterval, diagnostics)
	state.PeerServerPriorityIndex = internaltypes.Int64TypeOrNil(r.PeerServerPriorityIndex)
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.MaxUpdatesPerSecond = types.Int64Value(r.MaxUpdatesPerSecond)
	state.NumDeleteThreads = types.Int64Value(r.NumDeleteThreads)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a PeriodicGcPluginResponse object into the model struct
func readPeriodicGcPluginResponse(ctx context.Context, r *client.PeriodicGcPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("periodic-gc")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.InvokeGCDayOfWeek = internaltypes.GetStringSet(
		client.StringSliceEnumpluginInvokeGCDayOfWeekProp(r.InvokeGCDayOfWeek))
	state.InvokeGCTimeUtc = internaltypes.GetStringSet(r.InvokeGCTimeUtc)
	state.DelayAfterAlert = internaltypes.StringTypeOrNil(r.DelayAfterAlert, internaltypes.IsEmptyString(expectedValues.DelayAfterAlert))
	config.CheckMismatchedPDFormattedAttributes("delay_after_alert",
		expectedValues.DelayAfterAlert, state.DelayAfterAlert, diagnostics)
	state.DelayPostGC = internaltypes.StringTypeOrNil(r.DelayPostGC, internaltypes.IsEmptyString(expectedValues.DelayPostGC))
	config.CheckMismatchedPDFormattedAttributes("delay_post_gc",
		expectedValues.DelayPostGC, state.DelayPostGC, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a PeriodicGcPluginResponse object into the model struct
func readPeriodicGcPluginResponseDefault(ctx context.Context, r *client.PeriodicGcPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("periodic-gc")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.InvokeGCDayOfWeek = internaltypes.GetStringSet(
		client.StringSliceEnumpluginInvokeGCDayOfWeekProp(r.InvokeGCDayOfWeek))
	state.InvokeGCTimeUtc = internaltypes.GetStringSet(r.InvokeGCTimeUtc)
	state.DelayAfterAlert = internaltypes.StringTypeOrNil(r.DelayAfterAlert, internaltypes.IsEmptyString(expectedValues.DelayAfterAlert))
	config.CheckMismatchedPDFormattedAttributes("delay_after_alert",
		expectedValues.DelayAfterAlert, state.DelayAfterAlert, diagnostics)
	state.DelayPostGC = internaltypes.StringTypeOrNil(r.DelayPostGC, internaltypes.IsEmptyString(expectedValues.DelayPostGC))
	config.CheckMismatchedPDFormattedAttributes("delay_post_gc",
		expectedValues.DelayPostGC, state.DelayPostGC, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a PingOnePassThroughAuthenticationPluginResponse object into the model struct
func readPingOnePassThroughAuthenticationPluginResponse(ctx context.Context, r *client.PingOnePassThroughAuthenticationPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("ping-one-pass-through-authentication")
	state.Id = types.StringValue(r.Id)
	state.ApiURL = types.StringValue(r.ApiURL)
	state.AuthURL = types.StringValue(r.AuthURL)
	state.OAuthClientID = types.StringValue(r.OAuthClientID)
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.OAuthClientSecret = expectedValues.OAuthClientSecret
	state.OAuthClientSecretPassphraseProvider = internaltypes.StringTypeOrNil(r.OAuthClientSecretPassphraseProvider, internaltypes.IsEmptyString(expectedValues.OAuthClientSecretPassphraseProvider))
	state.EnvironmentID = types.StringValue(r.EnvironmentID)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, internaltypes.IsEmptyString(expectedValues.HttpProxyExternalServer))
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.TryLocalBind = internaltypes.BoolTypeOrNil(r.TryLocalBind)
	state.OverrideLocalPassword = internaltypes.BoolTypeOrNil(r.OverrideLocalPassword)
	state.UpdateLocalPassword = internaltypes.BoolTypeOrNil(r.UpdateLocalPassword)
	state.UpdateLocalPasswordDN = internaltypes.StringTypeOrNil(r.UpdateLocalPasswordDN, internaltypes.IsEmptyString(expectedValues.UpdateLocalPasswordDN))
	state.AllowLaxPassThroughAuthenticationPasswords = internaltypes.BoolTypeOrNil(r.AllowLaxPassThroughAuthenticationPasswords)
	state.IgnoredPasswordPolicyStateErrorCondition = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIgnoredPasswordPolicyStateErrorConditionProp(r.IgnoredPasswordPolicyStateErrorCondition))
	state.UserMappingLocalAttribute = internaltypes.GetStringSet(r.UserMappingLocalAttribute)
	state.UserMappingRemoteJSONField = internaltypes.GetStringSet(r.UserMappingRemoteJSONField)
	state.AdditionalUserMappingSCIMFilter = internaltypes.StringTypeOrNil(r.AdditionalUserMappingSCIMFilter, internaltypes.IsEmptyString(expectedValues.AdditionalUserMappingSCIMFilter))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a PingOnePassThroughAuthenticationPluginResponse object into the model struct
func readPingOnePassThroughAuthenticationPluginResponseDefault(ctx context.Context, r *client.PingOnePassThroughAuthenticationPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("ping-one-pass-through-authentication")
	state.Id = types.StringValue(r.Id)
	state.ApiURL = types.StringValue(r.ApiURL)
	state.AuthURL = types.StringValue(r.AuthURL)
	state.OAuthClientID = types.StringValue(r.OAuthClientID)
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.OAuthClientSecret = expectedValues.OAuthClientSecret
	state.OAuthClientSecretPassphraseProvider = internaltypes.StringTypeOrNil(r.OAuthClientSecretPassphraseProvider, internaltypes.IsEmptyString(expectedValues.OAuthClientSecretPassphraseProvider))
	state.EnvironmentID = types.StringValue(r.EnvironmentID)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, internaltypes.IsEmptyString(expectedValues.HttpProxyExternalServer))
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.TryLocalBind = internaltypes.BoolTypeOrNil(r.TryLocalBind)
	state.OverrideLocalPassword = internaltypes.BoolTypeOrNil(r.OverrideLocalPassword)
	state.UpdateLocalPassword = internaltypes.BoolTypeOrNil(r.UpdateLocalPassword)
	state.UpdateLocalPasswordDN = internaltypes.StringTypeOrNil(r.UpdateLocalPasswordDN, internaltypes.IsEmptyString(expectedValues.UpdateLocalPasswordDN))
	state.AllowLaxPassThroughAuthenticationPasswords = internaltypes.BoolTypeOrNil(r.AllowLaxPassThroughAuthenticationPasswords)
	state.IgnoredPasswordPolicyStateErrorCondition = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIgnoredPasswordPolicyStateErrorConditionProp(r.IgnoredPasswordPolicyStateErrorCondition))
	state.UserMappingLocalAttribute = internaltypes.GetStringSet(r.UserMappingLocalAttribute)
	state.UserMappingRemoteJSONField = internaltypes.GetStringSet(r.UserMappingRemoteJSONField)
	state.AdditionalUserMappingSCIMFilter = internaltypes.StringTypeOrNil(r.AdditionalUserMappingSCIMFilter, internaltypes.IsEmptyString(expectedValues.AdditionalUserMappingSCIMFilter))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a ChangelogPasswordEncryptionPluginResponse object into the model struct
func readChangelogPasswordEncryptionPluginResponseDefault(ctx context.Context, r *client.ChangelogPasswordEncryptionPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("changelog-password-encryption")
	state.Id = types.StringValue(r.Id)
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.ChangelogPasswordEncryptionKey = expectedValues.ChangelogPasswordEncryptionKey
	state.ChangelogPasswordEncryptionKeyPassphraseProvider = internaltypes.StringTypeOrNil(r.ChangelogPasswordEncryptionKeyPassphraseProvider, internaltypes.IsEmptyString(expectedValues.ChangelogPasswordEncryptionKeyPassphraseProvider))
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a ProcessingTimeHistogramPluginResponse object into the model struct
func readProcessingTimeHistogramPluginResponseDefault(ctx context.Context, r *client.ProcessingTimeHistogramPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("processing-time-histogram")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.HistogramCategoryBoundary = internaltypes.GetStringSet(r.HistogramCategoryBoundary)
	state.IncludeQueueTime = internaltypes.BoolTypeOrNil(r.IncludeQueueTime)
	state.SeparateMonitorEntryPerTrackedApplication = internaltypes.BoolTypeOrNil(r.SeparateMonitorEntryPerTrackedApplication)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a SearchShutdownPluginResponse object into the model struct
func readSearchShutdownPluginResponse(ctx context.Context, r *client.SearchShutdownPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("search-shutdown")
	state.Id = types.StringValue(r.Id)
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.Scope = types.StringValue(r.Scope.String())
	filterValues := []string{r.Filter}
	state.Filter = internaltypes.GetStringSet(filterValues)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.OutputFile = types.StringValue(r.OutputFile)
	state.PreviousFileExtension = internaltypes.StringTypeOrNil(r.PreviousFileExtension, internaltypes.IsEmptyString(expectedValues.PreviousFileExtension))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a SearchShutdownPluginResponse object into the model struct
func readSearchShutdownPluginResponseDefault(ctx context.Context, r *client.SearchShutdownPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("search-shutdown")
	state.Id = types.StringValue(r.Id)
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.Scope = types.StringValue(r.Scope.String())
	filterValues := []string{r.Filter}
	state.Filter = internaltypes.GetStringSet(filterValues)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.OutputFile = types.StringValue(r.OutputFile)
	state.PreviousFileExtension = internaltypes.StringTypeOrNil(r.PreviousFileExtension, internaltypes.IsEmptyString(expectedValues.PreviousFileExtension))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a PeriodicStatsLoggerPluginResponse object into the model struct
func readPeriodicStatsLoggerPluginResponse(ctx context.Context, r *client.PeriodicStatsLoggerPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("periodic-stats-logger")
	state.Id = types.StringValue(r.Id)
	state.LogInterval = types.StringValue(r.LogInterval)
	config.CheckMismatchedPDFormattedAttributes("log_interval",
		expectedValues.LogInterval, state.LogInterval, diagnostics)
	state.CollectionInterval = types.StringValue(r.CollectionInterval)
	config.CheckMismatchedPDFormattedAttributes("collection_interval",
		expectedValues.CollectionInterval, state.CollectionInterval, diagnostics)
	state.SuppressIfIdle = types.BoolValue(r.SuppressIfIdle)
	state.HeaderPrefixPerColumn = internaltypes.BoolTypeOrNil(r.HeaderPrefixPerColumn)
	state.EmptyInsteadOfZero = internaltypes.BoolTypeOrNil(r.EmptyInsteadOfZero)
	state.LinesBetweenHeader = types.Int64Value(r.LinesBetweenHeader)
	state.IncludedLDAPStat = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIncludedLDAPStatProp(r.IncludedLDAPStat))
	state.IncludedResourceStat = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIncludedResourceStatProp(r.IncludedResourceStat))
	state.HistogramFormat = types.StringValue(r.HistogramFormat.String())
	state.HistogramOpType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginHistogramOpTypeProp(r.HistogramOpType))
	state.PerApplicationLDAPStats = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginPeriodicStatsLoggerPerApplicationLDAPStatsProp(r.PerApplicationLDAPStats), internaltypes.IsEmptyString(expectedValues.PerApplicationLDAPStats))
	state.StatusSummaryInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginStatusSummaryInfoProp(r.StatusSummaryInfo), internaltypes.IsEmptyString(expectedValues.StatusSummaryInfo))
	state.LdapChangelogInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLdapChangelogInfoProp(r.LdapChangelogInfo), internaltypes.IsEmptyString(expectedValues.LdapChangelogInfo))
	state.GaugeInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginGaugeInfoProp(r.GaugeInfo), internaltypes.IsEmptyString(expectedValues.GaugeInfo))
	state.LogFileFormat = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLogFileFormatProp(r.LogFileFormat), internaltypes.IsEmptyString(expectedValues.LogFileFormat))
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.LocalDBBackendInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLocalDBBackendInfoProp(r.LocalDBBackendInfo), internaltypes.IsEmptyString(expectedValues.LocalDBBackendInfo))
	state.ReplicationInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginReplicationInfoProp(r.ReplicationInfo), internaltypes.IsEmptyString(expectedValues.ReplicationInfo))
	state.EntryCacheInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginEntryCacheInfoProp(r.EntryCacheInfo), internaltypes.IsEmptyString(expectedValues.EntryCacheInfo))
	state.HostInfo = internaltypes.GetStringSet(
		client.StringSliceEnumpluginHostInfoProp(r.HostInfo))
	state.IncludedLDAPApplication = internaltypes.GetStringSet(r.IncludedLDAPApplication)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a PeriodicStatsLoggerPluginResponse object into the model struct
func readPeriodicStatsLoggerPluginResponseDefault(ctx context.Context, r *client.PeriodicStatsLoggerPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("periodic-stats-logger")
	state.Id = types.StringValue(r.Id)
	state.LogInterval = types.StringValue(r.LogInterval)
	config.CheckMismatchedPDFormattedAttributes("log_interval",
		expectedValues.LogInterval, state.LogInterval, diagnostics)
	state.CollectionInterval = types.StringValue(r.CollectionInterval)
	config.CheckMismatchedPDFormattedAttributes("collection_interval",
		expectedValues.CollectionInterval, state.CollectionInterval, diagnostics)
	state.SuppressIfIdle = types.BoolValue(r.SuppressIfIdle)
	state.HeaderPrefixPerColumn = internaltypes.BoolTypeOrNil(r.HeaderPrefixPerColumn)
	state.EmptyInsteadOfZero = internaltypes.BoolTypeOrNil(r.EmptyInsteadOfZero)
	state.LinesBetweenHeader = types.Int64Value(r.LinesBetweenHeader)
	state.IncludedLDAPStat = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIncludedLDAPStatProp(r.IncludedLDAPStat))
	state.IncludedResourceStat = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIncludedResourceStatProp(r.IncludedResourceStat))
	state.HistogramFormat = types.StringValue(r.HistogramFormat.String())
	state.HistogramOpType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginHistogramOpTypeProp(r.HistogramOpType))
	state.PerApplicationLDAPStats = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginPeriodicStatsLoggerPerApplicationLDAPStatsProp(r.PerApplicationLDAPStats), internaltypes.IsEmptyString(expectedValues.PerApplicationLDAPStats))
	state.StatusSummaryInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginStatusSummaryInfoProp(r.StatusSummaryInfo), internaltypes.IsEmptyString(expectedValues.StatusSummaryInfo))
	state.LdapChangelogInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLdapChangelogInfoProp(r.LdapChangelogInfo), internaltypes.IsEmptyString(expectedValues.LdapChangelogInfo))
	state.GaugeInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginGaugeInfoProp(r.GaugeInfo), internaltypes.IsEmptyString(expectedValues.GaugeInfo))
	state.LogFileFormat = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLogFileFormatProp(r.LogFileFormat), internaltypes.IsEmptyString(expectedValues.LogFileFormat))
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.LocalDBBackendInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLocalDBBackendInfoProp(r.LocalDBBackendInfo), internaltypes.IsEmptyString(expectedValues.LocalDBBackendInfo))
	state.ReplicationInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginReplicationInfoProp(r.ReplicationInfo), internaltypes.IsEmptyString(expectedValues.ReplicationInfo))
	state.EntryCacheInfo = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginEntryCacheInfoProp(r.EntryCacheInfo), internaltypes.IsEmptyString(expectedValues.EntryCacheInfo))
	state.HostInfo = internaltypes.GetStringSet(
		client.StringSliceEnumpluginHostInfoProp(r.HostInfo))
	state.IncludedLDAPApplication = internaltypes.GetStringSet(r.IncludedLDAPApplication)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a PurgeExpiredDataPluginResponse object into the model struct
func readPurgeExpiredDataPluginResponse(ctx context.Context, r *client.PurgeExpiredDataPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("purge-expired-data")
	state.Id = types.StringValue(r.Id)
	state.DatetimeAttribute = types.StringValue(r.DatetimeAttribute)
	state.DatetimeJSONField = internaltypes.StringTypeOrNil(r.DatetimeJSONField, internaltypes.IsEmptyString(expectedValues.DatetimeJSONField))
	state.DatetimeFormat = types.StringValue(r.DatetimeFormat.String())
	state.CustomDatetimeFormat = internaltypes.StringTypeOrNil(r.CustomDatetimeFormat, internaltypes.IsEmptyString(expectedValues.CustomDatetimeFormat))
	state.CustomTimezone = internaltypes.StringTypeOrNil(r.CustomTimezone, internaltypes.IsEmptyString(expectedValues.CustomTimezone))
	state.ExpirationOffset = types.StringValue(r.ExpirationOffset)
	config.CheckMismatchedPDFormattedAttributes("expiration_offset",
		expectedValues.ExpirationOffset, state.ExpirationOffset, diagnostics)
	state.PurgeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginPurgeBehaviorProp(r.PurgeBehavior), internaltypes.IsEmptyString(expectedValues.PurgeBehavior))
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	filterValues := []string{}
	filterType := internaltypes.StringTypeOrNil(r.Filter, false)
	if !filterType.IsNull() {
		filterValues = append(filterValues, filterType.ValueString())
	}
	state.Filter = internaltypes.GetStringSet(filterValues)
	state.PollingInterval = types.StringValue(r.PollingInterval)
	config.CheckMismatchedPDFormattedAttributes("polling_interval",
		expectedValues.PollingInterval, state.PollingInterval, diagnostics)
	state.MaxUpdatesPerSecond = types.Int64Value(r.MaxUpdatesPerSecond)
	state.PeerServerPriorityIndex = internaltypes.Int64TypeOrNil(r.PeerServerPriorityIndex)
	state.NumDeleteThreads = types.Int64Value(r.NumDeleteThreads)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a PurgeExpiredDataPluginResponse object into the model struct
func readPurgeExpiredDataPluginResponseDefault(ctx context.Context, r *client.PurgeExpiredDataPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("purge-expired-data")
	state.Id = types.StringValue(r.Id)
	state.DatetimeAttribute = types.StringValue(r.DatetimeAttribute)
	state.DatetimeJSONField = internaltypes.StringTypeOrNil(r.DatetimeJSONField, internaltypes.IsEmptyString(expectedValues.DatetimeJSONField))
	state.DatetimeFormat = types.StringValue(r.DatetimeFormat.String())
	state.CustomDatetimeFormat = internaltypes.StringTypeOrNil(r.CustomDatetimeFormat, internaltypes.IsEmptyString(expectedValues.CustomDatetimeFormat))
	state.CustomTimezone = internaltypes.StringTypeOrNil(r.CustomTimezone, internaltypes.IsEmptyString(expectedValues.CustomTimezone))
	state.ExpirationOffset = types.StringValue(r.ExpirationOffset)
	config.CheckMismatchedPDFormattedAttributes("expiration_offset",
		expectedValues.ExpirationOffset, state.ExpirationOffset, diagnostics)
	state.PurgeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginPurgeBehaviorProp(r.PurgeBehavior), internaltypes.IsEmptyString(expectedValues.PurgeBehavior))
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	filterValues := []string{}
	filterType := internaltypes.StringTypeOrNil(r.Filter, false)
	if !filterType.IsNull() {
		filterValues = append(filterValues, filterType.ValueString())
	}
	state.Filter = internaltypes.GetStringSet(filterValues)
	state.PollingInterval = types.StringValue(r.PollingInterval)
	config.CheckMismatchedPDFormattedAttributes("polling_interval",
		expectedValues.PollingInterval, state.PollingInterval, diagnostics)
	state.MaxUpdatesPerSecond = types.Int64Value(r.MaxUpdatesPerSecond)
	state.PeerServerPriorityIndex = internaltypes.Int64TypeOrNil(r.PeerServerPriorityIndex)
	state.NumDeleteThreads = types.Int64Value(r.NumDeleteThreads)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a ChangeSubscriptionNotificationPluginResponse object into the model struct
func readChangeSubscriptionNotificationPluginResponseDefault(ctx context.Context, r *client.ChangeSubscriptionNotificationPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("change-subscription-notification")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a SubOperationTimingPluginResponse object into the model struct
func readSubOperationTimingPluginResponse(ctx context.Context, r *client.SubOperationTimingPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("sub-operation-timing")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.NumMostExpensivePhasesShown = types.Int64Value(r.NumMostExpensivePhasesShown)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a SubOperationTimingPluginResponse object into the model struct
func readSubOperationTimingPluginResponseDefault(ctx context.Context, r *client.SubOperationTimingPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("sub-operation-timing")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.NumMostExpensivePhasesShown = types.Int64Value(r.NumMostExpensivePhasesShown)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a ThirdPartyPluginResponse object into the model struct
func readThirdPartyPluginResponse(ctx context.Context, r *client.ThirdPartyPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a ThirdPartyPluginResponse object into the model struct
func readThirdPartyPluginResponseDefault(ctx context.Context, r *client.ThirdPartyPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a EncryptAttributeValuesPluginResponse object into the model struct
func readEncryptAttributeValuesPluginResponseDefault(ctx context.Context, r *client.EncryptAttributeValuesPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("encrypt-attribute-values")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.AttributeType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginAttributeTypeProp(r.AttributeType))
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a PassThroughAuthenticationPluginResponse object into the model struct
func readPassThroughAuthenticationPluginResponse(ctx context.Context, r *client.PassThroughAuthenticationPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("pass-through-authentication")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Server = internaltypes.GetStringSet(r.Server)
	state.TryLocalBind = types.BoolValue(r.TryLocalBind)
	state.OverrideLocalPassword = types.BoolValue(r.OverrideLocalPassword)
	state.UpdateLocalPassword = types.BoolValue(r.UpdateLocalPassword)
	state.AllowLaxPassThroughAuthenticationPasswords = internaltypes.BoolTypeOrNil(r.AllowLaxPassThroughAuthenticationPasswords)
	state.ServerAccessMode = types.StringValue(r.ServerAccessMode.String())
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.DnMap = internaltypes.GetStringSet(r.DnMap)
	state.BindDNPattern = internaltypes.StringTypeOrNil(r.BindDNPattern, internaltypes.IsEmptyString(expectedValues.BindDNPattern))
	state.SearchBaseDN = internaltypes.StringTypeOrNil(r.SearchBaseDN, internaltypes.IsEmptyString(expectedValues.SearchBaseDN))
	state.SearchFilterPattern = internaltypes.StringTypeOrNil(r.SearchFilterPattern, internaltypes.IsEmptyString(expectedValues.SearchFilterPattern))
	state.InitialConnections = types.Int64Value(r.InitialConnections)
	state.MaxConnections = types.Int64Value(r.MaxConnections)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a PassThroughAuthenticationPluginResponse object into the model struct
func readPassThroughAuthenticationPluginResponseDefault(ctx context.Context, r *client.PassThroughAuthenticationPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("pass-through-authentication")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Server = internaltypes.GetStringSet(r.Server)
	state.TryLocalBind = types.BoolValue(r.TryLocalBind)
	state.OverrideLocalPassword = types.BoolValue(r.OverrideLocalPassword)
	state.UpdateLocalPassword = types.BoolValue(r.UpdateLocalPassword)
	state.AllowLaxPassThroughAuthenticationPasswords = internaltypes.BoolTypeOrNil(r.AllowLaxPassThroughAuthenticationPasswords)
	state.ServerAccessMode = types.StringValue(r.ServerAccessMode.String())
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.DnMap = internaltypes.GetStringSet(r.DnMap)
	state.BindDNPattern = internaltypes.StringTypeOrNil(r.BindDNPattern, internaltypes.IsEmptyString(expectedValues.BindDNPattern))
	state.SearchBaseDN = internaltypes.StringTypeOrNil(r.SearchBaseDN, internaltypes.IsEmptyString(expectedValues.SearchBaseDN))
	state.SearchFilterPattern = internaltypes.StringTypeOrNil(r.SearchFilterPattern, internaltypes.IsEmptyString(expectedValues.SearchFilterPattern))
	state.InitialConnections = types.Int64Value(r.InitialConnections)
	state.MaxConnections = types.Int64Value(r.MaxConnections)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a DnMapperPluginResponse object into the model struct
func readDnMapperPluginResponse(ctx context.Context, r *client.DnMapperPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("dn-mapper")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.SourceDN = types.StringValue(r.SourceDN)
	state.TargetDN = types.StringValue(r.TargetDN)
	state.EnableAttributeMapping = types.BoolValue(r.EnableAttributeMapping)
	state.MapAttribute = internaltypes.GetStringSet(r.MapAttribute)
	state.EnableControlMapping = types.BoolValue(r.EnableControlMapping)
	state.AlwaysMapResponses = types.BoolValue(r.AlwaysMapResponses)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a DnMapperPluginResponse object into the model struct
func readDnMapperPluginResponseDefault(ctx context.Context, r *client.DnMapperPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("dn-mapper")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.SourceDN = types.StringValue(r.SourceDN)
	state.TargetDN = types.StringValue(r.TargetDN)
	state.EnableAttributeMapping = types.BoolValue(r.EnableAttributeMapping)
	state.MapAttribute = internaltypes.GetStringSet(r.MapAttribute)
	state.EnableControlMapping = types.BoolValue(r.EnableControlMapping)
	state.AlwaysMapResponses = types.BoolValue(r.AlwaysMapResponses)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a MonitorHistoryPluginResponse object into the model struct
func readMonitorHistoryPluginResponseDefault(ctx context.Context, r *client.MonitorHistoryPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("monitor-history")
	state.Id = types.StringValue(r.Id)
	state.LogInterval = types.StringValue(r.LogInterval)
	config.CheckMismatchedPDFormattedAttributes("log_interval",
		expectedValues.LogInterval, state.LogInterval, diagnostics)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.RetainFilesSparselyByAge = internaltypes.BoolTypeOrNil(r.RetainFilesSparselyByAge)
	state.Sanitize = internaltypes.BoolTypeOrNil(r.Sanitize)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a ReferralOnUpdatePluginResponse object into the model struct
func readReferralOnUpdatePluginResponse(ctx context.Context, r *client.ReferralOnUpdatePluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("referral-on-update")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.ReferralBaseURL = internaltypes.GetStringSet(r.ReferralBaseURL)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a ReferralOnUpdatePluginResponse object into the model struct
func readReferralOnUpdatePluginResponseDefault(ctx context.Context, r *client.ReferralOnUpdatePluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("referral-on-update")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.ReferralBaseURL = internaltypes.GetStringSet(r.ReferralBaseURL)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a SimpleToExternalBindPluginResponse object into the model struct
func readSimpleToExternalBindPluginResponse(ctx context.Context, r *client.SimpleToExternalBindPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("simple-to-external-bind")
	state.Id = types.StringValue(r.Id)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a SimpleToExternalBindPluginResponse object into the model struct
func readSimpleToExternalBindPluginResponseDefault(ctx context.Context, r *client.SimpleToExternalBindPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("simple-to-external-bind")
	state.Id = types.StringValue(r.Id)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a CustomPluginResponse object into the model struct
func readCustomPluginResponseDefault(ctx context.Context, r *client.CustomPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a SnmpSubagentPluginResponse object into the model struct
func readSnmpSubagentPluginResponse(ctx context.Context, r *client.SnmpSubagentPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("snmp-subagent")
	state.Id = types.StringValue(r.Id)
	state.ContextName = internaltypes.StringTypeOrNil(r.ContextName, internaltypes.IsEmptyString(expectedValues.ContextName))
	state.AgentxAddress = types.StringValue(r.AgentxAddress)
	state.AgentxPort = types.Int64Value(r.AgentxPort)
	state.NumWorkerThreads = internaltypes.Int64TypeOrNil(r.NumWorkerThreads)
	state.SessionTimeout = internaltypes.StringTypeOrNil(r.SessionTimeout, internaltypes.IsEmptyString(expectedValues.SessionTimeout))
	config.CheckMismatchedPDFormattedAttributes("session_timeout",
		expectedValues.SessionTimeout, state.SessionTimeout, diagnostics)
	state.ConnectRetryMaxWait = internaltypes.StringTypeOrNil(r.ConnectRetryMaxWait, internaltypes.IsEmptyString(expectedValues.ConnectRetryMaxWait))
	config.CheckMismatchedPDFormattedAttributes("connect_retry_max_wait",
		expectedValues.ConnectRetryMaxWait, state.ConnectRetryMaxWait, diagnostics)
	state.PingInterval = internaltypes.StringTypeOrNil(r.PingInterval, internaltypes.IsEmptyString(expectedValues.PingInterval))
	config.CheckMismatchedPDFormattedAttributes("ping_interval",
		expectedValues.PingInterval, state.PingInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a SnmpSubagentPluginResponse object into the model struct
func readSnmpSubagentPluginResponseDefault(ctx context.Context, r *client.SnmpSubagentPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("snmp-subagent")
	state.Id = types.StringValue(r.Id)
	state.ContextName = internaltypes.StringTypeOrNil(r.ContextName, internaltypes.IsEmptyString(expectedValues.ContextName))
	state.AgentxAddress = types.StringValue(r.AgentxAddress)
	state.AgentxPort = types.Int64Value(r.AgentxPort)
	state.NumWorkerThreads = internaltypes.Int64TypeOrNil(r.NumWorkerThreads)
	state.SessionTimeout = internaltypes.StringTypeOrNil(r.SessionTimeout, internaltypes.IsEmptyString(expectedValues.SessionTimeout))
	config.CheckMismatchedPDFormattedAttributes("session_timeout",
		expectedValues.SessionTimeout, state.SessionTimeout, diagnostics)
	state.ConnectRetryMaxWait = internaltypes.StringTypeOrNil(r.ConnectRetryMaxWait, internaltypes.IsEmptyString(expectedValues.ConnectRetryMaxWait))
	config.CheckMismatchedPDFormattedAttributes("connect_retry_max_wait",
		expectedValues.ConnectRetryMaxWait, state.ConnectRetryMaxWait, diagnostics)
	state.PingInterval = internaltypes.StringTypeOrNil(r.PingInterval, internaltypes.IsEmptyString(expectedValues.PingInterval))
	config.CheckMismatchedPDFormattedAttributes("ping_interval",
		expectedValues.PingInterval, state.PingInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a PasswordPolicyImportPluginResponse object into the model struct
func readPasswordPolicyImportPluginResponseDefault(ctx context.Context, r *client.PasswordPolicyImportPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("password-policy-import")
	state.Id = types.StringValue(r.Id)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.DefaultUserPasswordStorageScheme = internaltypes.GetStringSet(r.DefaultUserPasswordStorageScheme)
	state.DefaultAuthPasswordStorageScheme = internaltypes.GetStringSet(r.DefaultAuthPasswordStorageScheme)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a ProfilerPluginResponse object into the model struct
func readProfilerPluginResponseDefault(ctx context.Context, r *client.ProfilerPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("profiler")
	state.Id = types.StringValue(r.Id)
	state.ProfileSampleInterval = types.StringValue(r.ProfileSampleInterval)
	config.CheckMismatchedPDFormattedAttributes("profile_sample_interval",
		expectedValues.ProfileSampleInterval, state.ProfileSampleInterval, diagnostics)
	state.ProfileDirectory = types.StringValue(r.ProfileDirectory)
	state.EnableProfilingOnStartup = types.BoolValue(r.EnableProfilingOnStartup)
	state.ProfileAction = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginProfileActionProp(r.ProfileAction), internaltypes.IsEmptyString(expectedValues.ProfileAction))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a CleanUpInactivePingfederatePersistentSessionsPluginResponse object into the model struct
func readCleanUpInactivePingfederatePersistentSessionsPluginResponse(ctx context.Context, r *client.CleanUpInactivePingfederatePersistentSessionsPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("clean-up-inactive-pingfederate-persistent-sessions")
	state.Id = types.StringValue(r.Id)
	state.ExpirationOffset = types.StringValue(r.ExpirationOffset)
	config.CheckMismatchedPDFormattedAttributes("expiration_offset",
		expectedValues.ExpirationOffset, state.ExpirationOffset, diagnostics)
	state.PollingInterval = types.StringValue(r.PollingInterval)
	config.CheckMismatchedPDFormattedAttributes("polling_interval",
		expectedValues.PollingInterval, state.PollingInterval, diagnostics)
	state.PeerServerPriorityIndex = internaltypes.Int64TypeOrNil(r.PeerServerPriorityIndex)
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.MaxUpdatesPerSecond = types.Int64Value(r.MaxUpdatesPerSecond)
	state.NumDeleteThreads = types.Int64Value(r.NumDeleteThreads)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a CleanUpInactivePingfederatePersistentSessionsPluginResponse object into the model struct
func readCleanUpInactivePingfederatePersistentSessionsPluginResponseDefault(ctx context.Context, r *client.CleanUpInactivePingfederatePersistentSessionsPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("clean-up-inactive-pingfederate-persistent-sessions")
	state.Id = types.StringValue(r.Id)
	state.ExpirationOffset = types.StringValue(r.ExpirationOffset)
	config.CheckMismatchedPDFormattedAttributes("expiration_offset",
		expectedValues.ExpirationOffset, state.ExpirationOffset, diagnostics)
	state.PollingInterval = types.StringValue(r.PollingInterval)
	config.CheckMismatchedPDFormattedAttributes("polling_interval",
		expectedValues.PollingInterval, state.PollingInterval, diagnostics)
	state.PeerServerPriorityIndex = internaltypes.Int64TypeOrNil(r.PeerServerPriorityIndex)
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.MaxUpdatesPerSecond = types.Int64Value(r.MaxUpdatesPerSecond)
	state.NumDeleteThreads = types.Int64Value(r.NumDeleteThreads)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a ComposedAttributePluginResponse object into the model struct
func readComposedAttributePluginResponse(ctx context.Context, r *client.ComposedAttributePluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("composed-attribute")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	attributeTypeValues := []string{r.AttributeType}
	state.AttributeType = internaltypes.GetStringSet(attributeTypeValues)
	state.ValuePattern = internaltypes.GetStringSet(r.ValuePattern)
	state.MultipleValuePatternBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginMultipleValuePatternBehaviorProp(r.MultipleValuePatternBehavior), internaltypes.IsEmptyString(expectedValues.MultipleValuePatternBehavior))
	state.MultiValuedAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginMultiValuedAttributeBehaviorProp(r.MultiValuedAttributeBehavior), internaltypes.IsEmptyString(expectedValues.MultiValuedAttributeBehavior))
	state.TargetAttributeExistsDuringInitialPopulationBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginTargetAttributeExistsDuringInitialPopulationBehaviorProp(r.TargetAttributeExistsDuringInitialPopulationBehavior), internaltypes.IsEmptyString(expectedValues.TargetAttributeExistsDuringInitialPopulationBehavior))
	state.UpdateSourceAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdateSourceAttributeBehaviorProp(r.UpdateSourceAttributeBehavior), internaltypes.IsEmptyString(expectedValues.UpdateSourceAttributeBehavior))
	state.SourceAttributeRemovalBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginSourceAttributeRemovalBehaviorProp(r.SourceAttributeRemovalBehavior), internaltypes.IsEmptyString(expectedValues.SourceAttributeRemovalBehavior))
	state.UpdateTargetAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdateTargetAttributeBehaviorProp(r.UpdateTargetAttributeBehavior), internaltypes.IsEmptyString(expectedValues.UpdateTargetAttributeBehavior))
	state.IncludeBaseDN = internaltypes.GetStringSet(r.IncludeBaseDN)
	state.ExcludeBaseDN = internaltypes.GetStringSet(r.ExcludeBaseDN)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.ExcludeFilter = internaltypes.GetStringSet(r.ExcludeFilter)
	state.UpdatedEntryNewlyMatchesCriteriaBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdatedEntryNewlyMatchesCriteriaBehaviorProp(r.UpdatedEntryNewlyMatchesCriteriaBehavior), internaltypes.IsEmptyString(expectedValues.UpdatedEntryNewlyMatchesCriteriaBehavior))
	state.UpdatedEntryNoLongerMatchesCriteriaBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdatedEntryNoLongerMatchesCriteriaBehaviorProp(r.UpdatedEntryNoLongerMatchesCriteriaBehavior), internaltypes.IsEmptyString(expectedValues.UpdatedEntryNoLongerMatchesCriteriaBehavior))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a ComposedAttributePluginResponse object into the model struct
func readComposedAttributePluginResponseDefault(ctx context.Context, r *client.ComposedAttributePluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("composed-attribute")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	attributeTypeValues := []string{r.AttributeType}
	state.AttributeType = internaltypes.GetStringSet(attributeTypeValues)
	state.ValuePattern = internaltypes.GetStringSet(r.ValuePattern)
	state.MultipleValuePatternBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginMultipleValuePatternBehaviorProp(r.MultipleValuePatternBehavior), internaltypes.IsEmptyString(expectedValues.MultipleValuePatternBehavior))
	state.MultiValuedAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginMultiValuedAttributeBehaviorProp(r.MultiValuedAttributeBehavior), internaltypes.IsEmptyString(expectedValues.MultiValuedAttributeBehavior))
	state.TargetAttributeExistsDuringInitialPopulationBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginTargetAttributeExistsDuringInitialPopulationBehaviorProp(r.TargetAttributeExistsDuringInitialPopulationBehavior), internaltypes.IsEmptyString(expectedValues.TargetAttributeExistsDuringInitialPopulationBehavior))
	state.UpdateSourceAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdateSourceAttributeBehaviorProp(r.UpdateSourceAttributeBehavior), internaltypes.IsEmptyString(expectedValues.UpdateSourceAttributeBehavior))
	state.SourceAttributeRemovalBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginSourceAttributeRemovalBehaviorProp(r.SourceAttributeRemovalBehavior), internaltypes.IsEmptyString(expectedValues.SourceAttributeRemovalBehavior))
	state.UpdateTargetAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdateTargetAttributeBehaviorProp(r.UpdateTargetAttributeBehavior), internaltypes.IsEmptyString(expectedValues.UpdateTargetAttributeBehavior))
	state.IncludeBaseDN = internaltypes.GetStringSet(r.IncludeBaseDN)
	state.ExcludeBaseDN = internaltypes.GetStringSet(r.ExcludeBaseDN)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.ExcludeFilter = internaltypes.GetStringSet(r.ExcludeFilter)
	state.UpdatedEntryNewlyMatchesCriteriaBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdatedEntryNewlyMatchesCriteriaBehaviorProp(r.UpdatedEntryNewlyMatchesCriteriaBehavior), internaltypes.IsEmptyString(expectedValues.UpdatedEntryNewlyMatchesCriteriaBehavior))
	state.UpdatedEntryNoLongerMatchesCriteriaBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdatedEntryNoLongerMatchesCriteriaBehaviorProp(r.UpdatedEntryNoLongerMatchesCriteriaBehavior), internaltypes.IsEmptyString(expectedValues.UpdatedEntryNoLongerMatchesCriteriaBehavior))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a LdapResultCodeTrackerPluginResponse object into the model struct
func readLdapResultCodeTrackerPluginResponseDefault(ctx context.Context, r *client.LdapResultCodeTrackerPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("ldap-result-code-tracker")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a AttributeMapperPluginResponse object into the model struct
func readAttributeMapperPluginResponse(ctx context.Context, r *client.AttributeMapperPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("attribute-mapper")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.SourceAttribute = types.StringValue(r.SourceAttribute)
	state.TargetAttribute = types.StringValue(r.TargetAttribute)
	state.EnableControlMapping = types.BoolValue(r.EnableControlMapping)
	state.AlwaysMapResponses = types.BoolValue(r.AlwaysMapResponses)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a AttributeMapperPluginResponse object into the model struct
func readAttributeMapperPluginResponseDefault(ctx context.Context, r *client.AttributeMapperPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("attribute-mapper")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.SourceAttribute = types.StringValue(r.SourceAttribute)
	state.TargetAttribute = types.StringValue(r.TargetAttribute)
	state.EnableControlMapping = types.BoolValue(r.EnableControlMapping)
	state.AlwaysMapResponses = types.BoolValue(r.AlwaysMapResponses)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a DelayPluginResponse object into the model struct
func readDelayPluginResponse(ctx context.Context, r *client.DelayPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("delay")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Delay = types.StringValue(r.Delay)
	config.CheckMismatchedPDFormattedAttributes("delay",
		expectedValues.Delay, state.Delay, diagnostics)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a DelayPluginResponse object into the model struct
func readDelayPluginResponseDefault(ctx context.Context, r *client.DelayPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("delay")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Delay = types.StringValue(r.Delay)
	config.CheckMismatchedPDFormattedAttributes("delay",
		expectedValues.Delay, state.Delay, diagnostics)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a CleanUpExpiredPingfederatePersistentSessionsPluginResponse object into the model struct
func readCleanUpExpiredPingfederatePersistentSessionsPluginResponse(ctx context.Context, r *client.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("clean-up-expired-pingfederate-persistent-sessions")
	state.Id = types.StringValue(r.Id)
	state.PollingInterval = types.StringValue(r.PollingInterval)
	config.CheckMismatchedPDFormattedAttributes("polling_interval",
		expectedValues.PollingInterval, state.PollingInterval, diagnostics)
	state.PeerServerPriorityIndex = internaltypes.Int64TypeOrNil(r.PeerServerPriorityIndex)
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.MaxUpdatesPerSecond = types.Int64Value(r.MaxUpdatesPerSecond)
	state.NumDeleteThreads = types.Int64Value(r.NumDeleteThreads)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a CleanUpExpiredPingfederatePersistentSessionsPluginResponse object into the model struct
func readCleanUpExpiredPingfederatePersistentSessionsPluginResponseDefault(ctx context.Context, r *client.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("clean-up-expired-pingfederate-persistent-sessions")
	state.Id = types.StringValue(r.Id)
	state.PollingInterval = types.StringValue(r.PollingInterval)
	config.CheckMismatchedPDFormattedAttributes("polling_interval",
		expectedValues.PollingInterval, state.PollingInterval, diagnostics)
	state.PeerServerPriorityIndex = internaltypes.Int64TypeOrNil(r.PeerServerPriorityIndex)
	baseDNValues := []string{}
	baseDNType := internaltypes.StringTypeOrNil(r.BaseDN, false)
	if !baseDNType.IsNull() {
		baseDNValues = append(baseDNValues, baseDNType.ValueString())
	}
	state.BaseDN = internaltypes.GetStringSet(baseDNValues)
	state.MaxUpdatesPerSecond = types.Int64Value(r.MaxUpdatesPerSecond)
	state.NumDeleteThreads = types.Int64Value(r.NumDeleteThreads)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a GroovyScriptedPluginResponse object into the model struct
func readGroovyScriptedPluginResponse(ctx context.Context, r *client.GroovyScriptedPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a GroovyScriptedPluginResponse object into the model struct
func readGroovyScriptedPluginResponseDefault(ctx context.Context, r *client.GroovyScriptedPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a LastModPluginResponse object into the model struct
func readLastModPluginResponseDefault(ctx context.Context, r *client.LastModPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("last-mod")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.ExcludeAttribute = internaltypes.GetStringSet(r.ExcludeAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a PluggablePassThroughAuthenticationPluginResponse object into the model struct
func readPluggablePassThroughAuthenticationPluginResponse(ctx context.Context, r *client.PluggablePassThroughAuthenticationPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("pluggable-pass-through-authentication")
	state.Id = types.StringValue(r.Id)
	state.PassThroughAuthenticationHandler = types.StringValue(r.PassThroughAuthenticationHandler)
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.TryLocalBind = internaltypes.BoolTypeOrNil(r.TryLocalBind)
	state.OverrideLocalPassword = internaltypes.BoolTypeOrNil(r.OverrideLocalPassword)
	state.UpdateLocalPassword = internaltypes.BoolTypeOrNil(r.UpdateLocalPassword)
	state.UpdateLocalPasswordDN = internaltypes.StringTypeOrNil(r.UpdateLocalPasswordDN, internaltypes.IsEmptyString(expectedValues.UpdateLocalPasswordDN))
	state.AllowLaxPassThroughAuthenticationPasswords = internaltypes.BoolTypeOrNil(r.AllowLaxPassThroughAuthenticationPasswords)
	state.IgnoredPasswordPolicyStateErrorCondition = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIgnoredPasswordPolicyStateErrorConditionProp(r.IgnoredPasswordPolicyStateErrorCondition))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a PluggablePassThroughAuthenticationPluginResponse object into the model struct
func readPluggablePassThroughAuthenticationPluginResponseDefault(ctx context.Context, r *client.PluggablePassThroughAuthenticationPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("pluggable-pass-through-authentication")
	state.Id = types.StringValue(r.Id)
	state.PassThroughAuthenticationHandler = types.StringValue(r.PassThroughAuthenticationHandler)
	state.IncludedLocalEntryBaseDN = internaltypes.GetStringSet(r.IncludedLocalEntryBaseDN)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.TryLocalBind = internaltypes.BoolTypeOrNil(r.TryLocalBind)
	state.OverrideLocalPassword = internaltypes.BoolTypeOrNil(r.OverrideLocalPassword)
	state.UpdateLocalPassword = internaltypes.BoolTypeOrNil(r.UpdateLocalPassword)
	state.UpdateLocalPasswordDN = internaltypes.StringTypeOrNil(r.UpdateLocalPasswordDN, internaltypes.IsEmptyString(expectedValues.UpdateLocalPasswordDN))
	state.AllowLaxPassThroughAuthenticationPasswords = internaltypes.BoolTypeOrNil(r.AllowLaxPassThroughAuthenticationPasswords)
	state.IgnoredPasswordPolicyStateErrorCondition = internaltypes.GetStringSet(
		client.StringSliceEnumpluginIgnoredPasswordPolicyStateErrorConditionProp(r.IgnoredPasswordPolicyStateErrorCondition))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a ReferentialIntegrityPluginResponse object into the model struct
func readReferentialIntegrityPluginResponse(ctx context.Context, r *client.ReferentialIntegrityPluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("referential-integrity")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.AttributeType = internaltypes.GetStringSet(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.LogFile = internaltypes.StringTypeOrNil(r.LogFile, internaltypes.IsEmptyString(expectedValues.LogFile))
	state.UpdateInterval = internaltypes.StringTypeOrNil(r.UpdateInterval, internaltypes.IsEmptyString(expectedValues.UpdateInterval))
	config.CheckMismatchedPDFormattedAttributes("update_interval",
		expectedValues.UpdateInterval, state.UpdateInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a ReferentialIntegrityPluginResponse object into the model struct
func readReferentialIntegrityPluginResponseDefault(ctx context.Context, r *client.ReferentialIntegrityPluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("referential-integrity")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.AttributeType = internaltypes.GetStringSet(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.LogFile = internaltypes.StringTypeOrNil(r.LogFile, internaltypes.IsEmptyString(expectedValues.LogFile))
	state.UpdateInterval = internaltypes.StringTypeOrNil(r.UpdateInterval, internaltypes.IsEmptyString(expectedValues.UpdateInterval))
	config.CheckMismatchedPDFormattedAttributes("update_interval",
		expectedValues.UpdateInterval, state.UpdateInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Read a UniqueAttributePluginResponse object into the model struct
func readUniqueAttributePluginResponse(ctx context.Context, r *client.UniqueAttributePluginResponse, state *pluginResourceModel, expectedValues *pluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("unique-attribute")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Type = internaltypes.GetStringSet(r.Type)
	state.MultipleAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginMultipleAttributeBehaviorProp(r.MultipleAttributeBehavior), internaltypes.IsEmptyString(expectedValues.MultipleAttributeBehavior))
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.PreventConflictsWithSoftDeletedEntries = internaltypes.BoolTypeOrNil(r.PreventConflictsWithSoftDeletedEntries)
	filterValues := []string{}
	filterType := internaltypes.StringTypeOrNil(r.Filter, false)
	if !filterType.IsNull() {
		filterValues = append(filterValues, filterType.ValueString())
	}
	state.Filter = internaltypes.GetStringSet(filterValues)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSets(ctx, state)
}

// Read a UniqueAttributePluginResponse object into the model struct
func readUniqueAttributePluginResponseDefault(ctx context.Context, r *client.UniqueAttributePluginResponse, state *defaultPluginResourceModel, expectedValues *defaultPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("unique-attribute")
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Type = internaltypes.GetStringSet(r.Type)
	state.MultipleAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginMultipleAttributeBehaviorProp(r.MultipleAttributeBehavior), internaltypes.IsEmptyString(expectedValues.MultipleAttributeBehavior))
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.PreventConflictsWithSoftDeletedEntries = internaltypes.BoolTypeOrNil(r.PreventConflictsWithSoftDeletedEntries)
	filterValues := []string{}
	filterType := internaltypes.StringTypeOrNil(r.Filter, false)
	if !filterType.IsNull() {
		filterValues = append(filterValues, filterType.ValueString())
	}
	state.Filter = internaltypes.GetStringSet(filterValues)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePluginNilSetsDefault(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createPluginOperations(plan pluginResourceModel, state pluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.PassThroughAuthenticationHandler, state.PassThroughAuthenticationHandler, "pass-through-authentication-handler")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Type, state.Type, "type")
	operations.AddStringOperationIfNecessary(&ops, plan.MultipleAttributeBehavior, state.MultipleAttributeBehavior, "multiple-attribute-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddBoolOperationIfNecessary(&ops, plan.PreventConflictsWithSoftDeletedEntries, state.PreventConflictsWithSoftDeletedEntries, "prevent-conflicts-with-soft-deleted-entries")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdateInterval, state.UpdateInterval, "update-interval")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Delay, state.Delay, "delay")
	operations.AddStringOperationIfNecessary(&ops, plan.SourceAttribute, state.SourceAttribute, "source-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.TargetAttribute, state.TargetAttribute, "target-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ValuePattern, state.ValuePattern, "value-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.MultipleValuePatternBehavior, state.MultipleValuePatternBehavior, "multiple-value-pattern-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.MultiValuedAttributeBehavior, state.MultiValuedAttributeBehavior, "multi-valued-attribute-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.TargetAttributeExistsDuringInitialPopulationBehavior, state.TargetAttributeExistsDuringInitialPopulationBehavior, "target-attribute-exists-during-initial-population-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdateSourceAttributeBehavior, state.UpdateSourceAttributeBehavior, "update-source-attribute-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.SourceAttributeRemovalBehavior, state.SourceAttributeRemovalBehavior, "source-attribute-removal-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdateTargetAttributeBehavior, state.UpdateTargetAttributeBehavior, "update-target-attribute-behavior")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeBaseDN, state.IncludeBaseDN, "include-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeBaseDN, state.ExcludeBaseDN, "exclude-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeFilter, state.IncludeFilter, "include-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeFilter, state.ExcludeFilter, "exclude-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdatedEntryNewlyMatchesCriteriaBehavior, state.UpdatedEntryNewlyMatchesCriteriaBehavior, "updated-entry-newly-matches-criteria-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdatedEntryNoLongerMatchesCriteriaBehavior, state.UpdatedEntryNoLongerMatchesCriteriaBehavior, "updated-entry-no-longer-matches-criteria-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.ContextName, state.ContextName, "context-name")
	operations.AddStringOperationIfNecessary(&ops, plan.AgentxAddress, state.AgentxAddress, "agentx-address")
	operations.AddInt64OperationIfNecessary(&ops, plan.AgentxPort, state.AgentxPort, "agentx-port")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumWorkerThreads, state.NumWorkerThreads, "num-worker-threads")
	operations.AddStringOperationIfNecessary(&ops, plan.SessionTimeout, state.SessionTimeout, "session-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectRetryMaxWait, state.ConnectRetryMaxWait, "connect-retry-max-wait")
	operations.AddStringOperationIfNecessary(&ops, plan.PingInterval, state.PingInterval, "ping-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ReferralBaseURL, state.ReferralBaseURL, "referral-base-url")
	operations.AddStringOperationIfNecessary(&ops, plan.SourceDN, state.SourceDN, "source-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.TargetDN, state.TargetDN, "target-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableAttributeMapping, state.EnableAttributeMapping, "enable-attribute-mapping")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MapAttribute, state.MapAttribute, "map-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableControlMapping, state.EnableControlMapping, "enable-control-mapping")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlwaysMapResponses, state.AlwaysMapResponses, "always-map-responses")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Server, state.Server, "server")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.DatetimeAttribute, state.DatetimeAttribute, "datetime-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.DatetimeJSONField, state.DatetimeJSONField, "datetime-json-field")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerAccessMode, state.ServerAccessMode, "server-access-mode")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumMostExpensivePhasesShown, state.NumMostExpensivePhasesShown, "num-most-expensive-phases-shown")
	operations.AddStringOperationIfNecessary(&ops, plan.DatetimeFormat, state.DatetimeFormat, "datetime-format")
	operations.AddStringOperationIfNecessary(&ops, plan.CustomDatetimeFormat, state.CustomDatetimeFormat, "custom-datetime-format")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DnMap, state.DnMap, "dn-map")
	operations.AddStringOperationIfNecessary(&ops, plan.BindDNPattern, state.BindDNPattern, "bind-dn-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchBaseDN, state.SearchBaseDN, "search-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchFilterPattern, state.SearchFilterPattern, "search-filter-pattern")
	operations.AddInt64OperationIfNecessary(&ops, plan.InitialConnections, state.InitialConnections, "initial-connections")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxConnections, state.MaxConnections, "max-connections")
	operations.AddStringOperationIfNecessary(&ops, plan.CustomTimezone, state.CustomTimezone, "custom-timezone")
	operations.AddStringOperationIfNecessary(&ops, plan.ExpirationOffset, state.ExpirationOffset, "expiration-offset")
	operations.AddStringOperationIfNecessary(&ops, plan.PurgeBehavior, state.PurgeBehavior, "purge-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.LogInterval, state.LogInterval, "log-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.SuppressIfIdle, state.SuppressIfIdle, "suppress-if-idle")
	operations.AddBoolOperationIfNecessary(&ops, plan.HeaderPrefixPerColumn, state.HeaderPrefixPerColumn, "header-prefix-per-column")
	operations.AddBoolOperationIfNecessary(&ops, plan.EmptyInsteadOfZero, state.EmptyInsteadOfZero, "empty-instead-of-zero")
	operations.AddInt64OperationIfNecessary(&ops, plan.LinesBetweenHeader, state.LinesBetweenHeader, "lines-between-header")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedLDAPStat, state.IncludedLDAPStat, "included-ldap-stat")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedResourceStat, state.IncludedResourceStat, "included-resource-stat")
	operations.AddStringOperationIfNecessary(&ops, plan.HistogramFormat, state.HistogramFormat, "histogram-format")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.HistogramOpType, state.HistogramOpType, "histogram-op-type")
	operations.AddStringOperationIfNecessary(&ops, plan.Scope, state.Scope, "scope")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeAttribute, state.IncludeAttribute, "include-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.GaugeInfo, state.GaugeInfo, "gauge-info")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFileFormat, state.LogFileFormat, "log-file-format")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFile, state.LogFile, "log-file")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFilePermissions, state.LogFilePermissions, "log-file-permissions")
	operations.AddBoolOperationIfNecessary(&ops, plan.Append, state.Append, "append")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RotationPolicy, state.RotationPolicy, "rotation-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RotationListener, state.RotationListener, "rotation-listener")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RetentionPolicy, state.RetentionPolicy, "retention-policy")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.OutputFile, state.OutputFile, "output-file")
	operations.AddStringOperationIfNecessary(&ops, plan.PreviousFileExtension, state.PreviousFileExtension, "previous-file-extension")
	operations.AddStringOperationIfNecessary(&ops, plan.ApiURL, state.ApiURL, "api-url")
	operations.AddStringOperationIfNecessary(&ops, plan.AuthURL, state.AuthURL, "auth-url")
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthClientID, state.OAuthClientID, "oauth-client-id")
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthClientSecret, state.OAuthClientSecret, "oauth-client-secret")
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthClientSecretPassphraseProvider, state.OAuthClientSecretPassphraseProvider, "oauth-client-secret-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.EnvironmentID, state.EnvironmentID, "environment-id")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedLocalEntryBaseDN, state.IncludedLocalEntryBaseDN, "included-local-entry-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.PollingInterval, state.PollingInterval, "polling-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.TryLocalBind, state.TryLocalBind, "try-local-bind")
	operations.AddBoolOperationIfNecessary(&ops, plan.OverrideLocalPassword, state.OverrideLocalPassword, "override-local-password")
	operations.AddBoolOperationIfNecessary(&ops, plan.UpdateLocalPassword, state.UpdateLocalPassword, "update-local-password")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdateLocalPasswordDN, state.UpdateLocalPasswordDN, "update-local-password-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowLaxPassThroughAuthenticationPasswords, state.AllowLaxPassThroughAuthenticationPasswords, "allow-lax-pass-through-authentication-passwords")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IgnoredPasswordPolicyStateErrorCondition, state.IgnoredPasswordPolicyStateErrorCondition, "ignored-password-policy-state-error-condition")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.UserMappingLocalAttribute, state.UserMappingLocalAttribute, "user-mapping-local-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.UserMappingRemoteJSONField, state.UserMappingRemoteJSONField, "user-mapping-remote-json-field")
	operations.AddStringOperationIfNecessary(&ops, plan.AdditionalUserMappingSCIMFilter, state.AdditionalUserMappingSCIMFilter, "additional-user-mapping-scim-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.InvokeGCDayOfWeek, state.InvokeGCDayOfWeek, "invoke-gc-day-of-week")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.InvokeGCTimeUtc, state.InvokeGCTimeUtc, "invoke-gc-time-utc")
	operations.AddStringOperationIfNecessary(&ops, plan.DelayAfterAlert, state.DelayAfterAlert, "delay-after-alert")
	operations.AddStringOperationIfNecessary(&ops, plan.DelayPostGC, state.DelayPostGC, "delay-post-gc")
	operations.AddInt64OperationIfNecessary(&ops, plan.PeerServerPriorityIndex, state.PeerServerPriorityIndex, "peer-server-priority-index")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PluginType, state.PluginType, "plugin-type")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxUpdatesPerSecond, state.MaxUpdatesPerSecond, "max-updates-per-second")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumDeleteThreads, state.NumDeleteThreads, "num-delete-threads")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumThreads, state.NumThreads, "num-threads")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddInt64OperationIfNecessary(&ops, plan.LowerBound, state.LowerBound, "lower-bound")
	operations.AddInt64OperationIfNecessary(&ops, plan.UpperBound, state.UpperBound, "upper-bound")
	operations.AddStringOperationIfNecessary(&ops, plan.FilterPrefix, state.FilterPrefix, "filter-prefix")
	operations.AddStringOperationIfNecessary(&ops, plan.FilterSuffix, state.FilterSuffix, "filter-suffix")
	operations.AddStringOperationIfNecessary(&ops, plan.CollectionInterval, state.CollectionInterval, "collection-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.PerApplicationLDAPStats, state.PerApplicationLDAPStats, "per-application-ldap-stats")
	operations.AddStringOperationIfNecessary(&ops, plan.LdapChangelogInfo, state.LdapChangelogInfo, "ldap-changelog-info")
	operations.AddStringOperationIfNecessary(&ops, plan.StatusSummaryInfo, state.StatusSummaryInfo, "status-summary-info")
	operations.AddStringOperationIfNecessary(&ops, plan.LocalDBBackendInfo, state.LocalDBBackendInfo, "local-db-backend-info")
	operations.AddStringOperationIfNecessary(&ops, plan.ReplicationInfo, state.ReplicationInfo, "replication-info")
	operations.AddStringOperationIfNecessary(&ops, plan.EntryCacheInfo, state.EntryCacheInfo, "entry-cache-info")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.HostInfo, state.HostInfo, "host-info")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedLDAPApplication, state.IncludedLDAPApplication, "included-ldap-application")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddBoolOperationIfNecessary(&ops, plan.InvokeForInternalOperations, state.InvokeForInternalOperations, "invoke-for-internal-operations")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create any update operations necessary to make the state match the plan
func createPluginOperationsDefault(plan defaultPluginResourceModel, state defaultPluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.PassThroughAuthenticationHandler, state.PassThroughAuthenticationHandler, "pass-through-authentication-handler")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Type, state.Type, "type")
	operations.AddStringOperationIfNecessary(&ops, plan.MultipleAttributeBehavior, state.MultipleAttributeBehavior, "multiple-attribute-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddBoolOperationIfNecessary(&ops, plan.PreventConflictsWithSoftDeletedEntries, state.PreventConflictsWithSoftDeletedEntries, "prevent-conflicts-with-soft-deleted-entries")
	operations.AddStringOperationIfNecessary(&ops, plan.ProfileSampleInterval, state.ProfileSampleInterval, "profile-sample-interval")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeAttribute, state.ExcludeAttribute, "exclude-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdateInterval, state.UpdateInterval, "update-interval")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Delay, state.Delay, "delay")
	operations.AddStringOperationIfNecessary(&ops, plan.SourceAttribute, state.SourceAttribute, "source-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.TargetAttribute, state.TargetAttribute, "target-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.ProfileDirectory, state.ProfileDirectory, "profile-directory")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ValuePattern, state.ValuePattern, "value-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.MultipleValuePatternBehavior, state.MultipleValuePatternBehavior, "multiple-value-pattern-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.MultiValuedAttributeBehavior, state.MultiValuedAttributeBehavior, "multi-valued-attribute-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.TargetAttributeExistsDuringInitialPopulationBehavior, state.TargetAttributeExistsDuringInitialPopulationBehavior, "target-attribute-exists-during-initial-population-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdateSourceAttributeBehavior, state.UpdateSourceAttributeBehavior, "update-source-attribute-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.SourceAttributeRemovalBehavior, state.SourceAttributeRemovalBehavior, "source-attribute-removal-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdateTargetAttributeBehavior, state.UpdateTargetAttributeBehavior, "update-target-attribute-behavior")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeBaseDN, state.IncludeBaseDN, "include-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeBaseDN, state.ExcludeBaseDN, "exclude-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeFilter, state.IncludeFilter, "include-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeFilter, state.ExcludeFilter, "exclude-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdatedEntryNewlyMatchesCriteriaBehavior, state.UpdatedEntryNewlyMatchesCriteriaBehavior, "updated-entry-newly-matches-criteria-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdatedEntryNoLongerMatchesCriteriaBehavior, state.UpdatedEntryNoLongerMatchesCriteriaBehavior, "updated-entry-no-longer-matches-criteria-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableProfilingOnStartup, state.EnableProfilingOnStartup, "enable-profiling-on-startup")
	operations.AddStringOperationIfNecessary(&ops, plan.ProfileAction, state.ProfileAction, "profile-action")
	operations.AddStringOperationIfNecessary(&ops, plan.ContextName, state.ContextName, "context-name")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultUserPasswordStorageScheme, state.DefaultUserPasswordStorageScheme, "default-user-password-storage-scheme")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultAuthPasswordStorageScheme, state.DefaultAuthPasswordStorageScheme, "default-auth-password-storage-scheme")
	operations.AddStringOperationIfNecessary(&ops, plan.AgentxAddress, state.AgentxAddress, "agentx-address")
	operations.AddInt64OperationIfNecessary(&ops, plan.AgentxPort, state.AgentxPort, "agentx-port")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumWorkerThreads, state.NumWorkerThreads, "num-worker-threads")
	operations.AddStringOperationIfNecessary(&ops, plan.SessionTimeout, state.SessionTimeout, "session-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectRetryMaxWait, state.ConnectRetryMaxWait, "connect-retry-max-wait")
	operations.AddStringOperationIfNecessary(&ops, plan.PingInterval, state.PingInterval, "ping-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ReferralBaseURL, state.ReferralBaseURL, "referral-base-url")
	operations.AddStringOperationIfNecessary(&ops, plan.SourceDN, state.SourceDN, "source-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.TargetDN, state.TargetDN, "target-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableAttributeMapping, state.EnableAttributeMapping, "enable-attribute-mapping")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MapAttribute, state.MapAttribute, "map-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.RetainFilesSparselyByAge, state.RetainFilesSparselyByAge, "retain-files-sparsely-by-age")
	operations.AddBoolOperationIfNecessary(&ops, plan.Sanitize, state.Sanitize, "sanitize")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableControlMapping, state.EnableControlMapping, "enable-control-mapping")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlwaysMapResponses, state.AlwaysMapResponses, "always-map-responses")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Server, state.Server, "server")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionSettingsDefinitionID, state.EncryptionSettingsDefinitionID, "encryption-settings-definition-id")
	operations.AddStringOperationIfNecessary(&ops, plan.DatetimeAttribute, state.DatetimeAttribute, "datetime-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.DatetimeJSONField, state.DatetimeJSONField, "datetime-json-field")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerAccessMode, state.ServerAccessMode, "server-access-mode")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumMostExpensivePhasesShown, state.NumMostExpensivePhasesShown, "num-most-expensive-phases-shown")
	operations.AddStringOperationIfNecessary(&ops, plan.DatetimeFormat, state.DatetimeFormat, "datetime-format")
	operations.AddStringOperationIfNecessary(&ops, plan.CustomDatetimeFormat, state.CustomDatetimeFormat, "custom-datetime-format")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DnMap, state.DnMap, "dn-map")
	operations.AddStringOperationIfNecessary(&ops, plan.BindDNPattern, state.BindDNPattern, "bind-dn-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchBaseDN, state.SearchBaseDN, "search-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchFilterPattern, state.SearchFilterPattern, "search-filter-pattern")
	operations.AddInt64OperationIfNecessary(&ops, plan.InitialConnections, state.InitialConnections, "initial-connections")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxConnections, state.MaxConnections, "max-connections")
	operations.AddStringOperationIfNecessary(&ops, plan.CustomTimezone, state.CustomTimezone, "custom-timezone")
	operations.AddStringOperationIfNecessary(&ops, plan.ExpirationOffset, state.ExpirationOffset, "expiration-offset")
	operations.AddStringOperationIfNecessary(&ops, plan.PurgeBehavior, state.PurgeBehavior, "purge-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.LogInterval, state.LogInterval, "log-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.ChangelogPasswordEncryptionKey, state.ChangelogPasswordEncryptionKey, "changelog-password-encryption-key")
	operations.AddBoolOperationIfNecessary(&ops, plan.SuppressIfIdle, state.SuppressIfIdle, "suppress-if-idle")
	operations.AddBoolOperationIfNecessary(&ops, plan.HeaderPrefixPerColumn, state.HeaderPrefixPerColumn, "header-prefix-per-column")
	operations.AddBoolOperationIfNecessary(&ops, plan.EmptyInsteadOfZero, state.EmptyInsteadOfZero, "empty-instead-of-zero")
	operations.AddInt64OperationIfNecessary(&ops, plan.LinesBetweenHeader, state.LinesBetweenHeader, "lines-between-header")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedLDAPStat, state.IncludedLDAPStat, "included-ldap-stat")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedResourceStat, state.IncludedResourceStat, "included-resource-stat")
	operations.AddStringOperationIfNecessary(&ops, plan.HistogramFormat, state.HistogramFormat, "histogram-format")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.HistogramOpType, state.HistogramOpType, "histogram-op-type")
	operations.AddStringOperationIfNecessary(&ops, plan.Scope, state.Scope, "scope")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.HistogramCategoryBoundary, state.HistogramCategoryBoundary, "histogram-category-boundary")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeAttribute, state.IncludeAttribute, "include-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.GaugeInfo, state.GaugeInfo, "gauge-info")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFileFormat, state.LogFileFormat, "log-file-format")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFile, state.LogFile, "log-file")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFilePermissions, state.LogFilePermissions, "log-file-permissions")
	operations.AddBoolOperationIfNecessary(&ops, plan.Append, state.Append, "append")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RotationPolicy, state.RotationPolicy, "rotation-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RotationListener, state.RotationListener, "rotation-listener")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RetentionPolicy, state.RetentionPolicy, "retention-policy")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.OutputFile, state.OutputFile, "output-file")
	operations.AddStringOperationIfNecessary(&ops, plan.PreviousFileExtension, state.PreviousFileExtension, "previous-file-extension")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeQueueTime, state.IncludeQueueTime, "include-queue-time")
	operations.AddBoolOperationIfNecessary(&ops, plan.SeparateMonitorEntryPerTrackedApplication, state.SeparateMonitorEntryPerTrackedApplication, "separate-monitor-entry-per-tracked-application")
	operations.AddStringOperationIfNecessary(&ops, plan.ChangelogPasswordEncryptionKeyPassphraseProvider, state.ChangelogPasswordEncryptionKeyPassphraseProvider, "changelog-password-encryption-key-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.ApiURL, state.ApiURL, "api-url")
	operations.AddStringOperationIfNecessary(&ops, plan.AuthURL, state.AuthURL, "auth-url")
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthClientID, state.OAuthClientID, "oauth-client-id")
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthClientSecret, state.OAuthClientSecret, "oauth-client-secret")
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthClientSecretPassphraseProvider, state.OAuthClientSecretPassphraseProvider, "oauth-client-secret-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.EnvironmentID, state.EnvironmentID, "environment-id")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedLocalEntryBaseDN, state.IncludedLocalEntryBaseDN, "included-local-entry-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.PollingInterval, state.PollingInterval, "polling-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.TryLocalBind, state.TryLocalBind, "try-local-bind")
	operations.AddBoolOperationIfNecessary(&ops, plan.OverrideLocalPassword, state.OverrideLocalPassword, "override-local-password")
	operations.AddBoolOperationIfNecessary(&ops, plan.UpdateLocalPassword, state.UpdateLocalPassword, "update-local-password")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdateLocalPasswordDN, state.UpdateLocalPasswordDN, "update-local-password-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowLaxPassThroughAuthenticationPasswords, state.AllowLaxPassThroughAuthenticationPasswords, "allow-lax-pass-through-authentication-passwords")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IgnoredPasswordPolicyStateErrorCondition, state.IgnoredPasswordPolicyStateErrorCondition, "ignored-password-policy-state-error-condition")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.UserMappingLocalAttribute, state.UserMappingLocalAttribute, "user-mapping-local-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.UserMappingRemoteJSONField, state.UserMappingRemoteJSONField, "user-mapping-remote-json-field")
	operations.AddStringOperationIfNecessary(&ops, plan.AdditionalUserMappingSCIMFilter, state.AdditionalUserMappingSCIMFilter, "additional-user-mapping-scim-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.InvokeGCDayOfWeek, state.InvokeGCDayOfWeek, "invoke-gc-day-of-week")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.InvokeGCTimeUtc, state.InvokeGCTimeUtc, "invoke-gc-time-utc")
	operations.AddStringOperationIfNecessary(&ops, plan.DelayAfterAlert, state.DelayAfterAlert, "delay-after-alert")
	operations.AddStringOperationIfNecessary(&ops, plan.DelayPostGC, state.DelayPostGC, "delay-post-gc")
	operations.AddInt64OperationIfNecessary(&ops, plan.PeerServerPriorityIndex, state.PeerServerPriorityIndex, "peer-server-priority-index")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PluginType, state.PluginType, "plugin-type")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxUpdatesPerSecond, state.MaxUpdatesPerSecond, "max-updates-per-second")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumDeleteThreads, state.NumDeleteThreads, "num-delete-threads")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumThreads, state.NumThreads, "num-threads")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddInt64OperationIfNecessary(&ops, plan.LowerBound, state.LowerBound, "lower-bound")
	operations.AddInt64OperationIfNecessary(&ops, plan.UpperBound, state.UpperBound, "upper-bound")
	operations.AddStringOperationIfNecessary(&ops, plan.FilterPrefix, state.FilterPrefix, "filter-prefix")
	operations.AddStringOperationIfNecessary(&ops, plan.FilterSuffix, state.FilterSuffix, "filter-suffix")
	operations.AddStringOperationIfNecessary(&ops, plan.SampleInterval, state.SampleInterval, "sample-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.CollectionInterval, state.CollectionInterval, "collection-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.LdapInfo, state.LdapInfo, "ldap-info")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerInfo, state.ServerInfo, "server-info")
	operations.AddStringOperationIfNecessary(&ops, plan.PerApplicationLDAPStats, state.PerApplicationLDAPStats, "per-application-ldap-stats")
	operations.AddStringOperationIfNecessary(&ops, plan.LdapChangelogInfo, state.LdapChangelogInfo, "ldap-changelog-info")
	operations.AddStringOperationIfNecessary(&ops, plan.StatusSummaryInfo, state.StatusSummaryInfo, "status-summary-info")
	operations.AddBoolOperationIfNecessary(&ops, plan.GenerateCollectorFiles, state.GenerateCollectorFiles, "generate-collector-files")
	operations.AddStringOperationIfNecessary(&ops, plan.LocalDBBackendInfo, state.LocalDBBackendInfo, "local-db-backend-info")
	operations.AddStringOperationIfNecessary(&ops, plan.ReplicationInfo, state.ReplicationInfo, "replication-info")
	operations.AddStringOperationIfNecessary(&ops, plan.EntryCacheInfo, state.EntryCacheInfo, "entry-cache-info")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.HostInfo, state.HostInfo, "host-info")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedLDAPApplication, state.IncludedLDAPApplication, "included-ldap-application")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxUpdateFrequency, state.MaxUpdateFrequency, "max-update-frequency")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OperationType, state.OperationType, "operation-type")
	operations.AddBoolOperationIfNecessary(&ops, plan.InvokeForFailedBinds, state.InvokeForFailedBinds, "invoke-for-failed-binds")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxSearchResultEntriesToUpdate, state.MaxSearchResultEntriesToUpdate, "max-search-result-entries-to-update")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddBoolOperationIfNecessary(&ops, plan.InvokeForInternalOperations, state.InvokeForInternalOperations, "invoke-for-internal-operations")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a internal-search-rate plugin
func (r *pluginResource) CreateInternalSearchRatePlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddInternalSearchRatePluginRequest(plan.Id.ValueString(),
		[]client.EnuminternalSearchRatePluginSchemaUrn{client.ENUMINTERNALSEARCHRATEPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGININTERNAL_SEARCH_RATE},
		plan.BaseDN.Elements()[0].(types.String).ValueString(),
		plan.FilterPrefix.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalInternalSearchRatePluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddInternalSearchRatePluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readInternalSearchRatePluginResponse(ctx, addResponse.InternalSearchRatePluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a modifiable-password-policy-state plugin
func (r *pluginResource) CreateModifiablePasswordPolicyStatePlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddModifiablePasswordPolicyStatePluginRequest(plan.Id.ValueString(),
		[]client.EnummodifiablePasswordPolicyStatePluginSchemaUrn{client.ENUMMODIFIABLEPASSWORDPOLICYSTATEPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINMODIFIABLE_PASSWORD_POLICY_STATE},
		plan.Enabled.ValueBool())
	err := addOptionalModifiablePasswordPolicyStatePluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddModifiablePasswordPolicyStatePluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readModifiablePasswordPolicyStatePluginResponse(ctx, addResponse.ModifiablePasswordPolicyStatePluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a seven-bit-clean plugin
func (r *pluginResource) CreateSevenBitCleanPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddSevenBitCleanPluginRequest(plan.Id.ValueString(),
		[]client.EnumsevenBitCleanPluginSchemaUrn{client.ENUMSEVENBITCLEANPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINSEVEN_BIT_CLEAN},
		plan.Enabled.ValueBool())
	err := addOptionalSevenBitCleanPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddSevenBitCleanPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readSevenBitCleanPluginResponse(ctx, addResponse.SevenBitCleanPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a clean-up-expired-pingfederate-persistent-access-grants plugin
func (r *pluginResource) CreateCleanUpExpiredPingfederatePersistentAccessGrantsPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddCleanUpExpiredPingfederatePersistentAccessGrantsPluginRequest(plan.Id.ValueString(),
		[]client.EnumcleanUpExpiredPingfederatePersistentAccessGrantsPluginSchemaUrn{client.ENUMCLEANUPEXPIREDPINGFEDERATEPERSISTENTACCESSGRANTSPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINCLEAN_UP_EXPIRED_PINGFEDERATE_PERSISTENT_ACCESS_GRANTS},
		plan.Enabled.ValueBool())
	err := addOptionalCleanUpExpiredPingfederatePersistentAccessGrantsPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddCleanUpExpiredPingfederatePersistentAccessGrantsPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readCleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse(ctx, addResponse.CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a periodic-gc plugin
func (r *pluginResource) CreatePeriodicGcPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	var InvokeGCTimeUtcSlice []string
	plan.InvokeGCTimeUtc.ElementsAs(ctx, &InvokeGCTimeUtcSlice, false)
	addRequest := client.NewAddPeriodicGcPluginRequest(plan.Id.ValueString(),
		[]client.EnumperiodicGcPluginSchemaUrn{client.ENUMPERIODICGCPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINPERIODIC_GC},
		InvokeGCTimeUtcSlice,
		plan.Enabled.ValueBool())
	err := addOptionalPeriodicGcPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddPeriodicGcPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readPeriodicGcPluginResponse(ctx, addResponse.PeriodicGcPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a ping-one-pass-through-authentication plugin
func (r *pluginResource) CreatePingOnePassThroughAuthenticationPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	var UserMappingLocalAttributeSlice []string
	plan.UserMappingLocalAttribute.ElementsAs(ctx, &UserMappingLocalAttributeSlice, false)
	var UserMappingRemoteJSONFieldSlice []string
	plan.UserMappingRemoteJSONField.ElementsAs(ctx, &UserMappingRemoteJSONFieldSlice, false)
	addRequest := client.NewAddPingOnePassThroughAuthenticationPluginRequest(plan.Id.ValueString(),
		[]client.EnumpingOnePassThroughAuthenticationPluginSchemaUrn{client.ENUMPINGONEPASSTHROUGHAUTHENTICATIONPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINPING_ONE_PASS_THROUGH_AUTHENTICATION},
		plan.ApiURL.ValueString(),
		plan.AuthURL.ValueString(),
		plan.OAuthClientID.ValueString(),
		plan.EnvironmentID.ValueString(),
		UserMappingLocalAttributeSlice,
		UserMappingRemoteJSONFieldSlice,
		plan.Enabled.ValueBool())
	err := addOptionalPingOnePassThroughAuthenticationPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddPingOnePassThroughAuthenticationPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readPingOnePassThroughAuthenticationPluginResponse(ctx, addResponse.PingOnePassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a search-shutdown plugin
func (r *pluginResource) CreateSearchShutdownPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	scope, err := client.NewEnumpluginScopePropFromValue(plan.Scope.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for Scope", err.Error())
		return nil, err
	}
	addRequest := client.NewAddSearchShutdownPluginRequest(plan.Id.ValueString(),
		[]client.EnumsearchShutdownPluginSchemaUrn{client.ENUMSEARCHSHUTDOWNPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINSEARCH_SHUTDOWN},
		*scope,
		plan.Filter.Elements()[0].(types.String).ValueString(),
		plan.OutputFile.ValueString(),
		plan.Enabled.ValueBool())
	err = addOptionalSearchShutdownPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddSearchShutdownPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readSearchShutdownPluginResponse(ctx, addResponse.SearchShutdownPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a periodic-stats-logger plugin
func (r *pluginResource) CreatePeriodicStatsLoggerPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddPeriodicStatsLoggerPluginRequest(plan.Id.ValueString(),
		[]client.EnumperiodicStatsLoggerPluginSchemaUrn{client.ENUMPERIODICSTATSLOGGERPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINPERIODIC_STATS_LOGGER},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalPeriodicStatsLoggerPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddPeriodicStatsLoggerPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readPeriodicStatsLoggerPluginResponse(ctx, addResponse.PeriodicStatsLoggerPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a purge-expired-data plugin
func (r *pluginResource) CreatePurgeExpiredDataPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddPurgeExpiredDataPluginRequest(plan.Id.ValueString(),
		[]client.EnumpurgeExpiredDataPluginSchemaUrn{client.ENUMPURGEEXPIREDDATAPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINPURGE_EXPIRED_DATA},
		plan.DatetimeAttribute.ValueString(),
		plan.ExpirationOffset.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalPurgeExpiredDataPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddPurgeExpiredDataPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readPurgeExpiredDataPluginResponse(ctx, addResponse.PurgeExpiredDataPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a sub-operation-timing plugin
func (r *pluginResource) CreateSubOperationTimingPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddSubOperationTimingPluginRequest(plan.Id.ValueString(),
		[]client.EnumsubOperationTimingPluginSchemaUrn{client.ENUMSUBOPERATIONTIMINGPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINSUB_OPERATION_TIMING},
		plan.Enabled.ValueBool())
	err := addOptionalSubOperationTimingPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddSubOperationTimingPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readSubOperationTimingPluginResponse(ctx, addResponse.SubOperationTimingPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party plugin
func (r *pluginResource) CreateThirdPartyPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	var PluginTypeSlice []client.EnumpluginPluginTypeProp
	plan.PluginType.ElementsAs(ctx, &PluginTypeSlice, false)
	addRequest := client.NewAddThirdPartyPluginRequest(plan.Id.ValueString(),
		[]client.EnumthirdPartyPluginSchemaUrn{client.ENUMTHIRDPARTYPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool(),
		PluginTypeSlice)
	err := addOptionalThirdPartyPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddThirdPartyPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readThirdPartyPluginResponse(ctx, addResponse.ThirdPartyPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a pass-through-authentication plugin
func (r *pluginResource) CreatePassThroughAuthenticationPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	var ServerSlice []string
	plan.Server.ElementsAs(ctx, &ServerSlice, false)
	addRequest := client.NewAddPassThroughAuthenticationPluginRequest(plan.Id.ValueString(),
		[]client.EnumpassThroughAuthenticationPluginSchemaUrn{client.ENUMPASSTHROUGHAUTHENTICATIONPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINPASS_THROUGH_AUTHENTICATION},
		ServerSlice,
		plan.Enabled.ValueBool())
	err := addOptionalPassThroughAuthenticationPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddPassThroughAuthenticationPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readPassThroughAuthenticationPluginResponse(ctx, addResponse.PassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a dn-mapper plugin
func (r *pluginResource) CreateDnMapperPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddDnMapperPluginRequest(plan.Id.ValueString(),
		[]client.EnumdnMapperPluginSchemaUrn{client.ENUMDNMAPPERPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINDN_MAPPER},
		plan.SourceDN.ValueString(),
		plan.TargetDN.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalDnMapperPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddDnMapperPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readDnMapperPluginResponse(ctx, addResponse.DnMapperPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a referral-on-update plugin
func (r *pluginResource) CreateReferralOnUpdatePlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	var ReferralBaseURLSlice []string
	plan.ReferralBaseURL.ElementsAs(ctx, &ReferralBaseURLSlice, false)
	addRequest := client.NewAddReferralOnUpdatePluginRequest(plan.Id.ValueString(),
		[]client.EnumreferralOnUpdatePluginSchemaUrn{client.ENUMREFERRALONUPDATEPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINREFERRAL_ON_UPDATE},
		ReferralBaseURLSlice,
		plan.Enabled.ValueBool())
	err := addOptionalReferralOnUpdatePluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddReferralOnUpdatePluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readReferralOnUpdatePluginResponse(ctx, addResponse.ReferralOnUpdatePluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a simple-to-external-bind plugin
func (r *pluginResource) CreateSimpleToExternalBindPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddSimpleToExternalBindPluginRequest(plan.Id.ValueString(),
		[]client.EnumsimpleToExternalBindPluginSchemaUrn{client.ENUMSIMPLETOEXTERNALBINDPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINSIMPLE_TO_EXTERNAL_BIND},
		plan.Enabled.ValueBool())
	err := addOptionalSimpleToExternalBindPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddSimpleToExternalBindPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readSimpleToExternalBindPluginResponse(ctx, addResponse.SimpleToExternalBindPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a snmp-subagent plugin
func (r *pluginResource) CreateSnmpSubagentPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddSnmpSubagentPluginRequest(plan.Id.ValueString(),
		[]client.EnumsnmpSubagentPluginSchemaUrn{client.ENUMSNMPSUBAGENTPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINSNMP_SUBAGENT},
		plan.Enabled.ValueBool())
	err := addOptionalSnmpSubagentPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddSnmpSubagentPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readSnmpSubagentPluginResponse(ctx, addResponse.SnmpSubagentPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a clean-up-inactive-pingfederate-persistent-sessions plugin
func (r *pluginResource) CreateCleanUpInactivePingfederatePersistentSessionsPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddCleanUpInactivePingfederatePersistentSessionsPluginRequest(plan.Id.ValueString(),
		[]client.EnumcleanUpInactivePingfederatePersistentSessionsPluginSchemaUrn{client.ENUMCLEANUPINACTIVEPINGFEDERATEPERSISTENTSESSIONSPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINCLEAN_UP_INACTIVE_PINGFEDERATE_PERSISTENT_SESSIONS},
		plan.ExpirationOffset.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalCleanUpInactivePingfederatePersistentSessionsPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddCleanUpInactivePingfederatePersistentSessionsPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readCleanUpInactivePingfederatePersistentSessionsPluginResponse(ctx, addResponse.CleanUpInactivePingfederatePersistentSessionsPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a composed-attribute plugin
func (r *pluginResource) CreateComposedAttributePlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	var ValuePatternSlice []string
	plan.ValuePattern.ElementsAs(ctx, &ValuePatternSlice, false)
	addRequest := client.NewAddComposedAttributePluginRequest(plan.Id.ValueString(),
		[]client.EnumcomposedAttributePluginSchemaUrn{client.ENUMCOMPOSEDATTRIBUTEPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINCOMPOSED_ATTRIBUTE},
		plan.AttributeType.Elements()[0].(types.String).ValueString(),
		ValuePatternSlice,
		plan.Enabled.ValueBool())
	err := addOptionalComposedAttributePluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddComposedAttributePluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readComposedAttributePluginResponse(ctx, addResponse.ComposedAttributePluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a attribute-mapper plugin
func (r *pluginResource) CreateAttributeMapperPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddAttributeMapperPluginRequest(plan.Id.ValueString(),
		[]client.EnumattributeMapperPluginSchemaUrn{client.ENUMATTRIBUTEMAPPERPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINATTRIBUTE_MAPPER},
		plan.SourceAttribute.ValueString(),
		plan.TargetAttribute.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalAttributeMapperPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddAttributeMapperPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readAttributeMapperPluginResponse(ctx, addResponse.AttributeMapperPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a delay plugin
func (r *pluginResource) CreateDelayPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddDelayPluginRequest(plan.Id.ValueString(),
		[]client.EnumdelayPluginSchemaUrn{client.ENUMDELAYPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINDELAY},
		plan.Delay.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalDelayPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddDelayPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readDelayPluginResponse(ctx, addResponse.DelayPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a clean-up-expired-pingfederate-persistent-sessions plugin
func (r *pluginResource) CreateCleanUpExpiredPingfederatePersistentSessionsPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddCleanUpExpiredPingfederatePersistentSessionsPluginRequest(plan.Id.ValueString(),
		[]client.EnumcleanUpExpiredPingfederatePersistentSessionsPluginSchemaUrn{client.ENUMCLEANUPEXPIREDPINGFEDERATEPERSISTENTSESSIONSPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINCLEAN_UP_EXPIRED_PINGFEDERATE_PERSISTENT_SESSIONS},
		plan.Enabled.ValueBool())
	err := addOptionalCleanUpExpiredPingfederatePersistentSessionsPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddCleanUpExpiredPingfederatePersistentSessionsPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readCleanUpExpiredPingfederatePersistentSessionsPluginResponse(ctx, addResponse.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted plugin
func (r *pluginResource) CreateGroovyScriptedPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	var PluginTypeSlice []client.EnumpluginPluginTypeProp
	plan.PluginType.ElementsAs(ctx, &PluginTypeSlice, false)
	addRequest := client.NewAddGroovyScriptedPluginRequest(plan.Id.ValueString(),
		[]client.EnumgroovyScriptedPluginSchemaUrn{client.ENUMGROOVYSCRIPTEDPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool(),
		PluginTypeSlice)
	err := addOptionalGroovyScriptedPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddGroovyScriptedPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readGroovyScriptedPluginResponse(ctx, addResponse.GroovyScriptedPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a pluggable-pass-through-authentication plugin
func (r *pluginResource) CreatePluggablePassThroughAuthenticationPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	addRequest := client.NewAddPluggablePassThroughAuthenticationPluginRequest(plan.Id.ValueString(),
		[]client.EnumpluggablePassThroughAuthenticationPluginSchemaUrn{client.ENUMPLUGGABLEPASSTHROUGHAUTHENTICATIONPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINPLUGGABLE_PASS_THROUGH_AUTHENTICATION},
		plan.PassThroughAuthenticationHandler.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalPluggablePassThroughAuthenticationPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddPluggablePassThroughAuthenticationPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readPluggablePassThroughAuthenticationPluginResponse(ctx, addResponse.PluggablePassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a referential-integrity plugin
func (r *pluginResource) CreateReferentialIntegrityPlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	var AttributeTypeSlice []string
	plan.AttributeType.ElementsAs(ctx, &AttributeTypeSlice, false)
	addRequest := client.NewAddReferentialIntegrityPluginRequest(plan.Id.ValueString(),
		[]client.EnumreferentialIntegrityPluginSchemaUrn{client.ENUMREFERENTIALINTEGRITYPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINREFERENTIAL_INTEGRITY},
		AttributeTypeSlice,
		plan.Enabled.ValueBool())
	err := addOptionalReferentialIntegrityPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddReferentialIntegrityPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readReferentialIntegrityPluginResponse(ctx, addResponse.ReferentialIntegrityPluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a unique-attribute plugin
func (r *pluginResource) CreateUniqueAttributePlugin(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan pluginResourceModel) (*pluginResourceModel, error) {
	var TypeSlice []string
	plan.Type.ElementsAs(ctx, &TypeSlice, false)
	addRequest := client.NewAddUniqueAttributePluginRequest(plan.Id.ValueString(),
		[]client.EnumuniqueAttributePluginSchemaUrn{client.ENUMUNIQUEATTRIBUTEPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINUNIQUE_ATTRIBUTE},
		TypeSlice,
		plan.Enabled.ValueBool())
	err := addOptionalUniqueAttributePluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Plugin", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddUniqueAttributePluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Plugin", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pluginResourceModel
	readUniqueAttributePluginResponse(ctx, addResponse.UniqueAttributePluginResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *pluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *pluginResourceModel
	var err error
	if plan.ResourceType.ValueString() == "internal-search-rate" {
		state, err = r.CreateInternalSearchRatePlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "modifiable-password-policy-state" {
		state, err = r.CreateModifiablePasswordPolicyStatePlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "seven-bit-clean" {
		state, err = r.CreateSevenBitCleanPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "clean-up-expired-pingfederate-persistent-access-grants" {
		state, err = r.CreateCleanUpExpiredPingfederatePersistentAccessGrantsPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "periodic-gc" {
		state, err = r.CreatePeriodicGcPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "ping-one-pass-through-authentication" {
		state, err = r.CreatePingOnePassThroughAuthenticationPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "search-shutdown" {
		state, err = r.CreateSearchShutdownPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "periodic-stats-logger" {
		state, err = r.CreatePeriodicStatsLoggerPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "purge-expired-data" {
		state, err = r.CreatePurgeExpiredDataPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "sub-operation-timing" {
		state, err = r.CreateSubOperationTimingPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "pass-through-authentication" {
		state, err = r.CreatePassThroughAuthenticationPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "dn-mapper" {
		state, err = r.CreateDnMapperPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "referral-on-update" {
		state, err = r.CreateReferralOnUpdatePlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "simple-to-external-bind" {
		state, err = r.CreateSimpleToExternalBindPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "snmp-subagent" {
		state, err = r.CreateSnmpSubagentPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "clean-up-inactive-pingfederate-persistent-sessions" {
		state, err = r.CreateCleanUpInactivePingfederatePersistentSessionsPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "composed-attribute" {
		state, err = r.CreateComposedAttributePlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "attribute-mapper" {
		state, err = r.CreateAttributeMapperPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "delay" {
		state, err = r.CreateDelayPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "clean-up-expired-pingfederate-persistent-sessions" {
		state, err = r.CreateCleanUpExpiredPingfederatePersistentSessionsPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "groovy-scripted" {
		state, err = r.CreateGroovyScriptedPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "pluggable-pass-through-authentication" {
		state, err = r.CreatePluggablePassThroughAuthenticationPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "referential-integrity" {
		state, err = r.CreateReferentialIntegrityPlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.ResourceType.ValueString() == "unique-attribute" {
		state, err = r.CreateUniqueAttributePlugin(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, *state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultPluginResourceModel
	if plan.ResourceType.ValueString() == "last-access-time" {
		readLastAccessTimePluginResponseDefault(ctx, readResponse.LastAccessTimePluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "stats-collector" {
		readStatsCollectorPluginResponseDefault(ctx, readResponse.StatsCollectorPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "internal-search-rate" {
		readInternalSearchRatePluginResponseDefault(ctx, readResponse.InternalSearchRatePluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "modifiable-password-policy-state" {
		readModifiablePasswordPolicyStatePluginResponseDefault(ctx, readResponse.ModifiablePasswordPolicyStatePluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "seven-bit-clean" {
		readSevenBitCleanPluginResponseDefault(ctx, readResponse.SevenBitCleanPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "clean-up-expired-pingfederate-persistent-access-grants" {
		readCleanUpExpiredPingfederatePersistentAccessGrantsPluginResponseDefault(ctx, readResponse.CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "periodic-gc" {
		readPeriodicGcPluginResponseDefault(ctx, readResponse.PeriodicGcPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "ping-one-pass-through-authentication" {
		readPingOnePassThroughAuthenticationPluginResponseDefault(ctx, readResponse.PingOnePassThroughAuthenticationPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "changelog-password-encryption" {
		readChangelogPasswordEncryptionPluginResponseDefault(ctx, readResponse.ChangelogPasswordEncryptionPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "processing-time-histogram" {
		readProcessingTimeHistogramPluginResponseDefault(ctx, readResponse.ProcessingTimeHistogramPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "search-shutdown" {
		readSearchShutdownPluginResponseDefault(ctx, readResponse.SearchShutdownPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "periodic-stats-logger" {
		readPeriodicStatsLoggerPluginResponseDefault(ctx, readResponse.PeriodicStatsLoggerPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "purge-expired-data" {
		readPurgeExpiredDataPluginResponseDefault(ctx, readResponse.PurgeExpiredDataPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "change-subscription-notification" {
		readChangeSubscriptionNotificationPluginResponseDefault(ctx, readResponse.ChangeSubscriptionNotificationPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "sub-operation-timing" {
		readSubOperationTimingPluginResponseDefault(ctx, readResponse.SubOperationTimingPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "third-party" {
		readThirdPartyPluginResponseDefault(ctx, readResponse.ThirdPartyPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "encrypt-attribute-values" {
		readEncryptAttributeValuesPluginResponseDefault(ctx, readResponse.EncryptAttributeValuesPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "pass-through-authentication" {
		readPassThroughAuthenticationPluginResponseDefault(ctx, readResponse.PassThroughAuthenticationPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "dn-mapper" {
		readDnMapperPluginResponseDefault(ctx, readResponse.DnMapperPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "monitor-history" {
		readMonitorHistoryPluginResponseDefault(ctx, readResponse.MonitorHistoryPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "referral-on-update" {
		readReferralOnUpdatePluginResponseDefault(ctx, readResponse.ReferralOnUpdatePluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "simple-to-external-bind" {
		readSimpleToExternalBindPluginResponseDefault(ctx, readResponse.SimpleToExternalBindPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "custom" {
		readCustomPluginResponseDefault(ctx, readResponse.CustomPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "snmp-subagent" {
		readSnmpSubagentPluginResponseDefault(ctx, readResponse.SnmpSubagentPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "password-policy-import" {
		readPasswordPolicyImportPluginResponseDefault(ctx, readResponse.PasswordPolicyImportPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "profiler" {
		readProfilerPluginResponseDefault(ctx, readResponse.ProfilerPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "clean-up-inactive-pingfederate-persistent-sessions" {
		readCleanUpInactivePingfederatePersistentSessionsPluginResponseDefault(ctx, readResponse.CleanUpInactivePingfederatePersistentSessionsPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "composed-attribute" {
		readComposedAttributePluginResponseDefault(ctx, readResponse.ComposedAttributePluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "ldap-result-code-tracker" {
		readLdapResultCodeTrackerPluginResponseDefault(ctx, readResponse.LdapResultCodeTrackerPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "attribute-mapper" {
		readAttributeMapperPluginResponseDefault(ctx, readResponse.AttributeMapperPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "delay" {
		readDelayPluginResponseDefault(ctx, readResponse.DelayPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "clean-up-expired-pingfederate-persistent-sessions" {
		readCleanUpExpiredPingfederatePersistentSessionsPluginResponseDefault(ctx, readResponse.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "groovy-scripted" {
		readGroovyScriptedPluginResponseDefault(ctx, readResponse.GroovyScriptedPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "last-mod" {
		readLastModPluginResponseDefault(ctx, readResponse.LastModPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "pluggable-pass-through-authentication" {
		readPluggablePassThroughAuthenticationPluginResponseDefault(ctx, readResponse.PluggablePassThroughAuthenticationPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "referential-integrity" {
		readReferentialIntegrityPluginResponseDefault(ctx, readResponse.ReferentialIntegrityPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.ResourceType.ValueString() == "unique-attribute" {
		readUniqueAttributePluginResponseDefault(ctx, readResponse.UniqueAttributePluginResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createPluginOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.ResourceType.ValueString() == "last-access-time" {
			readLastAccessTimePluginResponseDefault(ctx, updateResponse.LastAccessTimePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "stats-collector" {
			readStatsCollectorPluginResponseDefault(ctx, updateResponse.StatsCollectorPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "internal-search-rate" {
			readInternalSearchRatePluginResponseDefault(ctx, updateResponse.InternalSearchRatePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "modifiable-password-policy-state" {
			readModifiablePasswordPolicyStatePluginResponseDefault(ctx, updateResponse.ModifiablePasswordPolicyStatePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "seven-bit-clean" {
			readSevenBitCleanPluginResponseDefault(ctx, updateResponse.SevenBitCleanPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "clean-up-expired-pingfederate-persistent-access-grants" {
			readCleanUpExpiredPingfederatePersistentAccessGrantsPluginResponseDefault(ctx, updateResponse.CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "periodic-gc" {
			readPeriodicGcPluginResponseDefault(ctx, updateResponse.PeriodicGcPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "ping-one-pass-through-authentication" {
			readPingOnePassThroughAuthenticationPluginResponseDefault(ctx, updateResponse.PingOnePassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "changelog-password-encryption" {
			readChangelogPasswordEncryptionPluginResponseDefault(ctx, updateResponse.ChangelogPasswordEncryptionPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "processing-time-histogram" {
			readProcessingTimeHistogramPluginResponseDefault(ctx, updateResponse.ProcessingTimeHistogramPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "search-shutdown" {
			readSearchShutdownPluginResponseDefault(ctx, updateResponse.SearchShutdownPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "periodic-stats-logger" {
			readPeriodicStatsLoggerPluginResponseDefault(ctx, updateResponse.PeriodicStatsLoggerPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "purge-expired-data" {
			readPurgeExpiredDataPluginResponseDefault(ctx, updateResponse.PurgeExpiredDataPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "change-subscription-notification" {
			readChangeSubscriptionNotificationPluginResponseDefault(ctx, updateResponse.ChangeSubscriptionNotificationPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "sub-operation-timing" {
			readSubOperationTimingPluginResponseDefault(ctx, updateResponse.SubOperationTimingPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "third-party" {
			readThirdPartyPluginResponseDefault(ctx, updateResponse.ThirdPartyPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "encrypt-attribute-values" {
			readEncryptAttributeValuesPluginResponseDefault(ctx, updateResponse.EncryptAttributeValuesPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "pass-through-authentication" {
			readPassThroughAuthenticationPluginResponseDefault(ctx, updateResponse.PassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "dn-mapper" {
			readDnMapperPluginResponseDefault(ctx, updateResponse.DnMapperPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "monitor-history" {
			readMonitorHistoryPluginResponseDefault(ctx, updateResponse.MonitorHistoryPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "referral-on-update" {
			readReferralOnUpdatePluginResponseDefault(ctx, updateResponse.ReferralOnUpdatePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "simple-to-external-bind" {
			readSimpleToExternalBindPluginResponseDefault(ctx, updateResponse.SimpleToExternalBindPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "custom" {
			readCustomPluginResponseDefault(ctx, updateResponse.CustomPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "snmp-subagent" {
			readSnmpSubagentPluginResponseDefault(ctx, updateResponse.SnmpSubagentPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "password-policy-import" {
			readPasswordPolicyImportPluginResponseDefault(ctx, updateResponse.PasswordPolicyImportPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "profiler" {
			readProfilerPluginResponseDefault(ctx, updateResponse.ProfilerPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "clean-up-inactive-pingfederate-persistent-sessions" {
			readCleanUpInactivePingfederatePersistentSessionsPluginResponseDefault(ctx, updateResponse.CleanUpInactivePingfederatePersistentSessionsPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "composed-attribute" {
			readComposedAttributePluginResponseDefault(ctx, updateResponse.ComposedAttributePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "ldap-result-code-tracker" {
			readLdapResultCodeTrackerPluginResponseDefault(ctx, updateResponse.LdapResultCodeTrackerPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "attribute-mapper" {
			readAttributeMapperPluginResponseDefault(ctx, updateResponse.AttributeMapperPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "delay" {
			readDelayPluginResponseDefault(ctx, updateResponse.DelayPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "clean-up-expired-pingfederate-persistent-sessions" {
			readCleanUpExpiredPingfederatePersistentSessionsPluginResponseDefault(ctx, updateResponse.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "groovy-scripted" {
			readGroovyScriptedPluginResponseDefault(ctx, updateResponse.GroovyScriptedPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "last-mod" {
			readLastModPluginResponseDefault(ctx, updateResponse.LastModPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "pluggable-pass-through-authentication" {
			readPluggablePassThroughAuthenticationPluginResponseDefault(ctx, updateResponse.PluggablePassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "referential-integrity" {
			readReferentialIntegrityPluginResponseDefault(ctx, updateResponse.ReferentialIntegrityPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "unique-attribute" {
			readUniqueAttributePluginResponseDefault(ctx, updateResponse.UniqueAttributePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *pluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state pluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.InternalSearchRatePluginResponse != nil {
		readInternalSearchRatePluginResponse(ctx, readResponse.InternalSearchRatePluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ModifiablePasswordPolicyStatePluginResponse != nil {
		readModifiablePasswordPolicyStatePluginResponse(ctx, readResponse.ModifiablePasswordPolicyStatePluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SevenBitCleanPluginResponse != nil {
		readSevenBitCleanPluginResponse(ctx, readResponse.SevenBitCleanPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse != nil {
		readCleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse(ctx, readResponse.CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PeriodicGcPluginResponse != nil {
		readPeriodicGcPluginResponse(ctx, readResponse.PeriodicGcPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PingOnePassThroughAuthenticationPluginResponse != nil {
		readPingOnePassThroughAuthenticationPluginResponse(ctx, readResponse.PingOnePassThroughAuthenticationPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SearchShutdownPluginResponse != nil {
		readSearchShutdownPluginResponse(ctx, readResponse.SearchShutdownPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PeriodicStatsLoggerPluginResponse != nil {
		readPeriodicStatsLoggerPluginResponse(ctx, readResponse.PeriodicStatsLoggerPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PurgeExpiredDataPluginResponse != nil {
		readPurgeExpiredDataPluginResponse(ctx, readResponse.PurgeExpiredDataPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SubOperationTimingPluginResponse != nil {
		readSubOperationTimingPluginResponse(ctx, readResponse.SubOperationTimingPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPluginResponse != nil {
		readThirdPartyPluginResponse(ctx, readResponse.ThirdPartyPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PassThroughAuthenticationPluginResponse != nil {
		readPassThroughAuthenticationPluginResponse(ctx, readResponse.PassThroughAuthenticationPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DnMapperPluginResponse != nil {
		readDnMapperPluginResponse(ctx, readResponse.DnMapperPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ReferralOnUpdatePluginResponse != nil {
		readReferralOnUpdatePluginResponse(ctx, readResponse.ReferralOnUpdatePluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SimpleToExternalBindPluginResponse != nil {
		readSimpleToExternalBindPluginResponse(ctx, readResponse.SimpleToExternalBindPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SnmpSubagentPluginResponse != nil {
		readSnmpSubagentPluginResponse(ctx, readResponse.SnmpSubagentPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CleanUpInactivePingfederatePersistentSessionsPluginResponse != nil {
		readCleanUpInactivePingfederatePersistentSessionsPluginResponse(ctx, readResponse.CleanUpInactivePingfederatePersistentSessionsPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ComposedAttributePluginResponse != nil {
		readComposedAttributePluginResponse(ctx, readResponse.ComposedAttributePluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AttributeMapperPluginResponse != nil {
		readAttributeMapperPluginResponse(ctx, readResponse.AttributeMapperPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DelayPluginResponse != nil {
		readDelayPluginResponse(ctx, readResponse.DelayPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CleanUpExpiredPingfederatePersistentSessionsPluginResponse != nil {
		readCleanUpExpiredPingfederatePersistentSessionsPluginResponse(ctx, readResponse.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedPluginResponse != nil {
		readGroovyScriptedPluginResponse(ctx, readResponse.GroovyScriptedPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PluggablePassThroughAuthenticationPluginResponse != nil {
		readPluggablePassThroughAuthenticationPluginResponse(ctx, readResponse.PluggablePassThroughAuthenticationPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ReferentialIntegrityPluginResponse != nil {
		readReferentialIntegrityPluginResponse(ctx, readResponse.ReferentialIntegrityPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UniqueAttributePluginResponse != nil {
		readUniqueAttributePluginResponse(ctx, readResponse.UniqueAttributePluginResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *defaultPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state defaultPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.LastAccessTimePluginResponse != nil {
		readLastAccessTimePluginResponseDefault(ctx, readResponse.LastAccessTimePluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.StatsCollectorPluginResponse != nil {
		readStatsCollectorPluginResponseDefault(ctx, readResponse.StatsCollectorPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ChangelogPasswordEncryptionPluginResponse != nil {
		readChangelogPasswordEncryptionPluginResponseDefault(ctx, readResponse.ChangelogPasswordEncryptionPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ProcessingTimeHistogramPluginResponse != nil {
		readProcessingTimeHistogramPluginResponseDefault(ctx, readResponse.ProcessingTimeHistogramPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ChangeSubscriptionNotificationPluginResponse != nil {
		readChangeSubscriptionNotificationPluginResponseDefault(ctx, readResponse.ChangeSubscriptionNotificationPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.EncryptAttributeValuesPluginResponse != nil {
		readEncryptAttributeValuesPluginResponseDefault(ctx, readResponse.EncryptAttributeValuesPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MonitorHistoryPluginResponse != nil {
		readMonitorHistoryPluginResponseDefault(ctx, readResponse.MonitorHistoryPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomPluginResponse != nil {
		readCustomPluginResponseDefault(ctx, readResponse.CustomPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PasswordPolicyImportPluginResponse != nil {
		readPasswordPolicyImportPluginResponseDefault(ctx, readResponse.PasswordPolicyImportPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ProfilerPluginResponse != nil {
		readProfilerPluginResponseDefault(ctx, readResponse.ProfilerPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdapResultCodeTrackerPluginResponse != nil {
		readLdapResultCodeTrackerPluginResponseDefault(ctx, readResponse.LdapResultCodeTrackerPluginResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LastModPluginResponse != nil {
		readLastModPluginResponseDefault(ctx, readResponse.LastModPluginResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *pluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan pluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state pluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.ResourceType.ValueString() == "internal-search-rate" {
			readInternalSearchRatePluginResponse(ctx, updateResponse.InternalSearchRatePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "modifiable-password-policy-state" {
			readModifiablePasswordPolicyStatePluginResponse(ctx, updateResponse.ModifiablePasswordPolicyStatePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "seven-bit-clean" {
			readSevenBitCleanPluginResponse(ctx, updateResponse.SevenBitCleanPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "clean-up-expired-pingfederate-persistent-access-grants" {
			readCleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse(ctx, updateResponse.CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "periodic-gc" {
			readPeriodicGcPluginResponse(ctx, updateResponse.PeriodicGcPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "ping-one-pass-through-authentication" {
			readPingOnePassThroughAuthenticationPluginResponse(ctx, updateResponse.PingOnePassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "search-shutdown" {
			readSearchShutdownPluginResponse(ctx, updateResponse.SearchShutdownPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "periodic-stats-logger" {
			readPeriodicStatsLoggerPluginResponse(ctx, updateResponse.PeriodicStatsLoggerPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "purge-expired-data" {
			readPurgeExpiredDataPluginResponse(ctx, updateResponse.PurgeExpiredDataPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "sub-operation-timing" {
			readSubOperationTimingPluginResponse(ctx, updateResponse.SubOperationTimingPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "third-party" {
			readThirdPartyPluginResponse(ctx, updateResponse.ThirdPartyPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "pass-through-authentication" {
			readPassThroughAuthenticationPluginResponse(ctx, updateResponse.PassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "dn-mapper" {
			readDnMapperPluginResponse(ctx, updateResponse.DnMapperPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "referral-on-update" {
			readReferralOnUpdatePluginResponse(ctx, updateResponse.ReferralOnUpdatePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "simple-to-external-bind" {
			readSimpleToExternalBindPluginResponse(ctx, updateResponse.SimpleToExternalBindPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "snmp-subagent" {
			readSnmpSubagentPluginResponse(ctx, updateResponse.SnmpSubagentPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "clean-up-inactive-pingfederate-persistent-sessions" {
			readCleanUpInactivePingfederatePersistentSessionsPluginResponse(ctx, updateResponse.CleanUpInactivePingfederatePersistentSessionsPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "composed-attribute" {
			readComposedAttributePluginResponse(ctx, updateResponse.ComposedAttributePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "attribute-mapper" {
			readAttributeMapperPluginResponse(ctx, updateResponse.AttributeMapperPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "delay" {
			readDelayPluginResponse(ctx, updateResponse.DelayPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "clean-up-expired-pingfederate-persistent-sessions" {
			readCleanUpExpiredPingfederatePersistentSessionsPluginResponse(ctx, updateResponse.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "groovy-scripted" {
			readGroovyScriptedPluginResponse(ctx, updateResponse.GroovyScriptedPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "pluggable-pass-through-authentication" {
			readPluggablePassThroughAuthenticationPluginResponse(ctx, updateResponse.PluggablePassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "referential-integrity" {
			readReferentialIntegrityPluginResponse(ctx, updateResponse.ReferentialIntegrityPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "unique-attribute" {
			readUniqueAttributePluginResponse(ctx, updateResponse.UniqueAttributePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *defaultPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan defaultPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state defaultPluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createPluginOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.ResourceType.ValueString() == "last-access-time" {
			readLastAccessTimePluginResponseDefault(ctx, updateResponse.LastAccessTimePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "stats-collector" {
			readStatsCollectorPluginResponseDefault(ctx, updateResponse.StatsCollectorPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "internal-search-rate" {
			readInternalSearchRatePluginResponseDefault(ctx, updateResponse.InternalSearchRatePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "modifiable-password-policy-state" {
			readModifiablePasswordPolicyStatePluginResponseDefault(ctx, updateResponse.ModifiablePasswordPolicyStatePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "seven-bit-clean" {
			readSevenBitCleanPluginResponseDefault(ctx, updateResponse.SevenBitCleanPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "clean-up-expired-pingfederate-persistent-access-grants" {
			readCleanUpExpiredPingfederatePersistentAccessGrantsPluginResponseDefault(ctx, updateResponse.CleanUpExpiredPingfederatePersistentAccessGrantsPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "periodic-gc" {
			readPeriodicGcPluginResponseDefault(ctx, updateResponse.PeriodicGcPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "ping-one-pass-through-authentication" {
			readPingOnePassThroughAuthenticationPluginResponseDefault(ctx, updateResponse.PingOnePassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "changelog-password-encryption" {
			readChangelogPasswordEncryptionPluginResponseDefault(ctx, updateResponse.ChangelogPasswordEncryptionPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "processing-time-histogram" {
			readProcessingTimeHistogramPluginResponseDefault(ctx, updateResponse.ProcessingTimeHistogramPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "search-shutdown" {
			readSearchShutdownPluginResponseDefault(ctx, updateResponse.SearchShutdownPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "periodic-stats-logger" {
			readPeriodicStatsLoggerPluginResponseDefault(ctx, updateResponse.PeriodicStatsLoggerPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "purge-expired-data" {
			readPurgeExpiredDataPluginResponseDefault(ctx, updateResponse.PurgeExpiredDataPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "change-subscription-notification" {
			readChangeSubscriptionNotificationPluginResponseDefault(ctx, updateResponse.ChangeSubscriptionNotificationPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "sub-operation-timing" {
			readSubOperationTimingPluginResponseDefault(ctx, updateResponse.SubOperationTimingPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "third-party" {
			readThirdPartyPluginResponseDefault(ctx, updateResponse.ThirdPartyPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "encrypt-attribute-values" {
			readEncryptAttributeValuesPluginResponseDefault(ctx, updateResponse.EncryptAttributeValuesPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "pass-through-authentication" {
			readPassThroughAuthenticationPluginResponseDefault(ctx, updateResponse.PassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "dn-mapper" {
			readDnMapperPluginResponseDefault(ctx, updateResponse.DnMapperPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "monitor-history" {
			readMonitorHistoryPluginResponseDefault(ctx, updateResponse.MonitorHistoryPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "referral-on-update" {
			readReferralOnUpdatePluginResponseDefault(ctx, updateResponse.ReferralOnUpdatePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "simple-to-external-bind" {
			readSimpleToExternalBindPluginResponseDefault(ctx, updateResponse.SimpleToExternalBindPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "custom" {
			readCustomPluginResponseDefault(ctx, updateResponse.CustomPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "snmp-subagent" {
			readSnmpSubagentPluginResponseDefault(ctx, updateResponse.SnmpSubagentPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "password-policy-import" {
			readPasswordPolicyImportPluginResponseDefault(ctx, updateResponse.PasswordPolicyImportPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "profiler" {
			readProfilerPluginResponseDefault(ctx, updateResponse.ProfilerPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "clean-up-inactive-pingfederate-persistent-sessions" {
			readCleanUpInactivePingfederatePersistentSessionsPluginResponseDefault(ctx, updateResponse.CleanUpInactivePingfederatePersistentSessionsPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "composed-attribute" {
			readComposedAttributePluginResponseDefault(ctx, updateResponse.ComposedAttributePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "ldap-result-code-tracker" {
			readLdapResultCodeTrackerPluginResponseDefault(ctx, updateResponse.LdapResultCodeTrackerPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "attribute-mapper" {
			readAttributeMapperPluginResponseDefault(ctx, updateResponse.AttributeMapperPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "delay" {
			readDelayPluginResponseDefault(ctx, updateResponse.DelayPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "clean-up-expired-pingfederate-persistent-sessions" {
			readCleanUpExpiredPingfederatePersistentSessionsPluginResponseDefault(ctx, updateResponse.CleanUpExpiredPingfederatePersistentSessionsPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "groovy-scripted" {
			readGroovyScriptedPluginResponseDefault(ctx, updateResponse.GroovyScriptedPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "last-mod" {
			readLastModPluginResponseDefault(ctx, updateResponse.LastModPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "pluggable-pass-through-authentication" {
			readPluggablePassThroughAuthenticationPluginResponseDefault(ctx, updateResponse.PluggablePassThroughAuthenticationPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "referential-integrity" {
			readReferentialIntegrityPluginResponseDefault(ctx, updateResponse.ReferentialIntegrityPluginResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.ResourceType.ValueString() == "unique-attribute" {
			readUniqueAttributePluginResponseDefault(ctx, updateResponse.UniqueAttributePluginResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *pluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state pluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PluginApi.DeletePluginExecute(r.apiClient.PluginApi.DeletePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Plugin", err, httpResp)
		return
	}
}

func (r *pluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPlugin(ctx, req, resp)
}

func (r *defaultPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPlugin(ctx, req, resp)
}

func importPlugin(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}