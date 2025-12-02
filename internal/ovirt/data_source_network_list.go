package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) networkListDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.networkListDataSourceRead,
		Schema: map[string]*schema.Schema{
			"networks": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "ID of the network.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the network.",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of the network.",
							Computed:    true,
						},
						"comment": {
							Type:        schema.TypeString,
							Description: "Comment of the network.",
							Computed:    true,
						},
						"vlan_id": {
							Type:        schema.TypeInt,
							Description: "VLAN ID of the network.",
							Computed:    true,
						},
						"dc_id": {
							Type:        schema.TypeString,
							Description: "ID of the datacenter to which the network belongs.",
							Computed:    true,
						},
					},
				},
			},
		},
		Description: `This data source retrieves a list of networks`,
	}
}

func (p *provider) networkListDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	allNetworks, err := client.ListNetworks()
	if err != nil {
		return errorToDiags("list all networks", err)
	}
	networkList := make([]map[string]interface{}, 0)
	for _, network := range allNetworks {
		networkMap := make(map[string]interface{}, 0)
		networkMap["id"] = string(network.ID())
		networkMap["name"] = network.Name()
		networkMap["description"] = network.Description()
		networkMap["comment"] = network.Comment()
		networkMap["vlan_id"] = network.VlanID()
		networkMap["dc_id"] = string(network.DatacenterID())
		networkList = append(networkList, networkMap)
	}
	if err := data.Set("networks", networkList); err != nil {
		return diag.FromErr(err)
	}
	data.SetId("networks")
	return nil
}
