package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) vnicListDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.vnicListDataSourceRead,
		Schema: map[string]*schema.Schema{
			"vnic_list": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "ID of the VNIC.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the VNIC.",
							Computed:    true,
						},
						"network_id": {
							Type:        schema.TypeString,
							Description: "ID of the network to which the VNIC is attached.",
							Computed:    true,
						},
					},
				},
			},
		},
		Description: `This data source retrieves a list of VNICs`,
	}
}

func (p *provider) vnicListDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	allVNICs, err := client.ListVNICProfiles()
	if err != nil {
		return errorToDiags("list all vNICs", err)
	}
	vnicList := make([]map[string]interface{}, 0)
	for _, vnic := range allVNICs {
		vnicMap := make(map[string]interface{}, 0)
		vnicMap["id"] = string(vnic.ID())
		vnicMap["name"] = vnic.Name()
		vnicMap["network_id"] = string(vnic.NetworkID())
		vnicList = append(vnicList, vnicMap)
	}

	if err := data.Set("vnic_list", vnicList); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("vnic_list")

	return nil
}
