package matchingrule_test

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

const testIdMatchingRule = "Case Exact Ordering Matching Rule"

// Attributes to test with. Add optional properties to test here if desired.
type matchingRuleTestModel struct {
	id      string
	enabled bool
}

func TestAccMatchingRule(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := matchingRuleTestModel{
		id:      testIdMatchingRule,
		enabled: false,
	}
	updatedResourceModel := matchingRuleTestModel{
		id:      testIdMatchingRule,
		enabled: true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccMatchingRuleResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedMatchingRuleAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccMatchingRuleResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedMatchingRuleAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccMatchingRuleResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_default_matching_rule." + resourceName,
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

func testAccMatchingRuleResource(resourceName string, resourceModel matchingRuleTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_default_matching_rule" "%[1]s" {
  type    = "ordering"
  id      = "%[2]s"
  enabled = %[3]t
}`, resourceName,
		resourceModel.id,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedMatchingRuleAttributes(config matchingRuleTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.MatchingRuleApi.GetMatchingRule(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Matching Rule"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.OrderingMatchingRuleResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}