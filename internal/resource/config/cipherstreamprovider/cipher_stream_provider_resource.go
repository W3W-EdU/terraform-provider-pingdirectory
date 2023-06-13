package cipherstreamprovider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
	_ resource.Resource                = &cipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &cipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &cipherStreamProviderResource{}
	_ resource.Resource                = &defaultCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &defaultCipherStreamProviderResource{}
)

// Create a Cipher Stream Provider resource
func NewCipherStreamProviderResource() resource.Resource {
	return &cipherStreamProviderResource{}
}

func NewDefaultCipherStreamProviderResource() resource.Resource {
	return &defaultCipherStreamProviderResource{}
}

// cipherStreamProviderResource is the resource implementation.
type cipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultCipherStreamProviderResource is the resource implementation.
type defaultCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *cipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cipher_stream_provider"
}

func (r *defaultCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_cipher_stream_provider"
}

// Configure adds the provider configured client to the resource.
func (r *cipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type cipherStreamProviderResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	LastUpdated                     types.String `tfsdk:"last_updated"`
	Notifications                   types.Set    `tfsdk:"notifications"`
	RequiredActions                 types.Set    `tfsdk:"required_actions"`
	Type                            types.String `tfsdk:"type"`
	ExtensionClass                  types.String `tfsdk:"extension_class"`
	ExtensionArgument               types.Set    `tfsdk:"extension_argument"`
	VaultExternalServer             types.String `tfsdk:"vault_external_server"`
	VaultServerBaseURI              types.Set    `tfsdk:"vault_server_base_uri"`
	VaultAuthenticationMethod       types.String `tfsdk:"vault_authentication_method"`
	VaultSecretPath                 types.String `tfsdk:"vault_secret_path"`
	VaultSecretFieldName            types.String `tfsdk:"vault_secret_field_name"`
	VaultEncryptionMetadataFile     types.String `tfsdk:"vault_encryption_metadata_file"`
	TrustStoreFile                  types.String `tfsdk:"trust_store_file"`
	TrustStorePin                   types.String `tfsdk:"trust_store_pin"`
	TrustStoreType                  types.String `tfsdk:"trust_store_type"`
	Pkcs11ProviderClass             types.String `tfsdk:"pkcs11_provider_class"`
	Pkcs11ProviderConfigurationFile types.String `tfsdk:"pkcs11_provider_configuration_file"`
	KeyStorePin                     types.String `tfsdk:"key_store_pin"`
	KeyStorePinFile                 types.String `tfsdk:"key_store_pin_file"`
	KeyStorePinEnvironmentVariable  types.String `tfsdk:"key_store_pin_environment_variable"`
	Pkcs11KeyStoreType              types.String `tfsdk:"pkcs11_key_store_type"`
	SslCertNickname                 types.String `tfsdk:"ssl_cert_nickname"`
	ConjurExternalServer            types.String `tfsdk:"conjur_external_server"`
	ConjurSecretRelativePath        types.String `tfsdk:"conjur_secret_relative_path"`
	PasswordFile                    types.String `tfsdk:"password_file"`
	WaitForPasswordFile             types.Bool   `tfsdk:"wait_for_password_file"`
	KeyVaultURI                     types.String `tfsdk:"key_vault_uri"`
	AzureAuthenticationMethod       types.String `tfsdk:"azure_authentication_method"`
	HttpProxyExternalServer         types.String `tfsdk:"http_proxy_external_server"`
	SecretName                      types.String `tfsdk:"secret_name"`
	EncryptedPassphraseFile         types.String `tfsdk:"encrypted_passphrase_file"`
	SecretID                        types.String `tfsdk:"secret_id"`
	SecretFieldName                 types.String `tfsdk:"secret_field_name"`
	SecretVersionID                 types.String `tfsdk:"secret_version_id"`
	SecretVersionStage              types.String `tfsdk:"secret_version_stage"`
	EncryptionMetadataFile          types.String `tfsdk:"encryption_metadata_file"`
	AwsExternalServer               types.String `tfsdk:"aws_external_server"`
	AwsAccessKeyID                  types.String `tfsdk:"aws_access_key_id"`
	AwsSecretAccessKey              types.String `tfsdk:"aws_secret_access_key"`
	AwsRegionName                   types.String `tfsdk:"aws_region_name"`
	KmsEncryptionKeyArn             types.String `tfsdk:"kms_encryption_key_arn"`
	Description                     types.String `tfsdk:"description"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *cipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	cipherStreamProviderSchema(ctx, req, resp, false)
}

func (r *defaultCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	cipherStreamProviderSchema(ctx, req, resp, true)
}

func cipherStreamProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Cipher Stream Provider.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Cipher Stream Provider resource. Options are ['amazon-key-management-service', 'amazon-secrets-manager', 'azure-key-vault', 'file-based', 'wait-for-passphrase', 'conjur', 'pkcs11', 'vault', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"amazon-key-management-service", "amazon-secrets-manager", "azure-key-vault", "file-based", "wait-for-passphrase", "conjur", "pkcs11", "vault", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Cipher Stream Provider.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Cipher Stream Provider. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"vault_external_server": schema.StringAttribute{
				Description: "An external server definition with information needed to connect and authenticate to the Vault server.",
				Optional:    true,
			},
			"vault_server_base_uri": schema.SetAttribute{
				Description: "The base URL needed to access the Vault server. The base URL should consist of the protocol (\"http\" or \"https\"), the server address (resolvable name or IP address), and the port number. For example, \"https://vault.example.com:8200/\".",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"vault_authentication_method": schema.StringAttribute{
				Description: "The mechanism used to authenticate to the Vault server.",
				Optional:    true,
			},
			"vault_secret_path": schema.StringAttribute{
				Description: "The path to the desired secret in the Vault service. This will be appended to the value of the base-url property for the associated Vault external server.",
				Optional:    true,
			},
			"vault_secret_field_name": schema.StringAttribute{
				Description: "The name of the field in the Vault secret record that contains the passphrase to use to generate the encryption key.",
				Optional:    true,
			},
			"vault_encryption_metadata_file": schema.StringAttribute{
				Description: "The path to a file that will hold metadata about the encryption performed by this Vault Cipher Stream Provider.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trust_store_file": schema.StringAttribute{
				Description: "The path to a file containing the information needed to trust the certificate presented by the Vault servers.",
				Optional:    true,
			},
			"trust_store_pin": schema.StringAttribute{
				Description: "The passphrase needed to access the contents of the trust store. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents.",
				Optional:    true,
				Sensitive:   true,
			},
			"trust_store_type": schema.StringAttribute{
				Description: "The store type for the specified trust store file. The value should likely be one of \"JKS\" or \"PKCS12\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pkcs11_provider_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java security provider class that implements support for interacting with PKCS #11 tokens.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pkcs11_provider_configuration_file": schema.StringAttribute{
				Description: "The path to the file to use to configure the security provider that implements support for interacting with PKCS #11 tokens.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key_store_pin": schema.StringAttribute{
				Description: "The clear-text user PIN needed to interact with the PKCS #11 token.",
				Optional:    true,
				Sensitive:   true,
			},
			"key_store_pin_file": schema.StringAttribute{
				Description: "The path to a file containing the user PIN needed to interact with the PKCS #11 token. The file must exist and must contain exactly one line with a clear-text representation of the PIN.",
				Optional:    true,
			},
			"key_store_pin_environment_variable": schema.StringAttribute{
				Description: "The name of an environment variable whose value is the user PIN needed to interact with the PKCS #11 token. The environment variable must be defined and must contain a clear-text representation of the PIN.",
				Optional:    true,
			},
			"pkcs11_key_store_type": schema.StringAttribute{
				Description: "The key store type to use when obtaining an instance of a key store for interacting with a PKCS #11 token.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_cert_nickname": schema.StringAttribute{
				Description: "The alias for the certificate in the PKCS #11 token that will be used to wrap the encryption key. The target certificate must exist in the PKCS #11 token, and it must have an RSA key pair because the JVM does not currently provide adequate key wrapping support for elliptic curve key pairs.  If you have also configured the server to use a PKCS #11 token for accessing listener certificates, we strongly recommend that you use a different certificate to protect the contents of the encryption settings database than you use for negotiating TLS sessions with clients. It is imperative that the certificate used by this PKCS11 Cipher Stream Provider remain constant for the life of the provider because if the certificate were to be replaced, then the contents of the encryption settings database could become inaccessible. Unlike with listener certificates used for TLS negotiation that need to be replaced on a regular basis, this PKCS11 Cipher Stream Provider does not consider the validity period for the associated certificate, and it will continue to function even after the certificate has expired.  If you need to rotate the certificate used to protect the server's encryption settings database, you should first install the desired new certificate in the PKCS #11 token under a different alias. Then, you should create a new instance of this PKCS11 Cipher Stream Provider that is configured to use that certificate, and that also uses a different value for the encryption-metadata-file because the information in that file is tied to the certificate used to generate it. Finally, you will need to update the global configuration so that the encryption-settings-cipher-stream-provider property references the new cipher stream provider rather than this one. The update to the global configuration must be done with the server online so that it can properly re-encrypt the contents of the encryption settings database with the correct key tied to the new certificate.",
				Optional:    true,
			},
			"conjur_external_server": schema.StringAttribute{
				Description: "An external server definition with information needed to connect and authenticate to the Conjur server.",
				Optional:    true,
			},
			"conjur_secret_relative_path": schema.StringAttribute{
				Description: "The portion of the path that follows the account name in the URI needed to obtain the secret passphrase to use to generate the encryption key. Any special characters in the path must be URL-encoded.",
				Optional:    true,
			},
			"password_file": schema.StringAttribute{
				Description: "The path to the file containing the password to use when generating ciphers.",
				Optional:    true,
			},
			"wait_for_password_file": schema.BoolAttribute{
				Description: "Indicates whether the server should wait for the password file to become available if it does not exist.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"key_vault_uri": schema.StringAttribute{
				Description: "The URI that identifies the Azure Key Vault from which the secret is to be retrieved.",
				Optional:    true,
			},
			"azure_authentication_method": schema.StringAttribute{
				Description: "The mechanism used to authenticate to the Azure service.",
				Optional:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "A reference to an HTTP proxy server that should be used for requests sent to the Azure service.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secret_name": schema.StringAttribute{
				Description: "The name of the secret to retrieve.",
				Optional:    true,
			},
			"encrypted_passphrase_file": schema.StringAttribute{
				Description: "The path to a file that will hold the encrypted passphrase used by this cipher stream provider.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secret_id": schema.StringAttribute{
				Description: "The Amazon Resource Name (ARN) or the user-friendly name of the secret to be retrieved.",
				Optional:    true,
			},
			"secret_field_name": schema.StringAttribute{
				Description: "The name of the JSON field whose value is the passphrase that will be used to generate the encryption key for protecting the contents of the encryption settings database.",
				Optional:    true,
			},
			"secret_version_id": schema.StringAttribute{
				Description: "The unique identifier for the version of the secret to be retrieved.",
				Optional:    true,
			},
			"secret_version_stage": schema.StringAttribute{
				Description: "The staging label for the version of the secret to be retrieved.",
				Optional:    true,
			},
			"encryption_metadata_file": schema.StringAttribute{
				Description: "The path to a file that will hold metadata about the encryption performed by this Amazon Secrets Manager Cipher Stream Provider.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"aws_external_server": schema.StringAttribute{
				Description: "The external server with information to use when interacting with the Amazon Key Management Service.",
				Optional:    true,
			},
			"aws_access_key_id": schema.StringAttribute{
				Description: "The access key ID that will be used if this cipher stream provider will authenticate to the Amazon Key Management Service using an access key rather than an IAM role associated with an EC2 instance.",
				Optional:    true,
			},
			"aws_secret_access_key": schema.StringAttribute{
				Description: "The secret access key that will be used if this cipher stream provider will authenticate to the Amazon Key Management Service using an access key rather than an IAM role associated with an EC2 instance.",
				Optional:    true,
				Sensitive:   true,
			},
			"aws_region_name": schema.StringAttribute{
				Description: "The name of the Amazon Web Services region that holds the encryption key. This is optional, and if it is not provided, then the server will attempt to determine the region from the key ARN.",
				Optional:    true,
			},
			"kms_encryption_key_arn": schema.StringAttribute{
				Description: "The Amazon resource name (ARN) for the KMS key that will be used to encrypt the contents of the passphrase file. This key must exist, and the AWS client must have access to encrypt and decrypt data using this key.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Cipher Stream Provider",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Cipher Stream Provider is enabled for use in the Directory Server.",
				Required:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"amazon-key-management-service", "amazon-secrets-manager", "azure-key-vault", "file-based", "wait-for-passphrase", "conjur", "pkcs11", "vault", "third-party"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *cipherStreamProviderResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCipherStreamProviderResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanCipherStreamProvider(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model cipherStreamProviderResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.SecretFieldName) && model.Type.ValueString() != "amazon-secrets-manager" {
		resp.Diagnostics.AddError("Attribute 'secret_field_name' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'secret_field_name', the 'type' attribute must be one of ['amazon-secrets-manager']")
	}
	if internaltypes.IsDefined(model.VaultEncryptionMetadataFile) && model.Type.ValueString() != "vault" {
		resp.Diagnostics.AddError("Attribute 'vault_encryption_metadata_file' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'vault_encryption_metadata_file', the 'type' attribute must be one of ['vault']")
	}
	if internaltypes.IsDefined(model.TrustStorePin) && model.Type.ValueString() != "vault" {
		resp.Diagnostics.AddError("Attribute 'trust_store_pin' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'trust_store_pin', the 'type' attribute must be one of ['vault']")
	}
	if internaltypes.IsDefined(model.KeyStorePinEnvironmentVariable) && model.Type.ValueString() != "pkcs11" {
		resp.Diagnostics.AddError("Attribute 'key_store_pin_environment_variable' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'key_store_pin_environment_variable', the 'type' attribute must be one of ['pkcs11']")
	}
	if internaltypes.IsDefined(model.SecretVersionID) && model.Type.ValueString() != "amazon-secrets-manager" {
		resp.Diagnostics.AddError("Attribute 'secret_version_id' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'secret_version_id', the 'type' attribute must be one of ['amazon-secrets-manager']")
	}
	if internaltypes.IsDefined(model.AwsRegionName) && model.Type.ValueString() != "amazon-key-management-service" {
		resp.Diagnostics.AddError("Attribute 'aws_region_name' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'aws_region_name', the 'type' attribute must be one of ['amazon-key-management-service']")
	}
	if internaltypes.IsDefined(model.WaitForPasswordFile) && model.Type.ValueString() != "file-based" {
		resp.Diagnostics.AddError("Attribute 'wait_for_password_file' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'wait_for_password_file', the 'type' attribute must be one of ['file-based']")
	}
	if internaltypes.IsDefined(model.PasswordFile) && model.Type.ValueString() != "file-based" {
		resp.Diagnostics.AddError("Attribute 'password_file' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'password_file', the 'type' attribute must be one of ['file-based']")
	}
	if internaltypes.IsDefined(model.Pkcs11ProviderConfigurationFile) && model.Type.ValueString() != "pkcs11" {
		resp.Diagnostics.AddError("Attribute 'pkcs11_provider_configuration_file' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'pkcs11_provider_configuration_file', the 'type' attribute must be one of ['pkcs11']")
	}
	if internaltypes.IsDefined(model.VaultExternalServer) && model.Type.ValueString() != "vault" {
		resp.Diagnostics.AddError("Attribute 'vault_external_server' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'vault_external_server', the 'type' attribute must be one of ['vault']")
	}
	if internaltypes.IsDefined(model.KmsEncryptionKeyArn) && model.Type.ValueString() != "amazon-key-management-service" {
		resp.Diagnostics.AddError("Attribute 'kms_encryption_key_arn' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'kms_encryption_key_arn', the 'type' attribute must be one of ['amazon-key-management-service']")
	}
	if internaltypes.IsDefined(model.KeyVaultURI) && model.Type.ValueString() != "azure-key-vault" {
		resp.Diagnostics.AddError("Attribute 'key_vault_uri' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'key_vault_uri', the 'type' attribute must be one of ['azure-key-vault']")
	}
	if internaltypes.IsDefined(model.TrustStoreType) && model.Type.ValueString() != "vault" {
		resp.Diagnostics.AddError("Attribute 'trust_store_type' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'trust_store_type', the 'type' attribute must be one of ['vault']")
	}
	if internaltypes.IsDefined(model.ExtensionArgument) && model.Type.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_argument' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_argument', the 'type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.KeyStorePin) && model.Type.ValueString() != "pkcs11" {
		resp.Diagnostics.AddError("Attribute 'key_store_pin' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'key_store_pin', the 'type' attribute must be one of ['pkcs11']")
	}
	if internaltypes.IsDefined(model.TrustStoreFile) && model.Type.ValueString() != "vault" {
		resp.Diagnostics.AddError("Attribute 'trust_store_file' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'trust_store_file', the 'type' attribute must be one of ['vault']")
	}
	if internaltypes.IsDefined(model.AzureAuthenticationMethod) && model.Type.ValueString() != "azure-key-vault" {
		resp.Diagnostics.AddError("Attribute 'azure_authentication_method' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'azure_authentication_method', the 'type' attribute must be one of ['azure-key-vault']")
	}
	if internaltypes.IsDefined(model.SecretVersionStage) && model.Type.ValueString() != "amazon-secrets-manager" {
		resp.Diagnostics.AddError("Attribute 'secret_version_stage' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'secret_version_stage', the 'type' attribute must be one of ['amazon-secrets-manager']")
	}
	if internaltypes.IsDefined(model.ConjurExternalServer) && model.Type.ValueString() != "conjur" {
		resp.Diagnostics.AddError("Attribute 'conjur_external_server' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'conjur_external_server', the 'type' attribute must be one of ['conjur']")
	}
	if internaltypes.IsDefined(model.Pkcs11ProviderClass) && model.Type.ValueString() != "pkcs11" {
		resp.Diagnostics.AddError("Attribute 'pkcs11_provider_class' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'pkcs11_provider_class', the 'type' attribute must be one of ['pkcs11']")
	}
	if internaltypes.IsDefined(model.AwsExternalServer) && model.Type.ValueString() != "amazon-key-management-service" && model.Type.ValueString() != "amazon-secrets-manager" {
		resp.Diagnostics.AddError("Attribute 'aws_external_server' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'aws_external_server', the 'type' attribute must be one of ['amazon-key-management-service', 'amazon-secrets-manager']")
	}
	if internaltypes.IsDefined(model.EncryptionMetadataFile) && model.Type.ValueString() != "amazon-secrets-manager" && model.Type.ValueString() != "azure-key-vault" && model.Type.ValueString() != "conjur" && model.Type.ValueString() != "pkcs11" {
		resp.Diagnostics.AddError("Attribute 'encryption_metadata_file' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'encryption_metadata_file', the 'type' attribute must be one of ['amazon-secrets-manager', 'azure-key-vault', 'conjur', 'pkcs11']")
	}
	if internaltypes.IsDefined(model.Pkcs11KeyStoreType) && model.Type.ValueString() != "pkcs11" {
		resp.Diagnostics.AddError("Attribute 'pkcs11_key_store_type' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'pkcs11_key_store_type', the 'type' attribute must be one of ['pkcs11']")
	}
	if internaltypes.IsDefined(model.VaultSecretFieldName) && model.Type.ValueString() != "vault" {
		resp.Diagnostics.AddError("Attribute 'vault_secret_field_name' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'vault_secret_field_name', the 'type' attribute must be one of ['vault']")
	}
	if internaltypes.IsDefined(model.HttpProxyExternalServer) && model.Type.ValueString() != "azure-key-vault" {
		resp.Diagnostics.AddError("Attribute 'http_proxy_external_server' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'http_proxy_external_server', the 'type' attribute must be one of ['azure-key-vault']")
	}
	if internaltypes.IsDefined(model.VaultSecretPath) && model.Type.ValueString() != "vault" {
		resp.Diagnostics.AddError("Attribute 'vault_secret_path' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'vault_secret_path', the 'type' attribute must be one of ['vault']")
	}
	if internaltypes.IsDefined(model.SecretName) && model.Type.ValueString() != "azure-key-vault" {
		resp.Diagnostics.AddError("Attribute 'secret_name' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'secret_name', the 'type' attribute must be one of ['azure-key-vault']")
	}
	if internaltypes.IsDefined(model.VaultServerBaseURI) && model.Type.ValueString() != "vault" {
		resp.Diagnostics.AddError("Attribute 'vault_server_base_uri' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'vault_server_base_uri', the 'type' attribute must be one of ['vault']")
	}
	if internaltypes.IsDefined(model.SecretID) && model.Type.ValueString() != "amazon-secrets-manager" {
		resp.Diagnostics.AddError("Attribute 'secret_id' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'secret_id', the 'type' attribute must be one of ['amazon-secrets-manager']")
	}
	if internaltypes.IsDefined(model.AwsAccessKeyID) && model.Type.ValueString() != "amazon-key-management-service" {
		resp.Diagnostics.AddError("Attribute 'aws_access_key_id' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'aws_access_key_id', the 'type' attribute must be one of ['amazon-key-management-service']")
	}
	if internaltypes.IsDefined(model.AwsSecretAccessKey) && model.Type.ValueString() != "amazon-key-management-service" {
		resp.Diagnostics.AddError("Attribute 'aws_secret_access_key' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'aws_secret_access_key', the 'type' attribute must be one of ['amazon-key-management-service']")
	}
	if internaltypes.IsDefined(model.VaultAuthenticationMethod) && model.Type.ValueString() != "vault" {
		resp.Diagnostics.AddError("Attribute 'vault_authentication_method' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'vault_authentication_method', the 'type' attribute must be one of ['vault']")
	}
	if internaltypes.IsDefined(model.EncryptedPassphraseFile) && model.Type.ValueString() != "amazon-key-management-service" {
		resp.Diagnostics.AddError("Attribute 'encrypted_passphrase_file' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'encrypted_passphrase_file', the 'type' attribute must be one of ['amazon-key-management-service']")
	}
	if internaltypes.IsDefined(model.SslCertNickname) && model.Type.ValueString() != "pkcs11" {
		resp.Diagnostics.AddError("Attribute 'ssl_cert_nickname' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'ssl_cert_nickname', the 'type' attribute must be one of ['pkcs11']")
	}
	if internaltypes.IsDefined(model.ConjurSecretRelativePath) && model.Type.ValueString() != "conjur" {
		resp.Diagnostics.AddError("Attribute 'conjur_secret_relative_path' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'conjur_secret_relative_path', the 'type' attribute must be one of ['conjur']")
	}
	if internaltypes.IsDefined(model.ExtensionClass) && model.Type.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_class' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_class', the 'type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.KeyStorePinFile) && model.Type.ValueString() != "pkcs11" {
		resp.Diagnostics.AddError("Attribute 'key_store_pin_file' not supported by pingdirectory_cipher_stream_provider resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'key_store_pin_file', the 'type' attribute must be one of ['pkcs11']")
	}
}

// Add optional fields to create request for amazon-key-management-service cipher-stream-provider
func addOptionalAmazonKeyManagementServiceCipherStreamProviderFields(ctx context.Context, addRequest *client.AddAmazonKeyManagementServiceCipherStreamProviderRequest, plan cipherStreamProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptedPassphraseFile) {
		addRequest.EncryptedPassphraseFile = plan.EncryptedPassphraseFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AwsExternalServer) {
		addRequest.AwsExternalServer = plan.AwsExternalServer.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AwsAccessKeyID) {
		addRequest.AwsAccessKeyID = plan.AwsAccessKeyID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AwsSecretAccessKey) {
		addRequest.AwsSecretAccessKey = plan.AwsSecretAccessKey.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AwsRegionName) {
		addRequest.AwsRegionName = plan.AwsRegionName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for amazon-secrets-manager cipher-stream-provider
func addOptionalAmazonSecretsManagerCipherStreamProviderFields(ctx context.Context, addRequest *client.AddAmazonSecretsManagerCipherStreamProviderRequest, plan cipherStreamProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SecretVersionID) {
		addRequest.SecretVersionID = plan.SecretVersionID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SecretVersionStage) {
		addRequest.SecretVersionStage = plan.SecretVersionStage.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionMetadataFile) {
		addRequest.EncryptionMetadataFile = plan.EncryptionMetadataFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for azure-key-vault cipher-stream-provider
func addOptionalAzureKeyVaultCipherStreamProviderFields(ctx context.Context, addRequest *client.AddAzureKeyVaultCipherStreamProviderRequest, plan cipherStreamProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpProxyExternalServer) {
		addRequest.HttpProxyExternalServer = plan.HttpProxyExternalServer.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionMetadataFile) {
		addRequest.EncryptionMetadataFile = plan.EncryptionMetadataFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for file-based cipher-stream-provider
func addOptionalFileBasedCipherStreamProviderFields(ctx context.Context, addRequest *client.AddFileBasedCipherStreamProviderRequest, plan cipherStreamProviderResourceModel) {
	if internaltypes.IsDefined(plan.WaitForPasswordFile) {
		addRequest.WaitForPasswordFile = plan.WaitForPasswordFile.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for wait-for-passphrase cipher-stream-provider
func addOptionalWaitForPassphraseCipherStreamProviderFields(ctx context.Context, addRequest *client.AddWaitForPassphraseCipherStreamProviderRequest, plan cipherStreamProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for conjur cipher-stream-provider
func addOptionalConjurCipherStreamProviderFields(ctx context.Context, addRequest *client.AddConjurCipherStreamProviderRequest, plan cipherStreamProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionMetadataFile) {
		addRequest.EncryptionMetadataFile = plan.EncryptionMetadataFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for pkcs11 cipher-stream-provider
func addOptionalPkcs11CipherStreamProviderFields(ctx context.Context, addRequest *client.AddPkcs11CipherStreamProviderRequest, plan cipherStreamProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Pkcs11ProviderClass) {
		addRequest.Pkcs11ProviderClass = plan.Pkcs11ProviderClass.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Pkcs11ProviderConfigurationFile) {
		addRequest.Pkcs11ProviderConfigurationFile = plan.Pkcs11ProviderConfigurationFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyStorePin) {
		addRequest.KeyStorePin = plan.KeyStorePin.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyStorePinFile) {
		addRequest.KeyStorePinFile = plan.KeyStorePinFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyStorePinEnvironmentVariable) {
		addRequest.KeyStorePinEnvironmentVariable = plan.KeyStorePinEnvironmentVariable.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Pkcs11KeyStoreType) {
		addRequest.Pkcs11KeyStoreType = plan.Pkcs11KeyStoreType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionMetadataFile) {
		addRequest.EncryptionMetadataFile = plan.EncryptionMetadataFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for vault cipher-stream-provider
func addOptionalVaultCipherStreamProviderFields(ctx context.Context, addRequest *client.AddVaultCipherStreamProviderRequest, plan cipherStreamProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VaultExternalServer) {
		addRequest.VaultExternalServer = plan.VaultExternalServer.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.VaultServerBaseURI) {
		var slice []string
		plan.VaultServerBaseURI.ElementsAs(ctx, &slice, false)
		addRequest.VaultServerBaseURI = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VaultAuthenticationMethod) {
		addRequest.VaultAuthenticationMethod = plan.VaultAuthenticationMethod.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VaultEncryptionMetadataFile) {
		addRequest.VaultEncryptionMetadataFile = plan.VaultEncryptionMetadataFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStoreFile) {
		addRequest.TrustStoreFile = plan.TrustStoreFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStorePin) {
		addRequest.TrustStorePin = plan.TrustStorePin.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStoreType) {
		addRequest.TrustStoreType = plan.TrustStoreType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for third-party cipher-stream-provider
func addOptionalThirdPartyCipherStreamProviderFields(ctx context.Context, addRequest *client.AddThirdPartyCipherStreamProviderRequest, plan cipherStreamProviderResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any sets that have a nil ElementType, to avoid a nil pointer when setting the state
func populateCipherStreamProviderNilSets(ctx context.Context, model *cipherStreamProviderResourceModel) {
	if model.VaultServerBaseURI.ElementType(ctx) == nil {
		model.VaultServerBaseURI = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
}

// Read a AmazonKeyManagementServiceCipherStreamProviderResponse object into the model struct
func readAmazonKeyManagementServiceCipherStreamProviderResponse(ctx context.Context, r *client.AmazonKeyManagementServiceCipherStreamProviderResponse, state *cipherStreamProviderResourceModel, expectedValues *cipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("amazon-key-management-service")
	state.Id = types.StringValue(r.Id)
	state.EncryptedPassphraseFile = types.StringValue(r.EncryptedPassphraseFile)
	state.AwsExternalServer = internaltypes.StringTypeOrNil(r.AwsExternalServer, internaltypes.IsEmptyString(expectedValues.AwsExternalServer))
	state.AwsAccessKeyID = internaltypes.StringTypeOrNil(r.AwsAccessKeyID, internaltypes.IsEmptyString(expectedValues.AwsAccessKeyID))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.AwsSecretAccessKey = expectedValues.AwsSecretAccessKey
	state.AwsRegionName = internaltypes.StringTypeOrNil(r.AwsRegionName, internaltypes.IsEmptyString(expectedValues.AwsRegionName))
	state.KmsEncryptionKeyArn = types.StringValue(r.KmsEncryptionKeyArn)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCipherStreamProviderNilSets(ctx, state)
}

// Read a AmazonSecretsManagerCipherStreamProviderResponse object into the model struct
func readAmazonSecretsManagerCipherStreamProviderResponse(ctx context.Context, r *client.AmazonSecretsManagerCipherStreamProviderResponse, state *cipherStreamProviderResourceModel, expectedValues *cipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("amazon-secrets-manager")
	state.Id = types.StringValue(r.Id)
	state.AwsExternalServer = types.StringValue(r.AwsExternalServer)
	state.SecretID = types.StringValue(r.SecretID)
	state.SecretFieldName = types.StringValue(r.SecretFieldName)
	state.SecretVersionID = internaltypes.StringTypeOrNil(r.SecretVersionID, internaltypes.IsEmptyString(expectedValues.SecretVersionID))
	state.SecretVersionStage = internaltypes.StringTypeOrNil(r.SecretVersionStage, internaltypes.IsEmptyString(expectedValues.SecretVersionStage))
	state.EncryptionMetadataFile = types.StringValue(r.EncryptionMetadataFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCipherStreamProviderNilSets(ctx, state)
}

// Read a AzureKeyVaultCipherStreamProviderResponse object into the model struct
func readAzureKeyVaultCipherStreamProviderResponse(ctx context.Context, r *client.AzureKeyVaultCipherStreamProviderResponse, state *cipherStreamProviderResourceModel, expectedValues *cipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("azure-key-vault")
	state.Id = types.StringValue(r.Id)
	state.KeyVaultURI = types.StringValue(r.KeyVaultURI)
	state.AzureAuthenticationMethod = types.StringValue(r.AzureAuthenticationMethod)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, internaltypes.IsEmptyString(expectedValues.HttpProxyExternalServer))
	state.SecretName = types.StringValue(r.SecretName)
	state.EncryptionMetadataFile = types.StringValue(r.EncryptionMetadataFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCipherStreamProviderNilSets(ctx, state)
}

// Read a FileBasedCipherStreamProviderResponse object into the model struct
func readFileBasedCipherStreamProviderResponse(ctx context.Context, r *client.FileBasedCipherStreamProviderResponse, state *cipherStreamProviderResourceModel, expectedValues *cipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based")
	state.Id = types.StringValue(r.Id)
	state.PasswordFile = types.StringValue(r.PasswordFile)
	state.WaitForPasswordFile = internaltypes.BoolTypeOrNil(r.WaitForPasswordFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCipherStreamProviderNilSets(ctx, state)
}

// Read a WaitForPassphraseCipherStreamProviderResponse object into the model struct
func readWaitForPassphraseCipherStreamProviderResponse(ctx context.Context, r *client.WaitForPassphraseCipherStreamProviderResponse, state *cipherStreamProviderResourceModel, expectedValues *cipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("wait-for-passphrase")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCipherStreamProviderNilSets(ctx, state)
}

// Read a ConjurCipherStreamProviderResponse object into the model struct
func readConjurCipherStreamProviderResponse(ctx context.Context, r *client.ConjurCipherStreamProviderResponse, state *cipherStreamProviderResourceModel, expectedValues *cipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("conjur")
	state.Id = types.StringValue(r.Id)
	state.ConjurExternalServer = types.StringValue(r.ConjurExternalServer)
	state.ConjurSecretRelativePath = types.StringValue(r.ConjurSecretRelativePath)
	state.EncryptionMetadataFile = types.StringValue(r.EncryptionMetadataFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCipherStreamProviderNilSets(ctx, state)
}

// Read a Pkcs11CipherStreamProviderResponse object into the model struct
func readPkcs11CipherStreamProviderResponse(ctx context.Context, r *client.Pkcs11CipherStreamProviderResponse, state *cipherStreamProviderResourceModel, expectedValues *cipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("pkcs11")
	state.Id = types.StringValue(r.Id)
	state.Pkcs11ProviderClass = internaltypes.StringTypeOrNil(r.Pkcs11ProviderClass, internaltypes.IsEmptyString(expectedValues.Pkcs11ProviderClass))
	state.Pkcs11ProviderConfigurationFile = internaltypes.StringTypeOrNil(r.Pkcs11ProviderConfigurationFile, internaltypes.IsEmptyString(expectedValues.Pkcs11ProviderConfigurationFile))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.KeyStorePin = expectedValues.KeyStorePin
	state.KeyStorePinFile = internaltypes.StringTypeOrNil(r.KeyStorePinFile, internaltypes.IsEmptyString(expectedValues.KeyStorePinFile))
	state.KeyStorePinEnvironmentVariable = internaltypes.StringTypeOrNil(r.KeyStorePinEnvironmentVariable, internaltypes.IsEmptyString(expectedValues.KeyStorePinEnvironmentVariable))
	state.Pkcs11KeyStoreType = internaltypes.StringTypeOrNil(r.Pkcs11KeyStoreType, internaltypes.IsEmptyString(expectedValues.Pkcs11KeyStoreType))
	state.SslCertNickname = types.StringValue(r.SslCertNickname)
	state.EncryptionMetadataFile = types.StringValue(r.EncryptionMetadataFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCipherStreamProviderNilSets(ctx, state)
}

// Read a VaultCipherStreamProviderResponse object into the model struct
func readVaultCipherStreamProviderResponse(ctx context.Context, r *client.VaultCipherStreamProviderResponse, state *cipherStreamProviderResourceModel, expectedValues *cipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("vault")
	state.Id = types.StringValue(r.Id)
	state.VaultExternalServer = internaltypes.StringTypeOrNil(r.VaultExternalServer, internaltypes.IsEmptyString(expectedValues.VaultExternalServer))
	state.VaultServerBaseURI = internaltypes.GetStringSet(r.VaultServerBaseURI)
	state.VaultAuthenticationMethod = internaltypes.StringTypeOrNil(r.VaultAuthenticationMethod, internaltypes.IsEmptyString(expectedValues.VaultAuthenticationMethod))
	state.VaultSecretPath = types.StringValue(r.VaultSecretPath)
	state.VaultSecretFieldName = types.StringValue(r.VaultSecretFieldName)
	state.VaultEncryptionMetadataFile = types.StringValue(r.VaultEncryptionMetadataFile)
	state.TrustStoreFile = internaltypes.StringTypeOrNil(r.TrustStoreFile, internaltypes.IsEmptyString(expectedValues.TrustStoreFile))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.TrustStorePin = expectedValues.TrustStorePin
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, internaltypes.IsEmptyString(expectedValues.TrustStoreType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCipherStreamProviderNilSets(ctx, state)
}

// Read a ThirdPartyCipherStreamProviderResponse object into the model struct
func readThirdPartyCipherStreamProviderResponse(ctx context.Context, r *client.ThirdPartyCipherStreamProviderResponse, state *cipherStreamProviderResourceModel, expectedValues *cipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCipherStreamProviderNilSets(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createCipherStreamProviderOperations(plan cipherStreamProviderResourceModel, state cipherStreamProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultExternalServer, state.VaultExternalServer, "vault-external-server")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.VaultServerBaseURI, state.VaultServerBaseURI, "vault-server-base-uri")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultAuthenticationMethod, state.VaultAuthenticationMethod, "vault-authentication-method")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultSecretPath, state.VaultSecretPath, "vault-secret-path")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultSecretFieldName, state.VaultSecretFieldName, "vault-secret-field-name")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultEncryptionMetadataFile, state.VaultEncryptionMetadataFile, "vault-encryption-metadata-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreFile, state.TrustStoreFile, "trust-store-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStorePin, state.TrustStorePin, "trust-store-pin")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreType, state.TrustStoreType, "trust-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.Pkcs11ProviderClass, state.Pkcs11ProviderClass, "pkcs11-provider-class")
	operations.AddStringOperationIfNecessary(&ops, plan.Pkcs11ProviderConfigurationFile, state.Pkcs11ProviderConfigurationFile, "pkcs11-provider-configuration-file")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePin, state.KeyStorePin, "key-store-pin")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePinFile, state.KeyStorePinFile, "key-store-pin-file")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePinEnvironmentVariable, state.KeyStorePinEnvironmentVariable, "key-store-pin-environment-variable")
	operations.AddStringOperationIfNecessary(&ops, plan.Pkcs11KeyStoreType, state.Pkcs11KeyStoreType, "pkcs11-key-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.SslCertNickname, state.SslCertNickname, "ssl-cert-nickname")
	operations.AddStringOperationIfNecessary(&ops, plan.ConjurExternalServer, state.ConjurExternalServer, "conjur-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.ConjurSecretRelativePath, state.ConjurSecretRelativePath, "conjur-secret-relative-path")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordFile, state.PasswordFile, "password-file")
	operations.AddBoolOperationIfNecessary(&ops, plan.WaitForPasswordFile, state.WaitForPasswordFile, "wait-for-password-file")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyVaultURI, state.KeyVaultURI, "key-vault-uri")
	operations.AddStringOperationIfNecessary(&ops, plan.AzureAuthenticationMethod, state.AzureAuthenticationMethod, "azure-authentication-method")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretName, state.SecretName, "secret-name")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptedPassphraseFile, state.EncryptedPassphraseFile, "encrypted-passphrase-file")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretID, state.SecretID, "secret-id")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretFieldName, state.SecretFieldName, "secret-field-name")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretVersionID, state.SecretVersionID, "secret-version-id")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretVersionStage, state.SecretVersionStage, "secret-version-stage")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionMetadataFile, state.EncryptionMetadataFile, "encryption-metadata-file")
	operations.AddStringOperationIfNecessary(&ops, plan.AwsExternalServer, state.AwsExternalServer, "aws-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.AwsAccessKeyID, state.AwsAccessKeyID, "aws-access-key-id")
	operations.AddStringOperationIfNecessary(&ops, plan.AwsSecretAccessKey, state.AwsSecretAccessKey, "aws-secret-access-key")
	operations.AddStringOperationIfNecessary(&ops, plan.AwsRegionName, state.AwsRegionName, "aws-region-name")
	operations.AddStringOperationIfNecessary(&ops, plan.KmsEncryptionKeyArn, state.KmsEncryptionKeyArn, "kms-encryption-key-arn")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a amazon-key-management-service cipher-stream-provider
func (r *cipherStreamProviderResource) CreateAmazonKeyManagementServiceCipherStreamProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan cipherStreamProviderResourceModel) (*cipherStreamProviderResourceModel, error) {
	addRequest := client.NewAddAmazonKeyManagementServiceCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumamazonKeyManagementServiceCipherStreamProviderSchemaUrn{client.ENUMAMAZONKEYMANAGEMENTSERVICECIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERAMAZON_KEY_MANAGEMENT_SERVICE},
		plan.KmsEncryptionKeyArn.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalAmazonKeyManagementServiceCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddAmazonKeyManagementServiceCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Cipher Stream Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state cipherStreamProviderResourceModel
	readAmazonKeyManagementServiceCipherStreamProviderResponse(ctx, addResponse.AmazonKeyManagementServiceCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a amazon-secrets-manager cipher-stream-provider
func (r *cipherStreamProviderResource) CreateAmazonSecretsManagerCipherStreamProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan cipherStreamProviderResourceModel) (*cipherStreamProviderResourceModel, error) {
	addRequest := client.NewAddAmazonSecretsManagerCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumamazonSecretsManagerCipherStreamProviderSchemaUrn{client.ENUMAMAZONSECRETSMANAGERCIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERAMAZON_SECRETS_MANAGER},
		plan.AwsExternalServer.ValueString(),
		plan.SecretID.ValueString(),
		plan.SecretFieldName.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalAmazonSecretsManagerCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddAmazonSecretsManagerCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Cipher Stream Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state cipherStreamProviderResourceModel
	readAmazonSecretsManagerCipherStreamProviderResponse(ctx, addResponse.AmazonSecretsManagerCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a azure-key-vault cipher-stream-provider
func (r *cipherStreamProviderResource) CreateAzureKeyVaultCipherStreamProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan cipherStreamProviderResourceModel) (*cipherStreamProviderResourceModel, error) {
	addRequest := client.NewAddAzureKeyVaultCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumazureKeyVaultCipherStreamProviderSchemaUrn{client.ENUMAZUREKEYVAULTCIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERAZURE_KEY_VAULT},
		plan.KeyVaultURI.ValueString(),
		plan.AzureAuthenticationMethod.ValueString(),
		plan.SecretName.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalAzureKeyVaultCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddAzureKeyVaultCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Cipher Stream Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state cipherStreamProviderResourceModel
	readAzureKeyVaultCipherStreamProviderResponse(ctx, addResponse.AzureKeyVaultCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a file-based cipher-stream-provider
func (r *cipherStreamProviderResource) CreateFileBasedCipherStreamProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan cipherStreamProviderResourceModel) (*cipherStreamProviderResourceModel, error) {
	addRequest := client.NewAddFileBasedCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumfileBasedCipherStreamProviderSchemaUrn{client.ENUMFILEBASEDCIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERFILE_BASED},
		plan.PasswordFile.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalFileBasedCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddFileBasedCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Cipher Stream Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state cipherStreamProviderResourceModel
	readFileBasedCipherStreamProviderResponse(ctx, addResponse.FileBasedCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a wait-for-passphrase cipher-stream-provider
func (r *cipherStreamProviderResource) CreateWaitForPassphraseCipherStreamProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan cipherStreamProviderResourceModel) (*cipherStreamProviderResourceModel, error) {
	addRequest := client.NewAddWaitForPassphraseCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumwaitForPassphraseCipherStreamProviderSchemaUrn{client.ENUMWAITFORPASSPHRASECIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERWAIT_FOR_PASSPHRASE},
		plan.Enabled.ValueBool())
	addOptionalWaitForPassphraseCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddWaitForPassphraseCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Cipher Stream Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state cipherStreamProviderResourceModel
	readWaitForPassphraseCipherStreamProviderResponse(ctx, addResponse.WaitForPassphraseCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a conjur cipher-stream-provider
func (r *cipherStreamProviderResource) CreateConjurCipherStreamProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan cipherStreamProviderResourceModel) (*cipherStreamProviderResourceModel, error) {
	addRequest := client.NewAddConjurCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumconjurCipherStreamProviderSchemaUrn{client.ENUMCONJURCIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERCONJUR},
		plan.ConjurExternalServer.ValueString(),
		plan.ConjurSecretRelativePath.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalConjurCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddConjurCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Cipher Stream Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state cipherStreamProviderResourceModel
	readConjurCipherStreamProviderResponse(ctx, addResponse.ConjurCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a pkcs11 cipher-stream-provider
func (r *cipherStreamProviderResource) CreatePkcs11CipherStreamProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan cipherStreamProviderResourceModel) (*cipherStreamProviderResourceModel, error) {
	addRequest := client.NewAddPkcs11CipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.Enumpkcs11CipherStreamProviderSchemaUrn{client.ENUMPKCS11CIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERPKCS11},
		plan.SslCertNickname.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalPkcs11CipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddPkcs11CipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Cipher Stream Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state cipherStreamProviderResourceModel
	readPkcs11CipherStreamProviderResponse(ctx, addResponse.Pkcs11CipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a vault cipher-stream-provider
func (r *cipherStreamProviderResource) CreateVaultCipherStreamProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan cipherStreamProviderResourceModel) (*cipherStreamProviderResourceModel, error) {
	addRequest := client.NewAddVaultCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumvaultCipherStreamProviderSchemaUrn{client.ENUMVAULTCIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERVAULT},
		plan.VaultSecretPath.ValueString(),
		plan.VaultSecretFieldName.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalVaultCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddVaultCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Cipher Stream Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state cipherStreamProviderResourceModel
	readVaultCipherStreamProviderResponse(ctx, addResponse.VaultCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party cipher-stream-provider
func (r *cipherStreamProviderResource) CreateThirdPartyCipherStreamProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan cipherStreamProviderResourceModel) (*cipherStreamProviderResourceModel, error) {
	addRequest := client.NewAddThirdPartyCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumthirdPartyCipherStreamProviderSchemaUrn{client.ENUMTHIRDPARTYCIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddThirdPartyCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Cipher Stream Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state cipherStreamProviderResourceModel
	readThirdPartyCipherStreamProviderResponse(ctx, addResponse.ThirdPartyCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *cipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *cipherStreamProviderResourceModel
	var err error
	if plan.Type.ValueString() == "amazon-key-management-service" {
		state, err = r.CreateAmazonKeyManagementServiceCipherStreamProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "amazon-secrets-manager" {
		state, err = r.CreateAmazonSecretsManagerCipherStreamProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "azure-key-vault" {
		state, err = r.CreateAzureKeyVaultCipherStreamProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "file-based" {
		state, err = r.CreateFileBasedCipherStreamProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "wait-for-passphrase" {
		state, err = r.CreateWaitForPassphraseCipherStreamProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "conjur" {
		state, err = r.CreateConjurCipherStreamProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "pkcs11" {
		state, err = r.CreatePkcs11CipherStreamProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "vault" {
		state, err = r.CreateVaultCipherStreamProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyCipherStreamProvider(ctx, req, resp, plan)
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
func (r *defaultCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state cipherStreamProviderResourceModel
	if plan.Type.ValueString() == "amazon-key-management-service" {
		readAmazonKeyManagementServiceCipherStreamProviderResponse(ctx, readResponse.AmazonKeyManagementServiceCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "amazon-secrets-manager" {
		readAmazonSecretsManagerCipherStreamProviderResponse(ctx, readResponse.AmazonSecretsManagerCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "azure-key-vault" {
		readAzureKeyVaultCipherStreamProviderResponse(ctx, readResponse.AzureKeyVaultCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "file-based" {
		readFileBasedCipherStreamProviderResponse(ctx, readResponse.FileBasedCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "wait-for-passphrase" {
		readWaitForPassphraseCipherStreamProviderResponse(ctx, readResponse.WaitForPassphraseCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "conjur" {
		readConjurCipherStreamProviderResponse(ctx, readResponse.ConjurCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "pkcs11" {
		readPkcs11CipherStreamProviderResponse(ctx, readResponse.Pkcs11CipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "vault" {
		readVaultCipherStreamProviderResponse(ctx, readResponse.VaultCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "third-party" {
		readThirdPartyCipherStreamProviderResponse(ctx, readResponse.ThirdPartyCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "amazon-key-management-service" {
			readAmazonKeyManagementServiceCipherStreamProviderResponse(ctx, updateResponse.AmazonKeyManagementServiceCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "amazon-secrets-manager" {
			readAmazonSecretsManagerCipherStreamProviderResponse(ctx, updateResponse.AmazonSecretsManagerCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "azure-key-vault" {
			readAzureKeyVaultCipherStreamProviderResponse(ctx, updateResponse.AzureKeyVaultCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based" {
			readFileBasedCipherStreamProviderResponse(ctx, updateResponse.FileBasedCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "wait-for-passphrase" {
			readWaitForPassphraseCipherStreamProviderResponse(ctx, updateResponse.WaitForPassphraseCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "conjur" {
			readConjurCipherStreamProviderResponse(ctx, updateResponse.ConjurCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "pkcs11" {
			readPkcs11CipherStreamProviderResponse(ctx, updateResponse.Pkcs11CipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "vault" {
			readVaultCipherStreamProviderResponse(ctx, updateResponse.VaultCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party" {
			readThirdPartyCipherStreamProviderResponse(ctx, updateResponse.ThirdPartyCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *cipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readCipherStreamProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state cipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.AmazonKeyManagementServiceCipherStreamProviderResponse != nil {
		readAmazonKeyManagementServiceCipherStreamProviderResponse(ctx, readResponse.AmazonKeyManagementServiceCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AmazonSecretsManagerCipherStreamProviderResponse != nil {
		readAmazonSecretsManagerCipherStreamProviderResponse(ctx, readResponse.AmazonSecretsManagerCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AzureKeyVaultCipherStreamProviderResponse != nil {
		readAzureKeyVaultCipherStreamProviderResponse(ctx, readResponse.AzureKeyVaultCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedCipherStreamProviderResponse != nil {
		readFileBasedCipherStreamProviderResponse(ctx, readResponse.FileBasedCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.WaitForPassphraseCipherStreamProviderResponse != nil {
		readWaitForPassphraseCipherStreamProviderResponse(ctx, readResponse.WaitForPassphraseCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConjurCipherStreamProviderResponse != nil {
		readConjurCipherStreamProviderResponse(ctx, readResponse.ConjurCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Pkcs11CipherStreamProviderResponse != nil {
		readPkcs11CipherStreamProviderResponse(ctx, readResponse.Pkcs11CipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.VaultCipherStreamProviderResponse != nil {
		readVaultCipherStreamProviderResponse(ctx, readResponse.VaultCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyCipherStreamProviderResponse != nil {
		readThirdPartyCipherStreamProviderResponse(ctx, readResponse.ThirdPartyCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *cipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateCipherStreamProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan cipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state cipherStreamProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "amazon-key-management-service" {
			readAmazonKeyManagementServiceCipherStreamProviderResponse(ctx, updateResponse.AmazonKeyManagementServiceCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "amazon-secrets-manager" {
			readAmazonSecretsManagerCipherStreamProviderResponse(ctx, updateResponse.AmazonSecretsManagerCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "azure-key-vault" {
			readAzureKeyVaultCipherStreamProviderResponse(ctx, updateResponse.AzureKeyVaultCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based" {
			readFileBasedCipherStreamProviderResponse(ctx, updateResponse.FileBasedCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "wait-for-passphrase" {
			readWaitForPassphraseCipherStreamProviderResponse(ctx, updateResponse.WaitForPassphraseCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "conjur" {
			readConjurCipherStreamProviderResponse(ctx, updateResponse.ConjurCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "pkcs11" {
			readPkcs11CipherStreamProviderResponse(ctx, updateResponse.Pkcs11CipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "vault" {
			readVaultCipherStreamProviderResponse(ctx, updateResponse.VaultCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party" {
			readThirdPartyCipherStreamProviderResponse(ctx, updateResponse.ThirdPartyCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *cipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state cipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProviderExecute(r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Cipher Stream Provider", err, httpResp)
		return
	}
}

func (r *cipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCipherStreamProvider(ctx, req, resp)
}

func (r *defaultCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCipherStreamProvider(ctx, req, resp)
}

func importCipherStreamProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}