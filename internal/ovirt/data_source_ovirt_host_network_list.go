package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/yosefsatrioaji/go-ovirt-client/v3"
)

func (p *provider) hostNetworkListDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.hostNetworkListDataSourceRead,
		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the network attachment",
						},
						"network_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the network",
						},
						"host_nic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the host NIC",
						},
					},
				},
			},
		},
	}
}

func (p *provider) hostNetworkListDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	hostID := data.Get("host_id").(string)
	hostNetworks, err := client.NetworkAttachmentList(ovirtclient.HostID(hostID))
	if err != nil {
		return diag.FromErr(err)
	}

	networks := make([]interface{}, len(hostNetworks))
	for i, network := range hostNetworks {
		networks[i] = map[string]interface{}{
			"id":          network.ID(),
			"network_id":  network.NetworkID(),
			"host_nic_id": network.HostNICID(),
		}
	}

	data.Set("networks", networks)
	data.SetId(hostID)
	return nil
}
