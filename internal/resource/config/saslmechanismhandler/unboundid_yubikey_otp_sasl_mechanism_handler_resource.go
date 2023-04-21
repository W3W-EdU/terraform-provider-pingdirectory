package saslmechanismhandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &unboundidYubikeyOtpSaslMechanismHandlerResource{}
	_ resource.ResourceWithConfigure   = &unboundidYubikeyOtpSaslMechanismHandlerResource{}
	_ resource.ResourceWithImportState = &unboundidYubikeyOtpSaslMechanismHandlerResource{}
)

// Create a Unboundid Yubikey Otp Sasl Mechanism Handler resource
func NewUnboundidYubikeyOtpSaslMechanismHandlerResource() resource.Resource {
	return &unboundidYubikeyOtpSaslMechanismHandlerResource{}
}

// unboundidYubikeyOtpSaslMechanismHandlerResource is the resource implementation.
type unboundidYubikeyOtpSaslMechanismHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *unboundidYubikeyOtpSaslMechanismHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_unboundid_yubikey_otp_sasl_mechanism_handler"
}

// Configure adds the provider configured client to the resource.
func (r *unboundidYubikeyOtpSaslMechanismHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type unboundidYubikeyOtpSaslMechanismHandlerResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	LastUpdated                     types.String `tfsdk:"last_updated"`
	Notifications                   types.Set    `tfsdk:"notifications"`
	RequiredActions                 types.Set    `tfsdk:"required_actions"`
	YubikeyClientID                 types.String `tfsdk:"yubikey_client_id"`
	YubikeyAPIKey                   types.String `tfsdk:"yubikey_api_key"`
	YubikeyAPIKeyPassphraseProvider types.String `tfsdk:"yubikey_api_key_passphrase_provider"`
	YubikeyValidationServerBaseURL  types.Set    `tfsdk:"yubikey_validation_server_base_url"`
	HttpProxyExternalServer         types.String `tfsdk:"http_proxy_external_server"`
	IdentityMapper                  types.String `tfsdk:"identity_mapper"`
	RequireStaticPassword           types.Bool   `tfsdk:"require_static_password"`
	KeyManagerProvider              types.String `tfsdk:"key_manager_provider"`
	TrustManagerProvider            types.String `tfsdk:"trust_manager_provider"`
	Description                     types.String `tfsdk:"description"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *unboundidYubikeyOtpSaslMechanismHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Unboundid Yubikey Otp Sasl Mechanism Handler.",
		Attributes: map[string]schema.Attribute{
			"yubikey_client_id": schema.StringAttribute{
				Description: "The client ID to include in requests to the YubiKey validation server. A client ID and API key may be obtained for free from https://upgrade.yubico.com/getapikey/.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"yubikey_api_key": schema.StringAttribute{
				Description: "The API key needed to verify signatures generated by the YubiKey validation server. A client ID and API key may be obtained for free from https://upgrade.yubico.com/getapikey/.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Sensitive: true,
			},
			"yubikey_api_key_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider to use to obtain the API key needed to verify signatures generated by the YubiKey validation server. A client ID and API key may be obtained for free from https://upgrade.yubico.com/getapikey/.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"yubikey_validation_server_base_url": schema.SetAttribute{
				Description: "The base URL of the validation server to use to verify one-time passwords. You should only need to change the value if you wish to use your own validation server instead of using one of the Yubico servers. The server must use the YubiKey Validation Protocol version 2.0.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "A reference to an HTTP proxy server that should be used for requests sent to the YubiKey validation service. Supported in PingDirectory product version 9.2.0.0+.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_mapper": schema.StringAttribute{
				Description: "The identity mapper that should be used to identify the user(s) targeted in the authentication and/or authorization identities contained in the bind request. This will only be used for \"u:\"-style identities.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"require_static_password": schema.BoolAttribute{
				Description: "Indicates whether a user will be required to provide a static password when authenticating via the UNBOUNDID-YUBIKEY-OTP SASL mechanism.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"key_manager_provider": schema.StringAttribute{
				Description: "Specifies which key manager provider should be used to obtain a client certificate to present to the validation server when performing HTTPS communication. This may be left undefined if communication will not be secured with HTTPS, or if there is no need to present a client certificate to the validation service.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "Specifies which trust manager provider should be used to determine whether to trust the certificate presented by the server when performing HTTPS communication. This may be left undefined if HTTPS communication is not needed, or if the validation service presents a certificate that is trusted by the default JVM configuration (which should be the case for the validation servers that Yubico provides, but may not be the case if an alternate validation server is configured).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this SASL Mechanism Handler",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the SASL mechanism handler is enabled for use.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Validate that any version restrictions are met in the plan
func (r *unboundidYubikeyOtpSaslMechanismHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingDirectory9200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model unboundidYubikeyOtpSaslMechanismHandlerResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsNonEmptyString(model.HttpProxyExternalServer) {
		resp.Diagnostics.AddError("Attribute 'http_proxy_external_server' not supported by PingDirectory version "+r.providerConfig.ProductVersion, "")
	}
}

// Read a UnboundidYubikeyOtpSaslMechanismHandlerResponse object into the model struct
func readUnboundidYubikeyOtpSaslMechanismHandlerResponse(ctx context.Context, r *client.UnboundidYubikeyOtpSaslMechanismHandlerResponse, state *unboundidYubikeyOtpSaslMechanismHandlerResourceModel, expectedValues *unboundidYubikeyOtpSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.YubikeyClientID = internaltypes.StringTypeOrNil(r.YubikeyClientID, true)
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.YubikeyAPIKey = expectedValues.YubikeyAPIKey
	state.YubikeyAPIKeyPassphraseProvider = internaltypes.StringTypeOrNil(r.YubikeyAPIKeyPassphraseProvider, true)
	state.YubikeyValidationServerBaseURL = internaltypes.GetStringSet(r.YubikeyValidationServerBaseURL)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, true)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.RequireStaticPassword = internaltypes.BoolTypeOrNil(r.RequireStaticPassword)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, true)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createUnboundidYubikeyOtpSaslMechanismHandlerOperations(plan unboundidYubikeyOtpSaslMechanismHandlerResourceModel, state unboundidYubikeyOtpSaslMechanismHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.YubikeyClientID, state.YubikeyClientID, "yubikey-client-id")
	operations.AddStringOperationIfNecessary(&ops, plan.YubikeyAPIKey, state.YubikeyAPIKey, "yubikey-api-key")
	operations.AddStringOperationIfNecessary(&ops, plan.YubikeyAPIKeyPassphraseProvider, state.YubikeyAPIKeyPassphraseProvider, "yubikey-api-key-passphrase-provider")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.YubikeyValidationServerBaseURL, state.YubikeyValidationServerBaseURL, "yubikey-validation-server-base-url")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireStaticPassword, state.RequireStaticPassword, "require-static-password")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyManagerProvider, state.KeyManagerProvider, "key-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustManagerProvider, state.TrustManagerProvider, "trust-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *unboundidYubikeyOtpSaslMechanismHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan unboundidYubikeyOtpSaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.GetSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Unboundid Yubikey Otp Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state unboundidYubikeyOtpSaslMechanismHandlerResourceModel
	readUnboundidYubikeyOtpSaslMechanismHandlerResponse(ctx, readResponse.UnboundidYubikeyOtpSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createUnboundidYubikeyOtpSaslMechanismHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Unboundid Yubikey Otp Sasl Mechanism Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readUnboundidYubikeyOtpSaslMechanismHandlerResponse(ctx, updateResponse.UnboundidYubikeyOtpSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *unboundidYubikeyOtpSaslMechanismHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state unboundidYubikeyOtpSaslMechanismHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.GetSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Unboundid Yubikey Otp Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readUnboundidYubikeyOtpSaslMechanismHandlerResponse(ctx, readResponse.UnboundidYubikeyOtpSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *unboundidYubikeyOtpSaslMechanismHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan unboundidYubikeyOtpSaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state unboundidYubikeyOtpSaslMechanismHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createUnboundidYubikeyOtpSaslMechanismHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Unboundid Yubikey Otp Sasl Mechanism Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readUnboundidYubikeyOtpSaslMechanismHandlerResponse(ctx, updateResponse.UnboundidYubikeyOtpSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *unboundidYubikeyOtpSaslMechanismHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *unboundidYubikeyOtpSaslMechanismHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
