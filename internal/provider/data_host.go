package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataHost() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"host_id", "unique_key"},
				Description:  "Identifier of the host to fetch.",
			},
			"unique_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"host_id", "unique_key"},
				Description:  "Unique key used to enforce host uniqueness within a namespace.",
			},
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Ed25519 public key in hex format for signing requests.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Optional:    true,
				Description: "Labels attached to the host.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		ReadContext: dataHostRead,
	}
}

func dataHostRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*grantoryClient)
	hostID := d.Get("host_id").(string)
	uniqueKey := d.Get("unique_key").(string)

	if hostID == "" && uniqueKey == "" {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "either host_id or unique_key must be specified",
		}}
	}
	if hostID != "" && uniqueKey != "" {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "cannot specify both host_id and unique_key",
		}}
	}

	var host apiHost
	if hostID != "" {
		var err error
		host, err = client.GetHost(ctx, hostID)
		if err != nil {
			if isNotFound(err) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}
	} else {
		hosts, err := client.ListHosts(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		found := false
		for _, h := range hosts {
			if h.UniqueKey == uniqueKey {
				host = h
				found = true
				break
			}
		}
		if !found {
			d.SetId("")
			return nil
		}
	}

	d.SetId(host.ID)
	if err := d.Set("host_id", host.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("unique_key", host.UniqueKey); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("public_key", host.PublicKey); err != nil {
		return diag.FromErr(err)
	}
	if host.Labels != nil {
		if err := d.Set("labels", flattenStringMap(host.Labels)); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}
