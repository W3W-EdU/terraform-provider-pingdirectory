package passwordgenerator_test

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

const testIdPasswordGenerator = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type passwordGeneratorTestModel struct {
	id                   string
	passwordCharacterSet []string
	passwordFormat       string
	enabled              bool
}

func TestAccPasswordGenerator(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := passwordGeneratorTestModel{
		id:                   testIdPasswordGenerator,
		passwordCharacterSet: []string{"set:abcdefghijklmnopqrstuvwxyz1234567890"},
		passwordFormat:       "set:15",
		enabled:              true,
	}
	updatedResourceModel := passwordGeneratorTestModel{
		id:                   testIdPasswordGenerator,
		passwordCharacterSet: []string{"set:abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		passwordFormat:       "set:20",
		enabled:              false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckPasswordGeneratorDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccPasswordGeneratorResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedPasswordGeneratorAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccPasswordGeneratorResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedPasswordGeneratorAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccPasswordGeneratorResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_password_generator." + resourceName,
				ImportStateId:     updatedResourceModel.id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated",
				},
			},
		},
	})
}

func testAccPasswordGeneratorResource(resourceName string, resourceModel passwordGeneratorTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_password_generator" "%[1]s" {
  type                   = "random"
  id                     = "%[2]s"
  password_character_set = %[3]s
  password_format        = "%[4]s"
  enabled                = %[5]t
}`, resourceName,
		resourceModel.id,
		acctest.StringSliceToTerraformString(resourceModel.passwordCharacterSet),
		resourceModel.passwordFormat,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedPasswordGeneratorAttributes(config passwordGeneratorTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.PasswordGeneratorApi.GetPasswordGenerator(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Password Generator"
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "password-character-set",
			config.passwordCharacterSet, response.RandomPasswordGeneratorResponse.PasswordCharacterSet)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "password-format",
			config.passwordFormat, response.RandomPasswordGeneratorResponse.PasswordFormat)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.RandomPasswordGeneratorResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckPasswordGeneratorDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.PasswordGeneratorApi.GetPasswordGenerator(ctx, testIdPasswordGenerator).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Password Generator", testIdPasswordGenerator)
	}
	return nil
}