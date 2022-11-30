package serverinstance

import (
	"context"
	"terraform-provider-pingdirectory/internal/utils"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdata-config-api-go-client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &authorizeServerInstanceResource{}
	_ resource.ResourceWithConfigure   = &authorizeServerInstanceResource{}
	_ resource.ResourceWithImportState = &authorizeServerInstanceResource{}
)

// Create a Authorize Server Instance resource
func NewAuthorizeServerInstanceResource() resource.Resource {
	return &authorizeServerInstanceResource{}
}

// authorizeServerInstanceResource is the resource implementation.
type authorizeServerInstanceResource struct {
	providerConfig utils.ProviderConfiguration
	apiClient      *client.APIClient
}

// authorizeServerInstanceResourceModel maps the resource schema data.
type authorizeServerInstanceResourceModel struct {
	ServerInstanceName        types.String `tfsdk:"server_instance_name"`
	ClusterName               types.String `tfsdk:"cluster_name"`
	ServerInstanceLocation    types.String `tfsdk:"server_instance_location"`
	Hostname                  types.String `tfsdk:"hostname"`
	ServerRoot                types.String `tfsdk:"server_root"`
	ServerVersion             types.String `tfsdk:"server_version"`
	InterServerCertificate    types.String `tfsdk:"inter_server_certificate"`
	LdapPort                  types.Int64  `tfsdk:"ldap_port"`
	LdapsPort                 types.Int64  `tfsdk:"ldaps_port"`
	HttpPort                  types.Int64  `tfsdk:"http_port"`
	HttpsPort                 types.Int64  `tfsdk:"https_port"`
	ReplicationPort           types.Int64  `tfsdk:"replication_port"`
	ReplicationServerID       types.Int64  `tfsdk:"replication_server_id"`
	ReplicationDomainServerID types.Set    `tfsdk:"replication_domain_server_id"`
	JmxPort                   types.Int64  `tfsdk:"jmx_port"`
	JmxsPort                  types.Int64  `tfsdk:"jmxs_port"`
	PreferredSecurity         types.String `tfsdk:"preferred_security"`
	StartTLSEnabled           types.Bool   `tfsdk:"start_tls_enabled"`
	BaseDN                    types.Set    `tfsdk:"base_dn"`
	MemberOfServerGroup       types.Set    `tfsdk:"member_of_server_group"`
	LastUpdated               types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *authorizeServerInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_server_instance"
}

// GetSchema defines the schema for the resource.
func (r *authorizeServerInstanceResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages a Authorize Server Instance.",
		Attributes: map[string]tfsdk.Attribute{
			// All are considered computed, since we are importing the existing server
			// instance from a server, rather than "creating" a server instance
			// like a typical Terraform resource.
			"server_instance_name": {
				Description: "The name of this Server Instance. The instance name needs to be unique if this server will be part of a topology of servers that are connected to each other. Once set, it may not be changed.",
				Type:        types.StringType,
				Required:    true,
			},
			"cluster_name": {
				Description: "The name of the cluster to which this Server Instance belongs. Server instances within the same cluster will share the same cluster-wide configuration.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"server_instance_location": {
				Description: "Specifies the location for the Server Instance.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"hostname": {
				Description: "The name of the host where this Server Instance is installed.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"server_root": {
				Description: "The file system path where this Server Instance is installed.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"server_version": {
				Description: "The version of the server.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"inter_server_certificate": {
				Description: "The public component of the certificate used by this instance to protect inter-server communication and to perform server-specific encryption. This will generally be managed by the server and should only be altered by administrators under explicit direction from Ping Identity support personnel.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"ldap_port": {
				Description: "The TCP port on which this server is listening for LDAP connections.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"ldaps_port": {
				Description: "The TCP port on which this server is listening for LDAP secure connections.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"http_port": {
				Description: "The TCP port on which this server is listening for HTTP connections.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"https_port": {
				Description: "The TCP port on which this server is listening for HTTPS connections.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"replication_port": {
				Description: "The replication TCP port.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"replication_server_id": {
				Description: "Specifies a unique identifier for the replication server on this server instance.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"replication_domain_server_id": {
				Description: "Specifies a unique identifier for the Directory Server within the replication domain.",
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
				Optional: true,
				Computed: true,
			},
			"jmx_port": {
				Description: "The TCP port on which this server is listening for JMX connections.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"jmxs_port": {
				Description: "The TCP port on which this server is listening for JMX secure connections.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"preferred_security": {
				Description: "Specifies the preferred mechanism to use for securing connections to the server.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"start_tls_enabled": {
				Description: "Indicates whether StartTLS is enabled on this server.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"base_dn": {
				Description: "The set of base DNs under the root DSE.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
				Computed: true,
			},
			"member_of_server_group": {
				Description: "The set of groups of which this server is a member.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
				Computed: true,
			},
			"last_updated": {
				Description: "Timestamp of the last Terraform update of the Server Instance.",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

// Configure adds the provider configured client to the resource.
func (r *authorizeServerInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(utils.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Create a new resource
// For server instances, create doesn't actually "create" anything - it "adopts" the existing
// server instance into management by terraform. This method reads the existing server instance
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *authorizeServerInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan authorizeServerInstanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	getResp, httpResp, err := r.apiClient.ServerInstanceApi.GetServerInstance(utils.BasicAuthContext(ctx, r.providerConfig), plan.ServerInstanceName.ValueString()).Execute()
	if err != nil {
		utils.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Instance", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := getResp.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read existing config
	var state authorizeServerInstanceResourceModel
	readAuthorizeServerInstanceResponse(getResp.AuthorizeServerInstanceResponse, &state)

	// Determine what changes need to be made to match the plan
	updateInstanceRequest := r.apiClient.ServerInstanceApi.UpdateServerInstance(utils.BasicAuthContext(ctx, r.providerConfig), plan.ServerInstanceName.ValueString())
	ops := createAuthorizeServerInstanceOperations(plan, state)

	if len(ops) > 0 {
		updateInstanceRequest = updateInstanceRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		utils.LogUpdateOperations(ctx, ops)
		instanceResp, httpResp, err := r.apiClient.ServerInstanceApi.UpdateServerInstanceExecute(updateInstanceRequest)
		if err != nil {
			utils.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Server Instance", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := instanceResp.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAuthorizeServerInstanceResponse(instanceResp.AuthorizeServerInstanceResponse, &plan)
		// Populate Computed attribute values
		plan.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		// Just put the initial read into the plan
		plan = state
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read a AuthorizeServerInstanceResponse object into the model struct.
// Use empty string for nils since everything is marked as computed.
func readAuthorizeServerInstanceResponse(r *client.AuthorizeServerInstanceResponse, state *authorizeServerInstanceResourceModel) {
	state.ServerInstanceName = types.StringValue(r.ServerInstanceName)
	state.ClusterName = types.StringValue(r.ClusterName)
	state.ServerInstanceLocation = utils.StringTypeOrNil(r.ServerInstanceLocation, true)
	state.Hostname = utils.StringTypeOrNil(r.Hostname, true)
	state.ServerRoot = utils.StringTypeOrNil(r.ServerRoot, true)
	state.ServerVersion = types.StringValue(r.ServerVersion)
	state.InterServerCertificate = utils.StringTypeOrNil(r.InterServerCertificate, true)
	state.LdapPort = utils.Int64TypeOrNil(r.LdapPort)
	state.LdapsPort = utils.Int64TypeOrNil(r.LdapsPort)
	state.HttpPort = utils.Int64TypeOrNil(r.HttpPort)
	state.HttpsPort = utils.Int64TypeOrNil(r.HttpsPort)
	state.ReplicationPort = utils.Int64TypeOrNil(r.ReplicationPort)
	state.ReplicationServerID = utils.Int64TypeOrNil(r.ReplicationServerID)
	state.ReplicationDomainServerID = utils.GetInt64Set(r.ReplicationDomainServerID)
	/*
		if r.ReplicationDomainServerID != nil {
			state.ReplicationDomainServerID = utils.GetInt64Set(*r.ReplicationDomainServerID)
		} else {
			state.ReplicationDomainServerID, _ = types.SetValue(types.Int64Type, []attr.Value{})
		}*/
	state.JmxPort = utils.Int64TypeOrNil(r.JmxPort)
	state.JmxsPort = utils.Int64TypeOrNil(r.JmxsPort)
	if r.PreferredSecurity != nil {
		state.PreferredSecurity = types.StringValue(string(*r.PreferredSecurity))
	} else {
		state.PreferredSecurity = types.StringValue("")
	}
	state.StartTLSEnabled = utils.BoolTypeOrNil(r.StartTLSEnabled)
	state.BaseDN = utils.GetStringSet(r.BaseDN)
	state.MemberOfServerGroup = utils.GetStringSet(r.MemberOfServerGroup)
}

// Read resource information
func (r *authorizeServerInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state authorizeServerInstanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverInstanceResponse, httpResp, err := r.apiClient.ServerInstanceApi.GetServerInstance(utils.BasicAuthContext(ctx, r.providerConfig), state.ServerInstanceName.ValueString()).Execute()
	if err != nil {
		utils.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Instance", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := serverInstanceResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAuthorizeServerInstanceResponse(serverInstanceResponse.AuthorizeServerInstanceResponse, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create any update operations necessary to make the state match the plan
func createAuthorizeServerInstanceOperations(plan authorizeServerInstanceResourceModel, state authorizeServerInstanceResourceModel) []client.Operation {
	var ops []client.Operation

	utils.AddStringOperationIfNecessary(&ops, plan.ClusterName, state.ClusterName, "cluster-name")
	utils.AddStringOperationIfNecessary(&ops, plan.ServerInstanceLocation, state.ServerInstanceLocation, "server-instance-location")
	utils.AddStringOperationIfNecessary(&ops, plan.Hostname, state.Hostname, "hostname")
	utils.AddStringOperationIfNecessary(&ops, plan.ServerRoot, state.ServerRoot, "server-root")
	utils.AddStringOperationIfNecessary(&ops, plan.ServerVersion, state.ServerVersion, "server-version")
	utils.AddStringOperationIfNecessary(&ops, plan.InterServerCertificate, state.InterServerCertificate, "inter-server-certificate")
	utils.AddInt64OperationIfNecessary(&ops, plan.LdapPort, state.LdapPort, "ldap-port")
	utils.AddInt64OperationIfNecessary(&ops, plan.LdapsPort, state.LdapsPort, "ldaps-port")
	utils.AddInt64OperationIfNecessary(&ops, plan.HttpPort, state.HttpPort, "http-port")
	utils.AddInt64OperationIfNecessary(&ops, plan.HttpsPort, state.HttpsPort, "https-port")
	utils.AddInt64OperationIfNecessary(&ops, plan.ReplicationPort, state.ReplicationPort, "replication-port")
	utils.AddInt64OperationIfNecessary(&ops, plan.ReplicationServerID, state.ReplicationServerID, "replication-server-id")
	utils.AddInt64SetOperationsIfNecessary(&ops, plan.ReplicationDomainServerID, state.ReplicationDomainServerID, "replication-domain-server-id")
	utils.AddInt64OperationIfNecessary(&ops, plan.JmxPort, state.JmxPort, "jmx-port")
	utils.AddInt64OperationIfNecessary(&ops, plan.JmxsPort, state.JmxsPort, "jmxs-port")
	utils.AddStringOperationIfNecessary(&ops, plan.PreferredSecurity, state.PreferredSecurity, "preferred-security")
	utils.AddBoolOperationIfNecessary(&ops, plan.StartTLSEnabled, state.StartTLSEnabled, "start-tls-enabled")
	utils.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	utils.AddStringSetOperationsIfNecessary(&ops, plan.MemberOfServerGroup, state.MemberOfServerGroup, "member-of-server-group")
	utils.AddStringOperationIfNecessary(&ops, plan.LastUpdated, state.LastUpdated, "last-updated")
	return ops
}

// Update a resource
func (r *authorizeServerInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan authorizeServerInstanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state authorizeServerInstanceResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ServerInstanceApi.UpdateServerInstance(utils.BasicAuthContext(ctx, r.providerConfig), plan.ServerInstanceName.ValueString())

	// Determine what update operations are necessary
	ops := createAuthorizeServerInstanceOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		utils.LogUpdateOperations(ctx, ops)

		serverInstanceResponse, httpResp, err := r.apiClient.ServerInstanceApi.UpdateServerInstanceExecute(updateRequest)
		if err != nil {
			utils.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Server Instance", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := serverInstanceResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAuthorizeServerInstanceResponse(serverInstanceResponse.AuthorizeServerInstanceResponse, &state)
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
// Terraform can't actually delete server instances, so this method does nothing.
// Terraform will just "forget" about the server instance config, and it can be managed elsewhere.
func (r *authorizeServerInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *authorizeServerInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to Name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}