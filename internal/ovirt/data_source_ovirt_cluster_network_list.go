package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/yosefsatrioaji/go-ovirt-client/v3"
)

func (p *provider) clusterNetworkListDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.clusterNetworkListDataSourceRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
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
							Description: "ID of the network",
						},
						"required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the network is required",
						},
					},
				},
			},
		},
	}
}

func (p *provider) clusterNetworkListDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	clusterID := data.Get("cluster_id").(string)
	clusterNetworks, err := client.ClusterNetworkList(ovirtclient.ClusterID(clusterID))
	if err != nil {
		return diag.FromErr(err)
	}

	networks := make([]interface{}, len(clusterNetworks))
	for i, network := range clusterNetworks {
		networks[i] = map[string]interface{}{
			"id":       network.NetworkID(),
			"required": network.Required(),
		}
	}

	data.Set("networks", networks)
	data.SetId(clusterID)
	return nil
}
