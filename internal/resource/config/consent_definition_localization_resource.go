package config

import (
	"context"
	"strings"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &consentDefinitionLocalizationResource{}
	_ resource.ResourceWithConfigure   = &consentDefinitionLocalizationResource{}
	_ resource.ResourceWithImportState = &consentDefinitionLocalizationResource{}
)

// Create a Consent Definition Localization resource
func NewConsentDefinitionLocalizationResource() resource.Resource {
	return &consentDefinitionLocalizationResource{}
}

// consentDefinitionLocalizationResource is the resource implementation.
type consentDefinitionLocalizationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *consentDefinitionLocalizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_consent_definition_localization"
}

// Configure adds the provider configured client to the resource.
func (r *consentDefinitionLocalizationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type consentDefinitionLocalizationResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	LastUpdated           types.String `tfsdk:"last_updated"`
	Notifications         types.Set    `tfsdk:"notifications"`
	RequiredActions       types.Set    `tfsdk:"required_actions"`
	ConsentDefinitionName types.String `tfsdk:"consent_definition_name"`
	Locale                types.String `tfsdk:"locale"`
	Version               types.String `tfsdk:"version"`
	TitleText             types.String `tfsdk:"title_text"`
	DataText              types.String `tfsdk:"data_text"`
	PurposeText           types.String `tfsdk:"purpose_text"`
}

// GetSchema defines the schema for the resource.
func (r *consentDefinitionLocalizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Consent Definition Localization.",
		Attributes: map[string]schema.Attribute{
			"consent_definition_name": schema.StringAttribute{
				Description: "Name of the parent Consent Definition",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"locale": schema.StringAttribute{
				Description: "The locale of this Consent Definition Localization.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.StringAttribute{
				Description: "The version of this Consent Definition Localization, using the format MAJOR.MINOR.",
				Required:    true,
			},
			"title_text": schema.StringAttribute{
				Description: "Localized text that may be used to provide a title or summary for a consent request or a granted consent.",
				Optional:    true,
			},
			"data_text": schema.StringAttribute{
				Description: "Localized text describing the data to be shared.",
				Required:    true,
			},
			"purpose_text": schema.StringAttribute{
				Description: "Localized text describing how the data is to be used.",
				Required:    true,
			},
		},
	}
	AddCommonSchema(&schema, false)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalConsentDefinitionLocalizationFields(ctx context.Context, addRequest *client.AddConsentDefinitionLocalizationRequest, plan consentDefinitionLocalizationResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TitleText) {
		stringVal := plan.TitleText.ValueString()
		addRequest.TitleText = &stringVal
	}
}

// Read a ConsentDefinitionLocalizationResponse object into the model struct
func readConsentDefinitionLocalizationResponse(ctx context.Context, r *client.ConsentDefinitionLocalizationResponse, state *consentDefinitionLocalizationResourceModel, expectedValues *consentDefinitionLocalizationResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ConsentDefinitionName = expectedValues.ConsentDefinitionName
	state.Locale = types.StringValue(r.Locale)
	state.Version = types.StringValue(r.Version)
	state.TitleText = internaltypes.StringTypeOrNil(r.TitleText, internaltypes.IsEmptyString(expectedValues.TitleText))
	state.DataText = types.StringValue(r.DataText)
	state.PurposeText = types.StringValue(r.PurposeText)
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createConsentDefinitionLocalizationOperations(plan consentDefinitionLocalizationResourceModel, state consentDefinitionLocalizationResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Locale, state.Locale, "locale")
	operations.AddStringOperationIfNecessary(&ops, plan.Version, state.Version, "version")
	operations.AddStringOperationIfNecessary(&ops, plan.TitleText, state.TitleText, "title-text")
	operations.AddStringOperationIfNecessary(&ops, plan.DataText, state.DataText, "data-text")
	operations.AddStringOperationIfNecessary(&ops, plan.PurposeText, state.PurposeText, "purpose-text")
	return ops
}

// Create a new resource
func (r *consentDefinitionLocalizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan consentDefinitionLocalizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddConsentDefinitionLocalizationRequest(plan.Locale.ValueString(),
		plan.Locale.ValueString(),
		plan.Version.ValueString(),
		plan.DataText.ValueString(),
		plan.PurposeText.ValueString())
	addOptionalConsentDefinitionLocalizationFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConsentDefinitionLocalizationApi.AddConsentDefinitionLocalization(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.ConsentDefinitionName.ValueString())
	apiAddRequest = apiAddRequest.AddConsentDefinitionLocalizationRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.ConsentDefinitionLocalizationApi.AddConsentDefinitionLocalizationExecute(apiAddRequest)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Consent Definition Localization", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state consentDefinitionLocalizationResourceModel
	readConsentDefinitionLocalizationResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *consentDefinitionLocalizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state consentDefinitionLocalizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConsentDefinitionLocalizationApi.GetConsentDefinitionLocalization(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Locale.ValueString(), state.ConsentDefinitionName.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Consent Definition Localization", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConsentDefinitionLocalizationResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *consentDefinitionLocalizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan consentDefinitionLocalizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state consentDefinitionLocalizationResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ConsentDefinitionLocalizationApi.UpdateConsentDefinitionLocalization(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Locale.ValueString(), plan.ConsentDefinitionName.ValueString())

	// Determine what update operations are necessary
	ops := createConsentDefinitionLocalizationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConsentDefinitionLocalizationApi.UpdateConsentDefinitionLocalizationExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Consent Definition Localization", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConsentDefinitionLocalizationResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *consentDefinitionLocalizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state consentDefinitionLocalizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ConsentDefinitionLocalizationApi.DeleteConsentDefinitionLocalizationExecute(r.apiClient.ConsentDefinitionLocalizationApi.DeleteConsentDefinitionLocalization(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Locale.ValueString(), state.ConsentDefinitionName.ValueString()))
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Consent Definition Localization", err, httpResp)
		return
	}
}

func (r *consentDefinitionLocalizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [consent-definition-name]/[consent-definition-localization-locale]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("consent_definition_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("locale"), split[1])...)
}
