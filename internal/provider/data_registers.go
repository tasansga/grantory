package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataRegisters() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels that each returned register entry must include.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"registers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Register entries that matched the supplied filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"register_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		ReadContext: dataRegistersRead,
	}
}

func dataRegistersRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*grantoryClient)
	opts := registerListOptions{
		Labels: expandStringMap(extractMap(d.Get("labels"))),
	}

	registers, err := client.listRegisters(ctx, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	values := make([]map[string]any, 0, len(registers))
	for _, reg := range registers {
		entry := map[string]any{
			"register_id": reg.ID,
			"host_id":     reg.HostID,
		}
		values = append(values, entry)
	}

	if err := d.Set("registers", values); err != nil {
		return diag.FromErr(err)
	}
	id, err := hashAsJSON(map[string]any{
		"labels":    opts.Labels,
		"registers": values,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return nil
}
