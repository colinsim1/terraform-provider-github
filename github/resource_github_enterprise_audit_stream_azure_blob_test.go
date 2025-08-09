package github

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGithubEnterpriseAuditStreamAzureBlob(t *testing.T) {
	keyID, keyIDexists := os.LookupEnv("GITHUB_AUDIT_STREAM_KEY_ID")
	encryptedSASUrl, encryptedSASUrlexists := os.LookupEnv("GITHUB_AUDIT_STREAM_ENCRYPTED_SAS_URL")

	t.Run("creates and imports an enterprise audit stream (Azure Blob) without error", func(t *testing.T) {

		if !keyIDexists {
			t.Skipf("%s environment variable is missing", keyID)
		}

		if !encryptedSASUrlexists {
			t.Skipf("%s environment variable is missing", encryptedSASUrl)
		}

		config := fmt.Sprintf(`
		    data "github_enterprise" "enterprise" {
			  slug = "%s"
		    }

		    resource "github_enterprise_audit_stream_azure_blob" "audit_stream" {
		      enterprise_slug = data.github_enterprise.enterprise.slug
		      enabled = true
		      key_id = "%s"
		      encrypted_sas_url = "%s"
		    }
		`, testEnterprise, keyID, encryptedSASUrl)

		checks := map[string]resource.TestCheckFunc{
			"after_create": resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(
					"github_enterprise_audit_stream_azure_blob.test", "enterprise_slug", enterprise,
				),
				resource.TestCheckResourceAttr(
					"github_enterprise_audit_stream_azure_blob.test", "enabled", "true",
				),
				// key_id should be set (Computed) after Create
				resource.TestCheckResourceAttrSet(
					"github_enterprise_audit_stream_azure_blob.test", "key_id",
				),
				// id should be set and look like "<enterprise>/<numeric_id>"
				resource.TestCheckResourceAttrSet(
					"github_enterprise_audit_stream_azure_blob.test", "id",
				),
			),
		}

		// Helper to derive import ID: "<enterprise_slug>/<stream_id>"
		importID := func(s *terraform.State) (string, error) {
			rs, ok := s.RootModule().Resources["github_enterprise_audit_stream_azure_blob.test"]
			if !ok {
				return "", fmt.Errorf("resource not found in state")
			}
			id := rs.Primary.ID
			if id == "" {
				return "", fmt.Errorf("id not set in state")
			}
			// Provider already uses "<enterprise>/<id>" format; just return it
			return id, nil
		}

		testCase := func(t *testing.T, mode string) {
			resource.Test(t, resource.TestCase{
				PreCheck:  func() { /* add any global prechecks here, e.g. auth */ },
				Providers: testAccProviders, // your provider map
				Steps: []resource.TestStep{
					{
						Config: config,
						Check:  checks["after_create"],
					},
					{
						ResourceName:      "github_enterprise_audit_stream_azure_blob.test",
						ImportState:       true,
						ImportStateIdFunc: importID,
						ImportStateVerify: true,
						// plaintext input should not be returned; ignore it during import verify
						ImportStateVerifyIgnore: []string{"encrypted_sas_url"},
					},
				},
			})
		}

		t.Run("with an enterprise account", func(t *testing.T) {
			testCase(t, enterprise)
		})
	})
}
