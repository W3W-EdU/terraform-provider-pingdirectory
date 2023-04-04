package accountstatusnotificationhandler

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
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &thirdPartyAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithConfigure   = &thirdPartyAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithImportState = &thirdPartyAccountStatusNotificationHandlerResource{}
	_ resource.Resource                = &defaultThirdPartyAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultThirdPartyAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithImportState = &defaultThirdPartyAccountStatusNotificationHandlerResource{}
)

// Create a Third Party Account Status Notification Handler resource
func NewThirdPartyAccountStatusNotificationHandlerResource() resource.Resource {
	return &thirdPartyAccountStatusNotificationHandlerResource{}
}

func NewDefaultThirdPartyAccountStatusNotificationHandlerResource() resource.Resource {
	return &defaultThirdPartyAccountStatusNotificationHandlerResource{}
}

// thirdPartyAccountStatusNotificationHandlerResource is the resource implementation.
type thirdPartyAccountStatusNotificationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultThirdPartyAccountStatusNotificationHandlerResource is the resource implementation.
type defaultThirdPartyAccountStatusNotificationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *thirdPartyAccountStatusNotificationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_third_party_account_status_notification_handler"
}

func (r *defaultThirdPartyAccountStatusNotificationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_third_party_account_status_notification_handler"
}

// Configure adds the provider configured client to the resource.
func (r *thirdPartyAccountStatusNotificationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultThirdPartyAccountStatusNotificationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type thirdPartyAccountStatusNotificationHandlerResourceModel struct {
	Id                                         types.String `tfsdk:"id"`
	LastUpdated                                types.String `tfsdk:"last_updated"`
	Notifications                              types.Set    `tfsdk:"notifications"`
	RequiredActions                            types.Set    `tfsdk:"required_actions"`
	ExtensionClass                             types.String `tfsdk:"extension_class"`
	ExtensionArgument                          types.Set    `tfsdk:"extension_argument"`
	Description                                types.String `tfsdk:"description"`
	Enabled                                    types.Bool   `tfsdk:"enabled"`
	Asynchronous                               types.Bool   `tfsdk:"asynchronous"`
	AccountCreationNotificationRequestCriteria types.String `tfsdk:"account_creation_notification_request_criteria"`
	AccountUpdateNotificationRequestCriteria   types.String `tfsdk:"account_update_notification_request_criteria"`
}

// GetSchema defines the schema for the resource.
func (r *thirdPartyAccountStatusNotificationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	thirdPartyAccountStatusNotificationHandlerSchema(ctx, req, resp, false)
}

func (r *defaultThirdPartyAccountStatusNotificationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	thirdPartyAccountStatusNotificationHandlerSchema(ctx, req, resp, true)
}

func thirdPartyAccountStatusNotificationHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Third Party Account Status Notification Handler.",
		Attributes: map[string]schema.Attribute{
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Account Status Notification Handler.",
				Required:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Account Status Notification Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Account Status Notification Handler",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Account Status Notification Handler is enabled. Only enabled handlers are invoked whenever a related event occurs in the server.",
				Required:    true,
			},
			"asynchronous": schema.BoolAttribute{
				Description: "Indicates whether the server should attempt to invoke this Account Status Notification Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver a message) will not delay processing for the operation that triggered the notification.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"account_creation_notification_request_criteria": schema.StringAttribute{
				Description: "A request criteria object that identifies which add requests should result in account creation notifications for this handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_update_notification_request_criteria": schema.StringAttribute{
				Description: "A request criteria object that identifies which modify and modify DN requests should result in account update notifications for this handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalThirdPartyAccountStatusNotificationHandlerFields(ctx context.Context, addRequest *client.AddThirdPartyAccountStatusNotificationHandlerRequest, plan thirdPartyAccountStatusNotificationHandlerResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountCreationNotificationRequestCriteria) {
		addRequest.AccountCreationNotificationRequestCriteria = plan.AccountCreationNotificationRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountUpdateNotificationRequestCriteria) {
		addRequest.AccountUpdateNotificationRequestCriteria = plan.AccountUpdateNotificationRequestCriteria.ValueStringPointer()
	}
}

// Read a ThirdPartyAccountStatusNotificationHandlerResponse object into the model struct
func readThirdPartyAccountStatusNotificationHandlerResponse(ctx context.Context, r *client.ThirdPartyAccountStatusNotificationHandlerResponse, state *thirdPartyAccountStatusNotificationHandlerResourceModel, expectedValues *thirdPartyAccountStatusNotificationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountCreationNotificationRequestCriteria))
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountUpdateNotificationRequestCriteria))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createThirdPartyAccountStatusNotificationHandlerOperations(plan thirdPartyAccountStatusNotificationHandlerResourceModel, state thirdPartyAccountStatusNotificationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountCreationNotificationRequestCriteria, state.AccountCreationNotificationRequestCriteria, "account-creation-notification-request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountUpdateNotificationRequestCriteria, state.AccountUpdateNotificationRequestCriteria, "account-update-notification-request-criteria")
	return ops
}

// Create a new resource
func (r *thirdPartyAccountStatusNotificationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan thirdPartyAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddThirdPartyAccountStatusNotificationHandlerRequest(plan.Id.ValueString(),
		[]client.EnumthirdPartyAccountStatusNotificationHandlerSchemaUrn{client.ENUMTHIRDPARTYACCOUNTSTATUSNOTIFICATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ACCOUNT_STATUS_NOTIFICATION_HANDLERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyAccountStatusNotificationHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AccountStatusNotificationHandlerApi.AddAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAccountStatusNotificationHandlerRequest(
		client.AddThirdPartyAccountStatusNotificationHandlerRequestAsAddAccountStatusNotificationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.AddAccountStatusNotificationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Third Party Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state thirdPartyAccountStatusNotificationHandlerResourceModel
	readThirdPartyAccountStatusNotificationHandlerResponse(ctx, addResponse.ThirdPartyAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultThirdPartyAccountStatusNotificationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan thirdPartyAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Third Party Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state thirdPartyAccountStatusNotificationHandlerResourceModel
	readThirdPartyAccountStatusNotificationHandlerResponse(ctx, readResponse.ThirdPartyAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createThirdPartyAccountStatusNotificationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Third Party Account Status Notification Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readThirdPartyAccountStatusNotificationHandlerResponse(ctx, updateResponse.ThirdPartyAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *thirdPartyAccountStatusNotificationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readThirdPartyAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultThirdPartyAccountStatusNotificationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readThirdPartyAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readThirdPartyAccountStatusNotificationHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state thirdPartyAccountStatusNotificationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Third Party Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readThirdPartyAccountStatusNotificationHandlerResponse(ctx, readResponse.ThirdPartyAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *thirdPartyAccountStatusNotificationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateThirdPartyAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultThirdPartyAccountStatusNotificationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateThirdPartyAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateThirdPartyAccountStatusNotificationHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan thirdPartyAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state thirdPartyAccountStatusNotificationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createThirdPartyAccountStatusNotificationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Third Party Account Status Notification Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readThirdPartyAccountStatusNotificationHandlerResponse(ctx, updateResponse.ThirdPartyAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultThirdPartyAccountStatusNotificationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *thirdPartyAccountStatusNotificationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state thirdPartyAccountStatusNotificationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.DeleteAccountStatusNotificationHandlerExecute(r.apiClient.AccountStatusNotificationHandlerApi.DeleteAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Third Party Account Status Notification Handler", err, httpResp)
		return
	}
}

func (r *thirdPartyAccountStatusNotificationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importThirdPartyAccountStatusNotificationHandler(ctx, req, resp)
}

func (r *defaultThirdPartyAccountStatusNotificationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importThirdPartyAccountStatusNotificationHandler(ctx, req, resp)
}

func importThirdPartyAccountStatusNotificationHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
