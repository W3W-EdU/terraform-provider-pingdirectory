package vaultauthenticationmethod

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	_ resource.Resource                = &vaultAuthenticationMethodResource{}
	_ resource.ResourceWithConfigure   = &vaultAuthenticationMethodResource{}
	_ resource.ResourceWithImportState = &vaultAuthenticationMethodResource{}
	_ resource.Resource                = &defaultVaultAuthenticationMethodResource{}
	_ resource.ResourceWithConfigure   = &defaultVaultAuthenticationMethodResource{}
	_ resource.ResourceWithImportState = &defaultVaultAuthenticationMethodResource{}
)

// Create a Vault Authentication Method resource
func NewVaultAuthenticationMethodResource() resource.Resource {
	return &vaultAuthenticationMethodResource{}
}

func NewDefaultVaultAuthenticationMethodResource() resource.Resource {
	return &defaultVaultAuthenticationMethodResource{}
}

// vaultAuthenticationMethodResource is the resource implementation.
type vaultAuthenticationMethodResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultVaultAuthenticationMethodResource is the resource implementation.
type defaultVaultAuthenticationMethodResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *vaultAuthenticationMethodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault_authentication_method"
}

func (r *defaultVaultAuthenticationMethodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_vault_authentication_method"
}

// Configure adds the provider configured client to the resource.
func (r *vaultAuthenticationMethodResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultVaultAuthenticationMethodResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type vaultAuthenticationMethodResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	LastUpdated        types.String `tfsdk:"last_updated"`
	Notifications      types.Set    `tfsdk:"notifications"`
	RequiredActions    types.Set    `tfsdk:"required_actions"`
	Type               types.String `tfsdk:"type"`
	Username           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	VaultRoleID        types.String `tfsdk:"vault_role_id"`
	VaultSecretID      types.String `tfsdk:"vault_secret_id"`
	LoginMechanismName types.String `tfsdk:"login_mechanism_name"`
	VaultAccessToken   types.String `tfsdk:"vault_access_token"`
	Description        types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *vaultAuthenticationMethodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	vaultAuthenticationMethodSchema(ctx, req, resp, false)
}

func (r *defaultVaultAuthenticationMethodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	vaultAuthenticationMethodSchema(ctx, req, resp, true)
}

func vaultAuthenticationMethodSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Vault Authentication Method.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Vault Authentication Method resource. Options are ['static-token', 'app-role', 'user-pass']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"static-token", "app-role", "user-pass"}...),
				},
			},
			"username": schema.StringAttribute{
				Description: "The username for the user to authenticate.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for the user to authenticate.",
				Optional:    true,
				Sensitive:   true,
			},
			"vault_role_id": schema.StringAttribute{
				Description: "The role ID for the AppRole to authenticate.",
				Optional:    true,
			},
			"vault_secret_id": schema.StringAttribute{
				Description: "The secret ID for the AppRole to authenticate.",
				Optional:    true,
				Sensitive:   true,
			},
			"login_mechanism_name": schema.StringAttribute{
				Description: "The name used when enabling the desired AppRole authentication mechanism in the Vault server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vault_access_token": schema.StringAttribute{
				Description: "The static token used to authenticate to the Vault server.",
				Optional:    true,
				Sensitive:   true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Vault Authentication Method",
				Optional:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"static-token", "app-role", "user-pass"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *vaultAuthenticationMethodResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanVaultAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultVaultAuthenticationMethodResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanVaultAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanVaultAuthenticationMethod(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model vaultAuthenticationMethodResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.VaultAccessToken) && model.Type.ValueString() != "static-token" {
		resp.Diagnostics.AddError("Attribute 'vault_access_token' not supported by pingdirectory_vault_authentication_method resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'vault_access_token', the 'type' attribute must be one of ['static-token']")
	}
	if internaltypes.IsDefined(model.LoginMechanismName) && model.Type.ValueString() != "app-role" && model.Type.ValueString() != "user-pass" {
		resp.Diagnostics.AddError("Attribute 'login_mechanism_name' not supported by pingdirectory_vault_authentication_method resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'login_mechanism_name', the 'type' attribute must be one of ['app-role', 'user-pass']")
	}
	if internaltypes.IsDefined(model.Password) && model.Type.ValueString() != "user-pass" {
		resp.Diagnostics.AddError("Attribute 'password' not supported by pingdirectory_vault_authentication_method resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'password', the 'type' attribute must be one of ['user-pass']")
	}
	if internaltypes.IsDefined(model.VaultRoleID) && model.Type.ValueString() != "app-role" {
		resp.Diagnostics.AddError("Attribute 'vault_role_id' not supported by pingdirectory_vault_authentication_method resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'vault_role_id', the 'type' attribute must be one of ['app-role']")
	}
	if internaltypes.IsDefined(model.VaultSecretID) && model.Type.ValueString() != "app-role" {
		resp.Diagnostics.AddError("Attribute 'vault_secret_id' not supported by pingdirectory_vault_authentication_method resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'vault_secret_id', the 'type' attribute must be one of ['app-role']")
	}
	if internaltypes.IsDefined(model.Username) && model.Type.ValueString() != "user-pass" {
		resp.Diagnostics.AddError("Attribute 'username' not supported by pingdirectory_vault_authentication_method resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'username', the 'type' attribute must be one of ['user-pass']")
	}
}

// Add optional fields to create request for static-token vault-authentication-method
func addOptionalStaticTokenVaultAuthenticationMethodFields(ctx context.Context, addRequest *client.AddStaticTokenVaultAuthenticationMethodRequest, plan vaultAuthenticationMethodResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for app-role vault-authentication-method
func addOptionalAppRoleVaultAuthenticationMethodFields(ctx context.Context, addRequest *client.AddAppRoleVaultAuthenticationMethodRequest, plan vaultAuthenticationMethodResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoginMechanismName) {
		addRequest.LoginMechanismName = plan.LoginMechanismName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for user-pass vault-authentication-method
func addOptionalUserPassVaultAuthenticationMethodFields(ctx context.Context, addRequest *client.AddUserPassVaultAuthenticationMethodRequest, plan vaultAuthenticationMethodResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoginMechanismName) {
		addRequest.LoginMechanismName = plan.LoginMechanismName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateVaultAuthenticationMethodUnknownValues(ctx context.Context, model *vaultAuthenticationMethodResourceModel) {
	if model.VaultAccessToken.IsUnknown() {
		model.VaultAccessToken = types.StringNull()
	}
	if model.VaultSecretID.IsUnknown() {
		model.VaultSecretID = types.StringNull()
	}
	if model.Password.IsUnknown() {
		model.Password = types.StringNull()
	}
}

// Read a StaticTokenVaultAuthenticationMethodResponse object into the model struct
func readStaticTokenVaultAuthenticationMethodResponse(ctx context.Context, r *client.StaticTokenVaultAuthenticationMethodResponse, state *vaultAuthenticationMethodResourceModel, expectedValues *vaultAuthenticationMethodResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("static-token")
	state.Id = types.StringValue(r.Id)
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.VaultAccessToken = expectedValues.VaultAccessToken
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVaultAuthenticationMethodUnknownValues(ctx, state)
}

// Read a AppRoleVaultAuthenticationMethodResponse object into the model struct
func readAppRoleVaultAuthenticationMethodResponse(ctx context.Context, r *client.AppRoleVaultAuthenticationMethodResponse, state *vaultAuthenticationMethodResourceModel, expectedValues *vaultAuthenticationMethodResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("app-role")
	state.Id = types.StringValue(r.Id)
	state.VaultRoleID = types.StringValue(r.VaultRoleID)
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.VaultSecretID = expectedValues.VaultSecretID
	state.LoginMechanismName = internaltypes.StringTypeOrNil(r.LoginMechanismName, internaltypes.IsEmptyString(expectedValues.LoginMechanismName))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVaultAuthenticationMethodUnknownValues(ctx, state)
}

// Read a UserPassVaultAuthenticationMethodResponse object into the model struct
func readUserPassVaultAuthenticationMethodResponse(ctx context.Context, r *client.UserPassVaultAuthenticationMethodResponse, state *vaultAuthenticationMethodResourceModel, expectedValues *vaultAuthenticationMethodResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("user-pass")
	state.Id = types.StringValue(r.Id)
	state.Username = types.StringValue(r.Username)
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.Password = expectedValues.Password
	state.LoginMechanismName = internaltypes.StringTypeOrNil(r.LoginMechanismName, internaltypes.IsEmptyString(expectedValues.LoginMechanismName))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVaultAuthenticationMethodUnknownValues(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createVaultAuthenticationMethodOperations(plan vaultAuthenticationMethodResourceModel, state vaultAuthenticationMethodResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Username, state.Username, "username")
	operations.AddStringOperationIfNecessary(&ops, plan.Password, state.Password, "password")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultRoleID, state.VaultRoleID, "vault-role-id")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultSecretID, state.VaultSecretID, "vault-secret-id")
	operations.AddStringOperationIfNecessary(&ops, plan.LoginMechanismName, state.LoginMechanismName, "login-mechanism-name")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultAccessToken, state.VaultAccessToken, "vault-access-token")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a static-token vault-authentication-method
func (r *vaultAuthenticationMethodResource) CreateStaticTokenVaultAuthenticationMethod(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan vaultAuthenticationMethodResourceModel) (*vaultAuthenticationMethodResourceModel, error) {
	addRequest := client.NewAddStaticTokenVaultAuthenticationMethodRequest(plan.Id.ValueString(),
		[]client.EnumstaticTokenVaultAuthenticationMethodSchemaUrn{client.ENUMSTATICTOKENVAULTAUTHENTICATIONMETHODSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VAULT_AUTHENTICATION_METHODSTATIC_TOKEN},
		plan.VaultAccessToken.ValueString())
	addOptionalStaticTokenVaultAuthenticationMethodFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VaultAuthenticationMethodApi.AddVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVaultAuthenticationMethodRequest(
		client.AddStaticTokenVaultAuthenticationMethodRequestAsAddVaultAuthenticationMethodRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VaultAuthenticationMethodApi.AddVaultAuthenticationMethodExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Vault Authentication Method", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state vaultAuthenticationMethodResourceModel
	readStaticTokenVaultAuthenticationMethodResponse(ctx, addResponse.StaticTokenVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a app-role vault-authentication-method
func (r *vaultAuthenticationMethodResource) CreateAppRoleVaultAuthenticationMethod(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan vaultAuthenticationMethodResourceModel) (*vaultAuthenticationMethodResourceModel, error) {
	addRequest := client.NewAddAppRoleVaultAuthenticationMethodRequest(plan.Id.ValueString(),
		[]client.EnumappRoleVaultAuthenticationMethodSchemaUrn{client.ENUMAPPROLEVAULTAUTHENTICATIONMETHODSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VAULT_AUTHENTICATION_METHODAPP_ROLE},
		plan.VaultRoleID.ValueString(),
		plan.VaultSecretID.ValueString())
	addOptionalAppRoleVaultAuthenticationMethodFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VaultAuthenticationMethodApi.AddVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVaultAuthenticationMethodRequest(
		client.AddAppRoleVaultAuthenticationMethodRequestAsAddVaultAuthenticationMethodRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VaultAuthenticationMethodApi.AddVaultAuthenticationMethodExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Vault Authentication Method", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state vaultAuthenticationMethodResourceModel
	readAppRoleVaultAuthenticationMethodResponse(ctx, addResponse.AppRoleVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a user-pass vault-authentication-method
func (r *vaultAuthenticationMethodResource) CreateUserPassVaultAuthenticationMethod(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan vaultAuthenticationMethodResourceModel) (*vaultAuthenticationMethodResourceModel, error) {
	addRequest := client.NewAddUserPassVaultAuthenticationMethodRequest(plan.Id.ValueString(),
		[]client.EnumuserPassVaultAuthenticationMethodSchemaUrn{client.ENUMUSERPASSVAULTAUTHENTICATIONMETHODSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VAULT_AUTHENTICATION_METHODUSER_PASS},
		plan.Username.ValueString(),
		plan.Password.ValueString())
	addOptionalUserPassVaultAuthenticationMethodFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VaultAuthenticationMethodApi.AddVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVaultAuthenticationMethodRequest(
		client.AddUserPassVaultAuthenticationMethodRequestAsAddVaultAuthenticationMethodRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VaultAuthenticationMethodApi.AddVaultAuthenticationMethodExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Vault Authentication Method", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state vaultAuthenticationMethodResourceModel
	readUserPassVaultAuthenticationMethodResponse(ctx, addResponse.UserPassVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *vaultAuthenticationMethodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vaultAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *vaultAuthenticationMethodResourceModel
	var err error
	if plan.Type.ValueString() == "static-token" {
		state, err = r.CreateStaticTokenVaultAuthenticationMethod(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "app-role" {
		state, err = r.CreateAppRoleVaultAuthenticationMethod(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "user-pass" {
		state, err = r.CreateUserPassVaultAuthenticationMethod(ctx, req, resp, plan)
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
func (r *defaultVaultAuthenticationMethodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vaultAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.VaultAuthenticationMethodApi.GetVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Vault Authentication Method", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state vaultAuthenticationMethodResourceModel
	if plan.Type.ValueString() == "static-token" {
		readStaticTokenVaultAuthenticationMethodResponse(ctx, readResponse.StaticTokenVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "app-role" {
		readAppRoleVaultAuthenticationMethodResponse(ctx, readResponse.AppRoleVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "user-pass" {
		readUserPassVaultAuthenticationMethodResponse(ctx, readResponse.UserPassVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.VaultAuthenticationMethodApi.UpdateVaultAuthenticationMethod(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createVaultAuthenticationMethodOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.VaultAuthenticationMethodApi.UpdateVaultAuthenticationMethodExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Vault Authentication Method", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "static-token" {
			readStaticTokenVaultAuthenticationMethodResponse(ctx, updateResponse.StaticTokenVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "app-role" {
			readAppRoleVaultAuthenticationMethodResponse(ctx, updateResponse.AppRoleVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "user-pass" {
			readUserPassVaultAuthenticationMethodResponse(ctx, updateResponse.UserPassVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
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
func (r *vaultAuthenticationMethodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readVaultAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultVaultAuthenticationMethodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readVaultAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readVaultAuthenticationMethod(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state vaultAuthenticationMethodResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.VaultAuthenticationMethodApi.GetVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Vault Authentication Method", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.StaticTokenVaultAuthenticationMethodResponse != nil {
		readStaticTokenVaultAuthenticationMethodResponse(ctx, readResponse.StaticTokenVaultAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AppRoleVaultAuthenticationMethodResponse != nil {
		readAppRoleVaultAuthenticationMethodResponse(ctx, readResponse.AppRoleVaultAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UserPassVaultAuthenticationMethodResponse != nil {
		readUserPassVaultAuthenticationMethodResponse(ctx, readResponse.UserPassVaultAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *vaultAuthenticationMethodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateVaultAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultVaultAuthenticationMethodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateVaultAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateVaultAuthenticationMethod(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan vaultAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state vaultAuthenticationMethodResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.VaultAuthenticationMethodApi.UpdateVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createVaultAuthenticationMethodOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.VaultAuthenticationMethodApi.UpdateVaultAuthenticationMethodExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Vault Authentication Method", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "static-token" {
			readStaticTokenVaultAuthenticationMethodResponse(ctx, updateResponse.StaticTokenVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "app-role" {
			readAppRoleVaultAuthenticationMethodResponse(ctx, updateResponse.AppRoleVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "user-pass" {
			readUserPassVaultAuthenticationMethodResponse(ctx, updateResponse.UserPassVaultAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultVaultAuthenticationMethodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *vaultAuthenticationMethodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vaultAuthenticationMethodResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.VaultAuthenticationMethodApi.DeleteVaultAuthenticationMethodExecute(r.apiClient.VaultAuthenticationMethodApi.DeleteVaultAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Vault Authentication Method", err, httpResp)
		return
	}
}

func (r *vaultAuthenticationMethodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importVaultAuthenticationMethod(ctx, req, resp)
}

func (r *defaultVaultAuthenticationMethodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importVaultAuthenticationMethod(ctx, req, resp)
}

func importVaultAuthenticationMethod(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}