package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) datacenterListDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.datacenterListDataSourceRead,
		Schema: map[string]*schema.Schema{
			"datacenters": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "ID of the datacenter.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the datacenter.",
							Computed:    true,
						},
					},
				},
			},
		},
		Description: `This data source retrieves a list of datacenters`,
	}
}

func (p *provider) datacenterListDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	allDatacenters, err := client.ListDatacenters()
	if err != nil {
		return errorToDiags("list all datacenters", err)
	}
	datacenterList := make([]map[string]interface{}, 0)
	for _, datacenter := range allDatacenters {
		datacenterMap := make(map[string]interface{}, 0)
		datacenterMap["id"] = string(datacenter.ID())
		datacenterMap["name"] = datacenter.Name()
		datacenterList = append(datacenterList, datacenterMap)
	}
	if err := data.Set("datacenters", datacenterList); err != nil {
		return diag.FromErr(err)
	}
	data.SetId("datacenters")
	return nil
}
