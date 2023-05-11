package azureauthenticationmethod_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

const testIdDefaultAzureAuthenticationMethod = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type defaultAzureAuthenticationMethodTestModel struct {
	id string
}

func TestAccDefaultAzureAuthenticationMethod(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := defaultAzureAuthenticationMethodTestModel{
		id: testIdDefaultAzureAuthenticationMethod,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckDefaultAzureAuthenticationMethodDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccDefaultAzureAuthenticationMethodResource(resourceName, initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccDefaultAzureAuthenticationMethodResource(resourceName, initialResourceModel),
				ResourceName:      "pingdirectory_default_azure_authentication_method." + resourceName,
				ImportStateId:     initialResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccDefaultAzureAuthenticationMethodResource(resourceName string, resourceModel defaultAzureAuthenticationMethodTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_azure_authentication_method" "%[1]s" {
  id = "%[2]s"
}`, resourceName,
		resourceModel.id)
}

// Test that any objects created by the test are destroyed
func testAccCheckDefaultAzureAuthenticationMethodDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.AzureAuthenticationMethodApi.GetAzureAuthenticationMethod(ctx, testIdDefaultAzureAuthenticationMethod).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Default Azure Authentication Method", testIdDefaultAzureAuthenticationMethod)
	}
	return nil
}