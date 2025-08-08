package github

import (
	"context"

	"github.com/google/go-github/v66/github"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceGitHubEnterpriseAuditStreamAzureBlobStorage() *schema.Resource {
	r := genBaseGitHubAuditStreamResource(flattenAuditStreamAzureBlobStorage, expandAuditStreamAzureBlobStorage)

	r.Schema["key_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Key ID obtained from the audit log stream key endpoint used to encrypt secrets.",
	}

	r.Schema["encrypted_sas_url"] = &schema.Schema{
		Type:     schema.TypeString,
		ForceNew: true,

		Required:     true,
		ValidateFunc: validation.IsURLWithHTTPS,
		Description:  "The URL of the Azure Storage Blob Container",
	}
	return r
}

func expandAuditStreamAzureBlobStorage(d *schema.ResourceData, m interface{}) (*github.AuditStream, string) {
	auditStream, enterprise := doBaseExpansion(d)
	ctx := context.Background()
	client := m.(*Owner).v3client
	key, _, err := client.Enterprise.GetAuditStreamKey(ctx, enterprise)
	if err != nil {
		return auditStream, "failed"
	}
	if key == nil || key.KeyID == "" {
		return auditStream, "failed"
	}
	auditStream.Enabled = d.Get("enabled").(bool)
	auditStream.StreamType = "Azure Blob Storage"
	auditStream.VendorSpecific = &map[string]string{
		"key_id":            key.KeyID,
		"encrypted_sas_url": d.Get("encrypted_sas_url").(string),
	}
	return auditStream, enterprise
}

func flattenAuditStreamAzureBlobStorage(d *schema.ResourceData, enterprise *string, auditStream *github.AuditStream) {
	doBaseFlattening(d, auditStream, enterprise)

	d.Set("key_id", (*auditStream.VendorSpecific)["key_id"])
	d.Set("encrypted_sas_url", (*auditStream.VendorSpecific)["encrypted_sas_url"])
}
