package logfilerotationlistener_test

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

const testIdSummarizeLogFileRotationListener = "MyId"

// Attributes to test with. Add optional properties to test here if desired.
type summarizeLogFileRotationListenerTestModel struct {
	id      string
	enabled bool
}

func TestAccSummarizeLogFileRotationListener(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := summarizeLogFileRotationListenerTestModel{
		id:      testIdSummarizeLogFileRotationListener,
		enabled: true,
	}
	updatedResourceModel := summarizeLogFileRotationListenerTestModel{
		id:      testIdSummarizeLogFileRotationListener,
		enabled: false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingdirectory": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckSummarizeLogFileRotationListenerDestroy,
		Steps: []resource.TestStep{
			{
				// Test basic resource.
				// Add checks for computed properties here if desired.
				Config: testAccSummarizeLogFileRotationListenerResource(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSummarizeLogFileRotationListenerAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccSummarizeLogFileRotationListenerResource(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSummarizeLogFileRotationListenerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSummarizeLogFileRotationListenerResource(resourceName, updatedResourceModel),
				ResourceName:      "pingdirectory_log_file_rotation_listener." + resourceName,
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

func testAccSummarizeLogFileRotationListenerResource(resourceName string, resourceModel summarizeLogFileRotationListenerTestModel) string {
	return fmt.Sprintf(`
resource "pingdirectory_log_file_rotation_listener" "%[1]s" {
  type    = "summarize"
  id      = "%[2]s"
  enabled = %[3]t
}`, resourceName,
		resourceModel.id,
		resourceModel.enabled)
}

// Test that the expected attributes are set on the PingDirectory server
func testAccCheckExpectedSummarizeLogFileRotationListenerAttributes(config summarizeLogFileRotationListenerTestModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LogFileRotationListenerApi.GetLogFileRotationListener(ctx, config.id).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		resourceType := "Summarize Log File Rotation Listener"
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled",
			config.enabled, response.SummarizeLogFileRotationListenerResponse.Enabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSummarizeLogFileRotationListenerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.LogFileRotationListenerApi.GetLogFileRotationListener(ctx, testIdSummarizeLogFileRotationListener).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Summarize Log File Rotation Listener", testIdSummarizeLogFileRotationListener)
	}
	return nil
}