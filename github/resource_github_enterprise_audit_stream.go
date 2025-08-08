package github

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-github/v66/github"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type flatFunc func(d *schema.ResourceData, enterprise *string, auditStream *github.AuditStream)
type expandFunc func(d *schema.ResourceData) (*github.AuditStream, string)

func genBaseGitHubAuditStreamResource(f flatFunc, e expandFunc) *schema.Resource {
	return &schema.Resource{
		Create: resourceGitHubEnterpriseAuditStreamCreate(f, e),
		Read:   resourceGitHubEnterpriseAuditStreamRead(f, e),
		// Update: resourceGitHubEnterpriseAuditStreamUpdate,z
		Delete: resourceGitHubEnterpriseAuditStreamDelete(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"enterprise_slug": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
				Description: "Whether the audit stream is enabled. " +
					"Defaults to true. If disabled, the stream will be created but paused.",
			},
		},
	}
}

// doBaseExpansion performs the expansion for the 'base' attributes that are defined in the schema, above
func doBaseExpansion(d *schema.ResourceData) (*github.AuditStream, string) {
	var auditStreamId *int
	parsedId, err := strconv.Atoi(d.Id())
	if err == nil {
		auditStreamId = &parsedId
	}

	enterprise := d.Get("enterprise_slug").(string)

	auditStream := &github.AuditStream{
		ID: auditStreamId,
	}

	return auditStream, enterprise
}

// doBaseFlattening performs the flattening for the 'base' attributes that are defined in the schema, above
func doBaseFlattening(d *schema.ResourceData, auditStream *github.AuditStream, enterprise *string) {
	d.SetId(strconv.Itoa(*auditStream.ID))
	d.Set("enterprise", enterprise)
}

func resourceGitHubEnterpriseAuditStreamCreate(flat flatFunc, expand expandFunc) func(d *schema.ResourceData, m interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		ctx := context.Background()
		client := m.(*Owner).v3client

		auditStream, enterprise := expand(d)

		out, _, err := client.Enterprise.CreateAuditStream(ctx, enterprise, auditStream)
		if err != nil {
			return fmt.Errorf("create audit log stream: %w", err)
		}

		d.SetId(fmt.Sprintf("%s/%d", enterprise, out.ID))
		return resourceGitHubEnterpriseAuditStreamRead(flat, expand)(d, m)
	}
}

func resourceGitHubEnterpriseAuditStreamRead(flat flatFunc, expand expandFunc) func(d *schema.ResourceData, m interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		ctx := context.Background()
		client := m.(*Owner).v3client

		streamId, err := strconv.Atoi(d.Id())
		if err != nil {
			return fmt.Errorf("Error parsing the audit stream ID from the Terraform resource data: %v", err)
		}

		enterprise := d.Get("enterprise_slug").(string)

		auditStream, _, err := client.Enterprise.GetAuditStream(ctx, enterprise, streamId)
		if err != nil {
			if ghErr, ok := err.(*github.ErrorResponse); ok && ghErr.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("read audit log stream: %w", err)
		}
		if auditStream == nil || auditStream.ID == nil {
			d.SetId("")
			return nil
		}

		flat(d, &enterprise, auditStream)

		return nil
	}
}

func resourceGitHubEnterpriseAuditStreamDelete() schema.DeleteFunc {
	return func(d *schema.ResourceData, m interface{}) error {
		ctx := context.Background()
		client := m.(*Owner).v3client

		streamId, err := strconv.Atoi(d.Id())
		if err != nil {
			return err
		}

		enterprise := d.Get("enterprise_slug").(string)

		_, err = client.Enterprise.DeleteAuditStream(ctx, enterprise, streamId)
		if err != nil {
			if ghErr, ok := err.(*github.ErrorResponse); ok && ghErr.Response.StatusCode == 404 {
				return nil
			}
			return fmt.Errorf("error deleting audit log stream: %w", err)
		}
		return nil
	}
}
